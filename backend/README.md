# Marsview 后端 API 服务

基于 Go 语言的后端 API 服务，为 Marsview 低代码设计平台提供数据存储和配置加载功能。

## 技术栈

- **Web 框架**：[Gin](https://github.com/gin-gonic/gin)
- **ORM**：[GORM](https://gorm.io/)
- **数据库**：SQLite（默认，可扩展为 MySQL）
- **认证**：JWT（[golang-jwt/jwt](https://github.com/golang-jwt/jwt)）
- **密码哈希**：bcrypt

## 快速启动

```bash
cd backend
go mod tidy
go run ./cmd/server
```

服务默认监听 `http://localhost:8000`。

### 环境变量

| 变量         | 默认值                      | 说明                  |
|--------------|-----------------------------|-----------------------|
| `PORT`       | `8000`                      | 服务监听端口          |
| `GIN_MODE`   | `debug`                     | Gin 运行模式          |
| `DB_DRIVER`  | `sqlite`                    | 数据库驱动            |
| `DB_DSN`     | `marsview.db`               | 数据库连接字符串      |
| `JWT_SECRET` | `marsview-secret-key-2024`  | JWT 签名密钥          |

### 默认管理员账号

首次启动时自动创建：

- 用户名：`admin@marsview.cc`
- 密码：`marsview2024`

## API 接口列表

### 用户
| 方法   | 路径             | 说明           | 需要认证 |
|--------|------------------|----------------|----------|
| POST   | /user/login      | 用户登录       | 否       |
| GET    | /user/info       | 获取当前用户   | 是       |
| GET    | /user/search     | 搜索用户       | 是       |

### 项目
| 方法   | 路径                          | 说明           | 需要认证 |
|--------|-------------------------------|----------------|----------|
| GET    | /admin/project/list           | 项目列表       | 是       |
| POST   | /admin/project/create         | 创建项目       | 是       |
| PUT    | /admin/project/update         | 更新项目       | 是       |
| DELETE | /admin/project/delete/:id     | 删除项目       | 是       |
| GET    | /admin/getProjectConfig       | 获取项目配置   | 是       |

### 页面
| 方法   | 路径                          | 说明               | 需要认证 |
|--------|-------------------------------|--------------------|----------|
| GET    | /admin/page/list              | 页面列表           | 是       |
| GET    | /admin/page/detail/:env/:id   | 页面详情           | 是       |
| POST   | /admin/page/create            | 创建页面           | 是       |
| POST   | /admin/page/update            | 更新页面数据       | 是       |
| DELETE | /admin/page/delete/:id        | 删除页面           | 是       |
| POST   | /admin/page/copy              | 复制页面           | 是       |
| POST   | /admin/page/publish           | 发布页面           | 是       |
| GET    | /admin/page/publishList       | 发布历史           | 是       |
| POST   | /admin/page/rollback          | 回滚页面           | 是       |

### 菜单
| 方法   | 路径                          | 说明           | 需要认证 |
|--------|-------------------------------|----------------|----------|
| GET    | /admin/menu/list/:projectId   | 菜单列表(树形) | 是       |
| POST   | /admin/menu/add               | 新增菜单       | 是       |
| PUT    | /admin/menu/update            | 更新菜单       | 是       |
| DELETE | /admin/menu/delete/:id        | 删除菜单       | 是       |
| POST   | /admin/menu/copy              | 复制菜单       | 是       |

### 角色
| 方法   | 路径                          | 说明           | 需要认证 |
|--------|-------------------------------|----------------|----------|
| GET    | /admin/role/list              | 角色分页列表   | 是       |
| GET    | /admin/role/all               | 全部角色       | 是       |
| POST   | /admin/role/create            | 创建角色       | 是       |
| DELETE | /admin/role/delete/:id        | 删除角色       | 是       |
| PUT    | /admin/role/update            | 更新角色信息   | 是       |
| PUT    | /admin/role/permissions       | 更新角色权限   | 是       |

### 项目用户管理
| 方法   | 路径                          | 说明           | 需要认证 |
|--------|-------------------------------|----------------|----------|
| GET    | /admin/user/list              | 项目用户列表   | 是       |
| POST   | /admin/user/add               | 添加用户       | 是       |
| DELETE | /admin/user/delete/:id        | 移除用户       | 是       |
| PUT    | /admin/user/update            | 更新用户角色   | 是       |

### 页面成员
| 方法   | 路径                          | 说明           | 需要认证 |
|--------|-------------------------------|----------------|----------|
| GET    | /admin/page/member/list       | 页面成员列表   | 是       |
| POST   | /admin/page/member/add        | 添加成员       | 是       |
| DELETE | /admin/page/member/delete/:id | 移除成员       | 是       |

## 统一响应格式

```json
{
  "code": 0,
  "data": {},
  "message": "success"
}
```

- `code = 0`：成功
- `code = 400`：请求参数错误
- `code = 10018`：未登录或 Token 过期

## 运行测试

```bash
cd backend
go test ./... -v
```

## 项目结构

```
backend/
├── cmd/server/main.go          # 启动入口
├── internal/
│   ├── config/config.go        # 配置加载
│   ├── middleware/auth.go      # JWT鉴权 & CORS
│   ├── model/model.go          # 数据库模型
│   ├── handler/                # HTTP处理器
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── page.go
│   │   ├── menu.go
│   │   ├── role.go
│   │   ├── project_user.go
│   │   └── page_member.go
│   └── router/router.go        # 路由注册
├── go.mod
└── go.sum
```
