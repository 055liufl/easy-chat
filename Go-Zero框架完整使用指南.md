# Go-Zero 框架完整使用指南

> 作者：木兮老师  
> 版本：v2.0  
> 更新日期：2026-03-17  
> 适用人群：Go-Zero 初学者、微服务开发者、后端工程师

---

## 目录

1. [Go-Zero 框架简介](#1-go-zero-框架简介)
2. [环境搭建](#2-环境搭建)
3. [API 服务开发](#3-api-服务开发)
4. [RPC 服务开发](#4-rpc-服务开发)
5. [配置管理](#5-配置管理)
6. [数据库操作](#6-数据库操作)
7. [中间件开发](#7-中间件开发)
8. [服务间通信](#8-服务间通信)
9. [缓存使用](#9-缓存使用)
10. [消息队列](#10-消息队列)
11. [日志和监控](#11-日志和监控)
12. [测试](#12-测试)
13. [部署](#13-部署)
14. [最佳实践](#14-最佳实践)
15. [常见问题](#15-常见问题)
16. [goctl 工具详解](#16-goctl-工具详解)
17. [完整项目示例](#17-完整项目示例)

---

## 1. Go-Zero 框架简介

### 1.1 什么是 Go-Zero？

Go-Zero 是一个集成了各种工程实践的 Web 和 RPC 框架。它由好未来（原学而思）开源，经过了大规模生产环境的验证。

**核心特性：**

- **高性能**：基于 Go 语言，天然支持高并发
- **微服务架构**：内置服务发现、负载均衡、熔断降级等功能
- **代码生成**：通过 goctl 工具自动生成代码，提高开发效率
- **API 网关**：支持 RESTful API 和 gRPC
- **中间件丰富**：内置日志、限流、熔断、链路追踪等中间件
- **配置管理**：支持本地配置和配置中心（ETCD）
- **数据库支持**：支持 MySQL、MongoDB、Redis 等

# Go-Zero 框架完整使用指南（续）

本文档是《Go-Zero框架完整使用指南.md》的补充，包含更多详细的使用示例和说明。

---

## 继续第1章内容

### 1.2 Go-Zero 的优势

**与其他框架对比：**

| 特性 | Go-Zero | Gin | Beego |
|------|---------|-----|-------|
| 代码生成 | ✅ 完整支持 | ❌ 不支持 | ⚠️ 部分支持 |
| 微服务 | ✅ 内置 | ❌ 需要自己实现 | ⚠️ 部分支持 |
| 服务发现 | ✅ 内置 | ❌ 需要自己实现 | ❌ 需要自己实现 |
| 熔断降级 | ✅ 内置 | ❌ 需要自己实现 | ❌ 需要自己实现 |
| 限流 | ✅ 内置 | ❌ 需要自己实现 | ⚠️ 部分支持 |
| 链路追踪 | ✅ 内置 | ❌ 需要自己实现 | ❌ 需要自己实现 |

### 1.3 Go-Zero 架构图

```
┌─────────────────────────────────────────────────────────┐
│                     客户端                               │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│                  API Gateway                             │
│  - 路由转发                                              │
│  - 认证鉴权                                              │
│  - 限流熔断                                              │
└────────────────────┬────────────────────────────────────┘
                     │
         ┌───────────┼───────────┐
         │           │           │
         ↓           ↓           ↓
    ┌────────┐  ┌────────┐  ┌────────┐
    │ API 1  │  │ API 2  │  │ API 3  │
    └───┬────┘  └───┬────┘  └───┬────┘
        │           │           │
        │    gRPC   │    gRPC   │
        ↓           ↓           ↓
    ┌────────┐  ┌────────┐  ┌────────┐
    │ RPC 1  │  │ RPC 2  │  │ RPC 3  │
    └───┬────┘  └───┬────┘  └───┬────┘
        │           │           │
        └───────────┼───────────┘
                    │
                    ↓
         ┌──────────────────────┐
         │   数据层（DB/Cache）  │
         └──────────────────────┘
```

---

## 2. 环境搭建

### 2.1 Go 环境安装

#### 2.1.1 Linux 安装

```bash
# 1. 下载 Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz

# 2. 解压到 /usr/local
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 3. 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# 4. 使配置生效
source ~/.bashrc

# 5. 验证安装
go version
# 输出：go version go1.21.0 linux/amd64
```

#### 2.1.2 配置 Go 环境

```bash
# 设置 Go 代理（国内用户必须）
go env -w GOPROXY=https://goproxy.cn,direct

# 开启 Go Modules
go env -w GO111MODULE=on

# 设置私有仓库（如果有）
go env -w GOPRIVATE=github.com/yourcompany/*

# 查看所有配置
go env
```

**重要配置说明：**

```
GOPROXY：Go 模块代理
  - 国内推荐：https://goproxy.cn,direct
  - 官方：https://proxy.golang.org,direct
  - 作用：加速依赖下载

GO111MODULE：Go Modules 开关
  - on：强制使用 Go Modules
  - off：使用 GOPATH 模式
  - auto：自动判断

GOPRIVATE：私有仓库
  - 设置后不会通过代理下载
  - 支持通配符
```

### 2.2 goctl 工具安装

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
goctl --version
# 输出：goctl version 1.6.0 linux/amd64

# 查看帮助
goctl --help
```

**goctl 命令概览：**

```bash
goctl api      # API 服务相关命令
goctl rpc      # RPC 服务相关命令
goctl model    # 数据库模型生成
goctl docker   # Docker 文件生成
goctl kube     # Kubernetes 配置生成
goctl template # 模板管理
goctl upgrade  # 升级 goctl
```

### 2.3 第一个 Hello World 项目

#### 步骤 1：创建项目目录

```bash
# 创建项目目录
mkdir hello-zero
cd hello-zero

# 初始化 Go Modules
go mod init hello-zero
```

#### 步骤 2：创建 API 定义文件

创建 `hello.api` 文件：

```go
// hello.api
syntax = "v1"

// API 信息
info(
    title: "Hello World API"
    desc: "第一个 Go-Zero 项目"
    author: "你的名字"
    version: "v1.0"
)

// 请求结构体
type HelloReq {
    Name string `json:"name"` // 用户名
}

// 响应结构体
type HelloResp {
    Message string `json:"message"` // 返回消息
}

// 服务定义
service hello {
    @handler hello
    post /hello (HelloReq) returns (HelloResp)
}
```

#### 步骤 3：生成代码

```bash
# 使用 goctl 生成代码
goctl api go -api hello.api -dir .

# 生成的目录结构：
# .
# ├── etc/
# │   └── hello.yaml          # 配置文件
# ├── internal/
# │   ├── config/
# │   │   └── config.go       # 配置结构体
# │   ├── handler/
# │   │   ├── hellohandler.go # HTTP 处理器
# │   │   └── routes.go       # 路由注册
# │   ├── logic/
# │   │   └── hellologic.go   # 业务逻辑
# │   ├── svc/
# │   │   └── servicecontext.go # 服务上下文
# │   └── types/
# │       └── types.go        # 数据结构
# ├── hello.api
# └── hello.go                # 主程序入口
```

#### 步骤 4：实现业务逻辑

打开 `internal/logic/hellologic.go`：

```go
package logic

import (
    "context"
    "fmt"

    "hello-zero/internal/svc"
    "hello-zero/internal/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type HelloLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
    return &HelloLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// Hello 处理 hello 请求
func (l *HelloLogic) Hello(req *types.HelloReq) (resp *types.HelloResp, err error) {
    // 打印日志
    l.Logger.Infof("收到请求，Name: %s", req.Name)

    // 构造响应
    resp = &types.HelloResp{
        Message: fmt.Sprintf("Hello, %s! 欢迎使用 Go-Zero!", req.Name),
    }

    return resp, nil
}
```

#### 步骤 5：运行服务

```bash
# 下载依赖
go mod tidy

# 运行服务
go run hello.go

# 输出：
# Starting server at 0.0.0.0:8888...
```

#### 步骤 6：测试接口

```bash
# 使用 curl 测试
curl -X POST http://localhost:8888/hello \
  -H "Content-Type: application/json" \
  -d '{"name":"张三"}'

# 响应：
# {"message":"Hello, 张三! 欢迎使用 Go-Zero!"}
```

**恭喜！你已经完成了第一个 Go-Zero 项目！** 🎉

---

## 3. API 服务开发详解

### 3.1 API 文件完整语法

#### 3.1.1 基本结构

```go
// 1. 语法版本（必须）
syntax = "v1"

// 2. API 信息（可选，但建议添加）
info(
    title: "用户服务 API"
    desc: "用户注册、登录、信息管理"
    author: "张三"
    email: "zhangsan@example.com"
    version: "v1.0"
)

// 3. 导入其他 API 文件（可选）
import "base.api"
import "user.api"

// 4. 数据结构定义
type User {
    Id       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 5. 服务定义
service user-api {
    @handler getUser
    get /user/:id (GetUserReq) returns (GetUserResp)
}
```

#### 3.1.2 数据类型详解

```go
type DataTypes {
    // 基本类型
    IntField    int     `json:"intField"`      // 整数
    Int8Field   int8    `json:"int8Field"`     // 8位整数
    Int16Field  int16   `json:"int16Field"`    // 16位整数
    Int32Field  int32   `json:"int32Field"`    // 32位整数
    Int64Field  int64   `json:"int64Field"`    // 64位整数
    
    UintField   uint    `json:"uintField"`     // 无符号整数
    Uint8Field  uint8   `json:"uint8Field"`    // 8位无符号整数
    Uint16Field uint16  `json:"uint16Field"`   // 16位无符号整数
    Uint32Field uint32  `json:"uint32Field"`   // 32位无符号整数
    Uint64Field uint64  `json:"uint64Field"`   // 64位无符号整数
    
    Float32Field float32 `json:"float32Field"` // 32位浮点数
    Float64Field float64 `json:"float64Field"` // 64位浮点数
    
    StringField string  `json:"stringField"`   // 字符串
    BoolField   bool    `json:"boolField"`     // 布尔值
    ByteField   byte    `json:"byteField"`     // 字节
    
    // 数组/切片
    IntArray    []int    `json:"intArray"`     // 整数数组
    StringArray []string `json:"stringArray"`  // 字符串数组
    UserArray   []User   `json:"userArray"`    // 结构体数组
    
    // Map
    MapField map[string]string `json:"mapField"` // Map
    MapField2 map[string]int   `json:"mapField2"` // Map
    
    // 嵌套结构体
    UserInfo User `json:"userInfo"`            // 嵌套用户信息
    
    // 指针类型（可选字段）
    OptionalField *string `json:"optionalField,optional"` // 可选字段
    OptionalInt   *int    `json:"optionalInt,optional"`   // 可选整数
}
```

#### 3.1.3 标签详解

```go
type TagExample {
    // json 标签：指定 JSON 字段名
    Id string `json:"id"`
    
    // optional 标签：可选字段（请求时可以不传）
    Nickname string `json:"nickname,optional"`
    
    // omitempty 标签：为空时不返回
    Avatar string `json:"avatar,omitempty"`
    
    // default 标签：默认值
    Status int `json:"status,default=1"`
    
    // options 标签：枚举值
    Gender string `json:"gender,options=male|female|unknown"`
    
    // range 标签：数值范围
    Age int `json:"age,range=[0:150]"`
    
    // path 标签：路径参数
    UserId string `path:"userId"`
    
    // form 标签：表单参数
    Username string `form:"username"`
    
    // header 标签：请求头参数
    Token string `header:"Authorization"`
}
```

#### 3.1.4 路由定义详解

```go
service user-api {
    // GET 请求 - 获取单个资源
    @handler getUser
    get /user/:id (GetUserReq) returns (GetUserResp)
    
    // GET 请求 - 获取列表
    @handler listUsers
    get /users (ListUsersReq) returns (ListUsersResp)
    
    // POST 请求 - 创建资源
    @handler createUser
    post /user (CreateUserReq) returns (CreateUserResp)
    
    // PUT 请求 - 更新资源
    @handler updateUser
    put /user/:id (UpdateUserReq) returns (UpdateUserResp)
    
    // DELETE 请求 - 删除资源
    @handler deleteUser
    delete /user/:id (DeleteUserReq) returns (DeleteUserResp)
    
    // PATCH 请求 - 部分更新
    @handler patchUser
    patch /user/:id (PatchUserReq) returns (PatchUserResp)
}
```

#### 3.1.5 路由分组和中间件

```go
// 不需要认证的接口
@server(
    prefix: /api/v1    // 路由前缀
    group: user        // 分组名称（生成的代码会放在 user 目录下）
)
service user-api {
    @doc "用户注册"
    @handler register
    post /register (RegisterReq) returns (RegisterResp)
    
    @doc "用户登录"
    @handler login
    post /login (LoginReq) returns (LoginResp)
}

// 需要 JWT 认证的接口
@server(
    prefix: /api/v1
    group: user
    jwt: Auth          // 启用 JWT 认证
)
service user-api {
    @doc "获取用户信息"
    @handler getUserInfo
    get /user/info (GetUserInfoReq) returns (GetUserInfoResp)
    
    @doc "更新用户信息"
    @handler updateUserInfo
    put /user/info (UpdateUserInfoReq) returns (UpdateUserInfoResp)
}

// 使用多个中间件
@server(
    prefix: /api/v1
    group: user
    middleware: Log, Auth, Limit  // 使用多个中间件（按顺序执行）
)
service user-api {
    @handler getUser
    get /user/:id (GetUserReq) returns (GetUserResp)
}

// 设置超时时间
@server(
    prefix: /api/v1
    group: user
    timeout: 3s        // 设置超时时间为 3 秒
)
service user-api {
    @handler slowOperation
    post /slow (SlowReq) returns (SlowResp)
}
```

---


## 3.3 完整的用户服务示例

### 3.3.1 用户注册功能

#### Handler 层实现

Handler 层负责接收 HTTP 请求，解析参数，调用 Logic 层处理业务逻辑，返回响应。

```go
// apps/user/api/internal/handler/user/registerhandler.go
package user

import (
    "net/http"
    
    "github.com/zeromicro/go-zero/rest/httpx"
    "imooc.com/easy-chat/apps/user/api/internal/logic/user"
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    "imooc.com/easy-chat/apps/user/api/internal/types"
)

// RegisterHandler 用户注册处理器
// 参数: svcCtx - 服务上下文，包含配置、依赖等
// 返回: http.HandlerFunc - HTTP 处理函数
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 定义请求结构体
        var req types.RegisterReq
        
        // 2. 解析请求参数（支持 JSON、Form、Path、Query 等）
        if err := httpx.Parse(r, &req); err != nil {
            // 参数解析失败，返回错误
            httpx.Error(w, err)
            return
        }
        
        // 3. 创建 Logic 实例
        l := user.NewRegisterLogic(r.Context(), svcCtx)
        
        // 4. 调用 Logic 层处理业务逻辑
        resp, err := l.Register(&req)
        if err != nil {
            // 业务处理失败，返回错误
            httpx.Error(w, err)
        } else {
            // 业务处理成功，返回 JSON 响应
            httpx.OkJson(w, resp)
        }
    }
}
```

#### Logic 层实现（API 层）

API 层的 Logic 负责参数转换和调用 RPC 服务。

```go
// apps/user/api/internal/logic/user/registerlogic.go
package user

import (
    "context"
    "github.com/jinzhu/copier"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    "imooc.com/easy-chat/apps/user/api/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

// RegisterLogic 注册业务逻辑
type RegisterLogic struct {
    logx.Logger                    // 日志记录器
    ctx    context.Context          // 上下文
    svcCtx *svc.ServiceContext      // 服务上下文
}

// NewRegisterLogic 创建注册逻辑实例
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
    return &RegisterLogic{
        Logger: logx.WithContext(ctx),  // 带上下文的日志
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// Register 注册方法
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 1. 调用 RPC 服务进行注册
    registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
        Phone:    req.Phone,      // 手机号
        Nickname: req.Nickname,   // 昵称
        Password: req.Password,   // 密码
        Avatar:   req.Avatar,     // 头像
        Sex:      int32(req.Sex), // 性别
    })
    if err != nil {
        return nil, err
    }
    
    // 2. 使用 copier 复制数据（避免手动赋值）
    var res types.RegisterResp
    copier.Copy(&res, registerResp)
    
    return &res, nil
}
```

#### Logic 层实现（RPC 层）

RPC 层的 Logic 负责实际的业务逻辑处理。

```go
// apps/user/rpc/internal/logic/registerlogic.go
package logic

import (
    "context"
    "database/sql"
    "errors"
    "imooc.com/easy-chat/apps/user/models"
    "imooc.com/easy-chat/pkg/ctxdata"
    "imooc.com/easy-chat/pkg/encrypt"
    "imooc.com/easy-chat/pkg/wuid"
    "time"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/logx"
)

var (
    // ErrPhoneIsRegister 手机号已注册错误
    ErrPhoneIsRegister = errors.New("手机号已经注册过")
)

// RegisterLogic 注册逻辑
type RegisterLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

// NewRegisterLogic 创建注册逻辑实例
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
    return &RegisterLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// Register 注册方法
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
    // 1. 验证用户是否已注册（根据手机号）
    userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
    if err != nil && err != models.ErrNotFound {
        // 数据库查询错误
        return nil, err
    }
    
    if userEntity != nil {
        // 用户已存在
        return nil, ErrPhoneIsRegister
    }
    
    // 2. 创建用户实体
    userEntity = &models.Users{
        Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource), // 生成分布式 ID
        Avatar:   in.Avatar,                                      // 头像
        Nickname: in.Nickname,                                    // 昵称
        Phone:    in.Phone,                                       // 手机号
        Sex: sql.NullInt64{                                       // 性别（可为空）
            Int64: int64(in.Sex),
            Valid: true,
        },
    }
    
    // 3. 加密密码
    if len(in.Password) > 0 {
        genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
        if err != nil {
            return nil, err
        }
        userEntity.Password = sql.NullString{
            String: string(genPassword),
            Valid:  true,
        }
    }
    
    // 4. 插入数据库
    _, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
    if err != nil {
        return nil, err
    }
    
    // 5. 生成 JWT Token
    now := time.Now().Unix()
    token, err := ctxdata.GetJwtToken(
        l.svcCtx.Config.Jwt.AccessSecret,  // JWT 密钥
        now,                                // 当前时间
        l.svcCtx.Config.Jwt.AccessExpire,  // 过期时间
        userEntity.Id,                      // 用户 ID
    )
    if err != nil {
        return nil, err
    }
    
    // 6. 返回响应
    return &user.RegisterResp{
        Token:  token,                                      // JWT Token
        Expire: now + l.svcCtx.Config.Jwt.AccessExpire,    // 过期时间戳
    }, nil
}
```

### 3.3.2 用户登录功能

#### API 定义

```api
// apps/user/api/user.api
type (
    // 登录请求
    LoginReq {
        Phone    string `json:"phone"`     // 手机号
        Password string `json:"password"`  // 密码
    }
    
    // 登录响应
    LoginResp {
        Token  string `json:"token"`   // JWT Token
        Expire int64  `json:"expire"`  // 过期时间戳
    }
)

@server(
    prefix: /api/user
    group: user
)
service user-api {
    @handler login
    post /login (LoginReq) returns (LoginResp)
}
```

#### Handler 层

```go
// apps/user/api/internal/handler/user/loginhandler.go
package user

import (
    "net/http"
    
    "github.com/zeromicro/go-zero/rest/httpx"
    "imooc.com/easy-chat/apps/user/api/internal/logic/user"
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    "imooc.com/easy-chat/apps/user/api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.LoginReq
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, err)
            return
        }
        
        l := user.NewLoginLogic(r.Context(), svcCtx)
        resp, err := l.Login(&req)
        if err != nil {
            httpx.Error(w, err)
        } else {
            httpx.OkJson(w, resp)
        }
    }
}
```

#### Logic 层（RPC）

```go
// apps/user/rpc/internal/logic/loginlogic.go
package logic

import (
    "context"
    "errors"
    "imooc.com/easy-chat/apps/user/models"
    "imooc.com/easy-chat/pkg/ctxdata"
    "imooc.com/easy-chat/pkg/encrypt"
    "time"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/logx"
)

var (
    // ErrPhoneNotRegister 手机号未注册
    ErrPhoneNotRegister = errors.New("手机号未注册")
    // ErrPasswordError 密码错误
    ErrPasswordError = errors.New("密码错误")
)

type LoginLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
    return &LoginLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
    // 1. 根据手机号查询用户
    userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
    if err != nil {
        if err == models.ErrNotFound {
            return nil, ErrPhoneNotRegister
        }
        return nil, err
    }
    
    // 2. 验证密码
    if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
        return nil, ErrPasswordError
    }
    
    // 3. 生成 Token
    now := time.Now().Unix()
    token, err := ctxdata.GetJwtToken(
        l.svcCtx.Config.Jwt.AccessSecret,
        now,
        l.svcCtx.Config.Jwt.AccessExpire,
        userEntity.Id,
    )
    if err != nil {
        return nil, err
    }
    
    return &user.LoginResp{
        Token:  token,
        Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
    }, nil
}
```

### 3.3.3 获取用户信息

#### API 定义

```api
type (
    // 获取用户信息请求
    DetailReq {
        UserId string `path:"userId"`  // 用户 ID（路径参数）
    }
    
    // 用户信息响应
    DetailResp {
        Id       string `json:"id"`       // 用户 ID
        Phone    string `json:"phone"`    // 手机号
        Nickname string `json:"nickname"` // 昵称
        Avatar   string `json:"avatar"`   // 头像
        Sex      int    `json:"sex"`      // 性别
    }
)

@server(
    prefix: /api/user
    group: user
    jwt: Auth  // 需要 JWT 认证
)
service user-api {
    @handler detail
    get /detail/:userId (DetailReq) returns (DetailResp)
}
```

#### Logic 层

```go
// apps/user/rpc/internal/logic/getuserinfologic.go
package logic

import (
    "context"
    "imooc.com/easy-chat/apps/user/models"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
    return &GetUserInfoLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
    // 1. 根据用户 ID 查询用户信息
    userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
    if err != nil {
        if err == models.ErrNotFound {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    
    // 2. 构造响应
    return &user.GetUserInfoResp{
        Id:       userEntity.Id,
        Phone:    userEntity.Phone,
        Nickname: userEntity.Nickname,
        Avatar:   userEntity.Avatar,
        Sex:      int32(userEntity.Sex.Int64),
    }, nil
}
```

### 3.3.4 ServiceContext 详解

ServiceContext 是服务上下文，用于存储配置、数据库连接、RPC 客户端等依赖。

```go
// apps/user/api/internal/svc/servicecontext.go
package svc

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/zrpc"
    "google.golang.org/grpc"
    "imooc.com/easy-chat/apps/user/api/internal/config"
    "imooc.com/easy-chat/apps/user/rpc/userclient"
)

// retryPolicy gRPC 重试策略配置
var retryPolicy = `{
    "methodConfig" : [{
        "name": [{
            "service": "user.User"
        }],
        "waitForReady": true,
        "retryPolicy": {
            "maxAttempts": 5,              // 最大重试次数
            "initialBackoff": "0.001s",    // 初始退避时间
            "maxBackoff": "0.002s",        // 最大退避时间
            "backoffMultiplier": 1.0,      // 退避倍数
            "retryableStatusCodes": ["UNKNOWN"]  // 可重试的状态码
        }
    }]
}`

// ServiceContext 服务上下文
type ServiceContext struct {
    Config config.Config      // 配置
    
    *redis.Redis              // Redis 客户端
    userclient.User           // User RPC 客户端
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        
        // 初始化 Redis 客户端
        Redis: redis.MustNewRedis(c.Redisx),
        
        // 初始化 User RPC 客户端（带重试策略）
        User: userclient.NewUser(
            zrpc.MustNewClient(
                c.UserRpc,
                zrpc.WithDialOption(
                    grpc.WithDefaultServiceConfig(retryPolicy),
                ),
            ),
        ),
    }
}
```

### 3.3.5 错误处理

#### 自定义错误

```go
// pkg/xerr/errors.go
package xerr

import "fmt"

// 错误码定义
const (
    OK                  = 0
    ServerError         = 1001  // 服务器错误
    ParamError          = 1002  // 参数错误
    DBError             = 1003  // 数据库错误
    
    UserNotFound        = 2001  // 用户不存在
    UserAlreadyExists   = 2002  // 用户已存在
    PasswordError       = 2003  // 密码错误
    TokenExpired        = 2004  // Token 过期
    TokenInvalid        = 2005  // Token 无效
)

// 错误消息映射
var message = map[int]string{
    OK:                  "success",
    ServerError:         "服务器错误",
    ParamError:          "参数错误",
    DBError:             "数据库错误",
    
    UserNotFound:        "用户不存在",
    UserAlreadyExists:   "用户已存在",
    PasswordError:       "密码错误",
    TokenExpired:        "Token 已过期",
    TokenInvalid:        "Token 无效",
}

// CodeError 自定义错误
type CodeError struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

// NewCodeError 创建自定义错误
func NewCodeError(code int, msg string) *CodeError {
    if msg == "" {
        msg = message[code]
    }
    return &CodeError{
        Code: code,
        Msg:  msg,
    }
}

// Error 实现 error 接口
func (e *CodeError) Error() string {
    return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// GetCode 获取错误码
func (e *CodeError) GetCode() int {
    return e.Code
}

// GetMsg 获取错误消息
func (e *CodeError) GetMsg() string {
    return e.Msg
}
```

#### 统一错误处理

```go
// pkg/xerr/errhandler.go
package xerr

import (
    "github.com/zeromicro/go-zero/rest/httpx"
    "net/http"
)

// ErrorHandler 统一错误处理
func ErrorHandler(err error) (int, interface{}) {
    // 判断是否为自定义错误
    if e, ok := err.(*CodeError); ok {
        return http.StatusOK, &ErrorResponse{
            Code: e.GetCode(),
            Msg:  e.GetMsg(),
        }
    }
    
    // 默认错误
    return http.StatusInternalServerError, &ErrorResponse{
        Code: ServerError,
        Msg:  err.Error(),
    }
}

// ErrorResponse 错误响应
type ErrorResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}
```

#### 在 main.go 中注册错误处理器

```go
// apps/user/api/user.go
package main

import (
    "flag"
    "fmt"
    
    "imooc.com/easy-chat/apps/user/api/internal/config"
    "imooc.com/easy-chat/apps/user/api/internal/handler"
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    "imooc.com/easy-chat/pkg/xerr"
    
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
    flag.Parse()
    
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()
    
    ctx := svc.NewServiceContext(c)
    handler.RegisterHandlers(server, ctx)
    
    // 注册统一错误处理器
    httpx.SetErrorHandler(xerr.ErrorHandler)
    
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
```

### 3.3.6 参数验证

#### 使用 validator 进行参数验证

```go
// apps/user/api/internal/types/types.go
package types

// RegisterReq 注册请求
type RegisterReq struct {
    Phone    string `json:"phone" validate:"required,len=11"`           // 手机号（必填，长度 11）
    Nickname string `json:"nickname" validate:"required,min=2,max=20"`  // 昵称（必填，2-20 字符）
    Password string `json:"password" validate:"required,min=6,max=20"`  // 密码（必填，6-20 字符）
    Avatar   string `json:"avatar"`                                     // 头像（可选）
    Sex      int    `json:"sex" validate:"oneof=0 1 2"`                 // 性别（0-未知，1-男，2-女）
}
```

#### 自定义参数验证中间件

```go
// pkg/middleware/validator.go
package middleware

import (
    "github.com/go-playground/validator/v10"
    "github.com/zeromicro/go-zero/rest/httpx"
    "imooc.com/easy-chat/pkg/xerr"
    "net/http"
)

var validate = validator.New()

// ValidatorMiddleware 参数验证中间件
func ValidatorMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 这里可以添加自定义验证逻辑
        next(w, r)
    }
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
    err := validate.Struct(s)
    if err != nil {
        // 返回第一个验证错误
        if errs, ok := err.(validator.ValidationErrors); ok {
            for _, e := range errs {
                return xerr.NewCodeError(xerr.ParamError, e.Error())
            }
        }
        return xerr.NewCodeError(xerr.ParamError, err.Error())
    }
    return nil
}
```

#### 在 Logic 中使用验证

```go
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 1. 参数验证
    if err := middleware.ValidateStruct(req); err != nil {
        return nil, err
    }
    
    // 2. 调用 RPC 服务
    registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
        Phone:    req.Phone,
        Nickname: req.Nickname,
        Password: req.Password,
        Avatar:   req.Avatar,
        Sex:      int32(req.Sex),
    })
    if err != nil {
        return nil, err
    }
    
    // 3. 返回响应
    var res types.RegisterResp
    copier.Copy(&res, registerResp)
    
    return &res, nil
}
```

---


## 4. RPC 服务开发

### 4.1 Proto 文件语法详解

Protocol Buffers (protobuf) 是 Google 开发的一种数据序列化协议，用于 RPC 通信。

#### 4.1.1 基本语法

```protobuf
// apps/user/rpc/user.proto
syntax = "proto3";  // 使用 proto3 语法

package user;       // 包名

// Go 包路径
option go_package = "./user";

// 用户服务定义
service User {
    // 注册方法
    rpc Register(RegisterReq) returns (RegisterResp);
    // 登录方法
    rpc Login(LoginReq) returns (LoginResp);
    // 获取用户信息
    rpc GetUserInfo(GetUserInfoReq) returns (GetUserInfoResp);
    // 查找用户
    rpc FindUser(FindUserReq) returns (FindUserResp);
}

// 注册请求
message RegisterReq {
    string phone = 1;      // 手机号（字段编号 1）
    string nickname = 2;   // 昵称（字段编号 2）
    string password = 3;   // 密码（字段编号 3）
    string avatar = 4;     // 头像（字段编号 4）
    int32 sex = 5;         // 性别（字段编号 5）
}

// 注册响应
message RegisterResp {
    string token = 1;   // JWT Token
    int64 expire = 2;   // 过期时间戳
}

// 登录请求
message LoginReq {
    string phone = 1;      // 手机号
    string password = 2;   // 密码
}

// 登录响应
message LoginResp {
    string token = 1;   // JWT Token
    int64 expire = 2;   // 过期时间戳
}

// 获取用户信息请求
message GetUserInfoReq {
    string id = 1;  // 用户 ID
}

// 获取用户信息响应
message GetUserInfoResp {
    string id = 1;         // 用户 ID
    string phone = 2;      // 手机号
    string nickname = 3;   // 昵称
    string avatar = 4;     // 头像
    int32 sex = 5;         // 性别
}

// 查找用户请求
message FindUserReq {
    string phone = 1;      // 手机号
    string nickname = 2;   // 昵称
}

// 查找用户响应
message FindUserResp {
    repeated UserInfo users = 1;  // 用户列表（repeated 表示数组）
}

// 用户信息
message UserInfo {
    string id = 1;         // 用户 ID
    string phone = 2;      // 手机号
    string nickname = 3;   // 昵称
    string avatar = 4;     // 头像
    int32 sex = 5;         // 性别
}
```

#### 4.1.2 数据类型

| Proto 类型 | Go 类型 | 说明 |
|-----------|---------|------|
| double | float64 | 双精度浮点数 |
| float | float32 | 单精度浮点数 |
| int32 | int32 | 32 位整数 |
| int64 | int64 | 64 位整数 |
| uint32 | uint32 | 无符号 32 位整数 |
| uint64 | uint64 | 无符号 64 位整数 |
| bool | bool | 布尔值 |
| string | string | 字符串 |
| bytes | []byte | 字节数组 |

#### 4.1.3 复杂类型

```protobuf
// 枚举类型
enum Sex {
    UNKNOWN = 0;  // 未知
    MALE = 1;     // 男
    FEMALE = 2;   // 女
}

// 嵌套消息
message User {
    string id = 1;
    string name = 2;
    Address address = 3;  // 嵌套消息
}

message Address {
    string province = 1;  // 省
    string city = 2;      // 市
    string detail = 3;    // 详细地址
}

// Map 类型
message UserMap {
    map<string, string> metadata = 1;  // 元数据（key-value）
}

// 数组类型
message UserList {
    repeated User users = 1;  // 用户列表
}

// Oneof（只能有一个字段有值）
message SearchRequest {
    oneof query {
        string phone = 1;
        string email = 2;
        string nickname = 3;
    }
}
```

### 4.2 使用 goctl 生成 RPC 代码

#### 4.2.1 生成 RPC 服务端代码

```bash
# 进入 RPC 目录
cd apps/user/rpc

# 生成代码
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 参数说明：
# --go_out=.           生成 protobuf 消息代码
# --go-grpc_out=.      生成 gRPC 服务代码
# --zrpc_out=.         生成 go-zero RPC 服务代码
```

#### 4.2.2 生成的目录结构

```
apps/user/rpc/
├── etc/
│   └── user.yaml           # 配置文件
├── internal/
│   ├── config/
│   │   └── config.go       # 配置结构体
│   ├── logic/
│   │   ├── registerlogic.go      # 注册逻辑
│   │   ├── loginlogic.go         # 登录逻辑
│   │   ├── getuserinfologic.go   # 获取用户信息逻辑
│   │   └── finduserlogic.go      # 查找用户逻辑
│   ├── server/
│   │   └── userserver.go   # RPC 服务器
│   └── svc/
│       └── servicecontext.go  # 服务上下文
├── user/
│   ├── user.pb.go          # protobuf 生成的消息代码
│   └── user_grpc.pb.go     # gRPC 生成的服务代码
├── userclient/
│   └── user.go             # RPC 客户端
├── user.go                 # 主程序入口
└── user.proto              # proto 文件
```

### 4.3 Server 端实现

#### 4.3.1 配置文件

```yaml
# apps/user/rpc/etc/user.yaml
Name: user.rpc              # 服务名称
ListenOn: 0.0.0.0:10001     # 监听地址

# Etcd 配置（服务注册与发现）
Etcd:
  Hosts:
    - 127.0.0.1:2379        # Etcd 地址
  Key: user.rpc             # 服务注册的 Key

# MySQL 配置
Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/easy_chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# Redis 配置
Redisx:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""

# JWT 配置
Jwt:
  AccessSecret: your-access-secret-key
  AccessExpire: 86400       # 24 小时

# 日志配置
Log:
  ServiceName: user.rpc
  Mode: console             # console 或 file
  Level: info               # debug, info, warn, error
  Encoding: plain           # plain 或 json
```

#### 4.3.2 配置结构体

```go
// apps/user/rpc/internal/config/config.go
package config

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/zrpc"
)

// Config RPC 配置
type Config struct {
    zrpc.RpcServerConf  // 继承 RPC 服务器配置
    
    // MySQL 配置
    Mysql struct {
        DataSource string  // 数据源
    }
    
    // Redis 配置
    Redisx redis.RedisConf
    
    // JWT 配置
    Jwt struct {
        AccessSecret string  // 密钥
        AccessExpire int64   // 过期时间（秒）
    }
}
```

#### 4.3.3 ServiceContext

```go
// apps/user/rpc/internal/svc/servicecontext.go
package svc

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "imooc.com/easy-chat/apps/user/models"
    "imooc.com/easy-chat/apps/user/rpc/internal/config"
)

// ServiceContext 服务上下文
type ServiceContext struct {
    Config config.Config  // 配置
    
    *redis.Redis          // Redis 客户端
    
    UsersModel models.UsersModel  // 用户 Model
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
    // 创建 MySQL 连接
    sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
    
    return &ServiceContext{
        Config: c,
        
        // 初始化 Redis
        Redis: redis.MustNewRedis(c.Redisx),
        
        // 初始化用户 Model
        UsersModel: models.NewUsersModel(sqlConn),
    }
}
```

#### 4.3.4 Server 实现

```go
// apps/user/rpc/internal/server/userserver.go
package server

import (
    "context"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/logic"
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
)

// UserServer RPC 服务器
type UserServer struct {
    svcCtx *svc.ServiceContext
    user.UnimplementedUserServer  // 嵌入未实现的服务器（向前兼容）
}

// NewUserServer 创建 RPC 服务器
func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
    return &UserServer{
        svcCtx: svcCtx,
    }
}

// Register 注册方法
func (s *UserServer) Register(ctx context.Context, in *user.RegisterReq) (*user.RegisterResp, error) {
    l := logic.NewRegisterLogic(ctx, s.svcCtx)
    return l.Register(in)
}

// Login 登录方法
func (s *UserServer) Login(ctx context.Context, in *user.LoginReq) (*user.LoginResp, error) {
    l := logic.NewLoginLogic(ctx, s.svcCtx)
    return l.Login(in)
}

// GetUserInfo 获取用户信息
func (s *UserServer) GetUserInfo(ctx context.Context, in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
    l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
    return l.GetUserInfo(in)
}

// FindUser 查找用户
func (s *UserServer) FindUser(ctx context.Context, in *user.FindUserReq) (*user.FindUserResp, error) {
    l := logic.NewFindUserLogic(ctx, s.svcCtx)
    return l.FindUser(in)
}
```

#### 4.3.5 主程序入口

```go
// apps/user/rpc/user.go
package main

import (
    "flag"
    "fmt"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/config"
    "imooc.com/easy-chat/apps/user/rpc/internal/server"
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/service"
    "github.com/zeromicro/go-zero/zrpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
    flag.Parse()
    
    // 1. 加载配置
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    // 2. 创建服务上下文
    ctx := svc.NewServiceContext(c)
    
    // 3. 创建 RPC 服务器
    s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
        // 注册服务
        user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
        
        // 开启反射（用于 grpcurl 等工具调试）
        if c.Mode == service.DevMode || c.Mode == service.TestMode {
            reflection.Register(grpcServer)
        }
    })
    defer s.Stop()
    
    fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
    s.Start()
}
```

### 4.4 Client 端调用

#### 4.4.1 生成客户端代码

客户端代码已经在生成 RPC 代码时自动生成：

```go
// apps/user/rpc/userclient/user.go
package userclient

