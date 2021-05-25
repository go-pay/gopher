package micro

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
)

func InitClient(cliName, version string, registry *EtcdRegistry, fn func(client.Client)) {
	var (
		s    micro.Service
		opts []micro.Option
	)
	opts = []micro.Option{
		micro.Name(cliName),
		//micro.WrapClient(LogClientWrap),
		micro.Version(version),
	}
	if registry != nil {
		opts = append(opts, micro.Registry(etcdRegistry(registry)))
	}
	//if broker != nil {
	//	opts = append(opts, micro.Broker(mqttBroker(broker)))
	//}

	s = micro.NewService(opts...)
	//s.Init()
	fn(s.Client())
}
