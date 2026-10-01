# Web Restic 系统概要设计

> 版本：v0.1（设计稿）
> 配套文档：《需求分析设计》`requirements-design.md`、《UI 原型设计》`ui-prototype.md`

## 1. 设计目标与原则

1. **轻量单实例**：后端单进程（Go Fiber）+ SQLite3，Docker Compose 一键部署，无外部中间件依赖。
2. **零安装目标主机**：restic 程序由服务端按平台托管上传，目标主机仅需 SSH 可达。
3. **异步执行**：备份/恢复是长任务，全部走"提交 → 执行队列 → 执行记录 → 日志推送"，HTTP 不阻塞。
4. **延续既有约定**：handler 子包装配、`auth.Require` 权限、POST-only 业务路由、`{ data: ... }` 响应、zap 日志、Swagger 注释、单数表名、`BaseModel` 继承、`CGO_ENABLED=0` 兼容 SQLite。

## 2. 总体架构

```
┌────────────────────────────────────────────────────────────────────┐
│                           前端 (Vue3 + antdv / vben)                 │
│   工作台 │ 仓库管理 │ 主机管理 │ 备份任务 │ 备份数据 │ 系统管理          │
└───────────────┬────────────────────────────────────────────────────┘
                │  POST /api/**（JWT Bearer，SSE 仅用于日志流）
┌───────────────▼────────────────────────────────────────────────────┐
│                         后端 (Go Fiber 单实例)                       │
│  ┌───────────────┐  ┌────────────────────────────────────────────┐  │
│  │ handler 装配层  │  │  basic / user / sys / restic(新增) 子包      │  │
│  └───────┬───────┘  └──────────────────────┬─────────────────────┘  │
│          │                                 │                        │
│  ┌───────▼─────────────────────────────────▼─────────────────────┐  │
│  │                    service 业务层                               │  │
│  │  UserService  RepoService  HostService  TaskService             │  │
│  │  SnapshotService  RestoreService                                │  │
│  └───────┬───────────────────────────┬───────────────────────────┘  │
│          │                           │                              │
│  ┌───────▼────────────┐   ┌──────────▼───────────────────────────┐  │
│  │ 基础设施             │   │ 执行引擎 (executor)                    │  │
│  │ auth(JWT) config    │   │  cron 调度器 → 任务队列 → SSH/SFTP 执行器 │  │
│  │ log(zap) db_(dao)   │   │  restic 配置生成器  二进制管理器(缓存/下载) │  │
│  │ crypto(凭据加密)      │   │  日志采集器 → 落库 + SSE 推送             │  │
│  └───────┬────────────┘   └───────────────────────────────────────┘  │
│          │                                                          │
│  ┌───────▼───────────────────────────────────────────────────────┐  │
│  │                SQLite (data/db.sqlite)                         │  │
│  │  user  sys_setting  backup_repo  host  backup_task             │  │
│  │  task_run  restore_record                                      │  │
│  └────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
                │ SSH/SFTP（执行 restic 命令、上传程序与配置）
┌───────────────▼──────────────────────────────┐   ┌──────────────────┐
│              目标主机 (多台)                    │   │   存储仓库 repo    │
│  restic backup / restore / snapshots ...     │──▶│ local/s3/sftp/... │
└──────────────────────────────────────────────┘   └──────────────────┘
```

数据流向：**目标主机 → 仓库**（备份）与 **仓库 → 目标主机**（恢复）由 restic 自身完成，服务端只负责编排、传输程序/配置与采集结果。

## 3. 技术选型

