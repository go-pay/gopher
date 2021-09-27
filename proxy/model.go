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
	ProxySchema SchemaType `json:"proxy_schema" yaml:"proxy_schema" toml:"proxy_schema"` // SchemaHTTP or SchemaHTTPS
	ProxyHost   string     `json:"proxy_host" yaml:"proxy_host" toml:"proxy_host"`       // 转发到的接口 Host
	ProxyPort   string     `json:"proxy_port" yaml:"proxy_port" toml:"proxy_port"`       // 转发到的接口 Port
	ServerPort  string     `json:"server_port" yaml:"server_port" toml:"server_port"`    // 代理转发服务启动的端口
	Key         string     `json:"key" yaml:"key" toml:"key"`                            // 简单的校验Key
}
