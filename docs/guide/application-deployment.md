# 应用部署指南

本文档说明如何在生产环境中部署 Go DDD Template 应用（**不是文档系统部署**）。

> 💡 **文档部署**：如需部署 VitePress 文档到 GitHub Pages，请查看 [文档部署指南](/guide/docs-deployment)

## 部署方式概览

本指南涵盖以下部署方式：

| 方式           | 适用场景                 | 复杂度   | 推荐度     |
| -------------- | ------------------------ | -------- | ---------- |
| Docker         | 单机部署、快速测试       | ⭐       | ⭐⭐⭐⭐⭐ |
| Docker Compose | 多服务编排、本地生产模拟 | ⭐⭐     | ⭐⭐⭐⭐   |
| Kubernetes     | 大规模、高可用、云原生   | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 二进制部署     | 简单环境、裸机           | ⭐       | ⭐⭐⭐     |

## Docker 部署

### 1. 创建 Dockerfile

在项目根目录创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.25.4-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /build/app \
    ./main.go

# 运行阶段
FROM alpine:latest

# 安装运行依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/app /app/app

# 复制配置文件（可选）
COPY --from=builder /build/configs /app/configs

# 设置文件所有权
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/bin/sh", "-c", "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"]

# 启动应用
CMD ["/app/app", "api"]
```

### 2. 构建镜像

```bash
# 构建镜像
docker build -t go-ddd-template:latest .

# 查看镜像
docker images go-ddd-template
```

### 3. 运行容器

```bash
# 运行容器（开发环境）
docker run -d \
  --name go-ddd-app \
  -p 8080:8080 \
  -e APP_SERVER_ADDR=:8080 \
  -e APP_DATA_PGSQL_URL=postgresql://postgres:postgres@host.docker.internal:5432/app \
  -e APP_DATA_REDIS_URL=redis://host.docker.internal:6379/0 \
  -e APP_JWT_SECRET=your-production-secret \
  go-ddd-template:latest

# 查看日志
docker logs -f go-ddd-app

# 停止容器
docker stop go-ddd-app

# 删除容器
docker rm go-ddd-app
```

## Docker Compose 部署

### 1. 创建 docker-compose.prod.yml

在项目根目录创建 `docker-compose.prod.yml`：

```yaml
version: "3.8"

services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:17-alpine
    container_name: go-ddd-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      POSTGRES_DB: app
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: go-ddd-redis
    restart: unless-stopped
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD:-}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  # Go 应用
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go-ddd-app
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      APP_SERVER_ADDR: ":8080"
      APP_SERVER_ENV: "production"
      APP_DATA_PGSQL_URL: "postgresql://postgres:${POSTGRES_PASSWORD:-postgres}@postgres:5432/app?sslmode=disable"
      APP_DATA_REDIS_URL: "redis://:${REDIS_PASSWORD:-}@redis:6379/0"
      APP_JWT_SECRET: ${JWT_SECRET}
      APP_JWT_ACCESS_TOKEN_EXPIRY: "15m"
      APP_JWT_REFRESH_TOKEN_EXPIRY: "7d"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--no-verbose",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 40s
    networks:
      - app-network

  # Nginx 反向代理（可选）
  nginx:
    image: nginx:alpine
    container_name: go-ddd-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app
    networks:
      - app-network

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  app-network:
    driver: bridge
```

### 2. 创建环境变量文件

创建 `.env.prod` 文件：

```bash
# 数据库配置
POSTGRES_PASSWORD=your-strong-password

# Redis 配置
REDIS_PASSWORD=your-redis-password

# JWT 配置
JWT_SECRET=your-very-secret-jwt-key-change-in-production
```

### 3. 启动服务

```bash
# 启动所有服务
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f app

# 停止服务
docker-compose -f docker-compose.prod.yml down

# 停止并删除数据卷
docker-compose -f docker-compose.prod.yml down -v
```

## Nginx 反向代理配置

创建 `nginx.conf`：

```nginx
events {
    worker_connections 1024;
}

