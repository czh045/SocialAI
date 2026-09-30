# Social AI：从零到就业的深度讲解

这份文档解释的是当前 `Z:\Program\java\Social` 下真实存在的代码。你已经有一个 Spring Boot + PostgreSQL 的 OnlineOrder 项目；Social AI 的价值在于让你再掌握一条不同的完整全栈技术链：

```text
React
  -> Axios / HTTP
Go + Gorilla Mux
  -> JWT Middleware
Service
  -> Elasticsearch + Local Storage / Google Cloud Storage
```

它实现的是一个图片、视频和 AI 生成图片的社交媒体原型：

1. 注册账号。
2. 登录取得 JWT。
3. 上传图片或视频并附带文字。
4. 以全部内容、关键词或用户搜索帖子。
5. 浏览图片与视频。
6. 删除自己的帖子。
7. 输入 prompt 调用 OpenAI 图像生成，再上传生成结果。

本项目能帮助你理解：React 状态、路由、文件上传、Go HTTP 服务、JWT、Elasticsearch、对象存储、云部署的基本串联方式。

先说一个就业角度的重要事实：

> 这个项目作为课程全栈项目已经有完整可演示链路；它不是可以直接拿去面对真实用户的社交平台。懂得指出安全、数据和部署差距，才是工程师的思考方式。

---

## 1. 目录和职责：先把地图记住

```text
Social/
├─ docker-compose.yml                 本机 Elasticsearch
├─ socialai/                          Go 后端
│  ├─ main.go                         后端入口
│  ├─ app.yaml                        Google App Engine 配置
│  ├─ conf/deploy.yml                 本地配置
│  ├─ backend/                        Elasticsearch、文件存储底层
│  ├─ constants/                      索引名称常量
│  ├─ handler/                        HTTP 路由、JWT、接口处理
│  ├─ model/                          Post、User 数据结构
│  ├─ service/                        搜索、保存、删除、用户业务
│  └─ util/config.go                  YAML 配置读取
└─ social-ai/                         React 前端
   ├─ package.json                    npm 依赖和命令
   ├─ .env                            本机 API 地址和可选 OpenAI key
   └─ src/
      ├─ index.js                     React 挂载入口
      ├─ constants.js                 token key、API 地址、搜索枚举
      ├─ components/                  页面和可复用 UI
      └─ styles/                      CSS
```

请先记住分层，而不是死记文件名：

| 层 | 作用 | Social AI 中的例子 |
| --- | --- | --- |
| Frontend | 展示 UI、收集输入、发请求 | `Login.js`、`Collection.js` |
| Router / Handler | 匹配 HTTP URL，校验请求，回响应 | `handler/router.go`、`handler/post.go` |
| Middleware | 在 Handler 前处理通用问题 | `handler/auth.go` 的 JWT |
| Service | 组合业务规则 | `service/post.go`、`service/user.go` |
| Backend Adapter | 封装外部系统 | `backend/elasticsearch.go`、`backend/storage.go` |
| Data store | 真实保存数据 | Elasticsearch、uploads 文件夹、GCS |

---

## 2. 正确的本地启动顺序

Social AI 需要三个进程/服务：

```text
Elasticsearch：数据搜索和帖子/用户文档
Go backend：API 和本地静态媒体
React frontend：浏览器页面
```

### 2.1 启动 Elasticsearch

在 `Z:\Program\java\Social`：

```powershell
docker compose up -d
docker compose ps
```

`docker-compose.yml` 中：

```yaml
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.17.25
    container_name: socialai-elasticsearch
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
    ports:
      - "9200:9200"
```

逐行理解：

- `image`：下载并运行 Elasticsearch 7.17.25。
- `container_name`：容器在 Docker Desktop 中显示为 `socialai-elasticsearch`。
- `discovery.type=single-node`：本机只跑一个节点，不加入集群。
- `xpack.security.enabled=false`：课程本机环境关闭 ES 登录认证，方便调试；云端不能照搬。
- `ES_JAVA_OPTS`：限制 Elasticsearch JVM 使用 512 MB 堆内存。
- `9200:9200`：电脑浏览器和 Go 后端通过 `localhost:9200` 访问 ES。
- named volume `socialai-es-data`：容器删掉后，索引数据仍可以保留。

验证：

```powershell
Invoke-RestMethod http://localhost:9200
```

能返回 Elasticsearch 的版本和集群信息，就说明数据库式的搜索服务已准备好。

### 2.2 启动 Go 后端

```powershell
cd Z:\Program\java\Social\socialai
go run .
```

默认监听 `8080`。如果 OnlineOrder 也占用 8080，请为 Social 使用 8081：

```powershell
$env:PORT = "8081"
$env:SOCIALAI_PUBLIC_BASE_URL = "http://localhost:8081/uploads"
go run .
```

这里有一个非常容易混淆的点：

```text
Go 后端的监听端口
必须和 React 请求的 API_BASE_URL 对上
本地媒体生成的 URL 也必须和 Go 后端端口对上
```

如果使用 8081，则 React `.env` 要是：

```env
REACT_APP_API_BASE_URL=http://localhost:8081
```

再重新启动 React。

### 2.3 启动 React 前端

```powershell
cd Z:\Program\java\Social\social-ai
npm.cmd install
npm.cmd start
```

浏览器打开：

```text
http://localhost:3000
```

构建生产静态文件：

```powershell
npm.cmd run build
```

---

## 3. 从一次“登录”开始理解整条链

用户在 Login 页面输入：

```text
username: alice
password: abc123
```

调用链：

