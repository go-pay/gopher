package web

import (
	"bytes"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/alog"
	"github.com/go-pay/gopher/util"
	"github.com/go-pay/gopher/xlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	OutputStdout OutputType = "stdout"
	OutputSLS    OutputType = "sls"
	OutputFile   OutputType = "file"

	_SlsTopic = "access_log"
)

var (
	// SLS client
	slsLogger *alog.Client
	// zip client
	zipLogger *zap.Logger
)

type OutputType string

type AccessConfig struct {
	AppName      string     `json:"app_name" yaml:"app_name" toml:"app_name"`                   // app name，记录日志标识
	OutputType   OutputType `json:"output_type" yaml:"output_type" toml:"output_type"`          // 日志输出类型：stdout、file、sls
	FilePath     string     `json:"file_path" yaml:"file_path" toml:"file_path"`                // 日志输出文件路径(输出类型为file时有效)
	SlsAccessKey string     `json:"sls_access_key" yaml:"sls_access_key" toml:"sls_access_key"` // 日志输出到阿里云SLS时的access_key
	SlsSecretKey string     `json:"sls_secret_key" yaml:"sls_secret_key" toml:"sls_secret_key"` // 日志输出到阿里云SLS时的secret_key
	SlsEndpoint  string     `json:"sls_endpoint" yaml:"sls_endpoint" toml:"sls_endpoint"`       // 日志输出到阿里云SLS时的endpoint
	SlsProject   string     `json:"sls_project" yaml:"sls_project" toml:"sls_project"`          // 日志输出到阿里云SLS时的project
	SlsLogStore  string     `json:"sls_log_store" yaml:"sls_log_store" toml:"sls_log_store"`    // 日志输出到阿里云SLS时的log_store
}

// AccessLog middleware for request and response body
func (g *GinEngine) AccessLog(ac *AccessConfig) gin.HandlerFunc {
	if ac == nil {
		ac = &AccessConfig{
			AppName:    "default",
			OutputType: OutputStdout,
			FilePath:   "./log/access.log",
		}
	}
	switch ac.OutputType {
	case OutputSLS:
		cfg := &alog.Config{
			AccessKey: ac.SlsAccessKey,
			SecretKey: ac.SlsSecretKey,
			Endpoint:  ac.SlsEndpoint,
			Project:   ac.SlsProject,
			LogStore:  ac.SlsLogStore,
		}
		alogger, err := alog.New(cfg)
		if err != nil || alogger == nil {
			xlog.Warnf("init sls logger failed, err: %v", err)
			return nil
		}
		slsLogger = alogger
	case OutputStdout:
	case OutputFile:
		if ac.FilePath == "" {
			ac.FilePath = "./log/access.log"
		}
		zipLogger = initZap(ac.FilePath)
	default:
		return nil
	}

	return func(c *gin.Context) {
		var (
			st        = time.Now()
			rHost     = c.Request.Host
			rUri      = c.Request.RequestURI
			rMethod   = c.Request.Method
			rHeader   = c.Request.Header.Clone()
			rClientIP = ClientIP(c.Request, rHeader)
			req       = map[string]any{}
			schema    = "http"
		)
		if c.Request.TLS != nil {
			schema = "https"
		}
		reqBs, err := ReadRequestBody(c.Request)
		if err != nil {
			return
		}
		_ = util.UnmarshalBytes(reqBs, &req)
		writer := responseWriter{
			ResponseWriter: c.Writer,
			b:              bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		defer func() {
			rbs := writer.b.Bytes()
			rsp := &CommonRsp{}
			_ = util.UnmarshalBytes(rbs, rsp)
			dataMap := map[string]any{
				"app_name":       ac.AppName,
				"client_ip":      rClientIP,
				"cost_ms":        time.Since(st).Milliseconds(),
				"host":           rHost,
				"method":         rMethod,
				"path":           rUri,
				"req_body":       req,
				"req_header":     rHeader,
				"result_code":    rsp.Code,
				"result_message": rsp.Message,
				"rsp_body":       rsp,
				"rsp_header":     c.Writer.Header(),
				"schema":         schema,
				"status_code":    c.Writer.Status(),
				"ts":             st.Unix(),
			}
			// output switch
			switch ac.OutputType {
			case OutputSLS:
				if slsLogger != nil {
					_ = slsLogger.Info(_SlsTopic, dataMap)
				}
			case OutputStdout:
				log.Printf("access_log: %s\n", util.MarshalBytes(dataMap))
			case OutputFile:
				if zipLogger != nil {
					fields := []zapcore.Field{
						zap.String("app_name", ac.AppName),
						zap.String("client_ip", rClientIP),
						zap.Int64("cost_ms", time.Since(st).Milliseconds()),
						zap.String("host", rHost),
						zap.String("method", rMethod),
						zap.String("path", rUri),
						zap.Any("req_body", req),
						zap.Any("req_header", rHeader),
						zap.Int("result_code", rsp.Code),
						zap.String("result_message", rsp.Message),
						zap.Any("rsp_body", rsp),
						zap.Any("rsp_header", c.Writer.Header()),
						zap.String("schema", schema),
						zap.Int("status_code", c.Writer.Status()),
						zap.Int64("ts", st.Unix()),
					}
					zipLogger.Info("access_log", fields...)
				}
			}
		}()
		c.Next()
	}
}

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	// 向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	// 完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// init zap
// fileName log path ./access_log
func initZap(fileName string) *zap.Logger {
	zConfig := zapcore.EncoderConfig{
		MessageKey: "",
		LevelKey:   "",
		TimeKey:    "",
		CallerKey:  "",
	}
	// io.Writer 使用 lumberjack
	infoWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024, //最大体积，单位M，超过则切割
		MaxBackups: 5,    //最大文件保留数，超过则删除最老的日志文件
		MaxAge:     30,   //最长保存时间30天
		Compress:   true, //是否压缩
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(zConfig), zapcore.AddSync(infoWriter), zap.InfoLevel), // 将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.InfoLevel))
}