import (
    "context"
    
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/zrpc"
)

type (
    // User RPC 客户端接口
    User interface {
        Register(ctx context.Context, in *user.RegisterReq, opts ...grpc.CallOption) (*user.RegisterResp, error)
        Login(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error)
        GetUserInfo(ctx context.Context, in *user.GetUserInfoReq, opts ...grpc.CallOption) (*user.GetUserInfoResp, error)
        FindUser(ctx context.Context, in *user.FindUserReq, opts ...grpc.CallOption) (*user.FindUserResp, error)
    }
    
    defaultUser struct {
        cli zrpc.Client
    }
)

// NewUser 创建 RPC 客户端
func NewUser(cli zrpc.Client) User {
    return &defaultUser{
        cli: cli,
    }
}

// Register 注册
func (m *defaultUser) Register(ctx context.Context, in *user.RegisterReq, opts ...grpc.CallOption) (*user.RegisterResp, error) {
    client := user.NewUserClient(m.cli.Conn())
    return client.Register(ctx, in, opts...)
}

// Login 登录
func (m *defaultUser) Login(ctx context.Context, in *user.LoginReq, opts ...grpc.CallOption) (*user.LoginResp, error) {
    client := user.NewUserClient(m.cli.Conn())
    return client.Login(ctx, in, opts...)
}

