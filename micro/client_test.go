package micro

import (
	"testing"
)

func TestClient(t *testing.T) {
	// 解开注释测试

	//var c proto.GreeterService
	//rg := &EtcdRegistry{Addrs: []string{"api.fumm.cc:2379"}}
	//
	//InitClient("service.client.clientName", "latest", rg, func(client client.Client) {
	//	c = proto.NewGreeterService("service.server.serverName", client)
	//})
	//count := 0
	//for {
	//	if count == 5 {
	//		return
	//	}
	//	time.Sleep(time.Second * 2)
	//	in := &proto.Request{Name: "Jerry"}
	//	rsp, err := c.Hello(context.Background(), in)
	//	if err != nil {
	//		xlog.Error(err)
	//		continue
	//	}
	//	xlog.Debug("rsp:", rsp.Msg)
	//	count++
	//}
}
