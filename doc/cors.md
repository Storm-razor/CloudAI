# CORS 作用、背景与配置详解

本文结合 `internal/middlerware/cors.go` 的实现，解释 CORS 的背景、工作机制、配置项含义与生产实践。

## 1. 背景：为什么需要 CORS

浏览器有同源策略（Same-Origin Policy）：默认不允许网页脚本跨域读取数据。只要协议/域名/端口任一不同，就被视为跨域。  
在前后端分离场景中，前端通常部署在 `https://app.example.com`，后端 API 在 `https://api.example.com`，这会触发浏览器限制。  
CORS 是一套浏览器层的跨域许可协议：**后端通过响应头声明允许跨域，浏览器决定是否放行**。

## 2. 工作机制：简单请求与预检请求

- **简单请求**：如 GET/POST 且请求头很简单。浏览器直接发请求，如果响应头中包含 `Access-Control-Allow-Origin` 且允许当前域名，浏览器才让 JS 访问响应内容。
- **预检请求（OPTIONS）**：带自定义头、JSON、PUT/DELETE 等，会先发 OPTIONS 询问“是否允许”。后端允许后，浏览器才发真实请求。

因此，CORS 本质上是**后端声明权限 + 浏览器执行限制**。

## 3. 你的代码如何配置 CORS

`SetupCORS` 读取配置并构造 `cors.Config`：

- `AllowOrigins`：允许跨域的来源白名单。  
  生产中通常写前端域名，如 `https://app.example.com`。  
  注意：若要携带 Cookie，不能用 `*`。
- `AllowMethods`：允许的 HTTP 方法，如 `GET,POST,PUT,DELETE,OPTIONS`。
- `AllowHeaders`：允许客户端携带的自定义请求头，如 `Authorization,Content-Type`。
- `ExposeHeaders`：允许前端 JS 读取的响应头，如 `X-Request-Id`。
- `AllowCredentials`：是否允许携带凭证（Cookie、Authorization）。  
  为 `true` 时，`AllowOrigins` 必须是明确域名，不能是 `*`。
- `MaxAge`：预检请求（OPTIONS）的缓存时间，减少预检次数。  
  这里配置为字符串并解析成 `time.Duration`，解析失败回退 12 小时。

## 4. 实际应用场景

- **前后端分离**：前端与后端不同域，必须配置 CORS，否则浏览器会拦截响应。
- **管理后台 + 公网 API**：后台域名是私有域，CORS 只放行后台域。
- **多环境部署**：测试、预发、生产域名不同，需要分别配置白名单。
- **微前端**：多个子应用不同域访问同一 API，需要允许多域名。

## 5. 常见坑与建议

- `AllowCredentials=true` 时 **不能** 配置 `AllowOrigins=*`。
- 预检请求如果未处理会导致前端报错“CORS preflight failed”。
- 如果响应头缺失 `Access-Control-Allow-Origin`，浏览器会直接拦截，即使后端正常返回 200。
- 推荐在生产环境使用**白名单域名**，避免任意跨域访问。

