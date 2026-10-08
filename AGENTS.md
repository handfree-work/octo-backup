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
- 后端 handler 按职责放在 `internal/handler/basic`、`internal/handler/user`、`internal/handler/sys` 子包，根 `handler` 包只负责统一装配。
- handler 文件按 `<module>-handler.go` 命名（例如 `user-handler.go`、`auth-handler.go`），单元测试文件使用对应的 `<module>-handler_test.go` 并与目标文件放在同一目录。
- 每个 `*_test.go` 文件必须与同目录下同名的源码文件一一对应（例如 `auth-handler_test.go` 对应 `auth-handler.go`）；测试辅助代码应放在已有的对应测试文件中，不得创建没有对应源码文件的孤立测试文件。
- 用户管理 handler 放在 `backend/internal/handler/sys/authority`，认证 handler 放在 `backend/internal/handler/basic`。
- `sys_setting` 是仅供内部使用的系统配置表，字段固定为 `id`、`key`、`setting`，不对外提供 CRUD 接口。

## 日志与权限

- 后端日志统一使用 zap，同时输出到终端和 `./logs/app.log`，文件日志应支持滚动切分。
- 控制台日志使用紧凑的 `YYYY-MM-DD HH:mm:ss.SSS [LEVEL]` 前缀：时间和级别带 ANSI 颜色，WARN 为黄色、ERROR 为红色；文件日志使用相同时间格式的无 ANSI JSON。
- 认证配置、JWT 校验、权限常量和中间件统一放在 `internal/base/web_`。
- 业务 POST 路由统一使用 `internal/base/web_` 的 `Route` 结构和 `RegisterRoutes` 注册；注册时传入认证配置，`Route.Permission` 为空时不添加权限中间件，非空时由 `web_.Require` 校验。新增路由字段应扩展通用 `Route` 结构。
- JWT 默认有效期为 7 天；生产环境通过 `JWT_SECRET` 提供密钥，不在代码中硬编码生产密钥。
- 后端关键业务流程必须记录结构化 zap 日志：至少包含开始、关键阶段、成功或失败及耗时；外部连接、文件上传、远程命令、下载和状态变更应记录可定位的上下文（ID、阶段、数量、路径等），密码、私钥、Token 等敏感信息不得写入日志。

## 错误处理

- 业务层、service、插件和基础设施在错误发生的位置立即使用 `internal/base/error_` 的 `NewTextError`、`NewCodeTextError` 或 `NewFormatError` 创建并记录 `CodedError`，错误消息必须包含足够的模块、操作和底层原因上下文。
- 查询或调用返回 `err` 时不得用“资源不存在”等新错误覆盖原始错误；如果能够确认 `err` 已经是包装过的业务错误，应直接返回 `err`，否则使用 `error_.NewApiError` 或 `error_.NewWrapError` 包装后返回，并保留原始错误及其上下文。
- handler 只负责将已有的 `CodedError` 转换为统一响应；不得在 handler 才首次包装业务错误，也不得重复记录已经编码的错误。普通 error 仅作为未覆盖路径的兜底。
- 插件 action 错误必须包含插件名和 action 名，底层连接、数据库和外部命令错误必须保留原始错误信息。
- 前端登录与注册直接调用 `/api/auth/login`、`/api/auth/register`，并按后端 `{ data: ... }` 响应保存 JWT 和用户角色；开发代理必须保留 `/api` 前缀。

## API 文档