| 层 | 选型 | 说明 |
| --- | --- | --- |
| 后端 | Go + Fiber v3 | 轻量 Web 框架，单实例 |
| ORM | GORM | SQLite 驱动（`modernc.org/sqlite` 等，支持 `CGO_ENABLED=0`） |
| 数据库 | SQLite3 | `./data/db.sqlite`，单文件持久化 |
| 认证 | JWT（HS256） | 沿用 `internal/auth`，默认 7 天，`JWT_SECRET` 环境注入 |
| SSH/SFTP | `golang.org/x/crypto/ssh` + SFTP 子协议 | 远程执行、文件上传 |
| 调度 | `robfig/cron/v3` | 任务 cron 调度（或等价轻量实现） |
| 日志 | zap | 终端紧凑彩色前缀 + `./logs/app.log` 滚动 JSON |
| 前端 | Vue3 + TS + antdv（vben 工程） | 已有骨架，新增业务视图 |
| 部署 | Docker Compose | backend + front + 数据卷 |

## 4. 后端模块划分（目录结构）

在现有结构上扩展：

```
backend/internal/
├── auth/            # JWT（已有，不动）
├── base/            # db_/log_/http_(SSE)/crypto_(新增凭据加密) 等（已有 + 扩展）
├── config/          # 配置加载（已有，新增 SecretKey 等字段）
├── handler/
│   ├── app.go       # 统一装配（新增 restic 子包注册）
│   ├── basic/       # 基础（已有）
│   ├── user/        # 认证与用户（已有）
│   ├── sys/         # 系统（已有）
│   └── restic/      # 新增：备份域 handler
│       ├── repo-handler.go      + repo-handler_test.go
│       ├── host-handler.go      + host-handler_test.go
│       ├── task-handler.go      + task-handler_test.go
│       ├── backup-handler.go    + backup-handler_test.go   # 快照/恢复/清理
│       └── run-handler.go       + run-handler_test.go      # 执行记录/日志
├── models/          # 新增：repo.go host.go task.go task_run.go restore_record.go
├── service/         # 新增：repo_service.go host_service.go task_service.go
│                    #       snapshot_service.go restore_service.go
├── executor/        # 新增：执行引擎
│   ├── ssh.go       # SSH 会话、SFTP 上传、远程命令执行
│   ├── config_gen.go# 按仓库类型生成 restic env/参数
│   ├── binary.go    # restic 程序下载与本地缓存（platform-arch）
│   ├── scheduler.go # cron 调度器 + 任务队列 + 主机互斥锁
│   └── runner.go    # 单次备份/恢复执行器（采集日志、写执行记录、SSE 推送）
└── svc/             # ServiceContext（扩展：SecretKey、Executor、Scheduler）
```

> 说明：`AGENTS.md` 目前列出的 handler 子包为 basic/user/sys（通用基础域）；备份是业务域，新增 `restic` 子包属于对既有约定的扩展，待评审确认后回写规范。

## 5. 数据库设计

> 约定：新表使用**单数表名**（显式 `TableName()`，同 `sys_setting`）；模型继承 `db_.BaseModel`（`id/created_at/updated_at`）+ `db_.LogicDelete`；时间戳为毫秒 int64；JSON 字段以 text 存储。

### 5.1 backup_repo（存储仓库）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | int64 PK | 主键 |
| name | varchar(100) | 仓库名称（唯一） |
| repo_type | varchar(50) | 后端类型：local / sftp / s3 / minio / azure / b2 / rest-server |
| address | varchar(500) | 仓库地址/URL，如 `/data/backup`、`s3:mybucket:/path`、`sftp:user@host:/backup` |
| bucket | varchar(200) | S3 bucket（按类型可选） |
| username | varchar(200) | 访问用户名 / S3 AccessKey（按类型可选） |
| password | text | 访问密钥，**AES 加密**（SecretKey、SFTP 密码等） |
| repo_password | text | restic 仓库密码，**AES 加密** |
| options | text | 附加参数 JSON：region、endpoint、hardlink 等 |
| status | varchar(20) | ok / error / unknown |
| last_check_at | int64 | 最近连接测试时间 |
| remark | varchar(500) | 备注 |
| created_at / updated_at / deleted_at | int64 | 审计与软删 |