http {
    upstream go_backend {
        server app:8080;
    }

    # 限流配置
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

    server {
        listen 80;
        server_name example.com;

        # 重定向到 HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name example.com;

        # SSL 证书
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        # SSL 配置
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers on;

        # 安全头
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # 日志
        access_log /var/log/nginx/access.log;
        error_log /var/log/nginx/error.log;

        # API 代理
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;

            proxy_pass http://go_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;

            # 超时配置
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # 健康检查
        location /health {
            proxy_pass http://go_backend;
            access_log off;
        }

        # 静态文件（如果有）
        location /docs/ {
            proxy_pass http://go_backend;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
        }
    }
}
```

## Kubernetes 部署

### 1. 创建命名空间

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-ddd
```

### 2. 创建 ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: go-ddd
data:
  APP_SERVER_ADDR: ":8080"
  APP_SERVER_ENV: "production"
  APP_JWT_ACCESS_TOKEN_EXPIRY: "15m"
  APP_JWT_REFRESH_TOKEN_EXPIRY: "7d"
```

### 3. 创建 Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
  namespace: go-ddd
type: Opaque
stringData:
  POSTGRES_PASSWORD: your-strong-password
  REDIS_PASSWORD: your-redis-password
  JWT_SECRET: your-very-secret-jwt-key
  APP_DATA_PGSQL_URL: postgresql://postgres:password@postgres-service:5432/app
  APP_DATA_REDIS_URL: redis://:password@redis-service:6379/0
```

### 4. 创建 Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-ddd-app
  namespace: go-ddd
  labels:
    app: go-ddd-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-ddd-app
  template:
    metadata:
      labels:
        app: go-ddd-app
    spec:
      containers:
        - name: app
          image: go-ddd-template:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              protocol: TCP
          envFrom:
            - configMapRef:
                name: app-config
            - secretRef:
                name: app-secret
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
```

### 5. 创建 Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: go-ddd-service
  namespace: go-ddd
spec:
  type: ClusterIP
  selector:
    app: go-ddd-app
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
```

### 6. 创建 Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-ddd-ingress
  namespace: go-ddd
  annotations:
    nginx.ingress.kubernetes.io/rate-limit: "100"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - api.example.com
      secretName: go-ddd-tls
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: go-ddd-service
                port:
                  number: 80
```

### 7. 部署到 Kubernetes

```bash
# 应用所有配置
kubectl apply -f k8s/

# 查看部署状态
kubectl get all -n go-ddd

# 查看 Pod 日志
kubectl logs -f -n go-ddd -l app=go-ddd-app

# 扩容
kubectl scale deployment go-ddd-app -n go-ddd --replicas=5

# 更新镜像
kubectl set image deployment/go-ddd-app app=go-ddd-template:v2.0.0 -n go-ddd

# 回滚
kubectl rollout undo deployment/go-ddd-app -n go-ddd

# 查看部署历史
kubectl rollout history deployment/go-ddd-app -n go-ddd
```

## 云平台部署

### AWS ECS

```bash
# 1. 推送镜像到 ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin your-account.dkr.ecr.us-west-2.amazonaws.com
docker tag go-ddd-template:latest your-account.dkr.ecr.us-west-2.amazonaws.com/go-ddd-template:latest
docker push your-account.dkr.ecr.us-west-2.amazonaws.com/go-ddd-template:latest

# 2. 创建 ECS 任务定义和服务（使用 AWS Console 或 CLI）
```

### Google Cloud Run

```bash
# 1. 推送镜像到 GCR
gcloud builds submit --tag gcr.io/your-project/go-ddd-template

# 2. 部署到 Cloud Run
gcloud run deploy go-ddd-app \
  --image gcr.io/your-project/go-ddd-template \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars APP_SERVER_ADDR=:8080 \
  --set-secrets APP_JWT_SECRET=jwt-secret:latest
```

### Azure Container Instances

```bash
# 1. 推送镜像到 ACR
az acr build --registry myregistry --image go-ddd-template:latest .

