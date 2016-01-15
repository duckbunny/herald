package registry

import (
	"github.com/duckbunny/consul"
	"github.com/duckbunny/etcd"
	"github.com/duckbunny/vulcand"
)

// Method to register all available heralds
func RegisterAll() {
	consul.Register()
	etcd.Register()
	vulcand.Register()
}
