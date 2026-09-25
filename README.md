# 明医电子病历管理系统（GB EMR）

面向中小型医疗机构的全栈电子病历（EMR）系统，提供患者建档、结构化病历、医嘱与电子处方、审签归档、检索统计和基础数据维护能力。

## Docker Compose 一键启动（推荐）

```bash
cp .env.example .env
docker compose up -d
```

首次构建会拉取 PostgreSQL、Node、Nginx 与 Go 基础镜像。查看运行状态：

```bash
docker compose ps
```

停止并保留数据库数据：

```bash
docker compose down
```

演示账号：`admin/admin123`、`doctor/doctor123`、`nurse/nurse123`。

## 访问地址

| 服务 | 地址 |
| --- | --- |
| Web 前端 | http://localhost:18930 |
| 后端健康检查 | http://localhost:19930/healthz |
| API 根路径 | http://localhost:19930/api/v1 |
| PostgreSQL | localhost:5432 |

## 主要功能

- **患者档案**：唯一档案编号；完整的基础信息、过敏史和既往病史；可按姓名、身份证号、手机号检索。
- **结构化病历**：门诊/住院类型，包含主诉、现病史、既往史、检查、诊断、方案和富文本补充内容；以时间线展示。
- **医嘱与处方**：长期/临时医嘱；电子处方包含药品、规格、剂量、频次、疗程，并跟踪待审核、已审核、已执行状态。
- **审签和留痕**：JWT 登录与管理员/医生/护士 RBAC；病历审核、审签归档与归档后的修改申请；关键创建动作写入审计日志。
- **检索与报表**：患者、科室、医生、日期、关键词多维病历检索；科室工作量与疾病谱报表。
- **系统管理**：管理员维护科室、账号、药品字典、ICD-10 诊断编码、病历模板库和操作审计。

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Ant Design + Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 15 |
| 鉴权 | JWT (`golang-jwt/jwt/v5`) + 角色权限控制 |
| 部署 | Docker Compose + Nginx |

## 本地开发（备选）

先在本机运行 PostgreSQL，并保证环境变量中的数据库连接可用：

```bash
cp .env.example .env
# 前端终端
cd frontend && npm install && npm run dev
# 后端终端
cd backend && go mod tidy && go run ./cmd/server
```

后端构建与测试：

```bash
cd backend
go build ./...
go test ./...
go vet ./...
```

前端检查与构建：

```bash
cd frontend
npm run check
npm run build
```

## API 清单

所有业务接口使用 `/api/v1` 前缀、`Authorization: Bearer <token>` 鉴权，并统一返回：

```json
{"code": 0, "message": "ok", "data": {}}
```

- `POST /auth/login`：登录；`GET /auth/me`：当前用户。
- `GET|POST|PUT /patients`：患者查询、创建和更新。
- `GET|POST /records`、`GET /records/:id`：病历检索、创建、详情。
- `POST /records/:id/review`、`POST /records/:id/change-requests`：审核/归档和归档修改申请。
- `POST /orders`、`GET /records/:id/orders`、`PUT /orders/:id/status`：医嘱。
- `POST /prescriptions`、`GET /records/:id/prescriptions`、`PUT /prescriptions/:id/status`：电子处方。
- `GET /reports/department-workload`、`GET /reports/disease-spectrum`：统计报表。
- `/admin/*`：科室、账号、药品、ICD-10、模板、审计日志（管理员）。

OpenAPI 摘要见 [`backend/api/openapi.yaml`](backend/api/openapi.yaml)。

## 目录结构

```text
.
├── frontend/                    # React 单页应用、页面、组件、API 客户端
├── backend/
│   ├── cmd/server/              # 启动装配
│   ├── internal/                # config/model/repository/service/handler/router/middleware/dto/constants
│   ├── api/openapi.yaml         # 接口文档摘要
│   ├── migrations/              # 迁移说明
│   └── deploy/                  # 部署说明
├── database/init.sql            # 初始化说明
├── docker-compose.yml           # 前后端和 PostgreSQL 编排
├── .env.example                 # 环境变量样例
└── README.md
```

## 环境变量

| 变量 | 默认/示例 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `gbemr` | Compose 项目与容器名前缀 |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | `gbemr` / `gbemr` / `gbemr_password` | PostgreSQL 连接凭据 |
| `JWT_SECRET` | 示例随机串 | JWT 签名密钥，生产必须替换 |
| `FRONTEND_PORT` | `18930` | Nginx 前端宿主机端口 |
| `BACKEND_PORT` | `19930` | Gin 后端宿主机端口 |

## Docker 部署说明与常见问题

- Compose 顶层固定 `name: gbemr`，同时 `.env` 指定 `COMPOSE_PROJECT_NAME=gbemr`，可在中文目录名中稳定运行。
- 数据库持久化至命名卷 `db_data`，不会绑定挂载到中文路径。若需完全清除测试数据：`docker compose down --volumes`。
- 前端只调用 `/api`，由 Nginx 反向代理至 `backend:8080`；无需在浏览器中配置本地后端地址。
- 若端口已被占用，请修改 `.env` 的 `FRONTEND_PORT` 或 `BACKEND_PORT` 后重新执行 `docker compose up -d`。
- 生产环境请设置强随机 `JWT_SECRET`、最小权限数据库账户，并在网关层补充 HTTPS、备份、审计留存与访问控制。

## License

MIT License。仅用于演示、学习和内部原型；实际医疗生产使用前应完成合规、安全、数据备份与隐私评估。