- 后端 Swagger UI 固定提供在 `/swagger/`，OpenAPI JSON 固定提供在 `/swagger/doc.json`。
- 新增或变更后端接口时，同步维护 handler 的 Swagger 注释并在 `backend` 目录执行 `swag init --generalInfo app.go --output docs --parseInternal`，提交生成的 `docs/` 文件。
- 除 Swagger 文档及静态资源外，前后端业务 API 统一使用 `POST`；同一资源的不同操作通过明确的动作子路径区分，不注册 `GET`、`PUT` 或 `DELETE` 业务路由。
- 业务接口统一使用 `POST` 和 `application/json`；响应固定为 `{code: 0, message: "成功", data: {}}`，`code=0` 表示成功，`code>0` 表示业务失败。无论业务成功或失败，HTTP 状态码统一返回 `200`；前端统一解析 `code` 和 `message`。
- 所有业务数据必须放在 `data` 字段内；分页数据统一使用 `{offset, limit, records, total}` 结构，service 层分页查询使用 `FindPage` 方法。
- 前端业务 API 文件就近放置在对应功能模块目录（例如用户模块使用 `views/sys/authority/user/api.ts`），不再集中放入全局 API 模块目录。
- 插件定义 metadata 统一由 Pinia `pluginDefineStore` 缓存并提供 `clear` 和 `reload`；插件实例记录通过 `/api/plugin/instance` API 查询，按 ID 回显使用批量简要接口，不把实例列表放入定义 store。
- 插件 metadata 字段可使用 `mergeScript` 与 fast-crud `compute` 按 `form.config` 动态生成组件参数；字段联动事件应同步依赖值，并在类型切换时清理失效引用。
- 登录、注册等认证 API 统一放置在 `views/framework/auth/api.ts`；用户管理模块 API 仅包含用户数据的增删改查。
- 前端登录、注册及用户密码变更请求在发送前使用 SHA-256 对密码做一次 hash；密码为空时不发送密码字段，表示不修改。
- 用户资源使用单数路径 `/api/user`；业务动作路径固定为 `/create`、`/update`、`/delete?id=`、`/info?id=`、`/list`、`/page`、`/batchDelete`，统一使用 POST。
- 存储仓库管理页面使用 `views/sys/repository`，对应后端资源路径为 `/api/repository`；仓库密码只写入不回显，更新时空密码表示不修改。

## 规范回顾

- 每次任务收尾时，回顾用户提出的要求，识别其中可跨任务复用、稳定且不与既有规范冲突的约束。
- 将这些共性要求精炼后更新本文件；一次性的任务细节、临时偏好和未经确认的推断不写入长期规范。
- 更新后检查规则是否清晰、可执行且没有重复或互相矛盾的表述。

## 开发偏好

- 功能或实现确认废弃后，应在同一变更中删除对应的死代码、兼容分支、配置和测试；不要为不再使用的路径保留“以后可能恢复”的代码。
- 实现前先沿调用链和相邻模块了解现有做法，优先复用已有组件、工具、类型和数据流；只做满足需求的最小改动，不顺手扩展范围。
- 避免没有明确收益的辅助函数、抽象层、依赖和配置；简单逻辑直接表达，尤其不要仅为命名或转发增加 `normalizeXxxx` 一类函数。
- 缺陷修复应定位并修复共享路径上的根因，同时检查其他调用方，避免只给单个页面或调用点打补丁。
- 前端交互按用户描述的完整操作流程实现，保留清楚的状态提示和错误反馈；新增或调整界面时遵循现有组件与视觉风格，不用无关的布局或样式改造扩大任务。
- 遇到可复现的问题，优先实际运行相关命令、测试或构建并据结果修复；交付时说明验证结果和未能验证的部分。
- 用户可见行为、接口或配置发生变化时，同步更新对应使用说明；文档只描述已实现且验证过的行为。
- 需求通常先给出目标，最终效果会在首轮实现后通过阅读代码和连续反馈逐步明确；首轮实现应主动检查完整调用链、数据流、边界状态、错误反馈、交互细节和视觉一致性，尽量一次覆盖可预见的遗漏。
- 对用户补充的每条反馈，先归纳其背后的通用规则和期望行为，再修改共享实现；不要只针对当前截图、当前数据或单个操作路径打补丁。
- 交付前应从用户实际操作顺序走查一遍：入口是否正确、状态是否连续、成功与失败是否分别反馈、完成后是否停留在合理上下文；这类走查优先于增加无关功能。
- 控制改动范围不等于拒绝抽象：当多个流程需要同一项有业务语义的能力时，应把能力下沉到其所有者并复用，避免复制流程；抽象必须承载稳定职责或变化点，不要只把几行代码换个名字塞进通用 `util`/helper。

### 新需求对齐与反过度设计

提出新需求时，先按“当前行为 → 第一个具体场景 → 验收结果 → 已确认的第二个变化点”理解问题，再决定代码边界。不要因为用户提到“以后还要支持”就提前实现完整平台；先交付第一个真实场景，同时把当前实现中已经明确的职责分开。

AI 只在答案会改变实现方案时提问，优先按以下顺序确认：