// GetUserInfo 获取用户信息
func (m *defaultUser) GetUserInfo(ctx context.Context, in *user.GetUserInfoReq, opts ...grpc.CallOption) (*user.GetUserInfoResp, error) {
    client := user.NewUserClient(m.cli.Conn())
    return client.GetUserInfo(ctx, in, opts...)
}

// FindUser 查找用户
func (m *defaultUser) FindUser(ctx context.Context, in *user.FindUserReq, opts ...grpc.CallOption) (*user.FindUserResp, error) {
    client := user.NewUserClient(m.cli.Conn())
    return client.FindUser(ctx, in, opts...)
}
```

#### 4.4.2 在 API 服务中调用 RPC

```go
// apps/user/api/internal/svc/servicecontext.go
package svc

import (
    "github.com/zeromicro/go-zero/zrpc"
    "imooc.com/easy-chat/apps/user/api/internal/config"
    "imooc.com/easy-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
    Config config.Config
    
    userclient.User  // User RPC 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        
        // 创建 RPC 客户端
        User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
    }
}
```

```go
// apps/user/api/internal/logic/user/registerlogic.go
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 调用 RPC 服务
    registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
        Phone:    req.Phone,
        Nickname: req.Nickname,
        Password: req.Password,
        Avatar:   req.Avatar,
        Sex:      int32(req.Sex),
    })
    if err != nil {
        return nil, err
    }
    
    var res types.RegisterResp
    copier.Copy(&res, registerResp)
    
    return &res, nil
}
```

### 4.5 服务注册与发现（Etcd）

#### 4.5.1 安装 Etcd

```bash
# 使用 Docker 安装 Etcd
docker run -d \
  --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  bitnami/etcd:latest
```

#### 4.5.2 RPC 服务端配置

```yaml
# apps/user/rpc/etc/user.yaml
Name: user.rpc
ListenOn: 0.0.0.0:10001

# Etcd 配置
Etcd:
  Hosts:
    - 127.0.0.1:2379  # Etcd 地址
  Key: user.rpc       # 服务注册的 Key
```

#### 4.5.3 RPC 客户端配置

```yaml
# apps/user/api/etc/user.yaml
Name: user-api
Host: 0.0.0.0
Port: 8888

# User RPC 配置
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379  # Etcd 地址
    Key: user.rpc       # 服务发现的 Key
```

#### 4.5.4 服务注册流程

1. RPC 服务启动时，自动将服务地址注册到 Etcd
2. Etcd 中存储的 Key 格式：`user.rpc/127.0.0.1:10001`
3. 服务会定期发送心跳，保持注册信息

#### 4.5.5 服务发现流程

1. API 服务启动时，从 Etcd 获取 RPC 服务地址列表
2. 使用负载均衡算法选择一个 RPC 服务地址
3. 监听 Etcd 变化，动态更新服务地址列表

#### 4.5.6 查看 Etcd 中的服务

```bash
# 安装 etcdctl
brew install etcd  # macOS
apt-get install etcd-client  # Ubuntu

# 查看所有 Key
etcdctl get --prefix ""

# 查看指定服务
etcdctl get --prefix "user.rpc"

# 输出示例：
# user.rpc/127.0.0.1:10001
# {"name":"user.rpc","addr":"127.0.0.1:10001"}
```

### 4.6 完整的 RPC 服务示例

#### 4.6.1 查找用户 Logic

```go
// apps/user/rpc/internal/logic/finduserlogic.go
package logic

import (
    "context"
    "imooc.com/easy-chat/apps/user/models"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
    return &FindUserLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// FindUser 查找用户
func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
    var (
        userEntities []*models.Users
        err          error
    )
    
    // 1. 根据条件查询用户
    if in.Phone != "" {
        // 根据手机号查询
        userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
        if err != nil && err != models.ErrNotFound {
            return nil, err
        }
        if userEntity != nil {
            userEntities = append(userEntities, userEntity)
        }
    } else if in.Nickname != "" {
        // 根据昵称模糊查询
        userEntities, err = l.svcCtx.UsersModel.FindByNickname(l.ctx, in.Nickname)
        if err != nil {
            return nil, err
        }
    }
    
    // 2. 构造响应
    var users []*user.UserInfo
    for _, userEntity := range userEntities {
        users = append(users, &user.UserInfo{
            Id:       userEntity.Id,
            Phone:    userEntity.Phone,
            Nickname: userEntity.Nickname,
            Avatar:   userEntity.Avatar,
            Sex:      int32(userEntity.Sex.Int64),
        })
    }
    
    return &user.FindUserResp{
        Users: users,
    }, nil
}
```

#### 4.6.2 测试 RPC 服务

```go
// apps/user/rpc/internal/logic/logic_test.go
package logic

import (
    "context"
    "testing"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/config"
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/zeromicro/go-zero/core/conf"
)

func TestRegisterLogic_Register(t *testing.T) {
    // 1. 加载配置
    var c config.Config
    conf.MustLoad("../../etc/user.yaml", &c)
    
    // 2. 创建服务上下文
    ctx := svc.NewServiceContext(c)
    
    // 3. 创建 Logic
    l := NewRegisterLogic(context.Background(), ctx)
    
    // 4. 调用方法
    resp, err := l.Register(&user.RegisterReq{
        Phone:    "13800138000",
        Nickname: "测试用户",
        Password: "123456",
        Avatar:   "https://example.com/avatar.jpg",
        Sex:      1,
    })
    
    // 5. 断言
    if err != nil {
        t.Errorf("Register failed: %v", err)
        return
    }
    
    if resp.Token == "" {
        t.Error("Token is empty")
    }
    
    t.Logf("Register success, token: %s, expire: %d", resp.Token, resp.Expire)
}
```

---


## 5. 配置管理

### 5.1 配置文件结构

Go-Zero 使用 YAML 格式的配置文件，支持环境变量替换和配置热更新。

#### 5.1.1 API 服务配置示例

```yaml
# apps/user/api/etc/user.yaml
Name: user-api              # 服务名称
Host: 0.0.0.0               # 监听地址
Port: 8888                  # 监听端口

# 日志配置
Log:
  ServiceName: user-api     # 服务名称
  Mode: console             # 日志模式：console（控制台）、file（文件）
  Level: info               # 日志级别：debug、info、warn、error
  Encoding: plain           # 日志编码：plain（纯文本）、json（JSON）
  Path: logs                # 日志文件路径（Mode 为 file 时有效）
  KeepDays: 7               # 日志保留天数

# JWT 配置
Auth:
  AccessSecret: your-access-secret-key-change-me  # JWT 密钥（必须修改）
  AccessExpire: 86400                             # Token 过期时间（秒），86400 = 24 小时

# Redis 配置
Redisx:
  Host: 127.0.0.1:6379      # Redis 地址
  Type: node                # Redis 类型：node（单节点）、cluster（集群）
  Pass: ""                  # Redis 密码
  Tls: false                # 是否启用 TLS

# User RPC 配置
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379      # Etcd 地址
    Key: user.rpc           # 服务发现的 Key
  Timeout: 5000             # 超时时间（毫秒）
  NonBlock: true            # 非阻塞模式

# 超时配置
Timeout: 30000              # 请求超时时间（毫秒）

# 限流配置
MaxConns: 10000             # 最大并发连接数
MaxBytes: 1048576           # 最大请求体大小（字节），1048576 = 1MB

# 跨域配置
CorsOrigins:
  - "*"                     # 允许的源（* 表示所有）
```

#### 5.1.2 RPC 服务配置示例

```yaml
# apps/user/rpc/etc/user.yaml
Name: user.rpc              # 服务名称
ListenOn: 0.0.0.0:10001     # 监听地址

# Etcd 配置（服务注册与发现）
Etcd:
  Hosts:
    - 127.0.0.1:2379        # Etcd 地址
  Key: user.rpc             # 服务注册的 Key

# 日志配置
Log:
  ServiceName: user.rpc
  Mode: console
  Level: info
  Encoding: plain

# MySQL 配置
Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/easy_chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# Redis 配置
Redisx:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""

# JWT 配置
Jwt:
  AccessSecret: your-access-secret-key-change-me
  AccessExpire: 86400

# 超时配置
Timeout: 5000               # RPC 调用超时时间（毫秒）

# 限流配置
CpuThreshold: 900           # CPU 使用率阈值（千分比），900 = 90%
```

### 5.2 服务配置详解

#### 5.2.1 基础配置

```go
// apps/user/api/internal/config/config.go
package config

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

// Config API 服务配置
type Config struct {
    rest.RestConf  // 继承 REST 配置
    
    // JWT 配置
    Auth struct {
        AccessSecret string  // JWT 密钥
        AccessExpire int64   // 过期时间（秒）
    }
    
    // Redis 配置
    Redisx redis.RedisConf
    
    // User RPC 配置
    UserRpc zrpc.RpcClientConf
}
```

#### 5.2.2 RestConf 配置项

```go
// RestConf REST 服务配置
type RestConf struct {
    service.ServiceConf           // 服务配置
    Host                string    // 监听地址
    Port                int       // 监听端口
    CertFile            string    // TLS 证书文件
    KeyFile             string    // TLS 密钥文件
    Verbose             bool      // 是否打印详细日志
    MaxConns            int       // 最大并发连接数
    MaxBytes            int64     // 最大请求体大小
    Timeout             int64     // 请求超时时间（毫秒）
    CpuThreshold        int64     // CPU 使用率阈值
    Signature           SignatureConf  // 签名配置
    Middlewares         []string       // 中间件列表
    CorsOrigins         []string       // 跨域配置
}
```

#### 5.2.3 ServiceConf 配置项

```go
// ServiceConf 服务配置
type ServiceConf struct {
    Name       string          // 服务名称
    Log        logx.LogConf    // 日志配置
    Mode       string          // 运行模式：dev、test、pre、pro
    MetricsUrl string          // Prometheus 指标地址
    Prometheus prometheus.Config  // Prometheus 配置
}
```

### 5.3 数据库配置

#### 5.3.1 MySQL 配置

```yaml
# MySQL 配置
Mysql:
  # 数据源格式：用户名:密码@tcp(地址:端口)/数据库名?参数
  DataSource: root:123456@tcp(127.0.0.1:3306)/easy_chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
  
  # 连接池配置（可选）
  MaxOpenConns: 100      # 最大打开连接数
  MaxIdleConns: 10       # 最大空闲连接数
  ConnMaxLifetime: 3600  # 连接最大生命周期（秒）
```

```go
// 配置结构体
type Config struct {
    zrpc.RpcServerConf
    
    Mysql struct {
        DataSource      string  // 数据源
        MaxOpenConns    int     // 最大打开连接数
        MaxIdleConns    int     // 最大空闲连接数
        ConnMaxLifetime int     // 连接最大生命周期
    }
}

// 创建 MySQL 连接
func NewServiceContext(c config.Config) *ServiceContext {
    sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
    
    // 设置连接池参数
    if c.Mysql.MaxOpenConns > 0 {
        sqlConn.SetMaxOpenConns(c.Mysql.MaxOpenConns)
    }
    if c.Mysql.MaxIdleConns > 0 {
        sqlConn.SetMaxIdleConns(c.Mysql.MaxIdleConns)
    }
    if c.Mysql.ConnMaxLifetime > 0 {
        sqlConn.SetConnMaxLifetime(time.Duration(c.Mysql.ConnMaxLifetime) * time.Second)
    }
    
    return &ServiceContext{
        Config:     c,
        UsersModel: models.NewUsersModel(sqlConn),
    }
}
```

#### 5.3.2 MongoDB 配置

```yaml
# MongoDB 配置
Mongo:
  Url: mongodb://127.0.0.1:27017  # MongoDB 地址
  Database: easy_chat              # 数据库名
  Username: ""                     # 用户名（可选）
  Password: ""                     # 密码（可选）
  MaxPoolSize: 100                 # 最大连接池大小
  MinPoolSize: 10                  # 最小连接池大小
```

```go
// 配置结构体
type Config struct {
    zrpc.RpcServerConf
    
    Mongo struct {
        Url         string  // MongoDB 地址
        Database    string  // 数据库名
        Username    string  // 用户名
        Password    string  // 密码
        MaxPoolSize int     // 最大连接池大小
        MinPoolSize int     // 最小连接池大小
    }
}

// 创建 MongoDB 连接
import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func NewServiceContext(c config.Config) *ServiceContext {
    // 创建客户端选项
    clientOptions := options.Client().ApplyURI(c.Mongo.Url)
    
    // 设置认证
    if c.Mongo.Username != "" && c.Mongo.Password != "" {
        clientOptions.SetAuth(options.Credential{
            Username: c.Mongo.Username,
            Password: c.Mongo.Password,
        })
    }
    
    // 设置连接池
    if c.Mongo.MaxPoolSize > 0 {
        clientOptions.SetMaxPoolSize(uint64(c.Mongo.MaxPoolSize))
    }
    if c.Mongo.MinPoolSize > 0 {
        clientOptions.SetMinPoolSize(uint64(c.Mongo.MinPoolSize))
    }
    
    // 连接 MongoDB
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        panic(err)
    }
    
    // 获取数据库
    db := client.Database(c.Mongo.Database)
    
    return &ServiceContext{
        Config:  c,
        MongoDB: db,
    }
}
```

#### 5.3.3 Redis 配置

```yaml
# Redis 单节点配置
Redisx:
  Host: 127.0.0.1:6379    # Redis 地址
  Type: node              # 类型：node（单节点）
  Pass: ""                # 密码
  Tls: false              # 是否启用 TLS
  DB: 0                   # 数据库编号

# Redis 集群配置
Redisx:
  Type: cluster           # 类型：cluster（集群）
  Hosts:
    - 127.0.0.1:7000
    - 127.0.0.1:7001
    - 127.0.0.1:7002
  Pass: ""
```

```go
// 配置结构体
import "github.com/zeromicro/go-zero/core/stores/redis"

type Config struct {
    rest.RestConf
    
    Redisx redis.RedisConf  // Redis 配置
}

// 创建 Redis 客户端
func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        Redis:  redis.MustNewRedis(c.Redisx),
    }
}

// 使用 Redis
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
    // 设置缓存
    err = l.svcCtx.Redis.Set("user:token:"+userId, token)
    if err != nil {
        return nil, err
    }
    
    // 设置过期时间
    err = l.svcCtx.Redis.Expire("user:token:"+userId, 86400)
    if err != nil {
        return nil, err
    }
    
    // 获取缓存
    token, err := l.svcCtx.Redis.Get("user:token:" + userId)
    if err != nil {
        return nil, err
    }
    
    return &types.LoginResp{Token: token}, nil
}
```

### 5.4 RPC 配置

#### 5.4.1 RPC 客户端配置

```yaml
# RPC 客户端配置
UserRpc:
  # Etcd 服务发现
  Etcd:
    Hosts:
      - 127.0.0.1:2379    # Etcd 地址
    Key: user.rpc         # 服务发现的 Key
  
  # 直连模式（不使用服务发现）
  # Endpoints:
  #   - 127.0.0.1:10001
  
  Timeout: 5000           # 超时时间（毫秒）
  NonBlock: true          # 非阻塞模式
```

```go
// 配置结构体
import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
    rest.RestConf
    
    UserRpc zrpc.RpcClientConf  // User RPC 配置
}

// 创建 RPC 客户端
import (
    "github.com/zeromicro/go-zero/zrpc"
    "imooc.com/easy-chat/apps/user/rpc/userclient"
)

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
    }
}
```

#### 5.4.2 RPC 服务端配置

```yaml
# RPC 服务端配置
Name: user.rpc
ListenOn: 0.0.0.0:10001   # 监听地址

# Etcd 服务注册
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: user.rpc

# 超时配置
Timeout: 5000             # 超时时间（毫秒）

# 限流配置
CpuThreshold: 900         # CPU 使用率阈值（千分比）
```

```go
// 配置结构体
import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
    zrpc.RpcServerConf  // 继承 RPC 服务器配置
    
    // 其他配置...
}
```

### 5.5 JWT 配置

```yaml
# JWT 配置
Auth:
  AccessSecret: your-access-secret-key-change-me  # JWT 密钥（必须修改）
  AccessExpire: 86400                             # Token 过期时间（秒）
```

```go
// 配置结构体
type Config struct {
    rest.RestConf
    
    Auth struct {
        AccessSecret string  // JWT 密钥
        AccessExpire int64   // 过期时间（秒）
    }
}

// 生成 Token
import (
    "github.com/golang-jwt/jwt/v4"
    "time"
)

func GenerateToken(secret string, expire int64, userId string) (string, error) {
    now := time.Now().Unix()
    
    // 创建 Claims
    claims := make(jwt.MapClaims)
    claims["exp"] = now + expire  // 过期时间
    claims["iat"] = now           // 签发时间
    claims["userId"] = userId     // 用户 ID
    
    // 创建 Token
    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims = claims
    
    // 签名
    return token.SignedString([]byte(secret))
}

// 解析 Token
func ParseToken(secret string, tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return "", err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userId := claims["userId"].(string)
        return userId, nil
    }
    
    return "", errors.New("invalid token")
}
```

### 5.6 日志配置

```yaml
# 日志配置
Log:
  ServiceName: user-api     # 服务名称
  Mode: console             # 日志模式：console（控制台）、file（文件）
  Level: info               # 日志级别：debug、info、warn、error
  Encoding: plain           # 日志编码：plain（纯文本）、json（JSON）
  Path: logs                # 日志文件路径（Mode 为 file 时有效）
  KeepDays: 7               # 日志保留天数
  Compress: true            # 是否压缩日志文件
  MaxBackups: 5             # 最大备份文件数
  MaxSize: 100              # 单个日志文件最大大小（MB）
```

```go
// 使用日志
import "github.com/zeromicro/go-zero/core/logx"

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // Debug 日志
    logx.Debugf("Register request: %+v", req)
    
    // Info 日志
    logx.Infof("User %s registered successfully", req.Phone)
    
    // Warn 日志
    logx.Warnf("User %s already exists", req.Phone)
    
    // Error 日志
    logx.Errorf("Register failed: %v", err)
    
    // 带上下文的日志
    logx.WithContext(l.ctx).Infof("User %s registered", req.Phone)
    
    return &types.RegisterResp{}, nil
}
```

### 5.7 配置热更新

Go-Zero 支持配置热更新，当配置文件发生变化时，会自动重新加载配置。

```go
// 启用配置热更新
import (
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
)

func main() {
    var c config.Config
    
    // 加载配置并监听变化
    conf.MustLoad(*configFile, &c, conf.UseEnv())
    
    // 监听配置变化
    conf.AddListener(func() {
        logx.Info("Config changed, reloading...")
        
        // 重新加载配置
        var newConfig config.Config
        conf.MustLoad(*configFile, &newConfig)
        
        // 更新配置
        c = newConfig
        
        logx.Info("Config reloaded successfully")
    })
    
    // 启动服务...
}
```

### 5.8 环境变量替换

配置文件支持使用环境变量：

```yaml
# 使用环境变量
Mysql:
  DataSource: ${MYSQL_DSN}  # 从环境变量 MYSQL_DSN 读取

Redisx:
  Host: ${REDIS_HOST:127.0.0.1:6379}  # 从环境变量读取，默认值为 127.0.0.1:6379
  Pass: ${REDIS_PASS:}                # 从环境变量读取，默认值为空

Auth:
  AccessSecret: ${JWT_SECRET}
  AccessExpire: ${JWT_EXPIRE:86400}
```

```bash
# 设置环境变量
export MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/easy_chat?charset=utf8mb4&parseTime=true"
export REDIS_HOST="127.0.0.1:6379"
export REDIS_PASS="your-redis-password"
export JWT_SECRET="your-jwt-secret-key"
export JWT_EXPIRE="86400"

