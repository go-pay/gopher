package proxy

type SchemaType string

const (
	SchemaHTTP  SchemaType = "http://"
	SchemaHTTPS SchemaType = "https://"

	HTTP_METHOD_GET  = "GET"
	HTTP_METHOD_POST = "POST"

	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
	CONTENT_TYPE_XML  = "application/xml"

	HEADER_CONTENT_TYPE = "Content-Type"
	HEADER_CONTENT_KEY  = "key"
)

// 配置文件
type Config struct {
	ProxySchema SchemaType // SchemaHTTP or SchemaHTTPS
	ProxyHost   string     // 转发到的接口 Host
	ProxyPort   string     // 转发到的接口 Port
	ServerPort  string     // 代理转发服务启动的端口
	Key         string     // 简单的校验Key
}