1. **首个场景**：这次必须先跑通哪一种输入、入口、执行环境和输出？什么结果才算完成？
2. **事实归属**：这个数据由哪个对象或系统拥有？是否已有唯一来源？新字段是事实、配置还是派生值？
3. **变化边界**：需求中哪些变化现在已经确定，哪些只是未来可能？是否存在第二个真实调用方或第二种实现，足以证明需要抽象？
4. **能力归属**：流程编排、数据采集、远程访问、外部工具、存储或恢复分别由谁负责？哪些能力需要被多个入口复用？
5. **契约和状态**：接口、数据格式、成功/失败/进行中/取消状态、日志和错误反馈是否有兼容要求？
6. **验证方式**：用户会按什么操作顺序验收？需要覆盖哪些边界、失败路径或恢复路径？

如果这些问题能从 README、现有代码或项目规范中确定，就直接读取并采用，不重复询问。若仍有未决项，只提会影响数据模型、职责边界、外部契约、安全或恢复语义的问题，通常不超过三个；其余细节采用现有模式并在实现前写明假设。

反过度设计的执行门槛：

- 只有一个具体实现时，优先使用现有对象或方法，不新增接口、注册表、工厂、插件体系或万能数据容器。
- 出现第二个真实实现或第二个真实调用方时，比较两条完整流程，抽取已经证明稳定且需要一致处理的部分；不要按想象中的第三、第四种实现设计。
- 不同维度只有在需求明确独立变化时才拆分。例如“数据库备份”先确认是逻辑 dump 还是原生快照、在哪里执行、如何恢复；若本次只支持 PostgreSQL dump，就先实现该路径，不提前做多数据库适配框架。
- 用户只要求一个垂直流程时，先完成最短可验证链路；完成后再检查重复、职责污染和扩展成本，必要时做一次有依据的下沉重构。
- 每次提出抽象时同时写出它解决的真实重复、变化点或一致性规则；写不出来就不抽象。

需求确认后的简短输出应包含：`首个场景`、`已确认的变化点`、`暂不设计的未来项`、`拟复用/新增的职责`、`验收路径`。这样既能让用户校准思路，也能防止把讨论中的可能性误当成当前需求。

### 典型案例复盘

以下案例记录的是“如何思考和演进”，不是要求以后保留当时的文件名、类型名或目录结构。后续实现应识别相同的职责、数据关系和变化点。

#### 案例一：从“SSH 主机 + Restic”扩展到文件、数据库和虚拟机备份

初始状态：项目以 Restic 管理远程备份。第一种执行方式是在 SSH 主机上运行 Restic，将主机文件路径备份到 SFTP 仓库（当前执行链路涉及 `BackupExecuteService`、`BackupProvider`/`ResticProvider` 和 `SshAccess`）。这只是“远程文件备份”的一条实现，不能据此认定 SSH + Restic 就是所有备份的通用模型。

用户提出的后续方向：还要支持数据库和虚拟机备份，同时复用计划、运行状态和日志；不同方法各自需要的采集、访问和恢复方式可能不同。

首轮实现的风险：把 SSH 配置、文件路径、Restic 二进制、仓库 URL 放进一个号称通用的备份接口，之后再把数据库连接、VM 标识和快照参数也塞进去。结果是接口越来越胖、实现里按类型分支。另一种过度设计是现在就造一个同时表示路径、dump 流和整机快照的万能 `BackupArtifact`，但还没有第二种实际流程证明这些输入应共享同一契约。

改进时按职责和变化轴划分：

- **计划执行生命周期**：读取计划、状态转换、统一结果与阶段日志；手动、定时和重试都进入同一执行用例（当前由 `BackupExecuteService` 编排，`BackupPlanService` 提供计划操作和触发入口）。
- **备份方法**：拥有特定数据的采集和一致性规则。文件方法处理目录与排除项；数据库方法处理数据库类型、认证、事务一致性及 dump/原生快照；虚拟机方法处理 hypervisor、VM 快照、磁盘导出与快照清理。三者不能仅靠增加一个 `type` 参数就视为同一种备份。
- **访问能力**：按需提供 SSH 命令/文件传输、数据库连接、hypervisor API 或本机进程执行。SSH 是当前方法用到的一种访问方式（历史实现为 `SshAccess.Platform`、`Upload`、`Execute`），不应成为数据库或 VM 方法的强制依赖。
- **Restic 能力**：负责 Restic 支持的输入、快照和恢复格式，以及仓库配置（当前由 `BackupProvider`/`ResticProvider` 承担）。SFTP 是 Restic 使用的仓库后端，不是数据库或 VM 的备份方法。只有先确认某种方法能提供 Restic 可接收的数据形态、并且恢复语义成立，才复用当前 Restic 流程。