### 5.2 host（目标主机）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | int64 PK | 主键 |
| name | varchar(100) | 主机名称（唯一） |
| hostname | varchar(200) | IP 或域名 |
| port | int | SSH 端口，默认 22 |
| username | varchar(100) | SSH 用户 |
| auth_type | varchar(20) | password / private_key |
| password | text | SSH 密码，AES 加密（auth_type=password） |
| private_key | text | 私钥内容，AES 加密（auth_type=private_key） |
| passphrase | text | 私钥口令，AES 加密（可选） |
| platform | varchar(20) | linux / windows / darwin / freebsd |
| arch | varchar(20) | amd64 / arm64 / 386 |
| restic_dir | varchar(500) | 主机上 restic 程序目录（默认 `~/.restic-web`） |
| work_dir | varchar(500) | 远程工作/临时目录 |
| status | varchar(20) | ok / error / unknown |
| last_check_at | int64 | 最近连通测试时间 |
| remark | varchar(500) | 备注 |

### 5.3 backup_task（备份任务）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | int64 PK | 主键 |
| name | varchar(100) | 任务名称 |
| host_id | bigint | 关联主机 |
| repo_id | bigint | 关联仓库 |
| source_paths | text | 备份源路径 JSON 数组，如 `["/etc","/var/www"]` |
| exclude_paths | text | 排除路径 JSON 数组 |
| exclude_file | varchar(500) | 排除规则文件（主机上的路径，可选） |
| schedule | varchar(100) | cron 表达式；空 = 仅手动 |
| enabled | bool | 是否启用调度 |
| keep_policy | text | 保留策略 JSON：`{keepLast, keepHourly, keepDaily, keepWeekly, keepMonthly, ...}` |
| tags | text | 标签 JSON 数组，透传 `--tag` |
| extra_args | text | restic 附加参数 JSON 数组 |
| last_status | varchar(20) | 最近执行结果：none / success / failed / running |
| last_run_at | int64 | 最近执行时间 |
| remark | varchar(500) | 备注 |

### 5.4 task_run（执行记录，备份与恢复共用）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | int64 PK | 主键 |
| task_id | bigint | 来源任务（恢复可为空） |
| run_type | varchar(20) | backup / restore |
| host_id / repo_id | bigint | 执行环境快照 |
| status | varchar(20) | pending / running / success / failed / canceled |
| snapshot_id | varchar(100) | restic 快照 ID（备份成功时写入） |
| trigger | varchar(20) | manual / cron |
| log_path | varchar(500) | 执行日志文件相对路径 |
| started_at / finished_at | int64 | 起止时间 |
| error_message | text | 失败原因 |

### 5.5 restore_record（恢复记录）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | int64 PK | 主键 |
| repo_id / host_id | bigint | 仓库与目标主机 |
| snapshot_id | varchar(100) | 恢复的快照 ID |
| source_path | varchar(500) | 恢复源路径（空 = 整个快照） |
| target_path | varchar(500) | 目标主机恢复路径 |
| status | varchar(20) | running / success / failed |
| files_restored | int64 | 恢复文件数（可空） |
| bytes_restored | int64 | 恢复字节数（可空） |
| log_path | varchar(500) | 日志路径 |
| started_at / finished_at | int64 | 起止时间 |
| error_message | text | 失败原因 |

### 5.6 索引与约束

- `backup_repo.name`、`host.name`、`backup_task.name` 唯一索引。
- `task_run(task_id, started_at)`、`restore_record(repo_id, started_at)` 组合索引，供历史查询。
- 外键逻辑关联（不强制 DB 级 FK，沿用项目现有做法），删除保护在 service 层校验。

## 6. 接口设计（API）

> 约定：业务接口统一 `POST` + 动作子路径；响应统一 `{ data: ... }`；权限用 `auth.Require` 声明（guest/login/admin/write/read）；每个 handler 同步维护 Swagger 注释并执行 `swag init`。

