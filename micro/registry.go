package micro

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3/registry"
)

//type Registry interface {
//	Address() []string
//}

func etcdRegistry(c *EtcdRegistry) (reg registry.Registry) {
	return etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = c.Address()
	})
}