伪代码只表达职责，不预定未来接口：

```text
run(plan):
    markRunning(plan)
    result = methodFor(plan).backup(plan)  # 文件、数据库或 VM 各自负责其采集规则
    markSuccessOrFailed(plan, result)

sshFileBackup:
    use SSH to prepare/run Restic on the source host

databaseBackup:
    create a consistent dump or native snapshot
    select a storage flow that can restore it

virtualMachineBackup:
    create a consistent VM snapshot and export disks/metadata
    clean up the snapshot after the backup outcome is known
```

最终组织方式不预先固定类图：共同生命周期只实现一次；每种备份方法拥有自己的配置、一致性、采集、清理和恢复规则；SSH、数据库连接和 hypervisor API 作为按需使用的访问能力；Restic 只复用到它确实支持且可恢复的数据输入上。新增第二种方法时，先对比它和现有方法的完整流程，只把真实共有且需要一致处理的能力下沉；那时再决定是否需要方法接口或共享输入契约。

好处：不会把数据库和 VM 参数硬塞进 SSH 文件备份；不会因预想所有未来场景而提前造万能抽象；真正共享的运行状态、错误和日志仍只有一套；类名/文件演进后，领域职责仍然清楚。

形成的规则：区分“计划生命周期、备份方法、访问能力、Restic 支持的输入/仓库配置、恢复语义”。抽象服务于真实共性与职责归属，不服务于当前类名；除非多个具体方法证明需要共享契约，否则不预造通用工件模型。

#### 案例二：用单一数据来源和业务约束解决仓库目录冲突

初始状态：来源实例可以配置多个来源路径，备份计划还准备保存来源路径；一个存储仓库又可能被多个计划使用，因此不同计划可能写入同一仓库位置。

用户需求：来源路径不要在计划中重复维护；同一仓库可以被多个计划使用，但每个计划需要有自己的目录；目录最好清楚、不能冲突、不能越界；空目录要有明确含义。

首轮实现容易出现的做法：在计划表再加一份 `sourcePath`，前端检查目录是否重复，后端只保存用户输入的字符串。这样来源修改后计划里的路径会失效，绕过前端也能写入冲突或危险路径。

用户要求改善：来源路径归来源实例管理；计划只保留仓库子目录；相同仓库不能冲突；需要考虑规范化和越界。

最终结果：来源实例（当时通过插件实例配置中的 `paths` 提供）是来源路径的唯一来源；计划实体（`BackupPlan`）只保存仓库子目录（`RepoTag`，输入字段为 `BackupPlanInput.RepoTag`）；服务层把反斜杠统一为斜杠、清理多余分隔符、拒绝绝对路径和 `..` 越界；数据库用“仓库标识 + 规范化子目录”做唯一约束；空子目录代表仓库根目录，并参与同样的冲突检查。

伪代码：

```text
savePlan(input):
    subPath = normalizeRelativePath(input.repoTag)
    require unique(repositoryId, subPath)
    persist(plan.repositoryId, subPath)

backup(plan):
    sourcePaths = loadSourceInstance(plan.sourceId)
    target = join(repository.path, plan.repoTag)
```

好处：只有一个地方维护来源路径；接口、页面和数据库都遵守同一约束；目录冲突和路径穿越不会依赖前端自觉；以后增加定时执行、恢复或迁移功能时可以直接复用同一数据模型。

形成的规则：能由数据模型和数据库约束表达的业务规则，不要只放在页面；同一事实只保留一个来源，派生关系在使用时组合。

#### 案例三：修改组件契约要贯穿包装层和使用层

初始状态：选择器组件使用自定义的 `value`/`update:value`，外层页面、包装组件和内部表格选择器之间的绑定命名不完全一致。

用户需求：统一改成 Vue 3 标准的 `modelValue`，并保证选择、清空、回显、监听和异步加载都正常。

首轮实现容易出现的做法：只把页面上的 `:value` 改成 `:model-value`，组件内部仍然声明旧 props 和旧事件。静态检查可能通过，但双向绑定或回显会在特定路径失效。

用户要求改善：包装组件和内部组件一起改，不能只改表面调用点。

最终结果：从页面到包装组件（当时的 `AccessSelector`）再到通用插件选择器（`PluginSelector`）或路径选择器（`PluginPathSelector`），统一使用 `modelValue` prop 和 `update:modelValue` 事件；初始化回显、用户选择、清空和异步查询都经过同一条契约。