### 6.1 认证与用户（已有）

| 接口 | 权限 | 说明 |
| --- | --- | --- |
| POST /api/auth/register | guest | 注册（首个用户为 admin） |
| POST /api/auth/login | guest | 登录，返回 `{token, expiresAt, user}` |
| POST /api/users/create | admin | 创建用户 |
| POST /api/users/list | read | 用户列表 |
| POST /api/users/:id/detail | read | 用户详情 |
| POST /api/users/:id/update | write | 更新（非 admin 仅限自己且不可改角色） |
| POST /api/users/:id/delete | admin | 删除用户 |

### 6.2 存储仓库（新增 `internal/handler/restic/repo-handler.go`）

| 接口 | 权限 | 说明 |
| --- | --- | --- |
| POST /api/repos/create | write | 创建仓库 |
| POST /api/repos/list | read | 仓库列表（分页/筛选） |
| POST /api/repos/:id/detail | read | 仓库详情（凭据掩码） |
| POST /api/repos/:id/update | write | 更新（凭据字段留空 = 不修改） |
| POST /api/repos/:id/delete | admin | 删除（有快照/运行中任务时拒绝） |
| POST /api/repos/:id/test | write | 连接测试（返回结果与错误信息） |

请求体示例（create/update）：

```json
{
  "name": "prod-s3",
  "repoType": "minio",
  "address": "s3:backup-bucket:/web",
  "bucket": "backup-bucket",
  "username": "minioadmin",
  "password": "********",
  "repoPassword": "repo-pass-123",
  "options": {"endpoint": "http://minio:9000", "region": "us-east-1"},
  "remark": "生产 MinIO"
}
```

### 6.3 主机管理（新增 `internal/handler/restic/host-handler.go`）

| 接口 | 权限 | 说明 |
| --- | --- | --- |
| POST /api/hosts/create | write | 创建主机 |
| POST /api/hosts/list | read | 主机列表（分页/筛选） |
| POST /api/hosts/:id/detail | read | 详情（凭据掩码） |
| POST /api/hosts/:id/update | write | 更新 |
| POST /api/hosts/:id/delete | admin | 删除（被任务引用时拒绝） |
| POST /api/hosts/:id/test | write | SSH 连通测试，返回平台探测结果 |

### 6.4 备份任务（新增 `internal/handler/restic/task-handler.go`）

| 接口 | 权限 | 说明 |
| --- | --- | --- |
| POST /api/tasks/create | write | 创建任务 |
| POST /api/tasks/list | read | 任务列表（分页/筛选/最近状态） |
| POST /api/tasks/:id/detail | read | 任务详情（含关联主机/仓库摘要） |
| POST /api/tasks/:id/update | write | 更新 |
| POST /api/tasks/:id/delete | admin | 删除（有运行中执行时拒绝） |
| POST /api/tasks/:id/enable | write | 启用调度 |
| POST /api/tasks/:id/disable | write | 停用调度 |
| POST /api/tasks/:id/run | write | 手动触发一次备份 |
| POST /api/tasks/:id/runs | read | 该任务的执行记录列表 |

### 6.5 备份数据（新增 `internal/handler/restic/backup-handler.go`、`run-handler.go`）

| 接口 | 权限 | 说明 |
| --- | --- | --- |
| POST /api/snapshots/list | read | 快照列表（repoId/hostId/时间范围筛选） |
| POST /api/snapshots/tree | read | 快照文件树（snapshotId + path，分页目录项） |
| POST /api/snapshots/restore | write | 提交恢复（snapshotId、sourcePath、hostId、targetPath），返回 runId |
| POST /api/snapshots/forget | write | 清理快照（snapshotIds 或按保留策略），返回 runId |
| POST /api/runs/list | read | 执行记录列表（taskId/runType/status 筛选） |
| POST /api/runs/:id/detail | read | 执行记录详情 |
| POST /api/runs/:id/logs | read | 执行日志（分页文本） |
| POST /api/runs/:id/stream | login | **SSE 实时日志流**（响应 `text/event-stream`；沿用 base/http_/sse.go 的 SseMessage 结构；属特殊流式接口，仍走 POST） |

