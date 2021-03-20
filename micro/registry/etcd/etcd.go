package etcd

import (
	"sync"

	"github.com/Augustu/go-draft/micro/registry"
	"go.etcd.io/etcd/clientv3"
)

var (
	prefix = "/draft/registry/"
)

type etcdRegistry struct {
	client  *clientv3.Client
	options registry.Options

	sync.RWMutex
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func NewRegistry(opts ...registry.RegisterOption) registry.Registry {
	// e := &etcdRegistry{
	// 	register: make(map[string]uint64),
	// 	leases:   make(map[string]clientv3.LeaseID),
	// }

	return nil
}

func (e *etcdRegistry) Init(opts ...registry.Option) error {
	return nil
}

func (e *etcdRegistry) Options() registry.Options {
	return e.options
}

func (e *etcdRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	return nil
}