伪代码：

```text
page:          modelValue <-> wrapper
wrapper:       modelValue <-> selector
selector:      modelValue -> selectedRecord
selector:      update:modelValue -> wrapper -> page
```

好处：跨层契约只有一个名字和一套语义；组件可以自由替换内部实现；不会出现“页面看起来改好了、内部状态没有同步”的隐蔽问题。

形成的规则：修改接口、事件或数据结构时，沿调用链一次检查生产者、消费者、包装层、监听、回显和清理逻辑，不能只修改最外层。

#### 案例四：先解决结构性问题，再修视觉细节

初始状态：备份计划新增了来源、计划、存储仓库三列数据流图，但静态布局下节点数量变化就会导致箭头方向、位置、颜色和动态效果异常。

用户需求：箭头要正确，图标颜色要一致，布局要灵动，数据流状态要能看懂。

首轮实现容易出现的做法：继续增加固定 `top/left`、伪元素和 CSS 动画，针对当前截图调整某一条线。

用户要求改善：寻找支持节点连接、箭头和自动布局的画布能力，而不是继续堆静态偏移。

最终结果：先采用能表达节点、边、方向和布局的结构化画布模型，再把运行中、成功、失败、禁用等状态映射为统一视觉样式；节点增删或尺寸变化时由布局和连接机制重新计算（当时页面视图命名为 `BackupPlanFlow`，单个可复用节点命名为 `PlanFlowNode`）。

好处：修复的是布局模型而不是单个截图；新增节点、不同窗口宽度和不同状态仍能保持可读；后续增加拖拽、缩放或更多关系时不需要重写定位 CSS。

形成的规则：当问题由节点关系、状态模型或布局机制引起时，先更换或修正结构；不要用更多局部样式掩盖结构缺陷。

#### 案例五：完整流程要可观察，成功失败要真实

初始状态：执行任务只打印开始和结束，Restic 下载地址按版本号拼接；下载失败后用户不知道失败在哪一步，甚至可能看到成功提示。

用户需求：展示整个执行过程；下载要按远程平台匹配真实 Release asset；失败必须明确失败，用户不能被误导。

首轮实现容易出现的做法：增加一条“执行中”日志，继续拼接猜测 URL，沿用成功分支的提示文本。

用户要求改善：记录计划读取、状态更新、平台探测、下载、上传、命令执行和最终结果；错误保留底层原因和阶段上下文。

最终结果：每个关键阶段都有结构化日志、对象标识、数量或路径和耗时；下载先读取 Release metadata 再精确匹配 asset（当时下载逻辑位于 `downloadResticBinary`，平台名称解析位于 `resticPlatform`）；任一步失败都会更新失败状态并返回真实错误，成功提示只在命令确认成功后出现（远程命令统一经 `SshAccess.Execute` 执行）。

好处：用户和开发者能定位失败阶段；外部下载和远程命令的行为可追踪；错误状态不会被成功文案掩盖；问题修复可以针对根因而不是猜测。

形成的规则：长流程按阶段设计日志和状态，错误发生在哪里就在哪里补充上下文；成功、失败、取消和进行中必须是互斥且真实的状态。

实现新需求时，先用上述案例问自己：事实由谁拥有？流程由谁编排？变化点在哪里？哪些能力会被另一个入口复用？约束能否下沉到业务层或数据库？用户从入口到完成会看到哪些状态？回答后再决定文件和类型如何组织。

#### 案例六：把 Restic 执行环境下沉到固定 executor

初始状态：Restic 同时需要支持两条真实链路：本机执行 Restic 查询远程仓库，以及通过 SSH 在来源主机执行 Restic 备份。旧实现把下载、上传、临时密钥、SSH 调用、仓库初始化和 Restic 命令混在 provider 中，调用方还需要传入不同的访问对象；Windows 私钥权限、远程二进制复用和失败诊断又让流程不断增加分支。

首个具体场景：创建一个 Restic 客户端，指定执行环境为 local 或 remote，然后完成环境准备、仓库初始化、备份、快照查询和恢复。验收结果是调用方只组装配置和业务请求，客户端能够执行完整流程；本地与远程执行使用同一组业务方法，SSH 免密由用户预先配置。