# 启动服务
go run user.go -f etc/user.yaml
```

---


## 6. 数据库操作

### 6.1 MySQL 集成

#### 6.1.1 安装依赖

```bash
go get -u github.com/zeromicro/go-zero/core/stores/sqlx
go get -u github.com/go-sql-driver/mysql
```

#### 6.1.2 创建数据库表

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS easy_chat DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE easy_chat;

-- 用户表
CREATE TABLE `users` (
  `id` varchar(24) NOT NULL COMMENT '用户ID',
  `phone` varchar(11) NOT NULL COMMENT '手机号',
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `password` varchar(255) DEFAULT NULL COMMENT '密码',
  `sex` tinyint(1) DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_nickname` (`nickname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### 6.1.3 使用 goctl 生成 Model 代码

```bash
# 方式1：从数据库生成
goctl model mysql datasource \
  -url="root:123456@tcp(127.0.0.1:3306)/easy_chat" \
  -table="users" \
  -dir="./apps/user/models" \
  -cache=true

# 方式2：从 DDL 文件生成
goctl model mysql ddl \
  -src="./apps/user/models/users.sql" \
  -dir="./apps/user/models" \
  -cache=true

# 参数说明：
# -url          数据库连接地址
# -table        表名（支持通配符，如 "user_*"）
# -dir          输出目录
# -cache        是否生成缓存代码
# -style        命名风格：gozero、go_zero、GoZero
```

#### 6.1.4 生成的 Model 代码

```go
// apps/user/models/usersmodel.go
package models

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    
    "github.com/zeromicro/go-zero/core/stores/builder"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlc"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/core/stringx"
)

var (
    usersFieldNames          = builder.RawFieldNames(&Users{})
    usersRows                = strings.Join(usersFieldNames, ",")
    usersRowsExpectAutoSet   = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
    usersRowsWithPlaceHolder = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
    
    cacheUsersIdPrefix    = "cache:users:id:"
    cacheUsersPhonePrefix = "cache:users:phone:"
)

type (
    usersModel interface {
        Insert(ctx context.Context, data *Users) (sql.Result, error)
        FindOne(ctx context.Context, id string) (*Users, error)
        FindOneByPhone(ctx context.Context, phone string) (*Users, error)
        Update(ctx context.Context, data *Users) error
        Delete(ctx context.Context, id string) error
    }
    
    defaultUsersModel struct {
        sqlc.CachedConn
        table string
    }
    
    Users struct {
        Id         string         `db:"id"`          // 用户ID
        Phone      string         `db:"phone"`       // 手机号
        Nickname   string         `db:"nickname"`    // 昵称
        Avatar     string         `db:"avatar"`      // 头像
        Password   sql.NullString `db:"password"`    // 密码
        Sex        sql.NullInt64  `db:"sex"`         // 性别
        CreateTime sql.NullTime   `db:"create_time"` // 创建时间
        UpdateTime sql.NullTime   `db:"update_time"` // 更新时间
    }
)

// newUsersModel 创建 Model 实例
func newUsersModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUsersModel {
    return &defaultUsersModel{
        CachedConn: sqlc.NewConn(conn, c),
        table:      "`users`",
    }
}

// Insert 插入数据
func (m *defaultUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
    query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        return conn.ExecCtx(ctx, query, data.Id, data.Phone, data.Nickname, data.Avatar, data.Password, data.Sex)
    }, usersIdKey(data.Id), usersPhoneKey(data.Phone))
    return ret, err
}

// FindOne 根据 ID 查询
func (m *defaultUsersModel) FindOne(ctx context.Context, id string) (*Users, error) {
    usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
    var resp Users
    err := m.QueryRowCtx(ctx, &resp, usersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
        query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
        return conn.QueryRowCtx(ctx, v, query, id)
    })
    switch err {
    case nil:
        return &resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

// FindOneByPhone 根据手机号查询
func (m *defaultUsersModel) FindOneByPhone(ctx context.Context, phone string) (*Users, error) {
    usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
    var resp Users
    err := m.QueryRowIndexCtx(ctx, &resp, usersPhoneKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
        query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
        if err := conn.QueryRowCtx(ctx, &resp, query, phone); err != nil {
            return nil, err
        }
        return resp.Id, nil
    }, m.queryPrimary)
    switch err {
    case nil:
        return &resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

// Update 更新数据
func (m *defaultUsersModel) Update(ctx context.Context, newData *Users) error {
    usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, newData.Id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
        return conn.ExecCtx(ctx, query, newData.Phone, newData.Nickname, newData.Avatar, newData.Password, newData.Sex, newData.Id)
    }, usersIdKey)
    return err
}

// Delete 删除数据
func (m *defaultUsersModel) Delete(ctx context.Context, id string) error {
    // 先查询数据（用于删除缓存）
    data, err := m.FindOne(ctx, id)
    if err != nil {
        return err
    }
    
    usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
    usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, data.Phone)
    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, id)
    }, usersIdKey, usersPhoneKey)
    return err
}

func (m *defaultUsersModel) formatPrimary(primary interface{}) string {
    return fmt.Sprintf("%s%v", cacheUsersIdPrefix, primary)
}

func (m *defaultUsersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
    query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
    return conn.QueryRowCtx(ctx, v, query, primary)
}

func usersIdKey(id string) string {
    return fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
}

func usersPhoneKey(phone string) string {
    return fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
}
```

#### 6.1.5 自定义 Model 方法

```go
// apps/user/models/usersmodel.go
package models

import (
    "context"
    "fmt"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
    // UsersModel 用户 Model 接口
    UsersModel interface {
        usersModel  // 继承生成的接口
        
        // 自定义方法
        FindByNickname(ctx context.Context, nickname string) ([]*Users, error)
        FindByIds(ctx context.Context, ids []string) ([]*Users, error)
        UpdateNickname(ctx context.Context, id string, nickname string) error
        Count(ctx context.Context) (int64, error)
    }
    
    customUsersModel struct {
        *defaultUsersModel
    }
)

// NewUsersModel 创建 Model 实例
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf) UsersModel {
    return &customUsersModel{
        defaultUsersModel: newUsersModel(conn, c),
    }
}

// FindByNickname 根据昵称模糊查询
func (m *customUsersModel) FindByNickname(ctx context.Context, nickname string) ([]*Users, error) {
    var resp []*Users
    query := fmt.Sprintf("select %s from %s where `nickname` like ? limit 100", usersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, "%"+nickname+"%")
    switch err {
    case nil:
        return resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

// FindByIds 根据 ID 列表批量查询
func (m *customUsersModel) FindByIds(ctx context.Context, ids []string) ([]*Users, error) {
    if len(ids) == 0 {
        return []*Users{}, nil
    }
    
    var resp []*Users
    query := fmt.Sprintf("select %s from %s where `id` in (?)", usersRows, m.table)
    query, args, err := sqlx.In(query, ids)
    if err != nil {
        return nil, err
    }
    
    err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
    switch err {
    case nil:
        return resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

// UpdateNickname 更新昵称
func (m *customUsersModel) UpdateNickname(ctx context.Context, id string, nickname string) error {
    usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("update %s set `nickname` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, nickname, id)
    }, usersIdKey)
    return err
}

// Count 统计用户数量
func (m *customUsersModel) Count(ctx context.Context) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query)
    return count, err
}
```

#### 6.1.6 事务操作

```go
// 事务示例
func (l *TransferLogic) Transfer(req *types.TransferReq) error {
    // 开启事务
    err := l.svcCtx.UsersModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // 1. 扣减转出用户余额
        _, err := session.ExecCtx(ctx, 
            "update users set balance = balance - ? where id = ? and balance >= ?",
            req.Amount, req.FromUserId, req.Amount)
        if err != nil {
            return err
        }
        
        // 2. 增加转入用户余额
        _, err = session.ExecCtx(ctx,
            "update users set balance = balance + ? where id = ?",
            req.Amount, req.ToUserId)
        if err != nil {
            return err
        }
        
        // 3. 记录转账记录
        _, err = session.ExecCtx(ctx,
            "insert into transfer_records (from_user_id, to_user_id, amount) values (?, ?, ?)",
            req.FromUserId, req.ToUserId, req.Amount)
        if err != nil {
            return err
        }
        
        return nil
    })
    
    return err
}
```

### 6.2 MongoDB 集成

#### 6.2.1 安装依赖

```bash
go get -u go.mongodb.org/mongo-driver/mongo
go get -u go.mongodb.org/mongo-driver/bson
```

#### 6.2.2 配置 MongoDB

```yaml
# apps/im/rpc/etc/im.yaml
Mongo:
  Url: mongodb://127.0.0.1:27017
  Database: easy_chat
```

```go
// apps/im/rpc/internal/config/config.go
type Config struct {
    zrpc.RpcServerConf
    
    Mongo struct {
        Url      string
        Database string
    }
}
```

#### 6.2.3 创建 MongoDB Model

```go
// apps/im/models/chatlogmodel.go
package models

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

const ChatLogCollectionName = "chat_log"

type (
    // ChatLogModel 聊天记录 Model 接口
    ChatLogModel interface {
        Insert(ctx context.Context, data *ChatLog) error
        FindOne(ctx context.Context, id string) (*ChatLog, error)
        FindByConversationId(ctx context.Context, conversationId string, page, pageSize int64) ([]*ChatLog, error)
        Update(ctx context.Context, data *ChatLog) error
        Delete(ctx context.Context, id string) error
    }
    
    defaultChatLogModel struct {
        conn       *mongo.Collection
    }
    
    // ChatLog 聊天记录
    ChatLog struct {
        ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
        ConversationId string             `bson:"conversation_id" json:"conversationId"`  // 会话 ID
        SendId         string             `bson:"send_id" json:"sendId"`                  // 发送者 ID
        RecvId         string             `bson:"recv_id" json:"recvId"`                  // 接收者 ID
        MsgType        int                `bson:"msg_type" json:"msgType"`                // 消息类型
        MsgContent     string             `bson:"msg_content" json:"msgContent"`          // 消息内容
        ChatType       int                `bson:"chat_type" json:"chatType"`              // 聊天类型：1-单聊，2-群聊
        SendTime       time.Time          `bson:"send_time" json:"sendTime"`              // 发送时间
        ReadRecords    []string           `bson:"read_records" json:"readRecords"`        // 已读记录
        CreateTime     time.Time          `bson:"create_time" json:"createTime"`
        UpdateTime     time.Time          `bson:"update_time" json:"updateTime"`
    }
)

// NewChatLogModel 创建 Model 实例
func NewChatLogModel(db *mongo.Database) ChatLogModel {
    return &defaultChatLogModel{
        conn: db.Collection(ChatLogCollectionName),
    }
}

// Insert 插入数据
func (m *defaultChatLogModel) Insert(ctx context.Context, data *ChatLog) error {
    if data.ID.IsZero() {
        data.ID = primitive.NewObjectID()
    }
    data.CreateTime = time.Now()
    data.UpdateTime = time.Now()
    
    _, err := m.conn.InsertOne(ctx, data)
    return err
}

// FindOne 根据 ID 查询
func (m *defaultChatLogModel) FindOne(ctx context.Context, id string) (*ChatLog, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    var data ChatLog
    err = m.conn.FindOne(ctx, bson.M{"_id": oid}).Decode(&data)
    if err != nil {
        return nil, err
    }
    
    return &data, nil
}

// FindByConversationId 根据会话 ID 查询（分页）
func (m *defaultChatLogModel) FindByConversationId(ctx context.Context, conversationId string, page, pageSize int64) ([]*ChatLog, error) {
    // 计算跳过的数量
    skip := (page - 1) * pageSize
    
    // 设置查询选项
    opts := options.Find().
        SetSort(bson.D{{"send_time", -1}}).  // 按发送时间倒序
        SetSkip(skip).
        SetLimit(pageSize)
    
    // 查询
    cursor, err := m.conn.Find(ctx, bson.M{"conversation_id": conversationId}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    // 解析结果
    var list []*ChatLog
    if err = cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    
    return list, nil
}

// Update 更新数据
func (m *defaultChatLogModel) Update(ctx context.Context, data *ChatLog) error {
    data.UpdateTime = time.Now()
    
    _, err := m.conn.UpdateOne(ctx,
        bson.M{"_id": data.ID},
        bson.M{"$set": data},
    )
    return err
}

// Delete 删除数据
func (m *defaultChatLogModel) Delete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    
    _, err = m.conn.DeleteOne(ctx, bson.M{"_id": oid})
    return err
}
```

### 6.3 Redis 集成

#### 6.3.1 基本操作

```go
// String 操作
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
    // 1. Set 设置值
    err = l.svcCtx.Redis.Set("user:token:"+userId, token)
    
    // 2. Setex 设置值并设置过期时间
    err = l.svcCtx.Redis.Setex("user:token:"+userId, token, 86400)
    
    // 3. Get 获取值
    token, err := l.svcCtx.Redis.Get("user:token:" + userId)
    
    // 4. Del 删除
    _, err = l.svcCtx.Redis.Del("user:token:" + userId)
    
    // 5. Exists 判断是否存在
    exists, err := l.svcCtx.Redis.Exists("user:token:" + userId)
    
    // 6. Expire 设置过期时间
    err = l.svcCtx.Redis.Expire("user:token:"+userId, 86400)
    
    // 7. Ttl 获取剩余过期时间
    ttl, err := l.svcCtx.Redis.Ttl("user:token:" + userId)
    
    return &types.LoginResp{Token: token}, nil
}
```

#### 6.3.2 Hash 操作

```go
// Hash 操作
func (l *UserLogic) SaveUserInfo(userId string, userInfo map[string]string) error {
    // 1. Hset 设置单个字段
    err := l.svcCtx.Redis.Hset("user:info:"+userId, "nickname", userInfo["nickname"])
    
    // 2. Hmset 设置多个字段
    err = l.svcCtx.Redis.Hmset("user:info:"+userId, userInfo)
    
    // 3. Hget 获取单个字段
    nickname, err := l.svcCtx.Redis.Hget("user:info:"+userId, "nickname")
    
    // 4. Hgetall 获取所有字段
    userInfo, err := l.svcCtx.Redis.Hgetall("user:info:" + userId)
    
    // 5. Hdel 删除字段
    _, err = l.svcCtx.Redis.Hdel("user:info:"+userId, "nickname")
    
    // 6. Hexists 判断字段是否存在
    exists, err := l.svcCtx.Redis.Hexists("user:info:"+userId, "nickname")
    
    return nil
}
```

#### 6.3.3 List 操作

```go
// List 操作
func (l *MessageLogic) SaveMessage(userId string, message string) error {
    // 1. Lpush 从左边插入
    _, err := l.svcCtx.Redis.Lpush("user:messages:"+userId, message)
    
    // 2. Rpush 从右边插入
    _, err = l.svcCtx.Redis.Rpush("user:messages:"+userId, message)
    
    // 3. Lpop 从左边弹出
    message, err := l.svcCtx.Redis.Lpop("user:messages:" + userId)
    
    // 4. Rpop 从右边弹出
    message, err = l.svcCtx.Redis.Rpop("user:messages:" + userId)
    
    // 5. Lrange 获取范围内的元素
    messages, err := l.svcCtx.Redis.Lrange("user:messages:"+userId, 0, 9)
    
    // 6. Llen 获取列表长度
    length, err := l.svcCtx.Redis.Llen("user:messages:" + userId)
    
    return nil
}
```

#### 6.3.4 Set 操作

```go
// Set 操作
func (l *FriendLogic) AddFriend(userId string, friendId string) error {
    // 1. Sadd 添加成员
    _, err := l.svcCtx.Redis.Sadd("user:friends:"+userId, friendId)
    
    // 2. Srem 删除成员
    _, err = l.svcCtx.Redis.Srem("user:friends:"+userId, friendId)
    
    // 3. Smembers 获取所有成员
    friends, err := l.svcCtx.Redis.Smembers("user:friends:" + userId)
    
    // 4. Sismember 判断是否是成员
    isMember, err := l.svcCtx.Redis.Sismember("user:friends:"+userId, friendId)
    
    // 5. Scard 获取成员数量
    count, err := l.svcCtx.Redis.Scard("user:friends:" + userId)
    
    return nil
}
```

#### 6.3.5 Sorted Set 操作

```go
// Sorted Set 操作
func (l *RankLogic) UpdateScore(userId string, score int64) error {
    // 1. Zadd 添加成员及分数
    _, err := l.svcCtx.Redis.Zadd("rank:score", score, userId)
    
    // 2. Zincrby 增加分数
    _, err = l.svcCtx.Redis.Zincrby("rank:score", 10, userId)
    
    // 3. Zscore 获取分数
    score, err := l.svcCtx.Redis.Zscore("rank:score", userId)
    
    // 4. Zrank 获取排名（从小到大）
    rank, err := l.svcCtx.Redis.Zrank("rank:score", userId)
    
    // 5. Zrevrank 获取排名（从大到小）
    rank, err = l.svcCtx.Redis.Zrevrank("rank:score", userId)
    
    // 6. Zrange 获取范围内的成员（从小到大）
    members, err := l.svcCtx.Redis.Zrange("rank:score", 0, 9)
    
    // 7. Zrevrange 获取范围内的成员（从大到小）
    members, err = l.svcCtx.Redis.Zrevrange("rank:score", 0, 9)
    
    // 8. Zrem 删除成员
    _, err = l.svcCtx.Redis.Zrem("rank:score", userId)
    
    return nil
}
```

#### 6.3.6 缓存策略

```go
// 缓存策略示例
func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    // 1. 先从缓存获取
    cacheKey := "user:info:" + req.UserId
    var userInfo types.GetUserInfoResp
    
    // 尝试从缓存获取
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil && cached != "" {
        // 缓存命中，反序列化
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 2. 缓存未命中，从数据库获取
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
    if err != nil {
        return nil, err
    }
    
    // 3. 构造响应
    userInfo = types.GetUserInfoResp{
        Id:       user.Id,
        Phone:    user.Phone,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
        Sex:      int(user.Sex.Int64),
    }
    
    // 4. 写入缓存
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)  // 缓存 1 小时
    
    return &userInfo, nil
}
```

---


## 7. 中间件开发

### 7.1 中间件原理

Go-Zero 的中间件基于 HTTP Handler 链式调用实现，每个中间件都是一个函数，接收下一个 Handler 作为参数，返回一个新的 Handler。

```go
// 中间件函数签名
type Middleware func(next http.HandlerFunc) http.HandlerFunc
```

### 7.2 JWT 认证中间件

#### 7.2.1 在 API 文件中配置 JWT

```api
// apps/user/api/user.api

// 需要 JWT 认证的接口
@server(
    prefix: /api/user
    group: user
    jwt: Auth  // 启用 JWT 认证
)
service user-api {
    @handler detail
    get /detail/:userId (DetailReq) returns (DetailResp)
    
    @handler updateInfo
    put /info (UpdateInfoReq) returns (UpdateInfoResp)
}
```

#### 7.2.2 配置 JWT 密钥

```yaml
# apps/user/api/etc/user.yaml
Auth:
  AccessSecret: your-access-secret-key-change-me  # JWT 密钥
  AccessExpire: 86400                             # 过期时间（秒）
```

#### 7.2.3 JWT 认证流程

1. 客户端在请求头中携带 Token：`Authorization: Bearer <token>`
2. 中间件解析 Token，验证签名和过期时间
3. 将用户 ID 等信息存入 Context
4. 在 Logic 中获取用户 ID

```go
// 在 Logic 中获取用户 ID
import "imooc.com/easy-chat/pkg/ctxdata"

