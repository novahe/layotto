package hello

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*HelloFactory)
	Create(name string) (HelloService, error)
}

type HelloFactory struct {
	Name          string
	FatcoryMethod func() HelloService
}

func NewHelloFactory(name string, f func() HelloService) *HelloFactory {
	return &HelloFactory{
		Name:          name,
		FatcoryMethod: f,
	}
}

type helloRegistry struct {
	stores map[string]func() HelloService
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName) // 添加服务信息
	return &helloRegistry{
		stores: make(map[string]func() HelloService),
		info:   info,
	}
}

func (r *helloRegistry) Register(fs ...*HelloFactory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FatcoryMethod
		r.info.RegisterComponent(ServiceName, f.Name) // 注册组件信息
	}
}

func (r *helloRegistry) Create(name string) (HelloService, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name) // 加载组件信息
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
