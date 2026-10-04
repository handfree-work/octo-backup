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

## 错误处理

- 业务层、service、插件和基础设施在错误发生的位置立即使用 `internal/base/error_` 的 `NewTextError`、`NewCodeTextError` 或 `NewFormatError` 创建并记录 `CodedError`，错误消息必须包含足够的模块、操作和底层原因上下文。
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
- 插件定义 metadata 统一由 Pinia `pluginDefineStore` 缓存并提供 `clear` 和 `reload`；插件实例记录通过 `plugin-instance` API 查询，按 ID 回显使用批量简要接口，不把实例列表放入定义 store。
- 插件 metadata 字段可使用 `mergeScript` 与 fast-crud `compute` 按 `form.config` 动态生成组件参数；字段联动事件应同步依赖值，并在类型切换时清理失效引用。
- 登录、注册等认证 API 统一放置在 `views/framework/auth/api.ts`；用户管理模块 API 仅包含用户数据的增删改查。
- 前端登录、注册及用户密码变更请求在发送前使用 SHA-256 对密码做一次 hash；密码为空时不发送密码字段，表示不修改。
- 用户资源使用单数路径 `/api/user`；业务动作路径固定为 `/create`、`/update`、`/delete?id=`、`/info?id=`、`/list`、`/page`、`/batchDelete`，统一使用 POST。
- 存储仓库管理页面使用 `views/sys/repository`，对应后端资源路径为 `/api/repository`；仓库密码只写入不回显，更新时空密码表示不修改。

## 规范回顾

- 每次任务收尾时，回顾用户提出的要求，识别其中可跨任务复用、稳定且不与既有规范冲突的约束。
- 将这些共性要求精炼后更新本文件；一次性的任务细节、临时偏好和未经确认的推断不写入长期规范。
- 更新后检查规则是否清晰、可执行且没有重复或互相矛盾的表述。

## 开发与验证

- 前端代码修改后使用 `front` 目录中的 Prettier 执行格式化，并使用 `prettier --check` 检查相关文件。
- 开发完成后不自动启动前端、后端或组合服务；由用户根据 VS Code 调试配置或项目命令自行启动和验证。
- 交付时说明已执行的自动化检查、构建结果以及未解决的既有问题。

## 前端代码风格

- 变量和函数使用有意义的完整命名，禁止无意义的简称和单字母变量。
- 禁止使用三元表达式；`if` 必须使用大括号，禁止将 `if` 写在单行。
- Vue 组件样式使用 `<style lang="less">`，不使用 `scoped`；根元素必须定义组件根 class，样式使用根 class 包裹子 class。
