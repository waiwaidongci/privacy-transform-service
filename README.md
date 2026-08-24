# 企业级隐私数据脱敏与策略执行服务

纯Go实现的结构化JSON隐私处理服务，按分类、字段路径和已发布策略执行遮罩、SHA-256哈希、确定性令牌化、截断、泛化和删除。令牌库、密钥服务、Redis、PostgreSQL和消息队列均保留接口边界；日志只记录请求ID、策略和结果摘要，不写入原始敏感值。

## 运行

需要Go 1.22或更高版本：

```bash
go test ./...
go run ./cmd/privacy-transform-service
```

默认监听 `:8083`。可通过 `PRIVACY_TRANSFORM_HTTP_ADDR`、`PRIVACY_TRANSFORM_ENVIRONMENT` 和 `PRIVACY_TRANSFORM_SHUTDOWN_SECONDS` 覆盖配置。`configs/config.yaml`提供配置示例；SQL迁移位于`migrations`，Docker文件位于`deploy`。

## API流程

```bash
curl -X POST localhost:8083/v1/privacy/classifications -H 'Content-Type: application/json' -d '{"id":"pii","name":"个人标识","level":3}'
curl -X POST localhost:8083/v1/privacy/policies -H 'Content-Type: application/json' -d '{"id":"customer-v1","name":"Customer policy","rules":[{"path":"email","classification":"pii","action":"mask","preserve":2},{"path":"phone","classification":"pii","action":"hash","salt":"v1"}]}'
curl -X POST localhost:8083/v1/privacy/policies/publish -H 'Content-Type: application/json' -d '{"id":"customer-v1"}'
curl -X POST localhost:8083/v1/privacy/transform -H 'Content-Type: application/json' -d '{"request_id":"r-1","policy_id":"customer-v1","data":{"email":"alice@example.com","phone":"13800138000","city":"Tokyo"}}'
curl -X POST localhost:8083/v1/privacy/batch -H 'Content-Type: application/json' -d '{"policy_id":"customer-v1","records":[{"email":"bob@example.com"},{"email":"carol@example.com"}]}'
```

策略发布和回滚会更新缓存版本；处理记录仅保存摘要，不持久化原始敏感值。服务收到SIGINT/SIGTERM后停止接收请求并优雅关闭。

## 目录

`cmd`为启动入口，`internal/domain`为纯领域模型和规则求值，`internal/application`为用例，`internal/adapter`为HTTP和缓存适配器，`internal/infrastructure`为配置、日志、指标和仓储实现，`api`为OpenAPI，`migrations`为数据库迁移，`deploy`和`scripts`提供部署与校验辅助。