首轮实现容易出现的做法：在每个业务方法中重新选择 executor；让 `ResticClient` 持有并操作 `SshAccess`；把二进制下载、上传、配置文件写入和仓库初始化继续放在客户端；通过 `executorFor(access)` 或带 `*SshAccess` 参数的方法保留第二套执行路径；为了兼容旧调用保留废弃 provider、私钥临时文件和诊断分支。这样看似抽象，实际职责仍然分散，调用方也无法只通过初始化确定执行环境。

最终结果：`NewResticClient` 初始化时固定 local 或 remote executor；`ResticClient` 负责业务流程编排和 Restic 参数，executor 负责本地/远程环境准备、二进制路径、仓库环境准备、命令执行和资源释放；下载实现独立到 `download.go`；远程 executor 负责 SSH 连接、二进制上传、`restic.info` 复用、仓库配置和初始化；本地 executor 负责本地缓存和进程执行；调用方分别构造 `ClientConfig`、业务 request、client，再调用方法。

伪代码：

```text
config = ClientConfig{Environment: remote, ...}
request = BackupRequest{...}
client = NewResticClient(serviceContext, config)
result = client.Backup(ctx, request)
```

明确不提前设计：只有一个 Restic 引擎时不增加通用 `BackupMethod` 注册表；没有第二个真实执行实现时不增加工厂层；暂时由用户配置 SSH 免密，不把密码、私钥临时文件或 Windows ACL 再塞回 ResticClient；不为了兼容旧调用保留废弃入口。

形成的规则：先沿调用链确认真实的执行环境和业务入口，再抽取稳定职责；执行环境的变化由 executor 封装，备份/恢复/快照等业务能力由客户端统一编排。构造客户端时确定环境，业务方法中禁止再次选择 executor。客户端文件中不得出现带 `SshAccess` 参数的环境操作方法；如果多个业务流程需要上传配置、准备二进制或释放连接，应下沉到 executor。独立且较大的下载、解析或协议适配流程放入职责明确的文件。多件有业务意义的动作必须拆成配置、请求、客户端、调用和错误处理等独立语句。

### 插件类型与核心能力契约

插件类型不是同一套业务接口的不同名称。所有插件只共享生命周期和元数据能力；每种插件类型必须拥有自己的领域接口，不能把所有行为退化成字符串 action 加 `map[string]any`。

#### 所有插件的通用接口

以下接口是职责契约示意，落地时应适配项目现有的 `plugin.Definition` 和注册流程；不要为了接口名称重写已经稳定的插件基础设施。所有插件都应能提供定义并校验配置：

```go
type Plugin interface {
    Definition() Definition
    ValidateConfig(config map[string]any) error
}
```

需要持有连接、临时目录或外部句柄的插件，再实现 `Close() error`。不要把 `Connect`、`Backup`、`Restore` 等具体业务方法放进所有插件的公共接口。

#### Access 插件：连接和访问外部系统

`Access` 只负责连接、连通性检查和释放资源；具体访问能力按需组合：

```go
type Access interface {
    Plugin
    Connect(ctx context.Context) error
    Test(ctx context.Context) error
    Close() error
}

type CommandExecutor interface {
    Execute(ctx context.Context, command string) ([]byte, error)
}

type FileTransfer interface {
    Upload(ctx context.Context, target string, data []byte, mode uint32) error
    Download(ctx context.Context, source string) ([]byte, error)
}
```

SSH 可以组合命令执行、文件传输、目录浏览和平台探测；数据库 Access 组合数据库查询、导出和导入；虚拟化平台 Access 组合虚拟机快照、磁盘导出和虚拟机创建。Access 不负责备份流程、Restic 仓库或数据库/虚拟机的业务编排。

#### Repository 插件：仓库配置和仓库访问

`Repository` 负责仓库地址、认证、后端配置和仓库校验；需要时提供 Restic 仓库配置、快照查询和删除能力：

```go
type Repository interface {
    Plugin
    Validate(ctx context.Context) error
    BuildResticConfig(ctx context.Context) (ResticRepositoryConfig, error)
}
```

SFTP、S3、MinIO 和本地目录是不同的 Repository 实现。Repository 不负责来源数据导出、虚拟机快照或数据库导入。

#### BackupSource 插件：来源数据的准备和恢复

`BackupSource` 拥有来源数据的一致性、导出、导入和临时资源清理规则。文件、数据库和虚拟机可以有不同的专用接口，不要强行塞进一个万能请求：

