# Mermaid 图表示例

本页面展示如何在 VitePress 文档中使用 Mermaid 图表。

## 流程图 (Flowchart)

```mermaid
flowchart TD
    A[开始] --> B{是否登录?}
    B -->|是| C[显示主页]
    B -->|否| D[跳转到登录页]
    D --> E[用户登录]
    E --> C
    C --> F[结束]
```

## 时序图 (Sequence Diagram)

```mermaid
sequenceDiagram
    participant 用户
    participant 前端
    participant 后端
    participant 数据库

    用户->>前端: 提交登录表单
    前端->>后端: POST /api/auth/login
    后端->>数据库: 查询用户信息
    数据库-->>后端: 返回用户数据
    后端-->>前端: 返回 JWT Token
    前端-->>用户: 登录成功
```

## 类图 (Class Diagram)

```mermaid
classDiagram
    class User {
        +ID uint
        +Username string
        +Email string
        +Password string
        +CreatedAt time
        +UpdatedAt time
        +DeletedAt time
        +Create()
        +Update()
        +Delete()
    }

    class UserRepository {
        <<interface>>
        +Create(user)
        +FindByID(id)
        +FindByUsername(username)
        +Update(user)
        +Delete(id)
    }

    class PostgresUserRepository {
        -db *gorm.DB
        +Create(user)
        +FindByID(id)
        +FindByUsername(username)
        +Update(user)
        +Delete(id)
    }

    UserRepository <|-- PostgresUserRepository
    PostgresUserRepository --> User
```

## 状态图 (State Diagram)

```mermaid
stateDiagram-v2
    [*] --> 未登录
    未登录 --> 登录中: 提交登录
    登录中 --> 已登录: 登录成功
    登录中 --> 未登录: 登录失败
    已登录 --> 未登录: 登出
    已登录 --> [*]
```

## ER 图 (Entity Relationship)

```mermaid
erDiagram
    User ||--o{ Post : creates
    User ||--o{ Comment : writes
    Post ||--o{ Comment : has

    User {
        uint id PK
        string username
        string email
        timestamp created_at
    }

    Post {
        uint id PK
        uint user_id FK
        string title
        text content
        timestamp created_at
    }

    Comment {
        uint id PK
        uint user_id FK
        uint post_id FK
        text content
        timestamp created_at
    }
```

## Git 分支图 (Gitgraph)

```mermaid
gitGraph
    commit
    commit
    branch develop
    checkout develop
    commit
    commit
    checkout main
    merge develop
    commit
    branch feature
    checkout feature
    commit
    checkout develop
    merge feature
    checkout main
    merge develop
```

## 甘特图 (Gantt)

```mermaid
gantt
    title 项目开发计划
    dateFormat YYYY-MM-DD
    section 设计阶段
    需求分析           :a1, 2025-01-01, 7d
    架构设计           :a2, after a1, 5d
    section 开发阶段
    后端开发           :b1, after a2, 15d
    前端开发           :b2, after a2, 15d
    section 测试阶段
    单元测试           :c1, after b1, 5d
    集成测试           :c2, after b2, 5d
    section 部署阶段
    上线部署           :d1, after c2, 2d
```

## 饼图 (Pie Chart)

```mermaid
pie title 技术栈占比
    "Go" : 45
    "Vue.js" : 30
    "PostgreSQL" : 15
    "Redis" : 10
```

## 思维导图 (Mindmap)

```mermaid
mindmap
  root((DDD 架构))
    表现层
      HTTP Handler
      路由
      中间件
    应用层
      用例
      DTO
      服务协调
    领域层
      实体
      值对象
      领域服务
      仓储接口
    基础设施层
      数据库实现
      缓存实现
      配置管理
```

## 时间线 (Timeline)

```mermaid
timeline
    title 项目发展历程
    2025-01 : 项目启动 : 完成架构设计
    2025-02 : 核心功能开发 : 用户管理 : JWT 认证
    2025-03 : 完善功能 : Redis 缓存 : 文档系统
    2025-04 : 测试与优化 : 单元测试 : 性能优化
    2025-05 : 正式发布 : v1.0.0
```

## 使用说明

在 Markdown 文件中使用 Mermaid 图表非常简单，只需使用标准的 Mermaid 代码块语法：

````markdown
```mermaid
flowchart LR
    A[开始] --> B[处理]
    B --> C[结束]
```
````

VitePress 会自动将其渲染为可交互的图表。

## 特性

- ✅ 支持所有 Mermaid 图表类型
- ✅ 自动适配亮色/暗色主题
- ✅ 响应式设计，适配移动端
- ✅ 与 VitePress 2.0 完美集成
- ✅ 标准 Markdown 代码块语法
- ✅ 保留换行符，语法正确

## 参考资料

- [Mermaid 官方文档](https://mermaid.js.org/)
- [Mermaid 在线编辑器](https://mermaid.live/)
