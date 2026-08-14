# 项目协作规范

## 语言

- 与用户沟通、代码注释、文档、日志和面向用户的界面文案，尽量使用中文。
- 技术专有名词、代码标识符、第三方库名称和命令可保留英文，以保证准确性与可检索性。

## TDD 开发

- 新功能、缺陷修复和重构遵循测试驱动开发：先编写或更新测试，再运行测试确认 RED，最后实现最小改动并确认 GREEN。
- 测试应覆盖正常路径、边界条件和明确的错误处理；优先为扫描、数据库和交互状态等核心逻辑添加自动化测试。
- 交付前至少执行相关测试；比如对 Go 代码默认执行 `go test ./...` 和 `go vet ./...`，无法执行时应说明原因。

## 产品目标

- 开始工作前阅读 `readme.md`，将其描述的项目目标作为设计和实现的边界。

## 技术栈

- 后端：Go Fiber + Go ORM （轻量级，单实例运行即可）
- 前端：Vue 3 + TypeScript + antdv
- 数据库：SQLite3
- SQLite 驱动需支持 `CGO_ENABLED=0` 构建，确保本地和 Docker 环境无需 C 编译器即可运行。
- 部署：Docker Compose

## 日志与权限

- 后端日志统一使用 zap，同时输出到终端和 `./logs/app.log`，文件日志应支持滚动切分。
- 控制台日志使用紧凑的 `YYYY-MM-DD HH:mm:ss.SSS [LEVEL]` 前缀：时间和级别带 ANSI 颜色，WARN 为黄色、ERROR 为红色；文件日志使用相同时间格式的无 ANSI JSON。
- 受保护路由在 handler 注册时使用 `auth.Require` 声明 `guest`、`login`、`admin`、`write` 或 `read` 权限，由认证中间件统一校验 JWT。
- JWT 默认有效期为 7 天；生产环境通过 `JWT_SECRET` 提供密钥，不在代码中硬编码生产密钥。
- 前端登录与注册直接调用 `/api/auth/login`、`/api/auth/register`，并按后端 `{ data: ... }` 响应保存 JWT 和用户角色；开发代理必须保留 `/api` 前缀。

## API 文档

- 后端 Swagger UI 固定提供在 `/swagger/`，OpenAPI JSON 固定提供在 `/swagger/doc.json`。
- 新增或变更后端接口时，同步维护 handler 的 Swagger 注释并在 `backend` 目录执行 `swag init --generalInfo app.go --output docs --parseInternal`，提交生成的 `docs/` 文件。
- 除 Swagger 文档及静态资源外，前后端业务 API 统一使用 `POST`；同一资源的不同操作通过明确的动作子路径区分，不注册 `GET`、`PUT` 或 `DELETE` 业务路由。

## 规范回顾

- 每次任务收尾时，回顾用户提出的要求，识别其中可跨任务复用、稳定且不与既有规范冲突的约束。
- 将这些共性要求精炼后更新本文件；一次性的任务细节、临时偏好和未经确认的推断不写入长期规范。
- 更新后检查规则是否清晰、可执行且没有重复或互相矛盾的表述。