### 6.6 通用响应/错误

```json
// 成功
{ "data": { ... } }
// 失败
{ "error": "错误信息" }
```

## 7. 核心模块设计

### 7.1 凭据加密（base/crypto_，新增）

- 对称加密 AES-256-GCM，密钥派生自环境变量 `SECRET_KEY`（base64，≥32 字节）；缺失时 dev 模式用内置默认值并打 WARN 日志，生产模式直接报错退出。
- 落库前加密、出库时解密；所有 handler 输出统一掩码（如 `******`）；日志记录器统一过滤敏感字段。

### 7.2 restic 配置生成（executor/config_gen.go）

按 `repo_type` 渲染 restic 环境变量到临时 `.env` 文件（随程序上传）：

| repo_type | 关键 env |
| --- | --- |
| local | `RESTIC_REPOSITORY=<address>` |
| sftp | `RESTIC_REPOSITORY=sftp:<user>@<host>:<path>` |
| s3 / minio | `RESTIC_REPOSITORY=s3:<bucket>:/<path>`、`AWS_ACCESS_KEY_ID`、`AWS_SECRET_ACCESS_KEY`、`AWS_ENDPOINT`、`AWS_DEFAULT_REGION` |
| azure | `RESTIC_REPOSITORY=azure:<container>:/<path>`、`AZURE_ACCOUNT_NAME`、`AZURE_ACCOUNT_KEY` |
| b2 | `RESTIC_REPOSITORY=b2:<bucket>:<path>`、`B2_ACCOUNT_ID`、`B2_ACCOUNT_KEY` |
| rest-server | `RESTIC_REPOSITORY=rest:<url>`、`RESTIC_REST_USERNAME`、`RESTIC_REST_PASSWORD` |

通用：`RESTIC_PASSWORD=<repo_password>`（或写入同目录 password 文件 + `RESTIC_PASSWORD_FILE`，更安全）。执行命令形如：

```bash
restic --repo <env> backup <src...> [--exclude ...] [--tag ...] [extra args]
restic --repo <env> restore <snapshot> --target <path> [--path <src>]
restic --repo <env> snapshots / ls <snapshot> / forget --keep-* / prune
```

### 7.3 restic 二进制管理（executor/binary.go）

- 本地缓存目录：`./data/restic/<platform>-<arch>/restic[.exe]`。
- 首次使用按 `host.platform + host.arch` 检查缓存，未命中则从配置的下载 URL 模板下载（`sys_setting` 键如 `restic.downloadUrl`，支持 `{version}/{platform}/{arch}` 占位），并校验 SHA256（`restic.downloadSha256`）。
- windows 平台二进制上传后使用 `.exe` 后缀；执行命令经 SSH shell 组装。

### 7.4 SSH/SFTP 执行器（executor/ssh.go）

- 会话复用：按主机维度维护连接池（LRU，闲置超时关闭）。
- `Exec`：打开 SSH session，`CombinedOutput`/逐行读取 stdout+stderr，回调把每行交给日志采集器。
- `Upload`：SFTP 上传文件（设置权限 0700/0600），上传前校验目标目录存在（不存在则 `mkdir -p`）。
- 平台差异：windows 主机命令经 `cmd /C` 或 PowerShell 包装，路径分隔符处理。

### 7.5 任务调度与执行队列（executor/scheduler.go、runner.go）

