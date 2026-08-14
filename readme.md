# Web Restic

一个基于Restic的远程备份管理系统，优雅的web界面，远程系统数据备份、恢复、管理等功能。

## API 文档

启动后端后可在 [Swagger UI](http://127.0.0.1:3000/swagger/) 查看和调试接口，OpenAPI JSON 位于 `/swagger/doc.json`。

后端接口发生变化时，在 `backend` 目录执行以下命令更新提交的文档：

```powershell
swag init --generalInfo app.go --output docs --parseInternal
```