func (l *DetailLogic) Detail(req *types.DetailReq) (*types.DetailResp, error) {
    // 从 Context 中获取用户 ID
    userId := ctxdata.GetUId(l.ctx)
    
    // 验证权限
    if userId != req.UserId {
        return nil, errors.New("无权访问")
    }
    
    // 查询用户信息
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &types.DetailResp{
        Id:       user.Id,
        Phone:    user.Phone,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
        Sex:      int(user.Sex.Int64),
    }, nil
}
```

#### 7.2.4 自定义 JWT 工具

```go
// pkg/ctxdata/jwt.go
package ctxdata

import (
    "context"
    "github.com/golang-jwt/jwt/v4"
    "time"
)

const CtxKeyJwtUserId = "userId"

// GetJwtToken 生成 JWT Token
func GetJwtToken(secret string, iat, seconds int64, userId string) (string, error) {
    claims := make(jwt.MapClaims)
    claims["exp"] = iat + seconds  // 过期时间
    claims["iat"] = iat            // 签发时间
    claims[CtxKeyJwtUserId] = userId  // 用户 ID
    
    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims = claims
    
    return token.SignedString([]byte(secret))
}

// GetUId 从 Context 中获取用户 ID
func GetUId(ctx context.Context) string {
    if uid, ok := ctx.Value(CtxKeyJwtUserId).(string); ok {
        return uid
    }
    return ""
}
```

### 7.3 日志中间件

#### 7.3.1 创建日志中间件

```go
// pkg/middleware/logmiddleware.go
package middleware

import (
    "fmt"
    "github.com/zeromicro/go-zero/core/logx"
    "net/http"
    "time"
)

// LogMiddleware 日志中间件
type LogMiddleware struct {
}

// NewLogMiddleware 创建日志中间件
func NewLogMiddleware() *LogMiddleware {
    return &LogMiddleware{}
}

// Handle 处理请求
func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 记录开始时间
        startTime := time.Now()
        
        // 记录请求信息
        logx.Infof("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
        
        // 创建响应包装器（用于捕获状态码）
        wrapper := &responseWrapper{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        // 调用下一个处理器
        next(wrapper, r)
        
        // 计算耗时
        duration := time.Since(startTime)
        
        // 记录响应信息
        logx.Infof("[%s] %s %s - %d - %v",
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
            wrapper.statusCode,
            duration,
        )
    }
}

// responseWrapper 响应包装器
type responseWrapper struct {
    http.ResponseWriter
    statusCode int
}

// WriteHeader 写入状态码
func (w *responseWrapper) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
```

#### 7.3.2 注册日志中间件

```go
// apps/user/api/internal/handler/routes.go
package handler

import (
    "net/http"
    
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    "imooc.com/easy-chat/apps/user/api/internal/handler/user"
    "imooc.com/easy-chat/pkg/middleware"
    
    "github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 注册全局中间件
    server.Use(middleware.NewLogMiddleware().Handle)
    
    server.AddRoutes(
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/register",
                Handler: user.RegisterHandler(serverCtx),
            },
            {
                Method:  http.MethodPost,
                Path:    "/login",
                Handler: user.LoginHandler(serverCtx),
            },
        },
        rest.WithPrefix("/api/user"),
    )
}
```

### 7.4 限流中间件

#### 7.4.1 创建限流中间件

```go
// apps/social/api/internal/middleware/limitmiddleware.go
package middleware

import (
    "net/http"
    "github.com/zeromicro/go-zero/core/limit"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// LimitMiddleware 限流中间件
type LimitMiddleware struct {
    limiter *limit.PeriodLimit  // 限流器
}

// NewLimitMiddleware 创建限流中间件
// 参数:
//   - redis: Redis 客户端
//   - seconds: 时间窗口（秒）
//   - quota: 配额（时间窗口内允许的请求数）
func NewLimitMiddleware(redis *redis.Redis, seconds, quota int) *LimitMiddleware {
    return &LimitMiddleware{
        limiter: limit.NewPeriodLimit(seconds, quota, redis, "limit"),
    }
}

// Handle 处理请求
func (m *LimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 获取限流 Key（可以使用 IP、用户 ID 等）
        key := r.RemoteAddr  // 使用 IP 作为限流 Key
        
        // 判断是否超过限流
        code, err := m.limiter.Take(key)
        if err != nil {
            httpx.Error(w, err)
            return
        }
        
        // code 说明：
        // limit.OverQuota: 超过配额
        // limit.Allowed: 允许通过
        // limit.HitQuota: 达到配额上限
        switch code {
        case limit.OverQuota:
            // 超过限流
            httpx.Error(w, errors.New("请求过于频繁，请稍后再试"))
            return
        case limit.Allowed:
            // 允许通过
            next(w, r)
        case limit.HitQuota:
            // 达到配额上限，但仍允许通过
            next(w, r)
        default:
            httpx.Error(w, errors.New("系统错误"))
            return
        }
    }
}
```

#### 7.4.2 使用限流中间件

```go
// apps/social/api/internal/svc/servicecontext.go
package svc

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "imooc.com/easy-chat/apps/social/api/internal/config"
    "imooc.com/easy-chat/apps/social/api/internal/middleware"
)

type ServiceContext struct {
    Config config.Config
    
    *redis.Redis
    
    // 限流中间件
    LimitMiddleware *middleware.LimitMiddleware
}

func NewServiceContext(c config.Config) *ServiceContext {
    rds := redis.MustNewRedis(c.Redisx)
    
    return &ServiceContext{
        Config: c,
        Redis:  rds,
        
        // 创建限流中间件：60 秒内最多 100 次请求
        LimitMiddleware: middleware.NewLimitMiddleware(rds, 60, 100),
    }
}
```

```go
// apps/social/api/internal/handler/routes.go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/friend/putin",
                Handler: friend.FriendPutinHandler(serverCtx),
            },
        },
        rest.WithPrefix("/api/social"),
        // 使用限流中间件
        rest.WithMiddlewares([]rest.Middleware{
            serverCtx.LimitMiddleware.Handle,
        }),
    )
}
```

### 7.5 幂等性中间件

#### 7.5.1 创建幂等性中间件

```go
// apps/social/api/internal/middleware/idempotencemiddleware.go
package middleware

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/rest/httpx"
    "io"
    "net/http"
    "time"
)

// IdempotenceMiddleware 幂等性中间件
type IdempotenceMiddleware struct {
    redis *redis.Redis
}

// NewIdempotenceMiddleware 创建幂等性中间件
func NewIdempotenceMiddleware(redis *redis.Redis) *IdempotenceMiddleware {
    return &IdempotenceMiddleware{
        redis: redis,
    }
}

// Handle 处理请求
func (m *IdempotenceMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 生成幂等性 Key
        idempotenceKey, err := m.generateKey(r)
        if err != nil {
            httpx.Error(w, err)
            return
        }
        
        // 2. 检查是否已处理过
        cacheKey := fmt.Sprintf("idempotence:%s", idempotenceKey)
        exists, err := m.redis.Exists(cacheKey)
        if err != nil {
            httpx.Error(w, err)
            return
        }
        
        if exists {
            // 已处理过，返回缓存的结果
            result, err := m.redis.Get(cacheKey)
            if err != nil {
                httpx.Error(w, err)
                return
            }
            
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(result))
            return
        }
        
        // 3. 设置处理标记（防止并发重复处理）
        ok, err := m.redis.Setnx(cacheKey+":lock", "1")
        if err != nil || !ok {
            httpx.Error(w, errors.New("请求正在处理中，请勿重复提交"))
            return
        }
        m.redis.Expire(cacheKey+":lock", 10)  // 锁定 10 秒
        
        // 4. 创建响应捕获器
        recorder := &responseRecorder{
            ResponseWriter: w,
            body:           []byte{},
        }
        
        // 5. 调用下一个处理器
        next(recorder, r)
        
        // 6. 缓存响应结果
        if recorder.statusCode == http.StatusOK {
            m.redis.Setex(cacheKey, string(recorder.body), 300)  // 缓存 5 分钟
        }
        
        // 7. 删除处理标记
        m.redis.Del(cacheKey + ":lock")
    }
}

// generateKey 生成幂等性 Key
func (m *IdempotenceMiddleware) generateKey(r *http.Request) (string, error) {
    // 读取请求体
    body, err := io.ReadAll(r.Body)
    if err != nil {
        return "", err
    }
    
    // 恢复请求体（供后续处理使用）
    r.Body = io.NopCloser(bytes.NewBuffer(body))
    
    // 生成 Key：Method + Path + Body 的 MD5
    h := md5.New()
    h.Write([]byte(r.Method))
    h.Write([]byte(r.URL.Path))
    h.Write(body)
    
    return hex.EncodeToString(h.Sum(nil)), nil
}

// responseRecorder 响应记录器
type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    body       []byte
}

// WriteHeader 写入状态码
func (r *responseRecorder) WriteHeader(statusCode int) {
    r.statusCode = statusCode
    r.ResponseWriter.WriteHeader(statusCode)
}

// Write 写入响应体
func (r *responseRecorder) Write(body []byte) (int, error) {
    r.body = body
    return r.ResponseWriter.Write(body)
}
```

### 7.6 跨域中间件

#### 7.6.1 配置跨域

```yaml
# apps/user/api/etc/user.yaml
CorsOrigins:
  - "*"  # 允许所有源（生产环境建议指定具体域名）
  # - "http://localhost:3000"
  # - "https://example.com"
```

#### 7.6.2 自定义跨域中间件

```go
// pkg/middleware/corsmiddleware.go
package middleware

import (
    "net/http"
)

// CorsMiddleware 跨域中间件
type CorsMiddleware struct {
    allowOrigins []string
}

// NewCorsMiddleware 创建跨域中间件
func NewCorsMiddleware(allowOrigins []string) *CorsMiddleware {
    return &CorsMiddleware{
        allowOrigins: allowOrigins,
    }
}

// Handle 处理请求
func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 设置跨域响应头
        origin := r.Header.Get("Origin")
        if m.isAllowedOrigin(origin) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Max-Age", "3600")
        
        // 处理 OPTIONS 预检请求
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        
        // 调用下一个处理器
        next(w, r)
    }
}

// isAllowedOrigin 判断是否允许的源
func (m *CorsMiddleware) isAllowedOrigin(origin string) bool {
    for _, allowOrigin := range m.allowOrigins {
        if allowOrigin == "*" || allowOrigin == origin {
            return true
        }
    }
    return false
}
```

### 7.7 自定义中间件示例

#### 7.7.1 签名验证中间件

```go
// pkg/middleware/signaturemiddleware.go
package middleware

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/zeromicro/go-zero/rest/httpx"
    "net/http"
    "sort"
    "strings"
)

// SignatureMiddleware 签名验证中间件
type SignatureMiddleware struct {
    secret string  // 签名密钥
}

// NewSignatureMiddleware 创建签名验证中间件
func NewSignatureMiddleware(secret string) *SignatureMiddleware {
    return &SignatureMiddleware{
        secret: secret,
    }
}

// Handle 处理请求
func (m *SignatureMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 获取签名
        signature := r.Header.Get("X-Signature")
        if signature == "" {
            httpx.Error(w, errors.New("缺少签名"))
            return
        }
        
        // 2. 获取时间戳
        timestamp := r.Header.Get("X-Timestamp")
        if timestamp == "" {
            httpx.Error(w, errors.New("缺少时间戳"))
            return
        }
        
        // 3. 验证时间戳（防止重放攻击）
        // TODO: 验证时间戳是否在有效期内
        
        // 4. 生成签名
        expectedSignature := m.generateSignature(r, timestamp)
        
        // 5. 验证签名
        if signature != expectedSignature {
            httpx.Error(w, errors.New("签名验证失败"))
            return
        }
        
        // 6. 调用下一个处理器
        next(w, r)
    }
}

// generateSignature 生成签名
func (m *SignatureMiddleware) generateSignature(r *http.Request, timestamp string) string {
    // 1. 获取所有查询参数
    params := make(map[string]string)
    for key, values := range r.URL.Query() {
        params[key] = values[0]
    }
    
    // 2. 添加时间戳
    params["timestamp"] = timestamp
    
    // 3. 按 Key 排序
    keys := make([]string, 0, len(params))
    for key := range params {
        keys = append(keys, key)
    }
    sort.Strings(keys)
    
    // 4. 拼接字符串
    var builder strings.Builder
    for _, key := range keys {
        builder.WriteString(key)
        builder.WriteString("=")
        builder.WriteString(params[key])
        builder.WriteString("&")
    }
    builder.WriteString("secret=")
    builder.WriteString(m.secret)
    
    // 5. 计算 MD5
    h := md5.New()
    h.Write([]byte(builder.String()))
    
    return hex.EncodeToString(h.Sum(nil))
}
```

### 7.8 中间件执行顺序

中间件的执行顺序是从外到内，响应顺序是从内到外：

```
请求 -> 中间件1 -> 中间件2 -> 中间件3 -> Handler -> 中间件3 -> 中间件2 -> 中间件1 -> 响应
```

```go
// 注册中间件
server.Use(middleware1.Handle)  // 最外层
server.Use(middleware2.Handle)
server.Use(middleware3.Handle)  // 最内层

// 或者使用 WithMiddlewares
server.AddRoutes(
    routes,
    rest.WithMiddlewares([]rest.Middleware{
        middleware1.Handle,
        middleware2.Handle,
        middleware3.Handle,
    }),
)
```

---


## 8. 服务间通信

### 8.1 RPC 调用示例

#### 8.1.1 API 服务调用 RPC 服务

```go
// apps/social/api/internal/logic/friend/friendlistlogic.go
package friend

import (
    "context"
    "imooc.com/easy-chat/apps/social/rpc/socialclient"
    
    "imooc.com/easy-chat/apps/social/api/internal/svc"
    "imooc.com/easy-chat/apps/social/api/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
    return &FriendListLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// FriendList 获取好友列表
func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
    // 1. 调用 Social RPC 服务获取好友列表
    friendListResp, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
        UserId: req.UserId,
    })
    if err != nil {
        return nil, err
    }
    
    // 2. 获取好友 ID 列表
    var friendIds []string
    for _, friend := range friendListResp.List {
        friendIds = append(friendIds, friend.FriendId)
    }
    
    // 3. 调用 User RPC 服务批量获取用户信息
    usersResp, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
        Ids: friendIds,
    })
    if err != nil {
        return nil, err
    }
    
    // 4. 构造响应
    var list []*types.Friend
    for _, user := range usersResp.Users {
        list = append(list, &types.Friend{
            Id:       user.Id,
            Nickname: user.Nickname,
            Avatar:   user.Avatar,
        })
    }
    
    return &types.FriendListResp{
        List: list,
    }, nil
}
```

#### 8.1.2 RPC 服务调用 RPC 服务

```go
// apps/im/rpc/internal/logic/setupuserconversationlogic.go
package logic

import (
    "context"
    "imooc.com/easy-chat/apps/social/rpc/socialclient"
    
    "imooc.com/easy-chat/apps/im/rpc/internal/svc"
    "imooc.com/easy-chat/apps/im/rpc/im"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type SetupUserConversationLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewSetupUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetupUserConversationLogic {
    return &SetupUserConversationLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// SetupUserConversation 建立用户会话
func (l *SetupUserConversationLogic) SetupUserConversation(in *im.SetupUserConversationReq) (*im.SetupUserConversationResp, error) {
    // 1. 调用 Social RPC 验证好友关系
    friendResp, err := l.svcCtx.Social.IsFriend(l.ctx, &socialclient.IsFriendReq{
        UserId:   in.SendId,
        FriendId: in.RecvId,
    })
    if err != nil {
        return nil, err
    }
    
    if !friendResp.IsFriend {
        return nil, errors.New("不是好友关系，无法发送消息")
    }
    
    // 2. 创建会话
    // ...
    
    return &im.SetupUserConversationResp{}, nil
}
```

### 8.2 服务发现

#### 8.2.1 基于 Etcd 的服务发现

```yaml
# RPC 服务端配置（注册服务）
Name: user.rpc
ListenOn: 0.0.0.0:10001

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: user.rpc  # 服务注册的 Key

# RPC 客户端配置（发现服务）
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc  # 服务发现的 Key
```

#### 8.2.2 直连模式（不使用服务发现）

```yaml
# RPC 客户端配置（直连模式）
UserRpc:
  Endpoints:
    - 127.0.0.1:10001
    - 127.0.0.1:10002  # 多个实例
  Timeout: 5000
```

### 8.3 负载均衡

Go-Zero 默认使用 gRPC 的负载均衡策略，支持以下几种：

#### 8.3.1 轮询（Round Robin）

```go
// 默认使用轮询策略
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
```

#### 8.3.2 自定义负载均衡

```go
// apps/user/api/internal/svc/servicecontext.go
import (
    "github.com/zeromicro/go-zero/zrpc"
    "google.golang.org/grpc/balancer/roundrobin"
)

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        
        // 使用轮询负载均衡
        User: userclient.NewUser(
            zrpc.MustNewClient(
                c.UserRpc,
                zrpc.WithDialOption(
                    grpc.WithDefaultServiceConfig(
                        fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, roundrobin.Name),
                    ),
                ),
            ),
        ),
    }
}
```

### 8.4 超时控制

#### 8.4.1 配置超时时间

```yaml
# RPC 客户端配置
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
  Timeout: 5000  # 超时时间（毫秒）
```

#### 8.4.2 在代码中设置超时

```go
import (
    "context"
    "time"
)

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 创建带超时的 Context
    ctx, cancel := context.WithTimeout(l.ctx, 3*time.Second)
    defer cancel()
    
    // 使用带超时的 Context 调用 RPC
    registerResp, err := l.svcCtx.User.Register(ctx, &user.RegisterReq{
        Phone:    req.Phone,
        Nickname: req.Nickname,
        Password: req.Password,
        Avatar:   req.Avatar,
        Sex:      int32(req.Sex),
    })
    if err != nil {
        return nil, err
    }
    
    return &types.RegisterResp{
        Token:  registerResp.Token,
        Expire: registerResp.Expire,
    }, nil
}
```

### 8.5 重试机制

#### 8.5.1 配置重试策略

```go
// apps/user/api/internal/svc/servicecontext.go
var retryPolicy = `{
    "methodConfig" : [{
        "name": [{
            "service": "user.User"
        }],
        "waitForReady": true,
        "retryPolicy": {
            "maxAttempts": 5,              // 最大重试次数
            "initialBackoff": "0.001s",    // 初始退避时间
            "maxBackoff": "0.002s",        // 最大退避时间
            "backoffMultiplier": 1.0,      // 退避倍数
            "retryableStatusCodes": ["UNKNOWN", "UNAVAILABLE"]  // 可重试的状态码
        }
    }]
}`

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        
        // 使用重试策略
        User: userclient.NewUser(
            zrpc.MustNewClient(
                c.UserRpc,
                zrpc.WithDialOption(
                    grpc.WithDefaultServiceConfig(retryPolicy),
                ),
            ),
        ),
    }
}
```

#### 8.5.2 自定义重试逻辑

```go
import (
    "context"
    "time"
)

// retryCall 重试调用
func retryCall(ctx context.Context, maxRetries int, fn func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        
        // 判断是否需要重试
        if !isRetryableError(err) {
            return err
        }
        
        // 等待一段时间后重试
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
        }
    }
    return err
}

// isRetryableError 判断是否可重试的错误
func isRetryableError(err error) bool {
    // 根据错误类型判断是否可重试
    // 例如：网络错误、超时错误等
    return true
}

