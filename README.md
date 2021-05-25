# go-util

### 项目 common 库

[![Golang](https://img.shields.io/badge/golang-1.13+-brightgreen.svg)](https://golang.google.cn)
[![GoDoc](https://img.shields.io/badge/doc-go.dev-informational.svg)](https://pkg.go.dev/github.com/iGoogle-ink/gotil)
[![Drone CI](https://cloud.drone.io/api/badges/iGoogle-ink/gotil/status.svg)](https://cloud.drone.io/iGoogle-ink/gotil)
[![GitHub Release](https://img.shields.io/github/v/release/iGoogle-ink/gotil)](https://github.com/iGoogle-ink/gotil/releases)
[![License](https://img.shields.io/github/license/iGoogle-ink/gopay)](https://www.apache.org/licenses/LICENSE-2.0)


### Install
```bash
go get github.com/iGoogle-ink/gotil
```

### 目录结构
```bash
.
├── LICENSE
├── README.md
├── aes
│   ├── aes_cbc_decrypt.go
│   ├── aes_cbc_encrypt.go
│   └── pkcs_padding.go
├── body_map.go
├── conf
│   ├── conf.go
│   ├── conf_test.go
│   ├── config.json
│   ├── config.toml
│   └── config.yaml
├── constants.go
├── convert.go
├── ecode
│   ├── common.go
│   └── ecode.go
├── errgroup
│   └── errgroup.go
├── go.mod
├── go.sum
├── limit
│   ├── limiter.go
│   └── limiter_test.go
├── lru
│   ├── cache.go
│   └── cache_test.go
├── orm
│   ├── model.go
│   ├── mysql_gorm.go
│   ├── mysql_xorm.go
│   ├── page.go
│   └── redis.go
├── proxy
│   ├── http.go
│   ├── model.go
│   ├── proxy_test.go
│   └── service.go
├── random.go
├── rate
│   └── rate.go
├── retry
│   ├── retry.go
│   └── retry_test.go
├── slice.go
├── string.go
├── verify.go
├── web
│   ├── gin.go
│   ├── gin_test.go
│   ├── page.go
│   ├── rsp.go
│   └── verify.go
├── ws
│   └── websocket_connection.go
├── xhttp
│   ├── client.go
│   ├── client_test.go
│   └── model.go
├── xlog
│   ├── debug_logger.go
│   ├── error_logger.go
│   ├── info_logger.go
│   ├── log.go
│   ├── log_test.go
│   ├── warn_logger.go
│   └── zap.go
├── xrsa
│   ├── rsa_decrypt.go
│   └── rsa_encrypt.go
└── xtime
    ├── parse_format.go
    └── xtime.go

```