```text
React Login.js
  -> axios POST http://localhost:8080/signin
  -> Go router 匹配 /signin
  -> signinHandler 解 JSON
  -> service.CheckUser
  -> Elasticsearch 查 user index
  -> 成功后 handler 签 JWT
  -> 返回纯文本 token
  -> React App.js 保存到 localStorage
  -> isLoggedIn = true
  -> Main.js 允许访问 /create 和 /collection
```

这条链要能讲清楚。它已经覆盖了全栈实习常见考点：

- 前端表单。
- POST JSON。
- REST API。
- 后端路由。
- 业务 Service。
- 外部数据存储。
- 登录认证。
- 浏览器本地存储。
- 前端路由守卫。

---

## 4. Go 后端从入口到路由

### 4.1 `main.go`：Go 程序的入口

核心逻辑：

```go
func main() {
    configPath := os.Getenv("SOCIALAI_CONFIG")
    if configPath == "" {
        configPath = "conf/deploy.yml"
    }

    config, err := util.LoadApplicationConfig(configPath)
    backend.InitElasticsearchBackend(config.ElasticsearchConfig)
    backend.InitMediaBackend(config)

    address := os.Getenv("PORT")
    if address == "" {
        address = "8080"
    }

    log.Fatal(http.ListenAndServe(":"+address, handler.InitRouter(config.TokenConfig)))
}
```

逐段理解：

```go
func main()
```

任何 Go 可执行程序的入口。和 Java 的 `public static void main` 是同一个地位。

```go
os.Getenv("SOCIALAI_CONFIG")
```

读取环境变量。部署时不同环境可以指定不同 YAML 文件，不需要改代码。

```go
util.LoadApplicationConfig(configPath)
```

读取 YAML。配置与代码分离是后端工程基本习惯。

```go
backend.InitElasticsearchBackend(...)
backend.InitMediaBackend(...)
```

后端启动前先确认两个依赖可用：

1. Elasticsearch 可以连接并已存在索引。
2. 文件保存位置已准备好：本地文件夹或 GCS。

```go
http.ListenAndServe(":"+address, ...)
```

Go 标准库启动 HTTP server。它会持续运行，等待请求；终端不返回提示符不是卡死。

### 4.2 `util/config.go`：YAML 如何变成 Go 对象

例如：

```go
type ElasticsearchInfo struct {
    Address  string `yaml:"address"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
```

这是 Go 的 struct：

- `Address string`：字段名和类型。
- `` `yaml:"address"` ``：tag，告诉 YAML 解析器把 `address:` 填入这个字段。

`ApplicationConfig` 聚合 Elasticsearch、storage、GCS、token 配置。

```go
content, err := os.ReadFile(configPath)
```

读文件。

```go
yaml.Unmarshal(content, &config)
```

把 YAML 字节内容解析并填入 `config`。`&config` 表示传地址，让函数能够修改这个变量。

### 4.3 `conf/deploy.yml`：本地环境配置

```yaml
elasticsearch:
  address: "http://localhost:9200"

storage:
  driver: "local"
  directory: "uploads"
  public_base_url: "http://localhost:8080/uploads"

token:
  secret: "replace-this-local-development-secret"
```

本机模式：

```text
帖子元数据 -> Elasticsearch
图片/视频实际字节 -> socialai/uploads/
浏览器访问媒体 -> Go 的 /uploads/... 静态路由
```

`deploy.gcp.example.yml` 是云端模板：

```text
帖子元数据 -> Elasticsearch on GCE
媒体 -> Google Cloud Storage
后端 -> Google App Engine
```

真实 secret、bucket 名、数据库密码不得提交到 Git。

---

## 5. Elasticsearch：这个项目为什么不像 OnlineOrder 一样使用 PostgreSQL

OnlineOrder 的主要问题是事务性关系数据：用户、购物车、订单项、外键。PostgreSQL 很适合。

Social AI 的主要需求包括：

- 按帖子文字关键词搜索。
- 按用户筛选帖子。
- 存帖子文档。

因此课程选择 Elasticsearch：

```text
user index：用户文档
post index：帖子文档
```

### 5.1 `constants/constants.go`

```go
const (
    UserIndex = "user"
    PostIndex = "post"
)
```

不要把 `"user"`、`"post"` 字符串复制到十多个文件里。常量可避免拼写错误，也方便将来改名。

### 5.2 `model/model.go`

```go
type Post struct {
    ID      string `json:"id"`
    User    string `json:"user"`
    Message string `json:"message"`
    URL     string `json:"url"`
    Type    string `json:"type"`
}
```

这是一篇帖子在 Go 和 Elasticsearch 中的文档结构：

- `ID`：UUID。
- `User`：发帖用户名。
- `Message`：文字说明，允许全文搜索。
- `URL`：媒体文件可访问地址。
- `Type`：`image` 或 `video`。

```go
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}
```

课程 User 模型提供了账户基础字段。注意：当前课程项目的 password 仍是明文，这是后文必须改的生产问题。

### 5.3 `backend/elasticsearch.go`

```go
options := []elastic.ClientOptionFunc{
    elastic.SetURL(config.Address),
    elastic.SetSniff(false),
    elastic.SetHealthcheck(false),
}
```

- `SetURL`：ES 地址。
- `SetSniff(false)`：不去发现 ES 集群其他节点。本机单节点若开 sniff，客户端可能连到不该连的地址。
- `SetHealthcheck(false)`：本地课程环境减少健康检查干扰。生产一般应谨慎评估，不能为了省事关闭真正需要的健康检测。

```go
client, err := elastic.NewClient(options...)
```

创建 ES 客户端。`...` 把 slice 展开为可变参数。

`createIndexIfMissing` 的逻辑：

```text
先 IndexExists
  -> 存在：直接返回
  -> 不存在：CreateIndex + mapping
