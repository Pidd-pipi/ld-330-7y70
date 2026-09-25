# 部署说明

Docker Compose 文件位于项目根目录。生产环境应替换 `.env` 中的 `JWT_SECRET` 与数据库密码，并通过受控迁移流程审核 `migrations/` 中的数据库变更。