// 使用示例
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    var registerResp *user.RegisterResp
    
    // 重试调用 RPC
    err = retryCall(l.ctx, 3, func() error {
        var err error
        registerResp, err = l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
            Phone:    req.Phone,
            Nickname: req.Nickname,
            Password: req.Password,
            Avatar:   req.Avatar,
            Sex:      int32(req.Sex),
        })
        return err
    })
    
    if err != nil {
        return nil, err
    }
    
    return &types.RegisterResp{
        Token:  registerResp.Token,
        Expire: registerResp.Expire,
    }, nil
}
```

---

## 9. 缓存使用

### 9.1 缓存策略

#### 9.1.1 Cache Aside（旁路缓存）

最常用的缓存策略，读取数据时先查缓存，缓存未命中再查数据库，然后更新缓存。

```go
func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    // 1. 先从缓存获取
    cacheKey := fmt.Sprintf("user:info:%s", req.UserId)
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil && cached != "" {
        // 缓存命中
        var userInfo types.GetUserInfoResp
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 2. 缓存未命中，从数据库获取
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
    if err != nil {
        return nil, err
    }
    
    // 3. 构造响应
    userInfo := &types.GetUserInfoResp{
        Id:       user.Id,
        Phone:    user.Phone,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
        Sex:      int(user.Sex.Int64),
    }
    
    // 4. 写入缓存
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)  // 缓存 1 小时
    
    return userInfo, nil
}
```

#### 9.1.2 Read Through（读穿透）

由缓存层负责从数据库加载数据，应用层只与缓存层交互。

```go
// 封装缓存读取逻辑
func (l *GetUserInfoLogic) getUserInfoWithCache(userId string) (*types.GetUserInfoResp, error) {
    cacheKey := fmt.Sprintf("user:info:%s", userId)
    
    // 尝试从缓存获取
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil && cached != "" {
        var userInfo types.GetUserInfoResp
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 缓存未命中，从数据库加载
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
    if err != nil {
        return nil, err
    }
    
    userInfo := &types.GetUserInfoResp{
        Id:       user.Id,
        Phone:    user.Phone,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
        Sex:      int(user.Sex.Int64),
    }
    
    // 写入缓存
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
    
    return userInfo, nil
}
```

#### 9.1.3 Write Through（写穿透）

更新数据时，先更新缓存，再由缓存层更新数据库。

```go
func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) error {
    // 1. 更新数据库
    err := l.svcCtx.UsersModel.Update(l.ctx, &models.Users{
        Id:       req.UserId,
        Nickname: req.Nickname,
        Avatar:   req.Avatar,
        Sex:      sql.NullInt64{Int64: int64(req.Sex), Valid: true},
    })
    if err != nil {
        return err
    }
    
    // 2. 删除缓存（让下次读取时重新加载）
    cacheKey := fmt.Sprintf("user:info:%s", req.UserId)
    l.svcCtx.Redis.Del(cacheKey)
    
    return nil
}
```

#### 9.1.4 Write Behind（写回）

更新数据时只更新缓存，由缓存层异步批量更新数据库。

```go
func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) error {
    // 1. 更新缓存
    cacheKey := fmt.Sprintf("user:info:%s", req.UserId)
    userInfo := map[string]interface{}{
        "id":       req.UserId,
        "nickname": req.Nickname,
        "avatar":   req.Avatar,
        "sex":      req.Sex,
    }
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
    
    // 2. 将更新操作放入队列（异步更新数据库）
    updateData := map[string]interface{}{
        "user_id":  req.UserId,
        "nickname": req.Nickname,
        "avatar":   req.Avatar,
        "sex":      req.Sex,
    }
    updateJson, _ := json.Marshal(updateData)
    l.svcCtx.Redis.Lpush("user:update:queue", string(updateJson))
    
    return nil
}
```

### 9.2 Redis 使用示例

#### 9.2.1 用户在线状态

```go
// 设置用户在线
func (l *UserOnlineLogic) SetUserOnline(userId string) error {
    // 使用 Set 存储在线用户
    _, err := l.svcCtx.Redis.Sadd("users:online", userId)
    if err != nil {
        return err
    }
    
    // 设置用户最后活跃时间
    now := time.Now().Unix()
    err = l.svcCtx.Redis.Hset("users:last:active", userId, fmt.Sprintf("%d", now))
    if err != nil {
        return err
    }
    
    return nil
}

// 设置用户离线
func (l *UserOfflineLogic) SetUserOffline(userId string) error {
    // 从在线用户集合中移除
    _, err := l.svcCtx.Redis.Srem("users:online", userId)
    return err
}

// 获取在线用户列表
func (l *GetOnlineUsersLogic) GetOnlineUsers() ([]string, error) {
    return l.svcCtx.Redis.Smembers("users:online")
}

// 判断用户是否在线
func (l *IsUserOnlineLogic) IsUserOnline(userId string) (bool, error) {
    return l.svcCtx.Redis.Sismember("users:online", userId)
}
```

#### 9.2.2 消息未读数

```go
// 增加未读数
func (l *IncrUnreadLogic) IncrUnread(userId string, conversationId string) error {
    key := fmt.Sprintf("user:%s:unread", userId)
    _, err := l.svcCtx.Redis.Hincrby(key, conversationId, 1)
    return err
}

// 获取未读数
func (l *GetUnreadLogic) GetUnread(userId string, conversationId string) (int64, error) {
    key := fmt.Sprintf("user:%s:unread", userId)
    unread, err := l.svcCtx.Redis.Hget(key, conversationId)
    if err != nil {
        return 0, err
    }
    return strconv.ParseInt(unread, 10, 64)
}

// 清空未读数
func (l *ClearUnreadLogic) ClearUnread(userId string, conversationId string) error {
    key := fmt.Sprintf("user:%s:unread", userId)
    _, err := l.svcCtx.Redis.Hdel(key, conversationId)
    return err
}

// 获取所有未读数
func (l *GetAllUnreadLogic) GetAllUnread(userId string) (map[string]int64, error) {
    key := fmt.Sprintf("user:%s:unread", userId)
    unreadMap, err := l.svcCtx.Redis.Hgetall(key)
    if err != nil {
        return nil, err
    }
    
    result := make(map[string]int64)
    for conversationId, unread := range unreadMap {
        count, _ := strconv.ParseInt(unread, 10, 64)
        result[conversationId] = count
    }
    
    return result, nil
}
```

### 9.3 缓存穿透、击穿、雪崩的解决方案

#### 9.3.1 缓存穿透

缓存穿透是指查询一个不存在的数据，缓存和数据库都没有，导致每次请求都会打到数据库。

**解决方案1：缓存空值**

```go
func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    cacheKey := fmt.Sprintf("user:info:%s", req.UserId)
    
    // 1. 从缓存获取
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil {
        if cached == "null" {
            // 缓存的空值
            return nil, errors.New("用户不存在")
        }
        
        var userInfo types.GetUserInfoResp
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 2. 从数据库获取
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
    if err != nil {
        if err == models.ErrNotFound {
            // 缓存空值，防止缓存穿透
            l.svcCtx.Redis.Setex(cacheKey, "null", 300)  // 缓存 5 分钟
        }
        return nil, err
    }
    
    // 3. 写入缓存
    userInfo := &types.GetUserInfoResp{
        Id:       user.Id,
        Phone:    user.Phone,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
        Sex:      int(user.Sex.Int64),
    }
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
    
    return userInfo, nil
}
```

**解决方案2：布隆过滤器**

```go
import "github.com/bits-and-blooms/bloom/v3"

// 初始化布隆过滤器
func NewServiceContext(c config.Config) *ServiceContext {
    // 创建布隆过滤器（预计 100 万个元素，误判率 0.01%）
    filter := bloom.NewWithEstimates(1000000, 0.0001)
    
    // 将所有用户 ID 加入布隆过滤器
    // TODO: 从数据库加载所有用户 ID
    
    return &ServiceContext{
        Config:      c,
        BloomFilter: filter,
    }
}

// 使用布隆过滤器
func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    // 1. 先用布隆过滤器判断
    if !l.svcCtx.BloomFilter.TestString(req.UserId) {
        // 布隆过滤器判断不存在，直接返回
        return nil, errors.New("用户不存在")
    }
    
    // 2. 从缓存获取
    // ...
    
    // 3. 从数据库获取
    // ...
}
```

#### 9.3.2 缓存击穿

缓存击穿是指一个热点 Key 过期，大量请求同时打到数据库。

**解决方案：互斥锁**

```go
import "github.com/zeromicro/go-zero/core/syncx"

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    cacheKey := fmt.Sprintf("user:info:%s", req.UserId)
    
    // 1. 从缓存获取
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil && cached != "" {
        var userInfo types.GetUserInfoResp
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 2. 使用互斥锁，防止缓存击穿
    lockKey := fmt.Sprintf("lock:user:info:%s", req.UserId)
    lock := syncx.NewSharedCalls()
    
    val, err := lock.Do(lockKey, func() (interface{}, error) {
        // 再次尝试从缓存获取（可能已被其他请求加载）
        cached, err := l.svcCtx.Redis.Get(cacheKey)
        if err == nil && cached != "" {
            var userInfo types.GetUserInfoResp
            json.Unmarshal([]byte(cached), &userInfo)
            return &userInfo, nil
        }
        
        // 从数据库获取
        user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
        if err != nil {
            return nil, err
        }
        
        userInfo := &types.GetUserInfoResp{
            Id:       user.Id,
            Phone:    user.Phone,
            Nickname: user.Nickname,
            Avatar:   user.Avatar,
            Sex:      int(user.Sex.Int64),
        }
        
        // 写入缓存
        data, _ := json.Marshal(userInfo)
        l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
        
        return userInfo, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return val.(*types.GetUserInfoResp), nil
}
```

#### 9.3.3 缓存雪崩

缓存雪崩是指大量缓存同时过期，导致大量请求打到数据库。

**解决方案：随机过期时间**

```go
import (
    "math/rand"
    "time"
)

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoResp, error) {
    // ...
    
    // 写入缓存时，设置随机过期时间
    baseExpire := 3600  // 基础过期时间 1 小时
    randomExpire := rand.Intn(300)  // 随机 0-5 分钟
    expire := baseExpire + randomExpire
    
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), expire)
    
    return userInfo, nil
}
```

---


## 10. 消息队列

### 10.1 Kafka 集成

#### 10.1.1 安装依赖

```bash
go get -u github.com/segmentio/kafka-go
```

#### 10.1.2 配置 Kafka

```yaml
# apps/task/mq/etc/task.yaml
Name: task.mq

# Kafka 配置
Kafka:
  Brokers:
    - 127.0.0.1:9092
  Group: task-group
  Topics:
    - chat-message      # 聊天消息
    - message-read      # 消息已读
    - group-message     # 群聊消息
```

```go
// apps/task/mq/internal/config/config.go
package config

type Config struct {
    Kafka struct {
        Brokers []string  // Kafka 地址列表
        Group   string    // 消费者组
        Topics  []string  // 主题列表
    }
}
```

### 10.2 生产者示例

#### 10.2.1 创建生产者

```go
// pkg/kafka/producer.go
package kafka

import (
    "context"
    "github.com/segmentio/kafka-go"
    "time"
)

// Producer Kafka 生产者
type Producer struct {
    writer *kafka.Writer
}

// NewProducer 创建生产者
func NewProducer(brokers []string, topic string) *Producer {
    return &Producer{
        writer: &kafka.Writer{
            Addr:         kafka.TCP(brokers...),
            Topic:        topic,
            Balancer:     &kafka.LeastBytes{},  // 负载均衡策略
            WriteTimeout: 10 * time.Second,
            ReadTimeout:  10 * time.Second,
            RequiredAcks: kafka.RequireOne,     // 确认级别
            Async:        false,                // 同步发送
        },
    }
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, key, value string) error {
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(key),
        Value: []byte(value),
        Time:  time.Now(),
    })
}

// SendMessages 批量发送消息
func (p *Producer) SendMessages(ctx context.Context, messages []kafka.Message) error {
    return p.writer.WriteMessages(ctx, messages...)
}

// Close 关闭生产者
func (p *Producer) Close() error {
    return p.writer.Close()
}
```

#### 10.2.2 使用生产者

```go
// apps/im/rpc/internal/logic/sendmessagelogic.go
package logic

import (
    "context"
    "encoding/json"
    "imooc.com/easy-chat/pkg/kafka"
    
    "imooc.com/easy-chat/apps/im/rpc/internal/svc"
    "imooc.com/easy-chat/apps/im/rpc/im"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
    return &SendMessageLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// SendMessage 发送消息
func (l *SendMessageLogic) SendMessage(in *im.SendMessageReq) (*im.SendMessageResp, error) {
    // 1. 构造消息
    message := map[string]interface{}{
        "conversation_id": in.ConversationId,
        "send_id":         in.SendId,
        "recv_id":         in.RecvId,
        "msg_type":        in.MsgType,
        "msg_content":     in.MsgContent,
        "chat_type":       in.ChatType,
        "send_time":       time.Now().Unix(),
    }
    
    // 2. 序列化消息
    messageJson, err := json.Marshal(message)
    if err != nil {
        return nil, err
    }
    
    // 3. 发送到 Kafka
    err = l.svcCtx.KafkaProducer.SendMessage(l.ctx, in.ConversationId, string(messageJson))
    if err != nil {
        return nil, err
    }
    
    return &im.SendMessageResp{
        MessageId: primitive.NewObjectID().Hex(),
    }, nil
}
```

### 10.3 消费者示例

#### 10.3.1 创建消费者

```go
// pkg/kafka/consumer.go
package kafka

import (
    "context"
    "github.com/segmentio/kafka-go"
    "time"
)

// Consumer Kafka 消费者
type Consumer struct {
    reader *kafka.Reader
}

// NewConsumer 创建消费者
func NewConsumer(brokers []string, group, topic string) *Consumer {
    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:        brokers,
            GroupID:        group,
            Topic:          topic,
            MinBytes:       10e3,            // 10KB
            MaxBytes:       10e6,            // 10MB
            CommitInterval: time.Second,     // 提交间隔
            StartOffset:    kafka.LastOffset, // 从最新位置开始消费
        }),
    }
}

// ConsumeMessages 消费消息
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func(message kafka.Message) error) error {
    for {
        // 读取消息
        message, err := c.reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        // 处理消息
        if err := handler(message); err != nil {
            // 处理失败，记录日志
            logx.Errorf("Handle message failed: %v", err)
            continue
        }
        
        // 提交偏移量
        if err := c.reader.CommitMessages(ctx, message); err != nil {
            logx.Errorf("Commit message failed: %v", err)
        }
    }
}

// Close 关闭消费者
func (c *Consumer) Close() error {
    return c.reader.Close()
}
```

#### 10.3.2 使用消费者

```go
// apps/task/mq/internal/handler/msgTransfer/msgChatTrasnfer.go
package msgTransfer

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
    "imooc.com/easy-chat/apps/task/mq/internal/svc"
    
    "github.com/zeromicro/go-zero/core/logx"
)

// MsgChatTransfer 聊天消息转发
type MsgChatTransfer struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