- 服务启动：加载所有 `enabled=true` 且有 `schedule` 的任务注册 cron；`cron` 命中与手动 `run` 统一提交到任务队列（buffered channel）。
- Worker 池（默认 1~2，可配置）：取任务 → 获取 `hostLock`（`sync.Map` 按 host_id 加锁，可重入拒绝）→ 执行备份流程 → 释放。
- 执行器状态：运行中写入 `task_run(status=running)`；完成后更新 `task_run` 与 `backup_task.last_status/last_run_at`；异常时写入 `error_message` 并保持任务启用。
- 重启恢复：启动时将残留 `running` 记录标记为 `failed`（原因：服务重启中断）。
- SSE 广播：`runner` 把日志行同时写入日志文件并 publish 到 `runId` 主题；前端 `POST /api/runs/:id/stream` 订阅，连接断开不影响执行。

### 7.6 快照与文件树（service/snapshot_service.go）

- 快照列表：`restic snapshots --json` 解析输出（ID、时间、主机、标签、路径、大小、文件数）。
- 文件树：`restic ls <snapshot> --json` 按 `--path` 过滤，服务端组织为树节点（懒加载目录项）。
- 快照与文件查询需要仓库连接：优先复用最近一次的 `task_run` 环境（主机+仓库+临时配置），未找到则临时生成。

## 8. 日志设计

- 系统日志：沿用 zap 双输出（终端紧凑彩色 + `./logs/app.log` 滚动 JSON），新增模块沿用 `log_.Info/Error` 等。
- 执行日志：按 `logs/runs/<runId>.log` 独立归档，同时内存保留最近 N 行供 SSE 追发；写入时统一脱敏（替换密码字段值）。

## 9. 部署设计（Docker Compose）

```yaml
services:
  backend:
    build: ./backend
    environment:
      - MODE=prod
      - JWT_SECRET=${JWT_SECRET}
      - SECRET_KEY=${SECRET_KEY}        # 凭据加密密钥
      - RESTIC_DOWNLOAD_URL=${RESTIC_DOWNLOAD_URL}
    volumes:
      - ./data:/app/data                # SQLite + restic 缓存
      - ./logs:/app/logs
    ports: ["3000:3000"]
  front:
    build: ./front
    ports: ["80:80"]
    depends_on: [backend]               # nginx 反代 /api → backend:3000
```

- 生产配置 `etc/prod.yaml` 扩展 `secretKey`、`restic.downloadUrl` 等字段（配置加载沿用 dev.yaml 兜底 + 模式覆盖）。
- 数据卷持久化 `data/`（db.sqlite + restic 二进制缓存）与 `logs/`。

## 10. 关键风险与应对

| 风险 | 应对 |
| --- | --- |
| 目标主机网络/SSH 不稳定导致执行失败 | 超时与重试策略（连接超时、命令超时）、失败记录可查、任务可手动重跑 |
| restic 二进制下载失败或平台缺失 | 下载失败重试；支持内网镜像 URL；`sys_setting` 可配版本与校验和 |
| 凭据泄露 | AES 加密落库、日志脱敏、前端掩码、生产密钥环境注入 |
| 大文件树/大量快照查询慢 | 文件树懒加载、快照列表分页、`ls --json` 流式解析 |
| 并发任务互相影响 | 主机级互斥锁；worker 池上限；任务队列缓冲 |
| 误删快照 | forget 二次确认；保留策略按配置执行；操作记录留痕 |

## 11. 实施拆分建议

1. **基础设施**：`SECRET_KEY` 加密模块、SSH/SFTP 执行器、restic 二进制管理（含下载缓存）——可独立单测（SSH 用测试桩/集成测试标记）。
2. **仓库/主机域**：模型 + service + handler + Swagger + 前端 CRUD 页面。
3. **任务域**：任务模型 + 调度器 + 执行队列 + runner + 日志采集/SSE。
4. **数据域**：快照列表/文件树/恢复/forget + 前端交互。
5. **收尾**：工作台统计、文档、Compose 验证、`AGENTS.md` 规范回写。
