# Koanf 配置系统重构完成 ✅

## 📋 完成的任务

- [x] 安装 koanf 依赖包
- [x] 更新 config.go 添加 koanf 标签
- [x] 创建配置加载器 loader.go
- [x] 修改 api.go 集成 koanf
- [x] 创建示例配置文件
- [x] 测试配置加载和优先级

## 🎯 实现的功能

### 1. 多层级配置系统

配置按以下优先级自动合并（从低到高）：

```
默认值 → 配置文件 → 环境变量 → 命令行参数
  ⭐      ⭐⭐       ⭐⭐⭐        ⭐⭐⭐⭐
```

### 2. 支持的配置来源

#### ① 默认值（代码中定义）
```go
Server: ServerConfig{
    Addr: "0.0.0.0:8080",
    Env:  "development",
}
```

#### ② 配置文件
```bash
# 自动搜索以下位置：
- config.yaml
- configs/config.yaml
```

#### ③ 环境变量
```bash
# 命名规则：APP_<SECTION>_<KEY>
APP_SERVER_ADDR=:8880
APP_SERVER_ENV=production
APP_DATABASE_HOST=db.example.com
```

#### ④ 命令行参数
```bash
./vmalert api --addr :9999
./vmalert api -a :8080
```

## 📁 项目结构

```
├── configs/
│   ├── config.example.yaml  # 配置模板
│   └── README.md            # 配置说明文档
├── internal/
│   └── config/
│       ├── config.go        # 配置结构（含 koanf 标签）
│       └── loader.go        # Koanf 加载器
└── internal/commands/api/
    └── api.go              # 集成 koanf 的 API 命令
```

## ✅ 测试结果

所有配置优先级测试通过：

| 测试场景 | 配置来源 | 结果 | 状态 |
|---------|---------|------|------|
| 默认值 | 代码默认 | `0.0.0.0:8080` | ✅ |
| 环境变量 | `APP_SERVER_ADDR=:7777` | `:7777` | ✅ |
| 命令行参数 | `--addr :9999` | `:9999` | ✅ |
| 优先级测试 | 环境变量 + 命令行 | 命令行胜出 | ✅ |

## 🚀 使用示例

### 场景 1: 开发环境（使用默认值）
```bash
./vmalert api
# → 0.0.0.0:8080
```

### 场景 2: 使用配置文件
```bash
cp configs/config.example.yaml config.yaml
# 编辑 config.yaml...
./vmalert api
# → 读取 config.yaml 中的配置
```

### 场景 3: Docker 容器（使用环境变量）
```bash
docker run -e APP_SERVER_ADDR=:8080 \
           -e APP_DATABASE_HOST=postgres \
           myapp
```

### 场景 4: 临时覆盖（命令行参数）
```bash
./vmalert api --addr :3000
# → 临时在 3000 端口运行
```

## 🔧 关键技术点

1. **Koanf Provider Chain**
   - `structs.Provider` - 默认值
   - `file.Provider` - YAML 文件
   - `env.Provider` - 环境变量
   - `posflag.Provider` - 命令行参数

2. **环境变量转换**
   ```
   APP_SERVER_ADDR → server.addr
   APP_DATABASE_HOST → database.host
   ```

3. **命令行参数桥接**
   - 使用 `cmd.IsSet()` 检测用户是否显式设置
   - 避免默认值覆盖环境变量

## 📚 相关文档

- [Koanf GitHub](https://github.com/knadh/koanf)
- [配置说明文档](configs/README.md)
- [示例配置文件](configs/config.example.yaml)

## 🎉 优势

相比之前的实现：

- ✅ **更灵活** - 支持 4 种配置来源
- ✅ **更轻量** - koanf 比 viper 小 10 倍
- ✅ **更清晰** - 配置优先级明确
- ✅ **更专业** - 符合 12-factor app 标准
- ✅ **类型安全** - 直接 Unmarshal 到结构体
- ✅ **易扩展** - 可添加远程配置、配置热重载等