```

Post 索引 mapping：

```json
{
  "id": {"type": "keyword"},
  "user": {"type": "keyword"},
  "message": {"type": "text"},
  "url": {"type": "keyword", "index": false},
  "type": {"type": "keyword", "index": false}
}
```

为什么 `message` 是 `text` 而 `user` 是 `keyword`？

- `text`：会分析、分词，适合关键词全文搜索。
- `keyword`：整体精确匹配，适合用户名、ID。
- `index: false`：URL、type 不作为搜索条件，用于返回即可。

`SaveToES`：

```go
backend.client.Index().
    Index(index).
    Id(id).
    BodyJson(document).
    Refresh("wait_for").
    Do(context.Background())
```

逐步意思：

1. 选择 index。
2. 用 ID 作为文档 ID。
3. 把 Go struct 转成 JSON。
4. `Refresh("wait_for")`：等待搜索可见，课程演示中上传后立即搜索就能看到。
5. `Do`：真正发网络请求。

注意 Go SDK 方法叫 `Id`，不是 `ID`。这是当前 ES v7 SDK 的实际 API。

`ReadFromES` 中：

```go
Size(100)
```

最多取 100 条。这是课程演示限制。真实产品必须有分页、cursor/search_after、排序和最大结果保护。

---

## 6. 媒体存储：帖子信息和图片字节不是一回事

这是全栈工程中非常关键的概念：

```text
Elasticsearch 存：
  帖子 ID、用户名、文字、媒体 URL、类型

本地磁盘或 GCS 存：
  PNG/JPG/MP4 的真实二进制文件
```

不能把大量视频字节直接塞进 Elasticsearch 文档。

### 6.1 `backend/storage.go` 的接口

```go
type MediaStorage interface {
    Save(reader io.Reader, objectName string) (string, error)
}
```

这是接口。含义是“不管存本地还是 GCS，只要能接收一个读取流和对象名，然后返回访问 URL，就符合要求”。

这正是依赖抽象的价值：

```text
handler/service 不关心存储实现
  -> MediaBackend.Save(...)
本地使用 localStorageBackend
云端使用 googleCloudStorageBackend
```

### 6.2 本地存储

```go
func (backend *localStorageBackend) Save(reader io.Reader, objectName string) (string, error) {
    path := filepath.Join(backend.directory, objectName)
    file, err := os.Create(path)
    defer file.Close()
    io.Copy(file, reader)
    return "/uploads/" + url.PathEscape(objectName), nil
}
```

- `io.Reader`：Go 用于“持续读取数据”的抽象，上传文件就是一个 reader。
- `filepath.Join`：正确拼文件路径，避免手写 `/` 或 `\`。
- `defer file.Close()`：函数返回前自动关闭文件，避免资源泄露。
- `io.Copy`：从上传流复制到磁盘文件。
- `url.PathEscape`：把文件名编码进 URL，避免特殊字符造成 URL 问题。

路由中：

```go
router.PathPrefix("/uploads/").
    Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
```

把请求 `/uploads/a.png` 映射到磁盘 `./uploads/a.png`。这是开发方便的静态文件服务。

### 6.3 Google Cloud Storage

`googleCloudStorageBackend` 使用：

```go
storage.NewClient(context.Background())
```

创建 GCS 客户端，云端通过 Google Application Default Credentials 认证。

```go
object := backend.client.Bucket(backend.bucket).Object(objectName)
writer := object.NewWriter(ctx)
io.Copy(writer, reader)
writer.Close()
```

表示把上传数据流写入 bucket 中一个对象。

课程版还给对象设置了 `storage.AllUsers` 读权限，所以浏览器能显示图片。它要求 GCS bucket 没有开启 uniform bucket-level access；真实项目还需考虑是否应公开、是否应改为签名 URL。

---

## 7. Service：搜索、保存、删除的业务动作

### 7.1 `service/post.go`

按用户查帖子：

```go
query := elastic.NewTermQuery("user", user)
```

用精确匹配查 `user` 字段，因为 user mapping 是 keyword。

按关键词查帖子：

```go
query := elastic.NewMatchQuery("message", keywords).Operator("AND")
if keywords == "" {
    query.ZeroTermsQuery("all")
}
```

- `MatchQuery`：全文搜索 message。
- `Operator("AND")`：多个词都需要满足，使结果更精确。
- 搜索词为空时 `ZeroTermsQuery("all")`：查全部帖子。

保存帖子：

```go
func SavePost(post *model.Post, file multipart.File, extension string) error {
    mediaURL, err := backend.MediaBackend.Save(file, post.ID+extension)
    post.URL = mediaURL
    return backend.ESBackend.SaveToES(post, constants.PostIndex, post.ID)
}
```

顺序是：

```text
先把文件存储成功
  -> 得到 URL
  -> URL 填入 Post
  -> 再把 Post 元数据存进 ES
```

若文件保存失败，函数会直接返回错误，ES 不会写一个指向不存在 URL 的帖子。这是一点很好的基本一致性意识。

删除帖子：

```go
query := elastic.NewBoolQuery().
    Must(elastic.NewTermQuery("id", id)).
    Must(elastic.NewTermQuery("user", user))