// NewMsgChatTransfer 创建聊天消息转发
func NewMsgChatTransfer(ctx context.Context, svcCtx *svc.ServiceContext) *MsgChatTransfer {
    return &MsgChatTransfer{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// Consume 消费消息
func (m *MsgChatTransfer) Consume() error {
    // 创建消费者
    consumer := kafka.NewConsumer(
        m.svcCtx.Config.Kafka.Brokers,
        m.svcCtx.Config.Kafka.Group,
        "chat-message",
    )
    defer consumer.Close()
    
    // 消费消息
    return consumer.ConsumeMessages(m.ctx, m.handleMessage)
}

// handleMessage 处理消息
func (m *MsgChatTransfer) handleMessage(message kafka.Message) error {
    // 1. 解析消息
    var msg map[string]interface{}
    if err := json.Unmarshal(message.Value, &msg); err != nil {
        return err
    }
    
    logx.Infof("Received message: %+v", msg)
    
    // 2. 保存消息到 MongoDB
    chatLog := &models.ChatLog{
        ConversationId: msg["conversation_id"].(string),
        SendId:         msg["send_id"].(string),
        RecvId:         msg["recv_id"].(string),
        MsgType:        int(msg["msg_type"].(float64)),
        MsgContent:     msg["msg_content"].(string),
        ChatType:       int(msg["chat_type"].(float64)),
        SendTime:       time.Unix(int64(msg["send_time"].(float64)), 0),
    }
    
    if err := m.svcCtx.ChatLogModel.Insert(m.ctx, chatLog); err != nil {
        return err
    }
    
    // 3. 推送消息给接收者（通过 WebSocket）
    // TODO: 实现消息推送逻辑
    
    return nil
}
```

#### 10.3.3 启动消费者

```go
// apps/task/mq/task.go
package main

import (
    "context"
    "flag"
    "fmt"
    
    "imooc.com/easy-chat/apps/task/mq/internal/config"
    "imooc.com/easy-chat/apps/task/mq/internal/handler/msgTransfer"
    "imooc.com/easy-chat/apps/task/mq/internal/svc"
    
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/task.yaml", "the config file")

func main() {
    flag.Parse()
    
    // 加载配置
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    // 创建服务上下文
    ctx := svc.NewServiceContext(c)
    
    // 创建消息转发器
    msgChatTransfer := msgTransfer.NewMsgChatTransfer(context.Background(), ctx)
    msgReadTransfer := msgTransfer.NewMsgReadTransfer(context.Background(), ctx)
    
    // 启动消费者
    go func() {
        if err := msgChatTransfer.Consume(); err != nil {
            logx.Errorf("MsgChatTransfer consume failed: %v", err)
        }
    }()
    
    go func() {
        if err := msgReadTransfer.Consume(); err != nil {
            logx.Errorf("MsgReadTransfer consume failed: %v", err)
        }
    }()
    
    fmt.Println("Task MQ started...")
    select {}  // 阻塞主线程
}
```

---

## 11. 日志和监控

### 11.1 日志配置

#### 11.1.1 日志级别

```yaml
# 日志配置
Log:
  ServiceName: user-api
  Mode: console           # console（控制台）、file（文件）
  Level: info             # debug、info、warn、error
  Encoding: plain         # plain（纯文本）、json（JSON）
  Path: logs              # 日志文件路径
  KeepDays: 7             # 日志保留天数
  Compress: true          # 是否压缩
  MaxBackups: 5           # 最大备份文件数
  MaxSize: 100            # 单个文件最大大小（MB）
```

#### 11.1.2 使用日志

```go
import "github.com/zeromicro/go-zero/core/logx"

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // Debug 日志（开发环境）
    logx.Debugf("Register request: %+v", req)
    
    // Info 日志（正常信息）
    logx.Infof("User %s registered successfully", req.Phone)
    
    // Warn 日志（警告信息）
    logx.Warnf("User %s already exists", req.Phone)
    
    // Error 日志（错误信息）
    logx.Errorf("Register failed: %v", err)
    
    // 带上下文的日志（包含 trace_id 等信息）
    logx.WithContext(l.ctx).Infof("User %s registered", req.Phone)
    
    // 带字段的日志
    logx.WithContext(l.ctx).WithFields(logx.Field("phone", req.Phone), logx.Field("nickname", req.Nickname)).Info("User registered")
    
    return &types.RegisterResp{}, nil
}
```

#### 11.1.3 自定义日志格式

```go
// pkg/log/logger.go
package log

import (
    "github.com/zeromicro/go-zero/core/logx"
    "time"
)

// InitLogger 初始化日志
func InitLogger(serviceName string) {
    logx.MustSetup(logx.LogConf{
        ServiceName: serviceName,
        Mode:        "file",
        Path:        "logs",
        Level:       "info",
        Compress:    true,
        KeepDays:    7,
        Encoding:    "json",
    })
    
    // 设置日志格式
    logx.AddGlobalFields(
        logx.Field("service", serviceName),
        logx.Field("timestamp", time.Now().Format("2006-01-02 15:04:05")),
    )
}
```

### 11.2 链路追踪

#### 11.2.1 配置链路追踪

```yaml
# 链路追踪配置
Telemetry:
  Name: user-api
  Endpoint: http://127.0.0.1:14268/api/traces  # Jaeger 地址
  Sampler: 1.0                                  # 采样率（0.0-1.0）
  Batcher: jaeger                               # 批处理器类型
```

#### 11.2.2 启用链路追踪

```go
// apps/user/api/user.go
package main

import (
    "flag"
    "fmt"
    
    "imooc.com/easy-chat/apps/user/api/internal/config"
    "imooc.com/easy-chat/apps/user/api/internal/handler"
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
    flag.Parse()
    
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()
    
    ctx := svc.NewServiceContext(c)
    handler.RegisterHandlers(server, ctx)
    
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
```

#### 11.2.3 在代码中使用链路追踪

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 创建 Span
    ctx, span := otel.Tracer("user-api").Start(l.ctx, "Register")
    defer span.End()
    
    // 添加属性
    span.SetAttributes(
        attribute.String("phone", req.Phone),
        attribute.String("nickname", req.Nickname),
    )
    
    // 调用 RPC（自动传递 trace 信息）
    registerResp, err := l.svcCtx.User.Register(ctx, &user.RegisterReq{
        Phone:    req.Phone,
        Nickname: req.Nickname,
        Password: req.Password,
        Avatar:   req.Avatar,
        Sex:      int32(req.Sex),
    })
    if err != nil {
        // 记录错误
        span.RecordError(err)
        return nil, err
    }
    
    return &types.RegisterResp{
        Token:  registerResp.Token,
        Expire: registerResp.Expire,
    }, nil
}
```

---

## 12. 测试

### 12.1 单元测试示例

#### 12.1.1 测试 Logic

```go
// apps/user/rpc/internal/logic/registerlogic_test.go
package logic

import (
    "context"
    "testing"
    
    "imooc.com/easy-chat/apps/user/rpc/internal/config"
    "imooc.com/easy-chat/apps/user/rpc/internal/svc"
    "imooc.com/easy-chat/apps/user/rpc/user"
    
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestRegisterLogic_Register(t *testing.T) {
    // 1. 加载配置
    var c config.Config
    conf.MustLoad("../../etc/user.yaml", &c)
    
    // 2. 创建服务上下文
    ctx := svc.NewServiceContext(c)
    
    // 3. 创建 Logic
    l := NewRegisterLogic(context.Background(), ctx)
    
    // 4. 测试用例
    tests := []struct {
        name    string
        req     *user.RegisterReq
        wantErr bool
    }{
        {
            name: "正常注册",
            req: &user.RegisterReq{
                Phone:    "13800138000",
                Nickname: "测试用户",
                Password: "123456",
                Avatar:   "https://example.com/avatar.jpg",
                Sex:      1,
            },
            wantErr: false,
        },
        {
            name: "手机号已注册",
            req: &user.RegisterReq{
                Phone:    "13800138000",
                Nickname: "测试用户2",
                Password: "123456",
                Avatar:   "https://example.com/avatar.jpg",
                Sex:      1,
            },
            wantErr: true,
        },
    }
    
    // 5. 执行测试
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := l.Register(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, resp.Token)
                assert.Greater(t, resp.Expire, int64(0))
            }
        })
    }
}
```

#### 12.1.2 测试 Model

```go
// apps/user/models/usersmodel_test.go
package models

import (
    "context"
    "database/sql"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestUsersModel_Insert(t *testing.T) {
    // 1. 创建数据库连接
    sqlConn := sqlx.NewMysql("root:123456@tcp(127.0.0.1:3306)/easy_chat_test?charset=utf8mb4&parseTime=true")
    
    // 2. 创建 Model
    model := NewUsersModel(sqlConn, nil)
    
    // 3. 插入数据
    user := &Users{
        Id:       "test_user_001",
        Phone:    "13800138001",
        Nickname: "测试用户",
        Avatar:   "https://example.com/avatar.jpg",
        Password: sql.NullString{String: "123456", Valid: true},
        Sex:      sql.NullInt64{Int64: 1, Valid: true},
    }
    
    _, err := model.Insert(context.Background(), user)
    assert.NoError(t, err)
    
    // 4. 查询数据
    found, err := model.FindOne(context.Background(), user.Id)
    assert.NoError(t, err)
    assert.Equal(t, user.Phone, found.Phone)
    assert.Equal(t, user.Nickname, found.Nickname)
    
    // 5. 清理数据
    err = model.Delete(context.Background(), user.Id)
    assert.NoError(t, err)
}
```

### 12.2 集成测试示例

#### 12.2.1 测试 API

```go
// apps/user/api/test/user_test.go
package test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "imooc.com/easy-chat/apps/user/api/internal/config"
    "imooc.com/easy-chat/apps/user/api/internal/handler"
    "imooc.com/easy-chat/apps/user/api/internal/svc"
    
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
)

func TestRegisterAPI(t *testing.T) {
    // 1. 加载配置
    var c config.Config
    conf.MustLoad("../etc/user.yaml", &c)
    
    // 2. 创建服务器
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()
    
    ctx := svc.NewServiceContext(c)
    handler.RegisterHandlers(server, ctx)
    
    // 3. 创建测试请求
    reqBody := map[string]interface{}{
        "phone":    "13800138002",
        "nickname": "测试用户",
        "password": "123456",
        "avatar":   "https://example.com/avatar.jpg",
        "sex":      1,
    }
    reqJson, _ := json.Marshal(reqBody)
    
    req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(reqJson))
    req.Header.Set("Content-Type", "application/json")
    
    // 4. 执行请求
    w := httptest.NewRecorder()
    server.ServeHTTP(w, req)
    
    // 5. 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    
    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NotEmpty(t, resp["token"])
    assert.Greater(t, resp["expire"], float64(0))
}
```

---


## 13. 部署

### 13.1 编译打包

#### 13.1.1 编译单个服务

```bash
# 编译 API 服务
cd apps/user/api
go build -o user-api user.go

# 编译 RPC 服务
cd apps/user/rpc
go build -o user-rpc user.go
```

#### 13.1.2 交叉编译

```bash
# 编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o user-api-linux user.go

# 编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o user-api.exe user.go

# 编译 macOS 版本
GOOS=darwin GOARCH=amd64 go build -o user-api-mac user.go
```

#### 13.1.3 编译脚本

```bash
#!/bin/bash
# build.sh

# 设置变量
APP_NAME="user-api"
VERSION="1.0.0"
BUILD_TIME=$(date +%Y%m%d%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD)

# 编译参数
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# 编译
echo "Building ${APP_NAME}..."
go build -ldflags "${LDFLAGS}" -o ${APP_NAME} user.go

echo "Build completed: ${APP_NAME}"
echo "Version: ${VERSION}"
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"
```

### 13.2 Docker 部署

#### 13.2.1 编写 Dockerfile

```dockerfile
# apps/user/api/Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-api apps/user/api/user.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/user-api .

# 复制配置文件
COPY apps/user/api/etc/user.yaml etc/

# 暴露端口
EXPOSE 8888

# 启动命令
CMD ["./user-api", "-f", "etc/user.yaml"]
```

#### 13.2.2 构建 Docker 镜像

```bash
# 构建镜像
docker build -t user-api:1.0.0 -f apps/user/api/Dockerfile .

# 查看镜像
docker images | grep user-api
```

#### 13.2.3 运行 Docker 容器

```bash
# 运行容器
docker run -d \
  --name user-api \
  -p 8888:8888 \
  -v /path/to/config:/root/etc \
  -v /path/to/logs:/root/logs \
  user-api:1.0.0

# 查看日志
docker logs -f user-api

# 停止容器
docker stop user-api

# 删除容器
docker rm user-api
```

#### 13.2.4 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  # MySQL
  mysql:
    image: mysql:8.0
    container_name: easy-chat-mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: easy_chat
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - easy-chat-network

  # Redis
  redis:
    image: redis:7-alpine
    container_name: easy-chat-redis
    ports:
      - "6379:6379"
    networks:
      - easy-chat-network

  # Etcd
  etcd:
    image: bitnami/etcd:latest
    container_name: easy-chat-etcd
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    ports:
      - "2379:2379"
      - "2380:2380"
    networks:
      - easy-chat-network

  # MongoDB
  mongodb:
    image: mongo:6
    container_name: easy-chat-mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongodb-data:/data/db
    networks:
      - easy-chat-network

  # Kafka
  kafka:
    image: bitnami/kafka:latest
    container_name: easy-chat-kafka
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
    ports:
      - "9092:9092"
    networks:
      - easy-chat-network

  # User RPC
  user-rpc:
    build:
      context: .
      dockerfile: apps/user/rpc/Dockerfile
    container_name: user-rpc
    ports:
      - "10001:10001"
    depends_on:
      - mysql
      - redis
      - etcd
    networks:
      - easy-chat-network

  # User API
  user-api:
    build:
      context: .
      dockerfile: apps/user/api/Dockerfile
    container_name: user-api
    ports:
      - "8888:8888"
    depends_on:
      - user-rpc
      - redis
    networks:
      - easy-chat-network

volumes:
  mysql-data:
  mongodb-data:

networks:
  easy-chat-network:
    driver: bridge
```

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f user-api

# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### 13.3 配置管理

#### 13.3.1 使用环境变量

```yaml
# apps/user/api/etc/user.yaml
Name: user-api
Host: ${HOST:0.0.0.0}
Port: ${PORT:8888}

Auth:
  AccessSecret: ${JWT_SECRET}
  AccessExpire: ${JWT_EXPIRE:86400}

Redisx:
  Host: ${REDIS_HOST:127.0.0.1:6379}
  Pass: ${REDIS_PASS:}

UserRpc:
  Etcd:
    Hosts:
      - ${ETCD_HOST:127.0.0.1:2379}
    Key: user.rpc
```

```bash
# 设置环境变量
export HOST=0.0.0.0
export PORT=8888
export JWT_SECRET=your-jwt-secret-key
export JWT_EXPIRE=86400
export REDIS_HOST=127.0.0.1:6379
export REDIS_PASS=your-redis-password
export ETCD_HOST=127.0.0.1:2379

# 启动服务
./user-api -f etc/user.yaml
```

#### 13.3.2 使用配置中心

```go
// 使用 Nacos 配置中心
import (
    "github.com/nacos-group/nacos-sdk-go/clients"
    "github.com/nacos-group/nacos-sdk-go/common/constant"
    "github.com/nacos-group/nacos-sdk-go/vo"
)

func loadConfigFromNacos() (string, error) {
    // 创建 Nacos 客户端
    sc := []constant.ServerConfig{
        *constant.NewServerConfig("127.0.0.1", 8848),
    }
    
    cc := constant.ClientConfig{
        NamespaceId:         "public",
        TimeoutMs:           5000,
        NotLoadCacheAtStart: true,
        LogDir:              "/tmp/nacos/log",
        CacheDir:            "/tmp/nacos/cache",
        LogLevel:            "info",
    }
    
    client, err := clients.CreateConfigClient(map[string]interface{}{
        "serverConfigs": sc,
        "clientConfig":  cc,
    })
    if err != nil {
        return "", err
    }
    
    // 获取配置
    content, err := client.GetConfig(vo.ConfigParam{
        DataId: "user-api.yaml",
        Group:  "DEFAULT_GROUP",
    })
    if err != nil {
        return "", err
    }
    
    return content, nil
}
```

---

## 14. 最佳实践

### 14.1 项目结构

```
easy-chat1/
├── apps/                    # 应用目录
│   ├── user/               # 用户服务
│   │   ├── api/            # API 服务
│   │   │   ├── etc/        # 配置文件
│   │   │   ├── internal/   # 内部代码
│   │   │   │   ├── config/     # 配置
│   │   │   │   ├── handler/    # 处理器
│   │   │   │   ├── logic/      # 业务逻辑
│   │   │   │   ├── svc/        # 服务上下文
│   │   │   │   ├── types/      # 类型定义
│   │   │   │   └── middleware/ # 中间件
│   │   │   ├── user.api    # API 定义
│   │   │   └── user.go     # 主程序
│   │   ├── rpc/            # RPC 服务
│   │   │   ├── etc/
│   │   │   ├── internal/
│   │   │   │   ├── config/
│   │   │   │   ├── logic/
│   │   │   │   ├── server/
│   │   │   │   └── svc/
│   │   │   ├── user.proto  # Proto 定义
│   │   │   └── user.go
│   │   └── models/         # 数据模型
│   ├── im/                 # IM 服务
│   ├── social/             # 社交服务
│   └── task/               # 任务服务
├── pkg/                    # 公共包
│   ├── ctxdata/           # 上下文数据
│   ├── encrypt/           # 加密
│   ├── kafka/             # Kafka
│   ├── middleware/        # 中间件
│   ├── wuid/              # 分布式 ID
│   └── xerr/              # 错误处理
├── go.mod
├── go.sum
└── README.md
```

### 14.2 代码规范

#### 14.2.1 命名规范

```go
// 1. 包名：小写，简短，有意义
package user

// 2. 文件名：小写，下划线分隔
// registerlogic.go
// user_model.go

// 3. 常量：大写，下划线分隔
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 5000
)

// 4. 变量：驼峰命名
var (
    userId     string
    userName   string
    userAge    int
)

// 5. 函数：驼峰命名，首字母大写表示导出
func GetUserInfo() {}
func getUserById() {}

// 6. 结构体：驼峰命名，首字母大写表示导出
type UserInfo struct {
    Id       string
    Nickname string
}

// 7. 接口：驼峰命名，以 er 结尾
type UserService interface {
    Register() error
    Login() error
}
```

#### 14.2.2 注释规范

```go
// Package user 用户服务
// 提供用户注册、登录、信息查询等功能
package user

// RegisterLogic 注册业务逻辑
// 负责处理用户注册相关的业务逻辑
type RegisterLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

// NewRegisterLogic 创建注册逻辑实例
// 参数:
//   - ctx: 上下文
//   - svcCtx: 服务上下文
// 返回:
//   - *RegisterLogic: 注册逻辑实例
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
    return &RegisterLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// Register 用户注册
// 功能:
//   1. 验证手机号是否已注册
//   2. 加密密码
//   3. 插入数据库
//   4. 生成 JWT Token
// 参数:
//   - req: 注册请求
// 返回:
//   - resp: 注册响应（包含 Token）
//   - err: 错误信息
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
    // 实现代码...
}
```

#### 14.2.3 错误处理规范

```go
// 1. 定义错误
var (
    ErrUserNotFound     = errors.New("用户不存在")
    ErrPasswordError    = errors.New("密码错误")
    ErrPhoneIsRegister  = errors.New("手机号已注册")
)

// 2. 使用自定义错误
import "imooc.com/easy-chat/pkg/xerr"

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    user, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, req.Phone)
    if err != nil {
        if err == models.ErrNotFound {
            return nil, xerr.NewCodeError(xerr.UserNotFound, "用户不存在")
        }
        return nil, xerr.NewCodeError(xerr.DBError, "数据库错误")
    }
    
    if !encrypt.ValidatePasswordHash(req.Password, user.Password.String) {
        return nil, xerr.NewCodeError(xerr.PasswordError, "密码错误")
    }
    
    // ...
}

// 3. 错误包装
import "fmt"

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
    _, err := l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
    if err != nil {
        return nil, fmt.Errorf("insert user failed: %w", err)
    }
    
    // ...
}
```

### 14.3 性能优化

#### 14.3.1 数据库优化

```go
// 1. 使用索引
// 在数据库表中创建索引
CREATE INDEX idx_phone ON users(phone);
CREATE INDEX idx_nickname ON users(nickname);

// 2. 批量查询
func (l *GetUsersLogic) GetUsers(userIds []string) ([]*types.UserInfo, error) {
    // 使用 IN 查询，而不是循环查询
    users, err := l.svcCtx.UsersModel.FindByIds(l.ctx, userIds)
    if err != nil {
        return nil, err
    }
    
    // ...
}

// 3. 分页查询
func (l *GetUserListLogic) GetUserList(page, pageSize int) ([]*types.UserInfo, error) {
    offset := (page - 1) * pageSize
    users, err := l.svcCtx.UsersModel.FindList(l.ctx, offset, pageSize)
    if err != nil {
        return nil, err
    }
    
    // ...
}

// 4. 使用连接池
sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
sqlConn.SetMaxOpenConns(100)  // 最大打开连接数
sqlConn.SetMaxIdleConns(10)   // 最大空闲连接数
sqlConn.SetConnMaxLifetime(time.Hour)  // 连接最大生命周期
```

#### 14.3.2 缓存优化

```go
// 1. 使用缓存
func (l *GetUserInfoLogic) GetUserInfo(userId string) (*types.UserInfo, error) {
    // 先从缓存获取
    cacheKey := fmt.Sprintf("user:info:%s", userId)
    cached, err := l.svcCtx.Redis.Get(cacheKey)
    if err == nil && cached != "" {
        var userInfo types.UserInfo
        json.Unmarshal([]byte(cached), &userInfo)
        return &userInfo, nil
    }
    
    // 从数据库获取
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    userInfo := &types.UserInfo{
        Id:       user.Id,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
    }
    data, _ := json.Marshal(userInfo)
    l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
    
    return userInfo, nil
}

// 2. 批量缓存
func (l *GetUsersLogic) GetUsers(userIds []string) ([]*types.UserInfo, error) {
    var result []*types.UserInfo
    var missIds []string
    
    // 批量从缓存获取
    for _, userId := range userIds {
        cacheKey := fmt.Sprintf("user:info:%s", userId)
        cached, err := l.svcCtx.Redis.Get(cacheKey)
        if err == nil && cached != "" {
            var userInfo types.UserInfo
            json.Unmarshal([]byte(cached), &userInfo)
            result = append(result, &userInfo)
        } else {
            missIds = append(missIds, userId)
        }
    }
    
    // 批量从数据库获取未命中的数据
    if len(missIds) > 0 {
        users, err := l.svcCtx.UsersModel.FindByIds(l.ctx, missIds)
        if err != nil {
            return nil, err
        }
        
        // 写入缓存
        for _, user := range users {
            userInfo := &types.UserInfo{
                Id:       user.Id,
                Nickname: user.Nickname,
                Avatar:   user.Avatar,
            }
            result = append(result, userInfo)
            
            cacheKey := fmt.Sprintf("user:info:%s", user.Id)
            data, _ := json.Marshal(userInfo)
            l.svcCtx.Redis.Setex(cacheKey, string(data), 3600)
        }
    }
    
    return result, nil
}
```

#### 14.3.3 并发优化

```go
// 1. 使用 goroutine 并发处理
func (l *GetFriendListLogic) GetFriendList(userId string) ([]*types.Friend, error) {
    // 获取好友 ID 列表
    friendIds, err := l.svcCtx.FriendModel.GetFriendIds(l.ctx, userId)
    if err != nil {
        return nil, err
    }
    
    // 并发获取好友信息
    var wg sync.WaitGroup
    var mu sync.Mutex
    var friends []*types.Friend
    
    for _, friendId := range friendIds {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            
            user, err := l.svcCtx.UsersModel.FindOne(l.ctx, id)
            if err != nil {
                return
            }
            
            mu.Lock()
            friends = append(friends, &types.Friend{
                Id:       user.Id,
                Nickname: user.Nickname,
                Avatar:   user.Avatar,
            })
            mu.Unlock()
        }(friendId)
    }
    
    wg.Wait()
    return friends, nil
}

// 2. 使用 errgroup 控制并发
import "golang.org/x/sync/errgroup"

