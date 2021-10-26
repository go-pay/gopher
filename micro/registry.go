package micro

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"go-micro.dev/v4/registry"
)

//type Registry interface {
//	Address() []string
//}

func etcdRegistry(c *EtcdRegistry) (reg registry.Registry) {
	return etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = c.Address()
	})
}