```

同时要求：

```text
帖子 ID 相等
AND
帖子所有者等于当前登录用户
```

所以即使恶意用户猜到别人的 post ID，也不能按课程后端删除别人的 ES 文档。

但注意：当前删除只删 ES 元数据，不删本地上传文件或 GCS 对象，会产生“孤儿文件”。这正是后续工程改进项。

### 7.2 `service/user.go`

注册：

```text
先按 username 查询 ES
  -> 已存在：返回 false
  -> 不存在：保存 User 文档
```

登录：

```go
elastic.NewBoolQuery().
    Must(elastic.NewTermQuery("username", username)).
    Must(elastic.NewTermQuery("password", password))
```

课程逻辑可用来理解搜索查询，但真实项目绝不能用“用户名 + 明文密码精确匹配 ES”做认证。原因后面详谈。

---

## 8. Handler、路由和 JWT：请求如何被保护

### 8.1 `handler/router.go`：所有路由的总表

公开接口：

```go
router.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)
router.HandleFunc("/signup", signupHandler).Methods(http.MethodPost)
router.HandleFunc("/signin", signinHandler).Methods(http.MethodPost)
```

需要登录的接口：

```go
protected := router.NewRoute().Subrouter()
protected.Use(func(next http.Handler) http.Handler {
    return withJWTAuth(next, signingKey)
})
protected.HandleFunc("/upload", uploadHandler).Methods(http.MethodPost)
protected.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
protected.HandleFunc("/post/{id}", deleteHandler).Methods(http.MethodDelete)
```

`protected.Use(...)` 的意思是：每一个 protected 路由真正进入 handler 前，先经过 `withJWTAuth`。

这叫 middleware：

```text
请求
  -> JWT Middleware
     -> token 合法：进入业务 handler
     -> token 缺失/错误/过期：401，业务 handler 不会执行
```

### 8.2 `handler/auth.go`：JWT 的验证

浏览器请求 header：

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

代码：

```go
header := r.Header.Get("Authorization")
if !strings.HasPrefix(header, "Bearer ") {
    http.Error(w, "Missing bearer token", http.StatusUnauthorized)
    return
}
```

先确认请求确实带 Bearer token。

```go
claims := &jwt.RegisteredClaims{}
token, err := jwt.ParseWithClaims(tokenString, claims, ...)
```

解析 token 中的标准 claims，包括 Subject 和过期时间。

```go
if !token.Valid || claims.Subject == "" {
    http.Error(..., http.StatusUnauthorized)
}
```

token 签名错误、过期或没有 Subject，都不能通过。

```go
ctx := context.WithValue(r.Context(), usernameContextKey, claims.Subject)
next.ServeHTTP(w, r.WithContext(ctx))
```

验证成功后，将 token 中的用户名放进 request context。后面的 upload/delete handler 不信任客户端随便传的 `username`，而从经过签名验证的 context 中取：

```go
username, ok := currentUsername(r)
```

### 8.3 JWT 是什么，不是什么

JWT 是一个由服务器签名的字符串。它通常包含：

```text
Header：使用何种签名算法
Payload：sub（用户名）、过期时间等 claims
Signature：用服务器 secret 验证内容没有被篡改
```

JWT 不是加密箱。它的 payload 默认可被客户端解码查看，不要放密码、银行卡号等敏感信息。

当前项目：

```text
登录成功
  -> token Subject = username
  -> 有效期 24 小时
  -> React localStorage 保存 token
  -> Axios 请求 header 带 Bearer token
```

---

## 9. 账号接口：`handler/user.go`

### 9.1 注册

```go
var user model.User
json.NewDecoder(r.Body).Decode(&user)
```

把浏览器 POST 的 JSON 解成 Go `User` struct。

```go
var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,30}$`)
```

限制用户名：

- 3 到 30 位。
- 只能是字母、数字、下划线。

这个限制是为了避免不合理用户名、查询歧义和部分输入风险。密码课程版只检查至少三位，生产项目需要更成熟策略。

### 9.2 登录和签发 token

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
    Subject: user.Username,
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
})
tokenString, err := token.SignedString(signingKey)
```

- `HS256`：HMAC SHA-256 对称签名。
- `Subject`：此 token 代表哪个用户。
- `ExpiresAt`：24 小时后失效。
- `SignedString`：用服务端 secret 签名。

返回内容类型为 `text/plain`，前端把 response.data 当 token 存储。

---

## 10. 帖子接口：`handler/post.go`

### 10.1 上传接口

客户端上传的是 `multipart/form-data`：

```text
message      普通文本字段
media_file   文件字段
```

后端：

```go
file, header, err := r.FormFile("media_file")
defer file.Close()
```

- `file`：文件内容流。
- `header.Filename`：原文件名。
- `defer file.Close()`：函数结束自动关闭上传流。

扩展名检查：

```go
extension := strings.ToLower(filepath.Ext(header.Filename))
mediaType, ok := mediaTypes[extension]
```

通过后缀把文件归为 image 或 video：

```text
.jpg .jpeg .png .gif .webp -> image
.mp4 .mov .avi .flv .wmv  -> video
```

用：

```go
uuid.NewString()
```

生成帖子 ID，避免不同文件同名覆盖。

最后：

```go
writeJSON(w, http.StatusOK, post)
```

把保存成功的帖子完整 JSON 返回给前端。

### 10.2 搜索接口

```text
GET /search
GET /search?user=alice
GET /search?keywords=sunset
```

逻辑：

```text
有 user 参数：按 user 查
否则：按 keywords 查
```

当前 `GET /search` 不带参数会列出全部帖子。

### 10.3 删除接口

```text
DELETE /post/{id}
```

`mux.Vars(r)["id"]` 从路径里取帖子 ID。

当前已认证用户名在 JWT context 中，传给 `service.DeletePost(id, username)`，由两条件 ES query 确保只能删自己的帖。

---

## 11. React 入口：页面如何在浏览器中启动

### 11.1 `src/index.js`

```jsx
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
```

逐项理解：

- `document.getElementById("root")`：在 `public/index.html` 中找 React 挂载点。
- `ReactDOM.createRoot`：React 18 的渲染入口。
- `<App />`：整个 UI 的根组件。
- `<BrowserRouter>`：启用前端路由，URL 换页面不需要浏览器重新请求一个 HTML 文件。
- `<React.StrictMode>`：开发期帮助发现副作用问题；在开发模式下某些 effect 可能被刻意多执行一次，别误认为 API 一定坏了。

### 11.2 `src/constants.js`

```jsx
export const TOKEN_KEY = "token";