# 2. 部署到 ACI
az container create \
  --resource-group myResourceGroup \
  --name go-ddd-app \
  --image myregistry.azurecr.io/go-ddd-template:latest \
  --dns-name-label go-ddd-app \
  --ports 8080
```

## 二进制部署

### 1. 构建二进制文件

```bash
# 本地构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

# 或使用 Task（如果配置了）
task build
```

### 2. 创建 systemd 服务

创建 `/etc/systemd/system/go-ddd-app.service`：

```ini
[Unit]
Description=Go DDD Application
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/go-ddd-app
ExecStart=/opt/go-ddd-app/app api
Restart=always
RestartSec=10

# 环境变量
Environment="APP_SERVER_ADDR=:8080"
Environment="APP_SERVER_ENV=production"
Environment="APP_DATA_PGSQL_URL=postgresql://postgres:password@localhost:5432/app"
Environment="APP_DATA_REDIS_URL=redis://localhost:6379/0"
Environment="APP_JWT_SECRET=your-secret"

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/go-ddd-app

[Install]
WantedBy=multi-user.target
```

### 3. 启动服务

```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start go-ddd-app

# 开机自启
sudo systemctl enable go-ddd-app

# 查看状态
sudo systemctl status go-ddd-app

# 查看日志
sudo journalctl -u go-ddd-app -f
```

## 生产环境最佳实践

### 1. 安全配置

```bash
# 使用强密码和密钥
APP_JWT_SECRET=$(openssl rand -base64 32)

# 启用 HTTPS
# 使用 Let's Encrypt 或购买 SSL 证书

# 限制访问（防火墙）
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. 监控和日志

```yaml
# 集成 Prometheus（在 docker-compose 中添加）
prometheus:
  image: prom/prometheus:latest
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
  ports:
    - "9090:9090"

# 集成 Grafana
grafana:
  image: grafana/grafana:latest
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
```

### 3. 备份策略

```bash
# PostgreSQL 备份
docker exec go-ddd-postgres pg_dump -U postgres app > backup_$(date +%Y%m%d).sql

# 自动备份（cron）
0 2 * * * /usr/local/bin/backup-db.sh
```

### 4. 性能优化

- 启用连接池（GORM 默认启用）
- 配置 Redis 缓存策略
- 使用 CDN 加速静态资源
- 启用 gzip 压缩
- 配置数据库索引

### 5. 高可用配置

- 使用负载均衡（如 Nginx、HAProxy、AWS ALB）
- 配置数据库主从复制
- 配置 Redis 哨兵或集群
- 使用容器编排（Kubernetes）

## 故障排查

### 检查应用日志

```bash
# Docker
docker logs -f go-ddd-app

# Kubernetes
kubectl logs -f -n go-ddd -l app=go-ddd-app

# systemd
journalctl -u go-ddd-app -f
```

### 检查健康状态

```bash
# 健康检查
curl http://localhost:8080/health

# 检查数据库连接
docker exec go-ddd-postgres psql -U postgres -c "SELECT 1"

# 检查 Redis 连接
docker exec go-ddd-redis redis-cli ping
```

### 常见问题

| 问题            | 原因           | 解决方案                     |
| --------------- | -------------- | ---------------------------- |
| 容器启动失败    | 环境变量未设置 | 检查 `.env` 文件和环境变量   |
| 数据库连接失败  | 数据库未就绪   | 使用 `depends_on` 和健康检查 |
| 502 Bad Gateway | 应用未启动     | 检查应用日志和健康检查       |
| 内存溢出        | 资源限制不足   | 增加 Docker/K8s 资源限制     |

## 相关链接

- [快速开始](/guide/getting-started) - 本地开发环境配置
- [配置系统](/guide/configuration) - 配置选项说明
- [PostgreSQL](/guide/postgresql) - 数据库配置
- [Redis](/guide/redis) - 缓存配置
- [文档部署](/guide/docs-deployment) - VitePress 文档部署
