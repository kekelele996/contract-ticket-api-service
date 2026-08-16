# 验证报告

## 项目信息

- 项目目录：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/生活服务主题项目提示词/cy-397`
- 项目类型：纯后端服务「合同模板生成与法律工单 API 服务」
- 技术栈：Go 1.22 + Gin + GORM + MySQL 8.0 + JWT + wkhtmltopdf

## 提示词 SHA-256

- 开始：`ca4a5efdf01a7003a301d55b8adf89edbe4065a80d82bc8d3a980add265f8480`
- 结束：`ca4a5efdf01a7003a301d55b8adf89edbe4065a80d82bc8d3a980add265f8480`
- 前后一致：是

## 顶层生成路径

- `backend/cmd/server/main.go`
- `backend/internal/config/`
- `backend/internal/constants/`
- `backend/internal/dto/`
- `backend/internal/handler/`
- `backend/internal/middleware/`
- `backend/internal/model/`
- `backend/internal/repository/`
- `backend/internal/router/`
- `backend/internal/service/`
- `backend/pkg/jwtutil/`
- `backend/templates/`（5 个合同模板文件）
- `backend/Dockerfile`
- `database/init.sql`
- `migrations/000001_init.sql`
- `api/openapi.yaml`
- `deploy/README.md`
- `docker-compose.yml`
- `.env.example`
- `.gitignore`
- `README.md`

## 本地构建与测试

在 `backend/` 下执行：

```bash
go mod tidy
go vet ./...
go build ./...
go test ./...
```

结果：

- `go mod tidy`：通过，`go.mod` 声明 `go 1.22.0`，包含 gin、gorm、gorm mysql driver、golang-jwt/v5、caarlos0/env/v11、x/crypto/bcrypt 等依赖。
- `go vet ./...`：通过。
- `go build ./...`：通过。
- `go test ./...`：通过，repository 与 service 均有表驱动测试。

## Docker Compose 验证

执行：

```bash
docker compose --env-file .env config --quiet
docker compose --env-file .env up -d --build --wait
```

结果：

- `docker compose config --quiet`：通过。
- `up -d --build --wait`：backend 与 mysql 均进入 healthy。
- backend 对外端口：`127.0.0.1:19412`。

### 容器内依赖

- `wkhtmltopdf 0.12.6` 可用。
- MySQL 8.0 数据表初始化成功，5 个合同模板与 6 条 FAQ 种子数据中文显示正常。

## 主要 API 流程 curl 证据

以下为关键调用摘要（完整命令见执行过程）：

1. `GET /healthz` 返回 `{"status":"ok"}`。
2. `POST /api/v1/auth/register` 创建 `verifyuser2`，返回 `code:0`。
3. `POST /api/v1/auth/login` 返回 JWT（长度 217）。
4. `GET /api/v1/templates` 返回 5 个模板：nda、cooperation、loan、labor、lease。
5. `POST /api/v1/contracts` 使用 lease 模板填充变量，生成草稿合同，`content_text` 包含“房屋租赁合同\n出租方（甲方）：张三…”。
6. `GET /api/v1/contracts` 用户合同库按用户隔离，返回 1 条草稿。
7. `GET /api/v1/contracts/3` 合同详情包含 `template_name: 房屋租赁合同`。
8. `POST /api/v1/templates/1/favorite` 收藏成功；`GET /api/v1/templates/favorites` 返回 1 条。
9. `POST /api/v1/contracts/3/submit`：draft → pending_signed。
10. `POST /api/v1/contracts/3/sign`：pending_signed → signed，记录签署时间。
11. `GET /api/v1/contracts/3/signers` 返回提交的 2 个签署方 + 实际签署方。
12. `GET /api/v1/contracts/3/export` 返回 `200 OK`，`Content-Type: application/pdf`，文件头 `%PDF-1.4`。
13. `POST /api/v1/tickets` 创建工单，状态 pending。
14. 用户回复后状态 pending → processing；律师回复后 processing → replied。
15. `PATCH /api/v1/tickets/1/status` 将 replied → closed。
16. `GET /api/v1/tickets` 列表返回 1 条 closed 工单。
17. `GET /api/v1/faqs?keyword=违约` 返回 2 条 FAQ。

### 错误处理验证

- 无 Token 访问受保护接口：`HTTP 401`，`{"code":40100,"message":"missing bearer token"}`。
- 密码错误登录：`HTTP 401`，`{"code":40100,"message":"invalid username or password"}`。
- 缺少必填变量生成合同：`HTTP 400`，`{"code":40001,"message":"missing required variables: party_b, address, amount, start_date, end_date, sign_date"}`。
- 非法工单类型：`HTTP 400`，`{"code":40001,"message":"invalid ticket type"}`。

## 清理验证

执行：

```bash
docker compose --env-file .env down -v --remove-orphans
```

- 端口 `19412`：已释放，无监听。
- `docker compose ps`：无运行中的 contractapi 服务。
- contractapi 相关容器/网络：无。
- contractapi 相关卷：已清理（含一个 2026-06-11 遗留的同名项目卷）。

## Git 信息

- 主提交：`2a0c33315b541cd1b0e5c713a751825f3ebdc4a7`
- 提交信息：`feat: 合同模板生成与法律工单 API 服务`
- 本地用户：`blueship581 <brysj.hhrhl.g@gmail.com>`

## 已知问题 / 偏差说明

1. 提示词给的 `backend/Dockerfile` 运行镜像为 `alpine:3.20` 并 `apk add wkhtmltopdf`，但 Alpine 3.20 仓库已无 `wkhtmltopdf` 包。为保证真实可用的 PDF 导出，运行镜像改为 `debian:bookworm-slim`，安装 `wkhtmltopdf` 与 `fonts-noto-cjk`；Go 多阶段构建保持不变。
2. 本机 Go 工具链为 1.25.1，曾将 `go.mod` 提升到 `go 1.25.0`；为兼容 Dockerfile 中的 `golang:1.22-alpine`，已将 `go.mod` 固定为 `go 1.22.0`，并降级 gin、x/crypto、x/sys、x/text、x/arch、protobuf 等到与 Go 1.22 兼容的版本。
3. 管理类接口（模板/FAQ 的新增、更新、删除）仅做 JWT 鉴权，未引入独立管理员角色，属于演示实现。