export const BASE_URL =
  process.env.REACT_APP_API_BASE_URL || "http://localhost:8080";

export const SEARCH_KEY = {
  all: 0,
  keyword: 1,
  user: 2,
};
```

三个目的：

- 集中管理 localStorage 的 key，避免每个组件写不同字符串。
- `BASE_URL` 支持本机或云端切换。
- `SEARCH_KEY` 让搜索类型不是散落的魔法数字。

React Create App 只有以 `REACT_APP_` 开头的环境变量才会注入前端构建结果。

### 11.3 `App.js`：登录状态的单一来源

```jsx
const [isLoggedIn, setIsLoggedIn] = useState(
  Boolean(localStorage.getItem(TOKEN_KEY))
);
```

`useState` 表示组件状态：

- 第一次加载时若 localStorage 有 token，状态为 true。
- 没 token，状态为 false。

```jsx
const loggedIn = (token) => {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token);
    setIsLoggedIn(true);
  }
};
```

登录成功后：

1. 保存 token，使刷新页面仍保持登录。
2. 更新 React state，使页面立刻重新渲染。

```jsx
const logout = () => {
  localStorage.removeItem(TOKEN_KEY);
  setIsLoggedIn(false);
};
```

登出做相反动作。

### 11.4 `Main.js`：前端路由守卫

路由表：

| URL | 登录前 | 登录后 |
| --- | --- | --- |
| `/` | 跳登录 | 跳 Create |
| `/login` | Login | 跳 Create |
| `/register` | Register | 跳 Create |
| `/create` | 跳登录 | Landing |
| `/collection` | 跳登录 | Collection |

例如：

```jsx
const showCollection = () =>
  isLoggedIn ? <Collection /> : <Navigate to="/login" replace />;
```

这是一层用户体验保护。真正安全仍由 Go 后端 JWT middleware 保证，因为任何人都可以在浏览器控制台手改前端状态，但不能伪造后端签名 token。

---

## 12. React 登录注册：前端表单怎样请求后端

### 12.1 `Login.js`

Ant Design 的：

```jsx
<Form onFinish={onFinish}>
```

表单校验通过后自动调用 `onFinish`。解构参数：

```jsx
const onFinish = ({ username, password }) => { ... }
```

表示从表单值对象中取两个字段。

Axios 请求：

```jsx
axios({
  method: "POST",
  url: `${BASE_URL}/signin`,
  data: { username, password },
  headers: { "Content-Type": "application/json" },
})
```

这对应：

```http
POST /signin
Content-Type: application/json

{"username":"alice","password":"abc123"}
```

成功：

```jsx
handleLoggedIn(response.data);
```

调用父组件传入的函数，令 App 存 token 并改变登录状态。

失败：

```jsx
message.error("Login failed. Check your username and password.");
```

用户看见友好提示；`console.error` 留给开发者调试。

### 12.2 `Register.js`

`Form.useForm()` 取得表单实例，`useNavigate()` 取得跳转函数。

确认密码规则：

```jsx
dependencies={["password"]}
```

意味着 confirm 字段依赖 password；password 改了时会重新验证 confirm。

```jsx
({ getFieldValue }) => ({
  validator(_, value) {
    if (!value || getFieldValue("password") === value) {
      return Promise.resolve();
    }
    return Promise.reject(new Error(...));
  },
})
```

这是 Ant Design 自定义异步校验器：

- 两次密码一致：返回 resolved Promise。
- 不一致：返回 rejected Promise，表单阻止提交。

前端校验是用户体验；后端校验才是安全边界。任何人都能绕过浏览器 JavaScript 直接发 HTTP。

---

## 13. 上传：React `FormData` 到 Go `multipart.File`

### 13.1 `PostForm.js`

```jsx
<Upload.Dragger name="files" beforeUpload={() => false} maxCount={1}>
```

- `Upload.Dragger`：可拖拽文件组件。
- `beforeUpload={() => false}`：阻止 Ant Design 自动上传。因为你要等用户点击 modal 的 Create 后，统一构造 `FormData` 发给自己后端。
- `maxCount={1}`：一次只选一个媒体文件。

`forwardRef`：

```jsx
export const PostForm = forwardRef((_, formRef) => ...)
```

使父组件 `CreatePostButton` 可以通过 ref 调用子表单的：

```jsx
this.postForm.validateFields()
this.postForm.resetFields()
```

### 13.2 `CreatePostButton.js`

这是 class component，和其他 function component 共存。你应能读懂两者：

| 写法 | 状态 | 生命周期/副作用 |
| --- | --- | --- |
| Function Component | `useState`、`useEffect` | Hook |
| Class Component | `this.state`、`this.setState` | 生命周期方法 |

在 `handleOk`：

```jsx
this.postForm.validateFields()
  .then(({ description, uploadPost }) => { ... })