```go
type BackupSource interface {
    Plugin
    Validate(ctx context.Context) error
    PrepareBackup(ctx context.Context) (BackupInput, error)
    PrepareRestore(ctx context.Context, input RestoreInput) error
}
```

主机文件来源直接提供路径和排除路径；数据库来源负责逻辑 dump 或原生快照以及导入；虚拟机来源负责一致性快照、磁盘和元数据导出、虚拟机创建和快照清理。BackupSource 可以调用 Access 和 Restic 核心，但不实现仓库内部逻辑。

典型流程：

```text
BackupSource 准备/导出
    -> ResticCore.Backup
    -> BackupSource 清理

ResticCore.Restore
    -> BackupSource 导入/创建资源
```

#### Restic：系统核心能力

Restic 负责所有来源可复用的文件系统备份和还原能力，不作为当前唯一实现的 `BackupMethod` 插件，也不通过一个只转发调用的插件包装：

```go
type ResticCore interface {
    Backup(ctx context.Context, request BackupRequest) (BackupResult, error)
    Restore(ctx context.Context, request RestoreRequest) (RestoreResult, error)
}
```

ResticCore 负责二进制准备、仓库环境、快照、加密、去重、备份、还原、命令执行和 Restic 错误；不负责数据库 dump、虚拟机快照、数据库导入或虚拟机创建。

只有出现第二种真实备份引擎（例如 Borg 或 Kopia），并且它与 Restic 需要稳定共享同一套生命周期契约时，才引入 `BackupMethod` 接口和注册机制。只有 Restic 一个实现时，不预造该插件类型。

#### 四类需求的职责落点

| 场景 | BackupSource 负责 | ResticCore 负责 |
| --- | --- | --- |
| 主机文件 | 选择主机路径和还原目录 | 备份/还原文件 |
| 虚拟机 | 创建快照、导出磁盘和元数据、创建虚拟机、清理快照 | 备份/还原导出文件 |
| 数据库直连 | 数据库导出和导入 | 备份/还原 dump 文件 |
| 数据库主机文件 | 登录数据库主机、导出和导入文件 | 备份/还原主机上的导出文件 |

判断是否新增接口时，先确认是否存在第二个真实实现和稳定的共享规则；如果只是当前 Restic 的一层转发，不新增插件或接口。

## 开发与验证

- 前端代码修改后使用 `front` 目录中的 Prettier 执行格式化，并使用 `prettier --check` 检查相关文件。
- 开发完成后不自动启动前端、后端或组合服务；由用户根据 VS Code 调试配置或项目命令自行启动和验证。
- 交付时说明已执行的自动化检查、构建结果以及未解决的既有问题。

## 前端代码风格

- 变量和函数使用有意义的完整命名，禁止无意义的简称和单字母变量。
- 禁止使用三元表达式；`if` 必须使用大括号，禁止将 `if` 写在单行。
- Vue 组件样式使用 `<style lang="less">`，不使用 `scoped`；根元素必须定义组件根 class，样式使用根 class 包裹子 class。
- 覆盖第三方 Vue 组件样式时，优先给组件增加语义明确的 class，并在自身根 class 下通过子 class 编写样式；只有无法通过 class 覆盖时才使用 `:deep`。

## 命名

- Go 标识符中的 ID 统一写作 `Id`，例如 `SourceId`、`RepositoryId`、`UserId` 和 `Id`；JSON 字段、数据库列名及外部配置键按现有协议保持不变。

## Go 注释规范

- 导出的类型、函数、方法和常量必须有以标识符名称开头的中文完整句子注释，说明职责和可观察行为。
- 非导出的复杂流程函数、协议适配函数和安全敏感逻辑应补充简短注释，说明原因、前置条件或幂等语义；简单转发和自解释代码不添加重复注释。
- 注释描述“为什么”和对调用方有用的契约，不逐行翻译实现；涉及错误、权限、资源释放、日志脱敏或兼容性时应明确边界。
- 新增和修改的 Go 注释必须使用中文；技术名词、函数名、命令和协议保留英文；示例命令使用代码块，避免把过期实现细节写入注释。

## 代码表达规范

- 禁止把多个有业务意义的动作压缩到同一行，尤其是对象初始化、方法调用、错误处理和返回值赋值的组合表达式；配置对象、请求对象、客户端初始化、业务调用和错误判断必须按职责拆成独立语句。
- Go 长调用必须拆分为可读的多行代码，禁止使用“构造配置 + 构造客户端 + 调用业务方法”这类链式单行表达式。