func (l *GetFriendListLogic) GetFriendList(userId string) ([]*types.Friend, error) {
    friendIds, err := l.svcCtx.FriendModel.GetFriendIds(l.ctx, userId)
    if err != nil {
        return nil, err
    }
    
    var mu sync.Mutex
    var friends []*types.Friend
    
    g, ctx := errgroup.WithContext(l.ctx)
    
    for _, friendId := range friendIds {
        id := friendId
        g.Go(func() error {
            user, err := l.svcCtx.UsersModel.FindOne(ctx, id)
            if err != nil {
                return err
            }
            
            mu.Lock()
            friends = append(friends, &types.Friend{
                Id:       user.Id,
                Nickname: user.Nickname,
                Avatar:   user.Avatar,
            })
            mu.Unlock()
            
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return friends, nil
}
```

---

## 15. 常见问题

### 15.1 常见错误及解决方案

#### 15.1.1 连接 Etcd 失败

```
错误信息: dial tcp 127.0.0.1:2379: connect: connection refused
```

**解决方案:**
1. 检查 Etcd 是否启动
2. 检查 Etcd 地址是否正确
3. 检查防火墙是否开放 2379 端口

```bash
# 启动 Etcd
docker run -d --name etcd -p 2379:2379 -p 2380:2380 -e ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd:latest

# 检查 Etcd 状态
etcdctl endpoint health
```

#### 15.1.2 RPC 调用超时

```
错误信息: rpc error: code = DeadlineExceeded desc = context deadline exceeded
```

**解决方案:**
1. 增加超时时间
2. 检查 RPC 服务是否正常
3. 检查网络是否正常

```yaml
# 增加超时时间
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
  Timeout: 10000  # 增加到 10 秒
```

#### 15.1.3 数据库连接失败

```
错误信息: Error 1045: Access denied for user 'root'@'localhost' (using password: YES)
```

**解决方案:**
1. 检查数据库用户名和密码
2. 检查数据库是否启动
3. 检查数据库权限

```bash
# 重置 MySQL 密码
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
```

#### 15.1.4 缓存未命中

```
错误信息: redis: nil
```

**解决方案:**
1. 检查 Redis 是否启动
2. 检查 Redis 地址是否正确
3. 检查缓存 Key 是否正确

```go
// 处理缓存未命中
cached, err := l.svcCtx.Redis.Get(cacheKey)
if err != nil || cached == "" {
    // 缓存未命中，从数据库获取
    // ...
}
```

### 15.2 性能问题排查

#### 15.2.1 慢查询排查

```bash
# 开启 MySQL 慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;  # 超过 1 秒的查询记录到慢查询日志

# 查看慢查询日志
tail -f /var/log/mysql/slow.log
```

#### 15.2.2 内存泄漏排查

```bash
# 使用 pprof 分析内存
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看内存分配
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap

# 查看 goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### 15.2.3 CPU 占用过高排查

```bash
# 使用 pprof 分析 CPU
go tool pprof http://localhost:6060/debug/pprof/profile

# 生成火焰图
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile
```

---


## 16. goctl 工具详解

### 16.1 goctl api 命令详解

#### 16.1.1 生成 API 代码

```bash
# 基本用法
goctl api go -api user.api -dir .

# 参数说明:
# -api      API 文件路径
# -dir      输出目录
# -style    命名风格: gozero, go_zero, GoZero
# -home     模板目录

# 指定命名风格
goctl api go -api user.api -dir . -style go_zero

# 使用自定义模板
goctl api go -api user.api -dir . -home ./template
```

#### 16.1.2 格式化 API 文件

```bash
# 格式化 API 文件
goctl api format -api user.api

# 格式化并输出到新文件
goctl api format -api user.api -o user_formatted.api

# 检查 API 文件语法
goctl api validate -api user.api
```

#### 16.1.3 生成 API 文档

```bash
# 生成 Markdown 文档
goctl api doc -api user.api -dir ./docs

# 生成 HTML 文档
goctl api doc -api user.api -dir ./docs -o html
```

#### 16.1.4 生成 API 插件

```bash
# 生成 TypeScript 类型定义
goctl api plugin -plugin goctl-ts -api user.api -dir ./types

# 生成 Dart 代码
goctl api plugin -plugin goctl-dart -api user.api -dir ./lib

# 生成 Java 代码
goctl api plugin -plugin goctl-java -api user.api -dir ./src
```

### 16.2 goctl rpc 命令详解

#### 16.2.1 生成 RPC 代码

```bash
# 基本用法
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 参数说明:
# --go_out        生成 protobuf 消息代码
# --go-grpc_out   生成 gRPC 服务代码
# --zrpc_out      生成 go-zero RPC 服务代码
# --style         命名风格
# --home          模板目录

# 指定命名风格
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=go_zero

# 使用自定义模板
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --home=./template
```

#### 16.2.2 生成 RPC 模板

```bash
# 生成 RPC 项目模板
goctl rpc template -o user.proto

# 生成的模板内容:
syntax = "proto3";

package user;
option go_package = "./user";

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

service User {
  rpc Ping(Request) returns(Response);
}
```

### 16.3 goctl model 命令详解

#### 16.3.1 从数据库生成 Model

```bash
# 从 MySQL 数据库生成 Model
goctl model mysql datasource \
  -url="root:123456@tcp(127.0.0.1:3306)/easy_chat" \
  -table="users" \
  -dir="./models" \
  -cache=true \
  -style=go_zero

# 参数说明:
# -url          数据库连接地址
# -table        表名（支持通配符，如 "user_*"）
# -dir          输出目录
# -cache        是否生成缓存代码
# -style        命名风格
# -home         模板目录
# -idea         是否生成 IDEA 配置

# 生成多个表
goctl model mysql datasource \
  -url="root:123456@tcp(127.0.0.1:3306)/easy_chat" \
  -table="users,friends,groups" \
  -dir="./models" \
  -cache=true

# 使用通配符
goctl model mysql datasource \
  -url="root:123456@tcp(127.0.0.1:3306)/easy_chat" \
  -table="user_*" \
  -dir="./models" \
  -cache=true
```

#### 16.3.2 从 DDL 文件生成 Model

```bash
# 从 DDL 文件生成 Model
goctl model mysql ddl \
  -src="./users.sql" \
  -dir="./models" \
  -cache=true \
  -style=go_zero

# users.sql 内容:
CREATE TABLE `users` (
  `id` varchar(24) NOT NULL COMMENT '用户ID',
  `phone` varchar(11) NOT NULL COMMENT '手机号',
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `password` varchar(255) DEFAULT NULL COMMENT '密码',
  `sex` tinyint(1) DEFAULT 0 COMMENT '性别',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_nickname` (`nickname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### 16.3.3 生成 MongoDB Model

```bash
# 生成 MongoDB Model
goctl model mongo \
  -type User \
  -dir ./models \
  -cache=true

# 参数说明:
# -type         类型名称
# -dir          输出目录
# -cache        是否生成缓存代码
# -style        命名风格
```

### 16.4 其他 goctl 命令

#### 16.4.1 生成 Dockerfile

```bash
# 生成 Dockerfile
goctl docker -go user.go

# 生成的 Dockerfile:
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/user user.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/user /app/user

CMD ["./user"]
```

#### 16.4.2 生成 Kubernetes 配置

```bash
# 生成 Kubernetes 配置
goctl kube deploy -name user-api -namespace default -image user-api:1.0.0 -o user-api.yaml

# 参数说明:
# -name         服务名称
# -namespace    命名空间
# -image        镜像名称
# -replicas     副本数
# -port         端口
# -nodePort     NodePort
# -requestCpu   CPU 请求
# -requestMem   内存请求
# -limitCpu     CPU 限制
# -limitMem     内存限制
# -o            输出文件
```

#### 16.4.3 升级 goctl

```bash
# 升级 goctl 到最新版本
goctl upgrade

# 查看 goctl 版本
goctl --version

# 查看帮助
goctl --help
```

#### 16.4.4 模板管理

```bash
# 初始化模板
goctl template init

# 清理模板
goctl template clean

# 更新模板
goctl template update

# 恢复默认模板
goctl template revert
```

---

## 17. 完整项目示例

### 17.1 基于 easy-chat1 项目的实际代码示例

#### 17.1.1 项目结构

```
easy-chat1/
├── apps/
│   ├── user/           # 用户服务
│   │   ├── api/        # 用户 API 服务
│   │   ├── rpc/        # 用户 RPC 服务
│   │   └── models/     # 用户数据模型
│   ├── im/             # IM 服务
│   │   ├── api/        # IM API 服务
│   │   ├── rpc/        # IM RPC 服务
│   │   ├── ws/         # WebSocket 服务
│   │   └── models/     # IM 数据模型
│   ├── social/         # 社交服务
│   │   ├── api/        # 社交 API 服务
│   │   ├── rpc/        # 社交 RPC 服务
│   │   └── models/     # 社交数据模型
│   └── task/           # 任务服务
│       └── mq/         # 消息队列消费者
├── pkg/                # 公共包
│   ├── ctxdata/       # 上下文数据
│   ├── encrypt/       # 加密
│   ├── kafka/         # Kafka
│   ├── middleware/    # 中间件
│   ├── wuid/          # 分布式 ID
│   └── xerr/          # 错误处理
├── go.mod
└── go.sum
```

### 17.2 用户服务完整示例

#### 17.2.1 API 定义

```api
// apps/user/api/user.api
syntax = "v1"

info(
    title: "用户服务 API"
    desc: "提供用户注册、登录、信息查询等功能"
    author: "easy-chat"
    version: "1.0.0"
)

// 注册请求
type RegisterReq {
    Phone    string `json:"phone"`     // 手机号
    Nickname string `json:"nickname"`  // 昵称
    Password string `json:"password"`  // 密码
    Avatar   string `json:"avatar"`    // 头像
    Sex      int    `json:"sex"`       // 性别：0-未知，1-男，2-女
}

// 注册响应
type RegisterResp {
    Token  string `json:"token"`   // JWT Token
    Expire int64  `json:"expire"`  // 过期时间戳
}

// 登录请求
type LoginReq {
    Phone    string `json:"phone"`     // 手机号
    Password string `json:"password"`  // 密码
}

// 登录响应
type LoginResp {
    Token  string `json:"token"`   // JWT Token
    Expire int64  `json:"expire"`  // 过期时间戳
}

// 获取用户信息请求
type DetailReq {
    UserId string `path:"userId"`  // 用户 ID
}

// 用户信息响应
type DetailResp {
    Id       string `json:"id"`       // 用户 ID
    Phone    string `json:"phone"`    // 手机号
    Nickname string `json:"nickname"` // 昵称
    Avatar   string `json:"avatar"`   // 头像
    Sex      int    `json:"sex"`      // 性别
}

// 不需要认证的接口
@server(
    prefix: /api/user
    group: user
)
service user-api {
    @doc "用户注册"
    @handler register
    post /register (RegisterReq) returns (RegisterResp)
    
    @doc "用户登录"
    @handler login
    post /login (LoginReq) returns (LoginResp)
}

// 需要 JWT 认证的接口
@server(
    prefix: /api/user
    group: user
    jwt: Auth
)
service user-api {
    @doc "获取用户信息"
    @handler detail
    get /detail/:userId (DetailReq) returns (DetailResp)
}
```

#### 17.2.2 Proto 定义

```protobuf
// apps/user/rpc/user.proto
syntax = "proto3";

package user;
option go_package = "./user";

// 用户服务
service User {
    // 注册
    rpc Register(RegisterReq) returns (RegisterResp);
    // 登录
    rpc Login(LoginReq) returns (LoginResp);
    // 获取用户信息
    rpc GetUserInfo(GetUserInfoReq) returns (GetUserInfoResp);
    // 查找用户
    rpc FindUser(FindUserReq) returns (FindUserResp);
}

// 注册请求
message RegisterReq {
    string phone = 1;
    string nickname = 2;
    string password = 3;
    string avatar = 4;
    int32 sex = 5;
}

// 注册响应
message RegisterResp {
    string token = 1;
    int64 expire = 2;
}

// 登录请求
message LoginReq {
    string phone = 1;
    string password = 2;
}

// 登录响应
message LoginResp {
    string token = 1;
    int64 expire = 2;
}

// 获取用户信息请求
message GetUserInfoReq {
    string id = 1;
}

// 获取用户信息响应
message GetUserInfoResp {
    string id = 1;
    string phone = 2;
    string nickname = 3;
    string avatar = 4;
    int32 sex = 5;
}

// 查找用户请求
message FindUserReq {
    string phone = 1;
    string nickname = 2;
    repeated string ids = 3;
}

// 查找用户响应
message FindUserResp {
    repeated UserInfo users = 1;
}

// 用户信息
message UserInfo {
    string id = 1;
    string phone = 2;
    string nickname = 3;
    string avatar = 4;
    int32 sex = 5;
}
```

### 17.3 IM 服务完整示例

#### 17.3.1 WebSocket 连接管理

```go
// apps/im/ws/internal/handler/user/user.go
package user

import (
    "context"
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/zeromicro/go-zero/core/logx"
    "imooc.com/easy-chat/apps/im/ws/internal/svc"
    "net/http"
    "sync"
)

// UserHandler WebSocket 用户处理器
type UserHandler struct {
    svcCtx *svc.ServiceContext
    
    // 用户连接管理
    connections map[string]*websocket.Conn  // userId -> conn
    mu          sync.RWMutex
}

// NewUserHandler 创建用户处理器
func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
    return &UserHandler{
        svcCtx:      svcCtx,
        connections: make(map[string]*websocket.Conn),
    }
}

// HandleWebSocket 处理 WebSocket 连接
func (h *UserHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // 1. 升级为 WebSocket 连接
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true  // 允许所有来源
        },
    }
    
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logx.Errorf("Upgrade to websocket failed: %v", err)
        return
    }
    defer conn.Close()
    
    // 2. 获取用户 ID（从 Token 中解析）
    token := r.URL.Query().Get("token")
    userId, err := h.parseToken(token)
    if err != nil {
        logx.Errorf("Parse token failed: %v", err)
        return
    }
    
    // 3. 保存连接
    h.addConnection(userId, conn)
    defer h.removeConnection(userId)
    
    // 4. 设置用户在线状态
    h.setUserOnline(userId)
    defer h.setUserOffline(userId)
    
    // 5. 处理消息
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            logx.Errorf("Read message failed: %v", err)
            break
        }
        
        // 处理消息
        h.handleMessage(userId, message)
    }
}

// addConnection 添加连接
func (h *UserHandler) addConnection(userId string, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.connections[userId] = conn
}

// removeConnection 移除连接
func (h *UserHandler) removeConnection(userId string) {
    h.mu.Lock()
    defer h.mu.Unlock()
    delete(h.connections, userId)
}

// getConnection 获取连接
func (h *UserHandler) getConnection(userId string) (*websocket.Conn, bool) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    conn, ok := h.connections[userId]
    return conn, ok
}

// SendMessage 发送消息给用户
func (h *UserHandler) SendMessage(userId string, message interface{}) error {
    conn, ok := h.getConnection(userId)
    if !ok {
        return errors.New("user not online")
    }
    
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    return conn.WriteMessage(websocket.TextMessage, data)
}

// handleMessage 处理消息
func (h *UserHandler) handleMessage(userId string, message []byte) {
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        logx.Errorf("Unmarshal message failed: %v", err)
        return
    }
    
    msgType := msg["type"].(string)
    
    switch msgType {
    case "chat":
        // 处理聊天消息
        h.handleChatMessage(userId, msg)
    case "read":
        // 处理已读消息
        h.handleReadMessage(userId, msg)
    case "heartbeat":
        // 处理心跳消息
        h.handleHeartbeat(userId)
    default:
        logx.Warnf("Unknown message type: %s", msgType)
    }
}

// setUserOnline 设置用户在线
func (h *UserHandler) setUserOnline(userId string) {
    h.svcCtx.Redis.Sadd("users:online", userId)
    h.svcCtx.Redis.Hset("users:last:active", userId, fmt.Sprintf("%d", time.Now().Unix()))
}

// setUserOffline 设置用户离线
func (h *UserHandler) setUserOffline(userId string) {
    h.svcCtx.Redis.Srem("users:online", userId)
}
```

### 17.4 Social 服务完整示例

#### 17.4.1 好友申请流程

```go
// apps/social/api/internal/logic/friend/friendputinlogic.go
package friend

import (
    "context"
    "imooc.com/easy-chat/apps/social/rpc/socialclient"
    "imooc.com/easy-chat/pkg/ctxdata"
    
    "imooc.com/easy-chat/apps/social/api/internal/svc"
    "imooc.com/easy-chat/apps/social/api/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type FriendPutinLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewFriendPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutinLogic {
    return &FriendPutinLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// FriendPutin 好友申请
func (l *FriendPutinLogic) FriendPutin(req *types.FriendPutinReq) (resp *types.FriendPutinResp, err error) {
    // 1. 获取当前用户 ID
    userId := ctxdata.GetUId(l.ctx)
    
    // 2. 验证是否已经是好友
    isFriendResp, err := l.svcCtx.Social.IsFriend(l.ctx, &socialclient.IsFriendReq{
        UserId:   userId,
        FriendId: req.FriendId,
    })
    if err != nil {
        return nil, err
    }
    
    if isFriendResp.IsFriend {
        return nil, errors.New("已经是好友关系")
    }
    
    // 3. 验证是否已经申请过
    existsResp, err := l.svcCtx.Social.FriendPutinExists(l.ctx, &socialclient.FriendPutinExistsReq{
        UserId:   userId,
        FriendId: req.FriendId,
    })
    if err != nil {
        return nil, err
    }
    
    if existsResp.Exists {
        return nil, errors.New("已经申请过，请等待对方处理")
    }
    
    // 4. 创建好友申请
    _, err = l.svcCtx.Social.FriendPutin(l.ctx, &socialclient.FriendPutinReq{
        UserId:   userId,
        FriendId: req.FriendId,
        ReqMsg:   req.ReqMsg,
    })
    if err != nil {
        return nil, err
    }
    
    // 5. 发送通知给对方（通过 WebSocket）
    // TODO: 实现通知逻辑
    
    return &types.FriendPutinResp{}, nil
}
```

---

## 总结

本文档详细介绍了 Go-Zero 框架的完整使用方法，包括：

1. **框架介绍**: Go-Zero 的特点、架构、核心组件
2. **快速开始**: 环境搭建、创建项目、运行服务
3. **API 服务开发**: API 语法、Handler、Logic、ServiceContext、错误处理、参数验证
4. **RPC 服务开发**: Proto 语法、代码生成、Server 实现、Client 调用、服务发现
5. **配置管理**: 配置文件结构、数据库配置、RPC 配置、JWT 配置、日志配置
6. **数据库操作**: MySQL、MongoDB、Redis 集成及使用
7. **中间件开发**: JWT 认证、日志、限流、幂等性、跨域等中间件
8. **服务间通信**: RPC 调用、服务发现、负载均衡、超时控制、重试机制
9. **缓存使用**: 缓存策略、Redis 使用、缓存穿透/击穿/雪崩解决方案
10. **消息队列**: Kafka 集成、生产者、消费者
11. **日志和监控**: 日志配置、链路追踪
12. **测试**: 单元测试、集成测试
13. **部署**: 编译打包、Docker 部署、配置管理
14. **最佳实践**: 项目结构、代码规范、性能优化
15. **常见问题**: 错误排查、性能问题排查
16. **goctl 工具**: API、RPC、Model 代码生成
17. **完整项目示例**: 基于 easy-chat1 项目的实际代码

通过本文档，你可以全面掌握 Go-Zero 框架的使用，快速开发高性能的微服务应用。

---

**文档版本**: 1.0.0  
**最后更新**: 2024-01-01  
**作者**: easy-chat 团队