```

先验证必填字段。再从 AntD 文件对象中取：

```jsx
const { type, originFileObj } = uploadPost[0];
```

- `type`：浏览器报告的 MIME 类型，例如 `image/png`。
- `originFileObj`：真实 File 对象。

构造：

```jsx
const formData = new FormData();
formData.append("message", description);
formData.append("media_file", originFileObj);
```

字段名必须准确对应 Go：

```go
r.FormValue("message")
r.FormFile("media_file")
```

Axios header 只传 Authorization：

```jsx
headers: {
  Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
}
```

不要手动写 multipart 的 `Content-Type` 和 boundary；浏览器会替你正确生成。

上传成功后：

```jsx
this.props.onShowPost(postType);
```

通知 Collection 切换到图片或视频 tab，并刷新帖子列表。

---

## 14. Collection：React 状态、Effect 与搜索

`Collection.js` 里三个关键状态：

```jsx
const [posts, setPosts] = useState([]);
const [activeTab, setActiveTab] = useState("image");
const [searchOption, setSearchOption] = useState({
  type: SEARCH_KEY.all,
  keyword: "",
});
```

它们分别代表：

- `posts`：后端返回的帖子。
- `activeTab`：当前显示图片还是视频。
- `searchOption`：当前过滤条件。

核心：

```jsx
useEffect(() => {
  ...
}, [searchOption]);
```

意思是：每次 `searchOption` 改变时，执行一次请求。不是每次任何页面重渲染都请求。

请求 URL 按搜索类型变化：

```text
all      -> GET /search
user     -> GET /search?user=alice
keyword  -> GET /search?keywords=sunset
```

`encodeURIComponent(keyword)` 是必须养成的习惯。它把空格、中文、`&` 等特殊字符安全放进 URL 查询参数。

`renderPosts`：

```text
先 posts.filter(post.type === type)
image -> 转换为 PhotoGallery 的 images 格式
video -> 用 HTML <video controls> 显示
```

这里也体现 Adapter 思路：后端统一返回 `Post`，UI 图库组件有自己的 photo 数据要求，前端做一次格式转换。

### `SearchBar.js`

SearchBar 并不自己请求 API，而是：

```jsx
handleSearch({ type: searchType, keyword: value.trim() });
```

把筛选状态交给父组件 Collection。Collection 状态改变，effect 自动请求。

这种设计叫“状态提升”：

```text
SearchBar 负责用户输入
Collection 拥有搜索条件和帖子数据
```

这样子组件职责简单，后期更容易维护。

---

## 15. 图片展示和删除：`PhotoGallery.js`

组件的两个 state：

```jsx
const [galleryImages, setGalleryImages] = useState(images);
const [index, setIndex] = useState(-1);
```

- `galleryImages`：当前页面展示的图片列表。
- `index=-1`：Lightbox 未打开；大于等于 0 表示打开第几张图片。

```jsx
useEffect(() => setGalleryImages(images), [images]);
```

父组件异步拿到新 images 后，子组件同步自己的显示状态。

删除：

```jsx
axios({
  method: "DELETE",
  url: `${BASE_URL}/post/${postId}`,
  headers: {
    Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
  },
})
```

成功后：

```jsx
setGalleryImages((current) =>
  current.filter((image) => image.postId !== postId)
);
```

使用函数式 state 更新。React 可能批量处理 state 更新，基于最新 `current` 计算比直接用闭包旧变量更稳。

`useMemo`：

```jsx
const imageArr = useMemo(() => ..., [galleryImages]);
```

只在 `galleryImages` 改变时重建图库格式的数据，减少不必要计算。这里数据量小，性能差异有限，但它很好地展示了 React memoization 的概念。

---

## 16. AI 生成图片：`Landing.js`

页面行为：

```text
用户输入 prompt
  -> 点击箭头
  -> OpenAI 图像 API 返回 base64 PNG
  -> React 用 data:image/png;base64,... 预览
  -> 用户点击上传
  -> base64 转 Blob / File
  -> 复用 POST /upload 保存到 Social AI
```

### 16.1 `useState`

```jsx
const [inputValue, setInputValue] = useState("");
const [isGeneratingImage, setIsGeneratingImage] = useState(false);
const [generatedImageUrl, setGeneratedImageUrl] = useState("");
```

- inputValue：prompt。
- isGeneratingImage：展示 loading 遮罩，避免用户不停重复点击。
- generatedImageUrl：base64 data URL。

### 16.2 调用图像 API

```jsx
const openai = new OpenAI({
  apiKey,
  dangerouslyAllowBrowser: true,
});
```

`dangerouslyAllowBrowser` 这个名字已经在提醒你：直接在浏览器使用 API key 有安全风险。课程演示可用自己的临时、额度受控 key；生产不能这么做。

生产正确结构：

```text
React
  -> POST /generate-image (你自己的后端，用户身份、限流、内容审核)
  -> 后端从服务器环境变量读取 OPENAI_API_KEY
  -> OpenAI API
  -> 后端返回安全的图片 URL / job ID
```

绝不把长期 API key 打入 React build。任何拿到网页的人都能在浏览器开发者工具中找到它。

### 16.3 为什么要先转 Blob 再上传

生成图片是 data URL：

```text
data:image/png;base64,......
```

而现有 `/upload` 接口期望 multipart 文件。代码：

```jsx
const response = await fetch(generatedImageUrl);
const blob = await response.blob();
formData.append(
  "media_file",
  new File([blob], "ai-generated.png", { type: "image/png" })
);
```

把内存中的 data URL 变成浏览器 File 对象，之后就能复用正常的上传接口。复用既有接口而不是另建“AI 上传专用接口”，是简洁的设计。

---

## 17. CORS：为什么 React 3000 能请求 Go 8080

浏览器把：

```text
http://localhost:3000
http://localhost:8080
```

视为不同 origin，因为端口不同。

Go 路由最后包了一层：

```go
return handlers.CORS(
    handlers.AllowedOrigins([]string{"*"}),
    handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
    handlers.AllowedMethods([]string{
        http.MethodGet,
        http.MethodPost,
        http.MethodDelete,
        http.MethodOptions,
    }),
)(router)
```

这使浏览器预检和正式请求被允许。

课程本机环境用 `"*"` 很省事。生产环境应明确只允许真实前端域名，例如：

```text
https://app.example.com
```

否则任意站点都可从浏览器对你的 API 发跨域请求。

---

## 18. API 契约：用 Postman/终端独立验证后端

健康检查：

```powershell
curl.exe -i http://localhost:8080/healthz
```

注册：

```powershell
curl.exe -i -X POST http://localhost:8080/signup `
  -H "Content-Type: application/json" `
  -d "{\"username\":\"alice_01\",\"password\":\"abc123\"}"
```

登录：

```powershell
curl.exe -i -X POST http://localhost:8080/signin `
  -H "Content-Type: application/json" `
  -d "{\"username\":\"alice_01\",\"password\":\"abc123\"}"
```

响应正文是一串 JWT。复制为 PowerShell 变量：

```powershell
$token = "粘贴登录响应的完整 token"
```

查看所有帖子：

```powershell
curl.exe -i http://localhost:8080/search `
  -H "Authorization: Bearer $token"
```

按用户：

```powershell
curl.exe -i "http://localhost:8080/search?user=alice_01" `
  -H "Authorization: Bearer $token"
```

上传：

```powershell
curl.exe -i -X POST http://localhost:8080/upload `
  -H "Authorization: Bearer $token" `
  -F "message=My first post" `
  -F "media_file=@C:\path\to\image.png"
```

注意：

- `-F` 才是 multipart form。
- 不带 Authorization，上传/搜索/删除返回 401 是正确的。
- `username` 不需要由客户端作为上传字段传入，服务端从 JWT 获得。

---

## 19. Social AI 与 OnlineOrder 的核心差异

| 对比 | OnlineOrder | Social AI |
| --- | --- | --- |
| 后端语言 | Java | Go |
| Web 框架 | Spring Boot / Spring MVC | net/http + Gorilla Mux |
| 主数据 | PostgreSQL 关系表 | Elasticsearch 文档 |
| 登录 | Spring Security Session Cookie | JWT Bearer Token |
| 典型核心 | 事务、外键、购物车 | 上传、搜索、媒体 URL |
| 缓存 | Caffeine | 当前没有独立缓存层 |
| 文件 | 外链菜单图片 | 本地 uploads / GCS |
| 部署课程 | Docker/RDS/ECR/ECS | GCE/Elasticsearch/GCS/App Engine |

你不需要把两套技术混为一谈。要理解选型：

```text
强事务、复杂关系、订单金额
  -> PostgreSQL 更自然

全文搜索、文档搜索、社交内容检索
  -> Elasticsearch 很擅长

实际大型产品常常是 PostgreSQL/MySQL 作为事实主库，
Elasticsearch 作为由主库同步而来的搜索索引，
而不是把所有用户认证和核心事实都交给 Elasticsearch。
```

---

## 20. 当前项目的生产风险与升级路线

这张表是最值得你认真记住的。它能让你在面试中体现出不会盲目把 demo 当生产系统。

| 优先级 | 当前问题 | 风险 | 建议改法 |
| --- | --- | --- | --- |
| P0 | 用户密码明文存 ES | 数据泄露即泄露密码 | bcrypt/Argon2 哈希；用户放 PostgreSQL |
| P0 | OpenAI key 在浏览器 | key 可被任何访问网页者拿走 | 图像生成移到后端，key 存服务端环境变量 |
| P0 | ES 本机关闭安全且可公开端口 | 任何人可能读写数据 | 启用 TLS、账号、VPC/私网访问 |
| P0 | CORS 允许 `*` | 任意来源可调用 API | 白名单真实前端域名 |
| P0 | 文件类型主要看扩展名 | 改名文件可伪装类型 | 检测 MIME/魔数，限制文件大小，病毒扫描 |
| P1 | 删除只删 ES 文档 | 本地/GCS 留下孤儿文件 | 事务式工作流或异步 cleanup job |
| P1 | 搜索最多 100 条且无分页 | 数据大时不完整 | pagination/search_after、排序 |
| P1 | JWT 只用长 access token | 无法良好登出/续期 | 短 token + refresh token + denylist 策略 |
| P1 | 无速率限制 | 登录、上传、AI 容易被滥用 | rate limit、配额、验证码、审计 |
| P1 | 无内容审核 | 生成/上传内容可能违规 | 审核服务、举报、人工处理流程 |
| P1 | 测试很少 | 改动容易回归 | Go handler 测试、前端测试、E2E |
| P2 | 无观测性 | 不知道慢在哪里/失败多少 | structured logs、metrics、tracing |
| P2 | 无 CI/CD | 发布全手工 | GitHub Actions、镜像、环境分离 |

### 生产级数据架构应如何演进

合理的后续方向：

```text
PostgreSQL
  - users
  - posts
  - permissions
  - likes/comments/follows
  - password hash
  - 事务事实数据

Object Storage (S3 / GCS)
  - 图片、视频真实字节

Elasticsearch / OpenSearch
  - 从 PostgreSQL 同步的可搜索帖子索引

Redis
  - rate limit、session/token denylist、热门缓存
```

这不表示你现在要一口气全部重构。它表示你知道系统规模变大后的正确方向。

---

## 21. 面试时如何讲 Social AI

### 60 秒版本

> 我完成了一个 React、Go 和 Elasticsearch 的社交媒体全栈项目。前端用 React Router 实现登录、注册、创作和内容集合页面，通过 Axios 发起 JSON 和 multipart 上传请求；Go 后端使用 Gorilla Mux 组织 REST 路由，用 JWT middleware 保护上传、搜索和删除接口，并把认证用户写入 request context。帖子元数据保存到 Elasticsearch，支持按用户精确查询和按 message 全文搜索，图片视频文件在本地开发时存入 uploads，云端设计可切换到 Google Cloud Storage。项目还接入了图像生成 API，将生成的 base64 图片转换为 File 后复用上传接口。当前我清楚地识别了生产版需要将明文密码改为哈希、将 OpenAI key 移到后端、增加文件校验、分页和限流等问题。

### 高频问题与回答方向

**Q：为什么用 JWT？**

> 前端与后端分离、可能跨域部署时，Bearer token 便于客户端在请求头携带身份。后端验证签名和过期时间后，将 subject 放到 request context，业务 handler 不信任客户端自行提交的用户名。

**Q：JWT 和 session 有什么区别？**

> OnlineOrder 使用服务端 session，浏览器带 cookie，服务端保存会话状态；Social AI 使用 JWT，客户端携带签名 token，后端可无状态验证。JWT 适合 API/多端，但撤销、续期和泄露风险需要更谨慎设计。

**Q：为什么 Post 元数据与媒体文件分开？**

> ES 适合索引、搜索结构化文本元数据，GCS/对象存储适合大文件。把视频字节塞入 ES 会破坏搜索集群性能和成本结构。

**Q：如何保证只能删除自己的帖子？**

> JWT middleware 解析 token 获得认证 username，删除 service 用 `id AND user` 构建 ES bool query。客户端 URL 中只有 post id，不能通过自报用户名越权。

**Q：前端如何在上传时带文件和文字？**

> 用 `FormData`，append `message` 和 `media_file`。浏览器自动编码 multipart boundary，Go 用 `FormValue` 和 `FormFile` 接收。

**Q：这个项目最大的安全问题是什么？**

> 密码明文存储和浏览器端 OpenAI key 都是 P0 问题。生产应使用强密码哈希、服务端密钥管理、HTTPS、严格 CORS、限流和文件校验。

---

## 22. 你应该亲手做的练习

### 第一阶段：能独立运行、理解每一条请求

1. 停掉 Elasticsearch，再启动 Go，观察错误并解释为什么。
2. 用 curl 或 Postman 注册、登录、复制 JWT、上传、搜索、删除。
3. 在浏览器 DevTools Network 看 `/signin`、`/upload` 的 Request Headers 和 Request Payload。
4. 打开 `http://localhost:9200/_cat/indices?v` 看 user 和 post 索引。
5. 在 `socialai/uploads` 看上传文件；理解 ES 只存 URL。

### 第二阶段：做实用功能

1. 给帖子添加 `created_at`，搜索结果按时间倒序。
2. 添加文件大小限制。
3. 新增只展示“我的帖子”的页面按钮。
4. 给 Search 加分页参数 `page`、`size`。
5. 在 Go 写 handler 单元测试，例如无 token 必须返回 401。

### 第三阶段：做成更可信的个人作品

1. 用户迁移到 PostgreSQL，密码改 bcrypt/Argon2。
2. ES 仅作 post 的搜索索引。
3. OpenAI 图像生成改成后端代理接口，增加每用户配额。
4. 使用 GCS/S3、签名 URL 和 CDN。
5. 加 Dockerfile、CI、部署说明、监控和错误响应标准。

---

## 23. 两个项目放在一起，你已经学到了什么

你不是只做了两个“页面项目”。你已经接触到：

```text
后端：
  Java、Spring Boot、Go、REST、HTTP、JSON、JWT、Session

数据：
  PostgreSQL、表设计、主键外键、事务、SQL、Elasticsearch、全文搜索

前端：
  React、组件、props、state、useEffect、Router、Axios、FormData

工程：
  Gradle、npm、Docker、环境变量、配置文件、云存储、云部署基础

质量：
  单元测试、缓存、认证、错误排查、生产风险评估
```

不过“学到”要转化成“能就业”，还需要反复做三件事：

1. 不看代码，用自己的话讲清请求链。
2. 自己改一个功能并写测试。
3. 记录错误和解决过程，形成排障能力。

---

## 24. 最后用一句话记住 Social AI

Social AI 的核心不是“会调用一个 AI 图片接口”，而是：

```text
React 管理页面、路由和用户状态
  -> JWT 放在请求头证明身份
  -> Go middleware 验证身份
  -> handler/service 处理上传、搜索、删除
  -> Elasticsearch 保存可搜索的帖子文档
  -> 本地磁盘或 GCS 保存真正的媒体文件
  -> 云配置让同一套应用可以从本机走向部署环境
```

你能把这条链、它和 OnlineOrder 的区别、以及它的生产缺口讲清楚，就不只是“跟着课件做过 React/Go”，而是已经具备全栈工程学习和迭代的基本能力。
