package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container é um simples DI container
type Container struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewContainer cria uma nova instância do container
func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register registra um serviço no container
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// RegisterSingleton registra um singleton no container
func (c *Container) RegisterSingleton(name string, factory func() interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Lazy initialization - o serviço só é criado quando solicitado
	c.services[name] = &singleton{
		factory: factory,
		once:    &sync.Once{},
	}
}

// Get obtém um serviço do container
func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	// Se for um singleton, inicializa se necessário
	if s, ok := service.(*singleton); ok {
		s.once.Do(func() {
			s.instance = s.factory()
		})
		return s.instance, nil
	}

	return service, nil
}

// MustGet obtém um serviço do container ou entra em pânico
func (c *Container) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// GetAs obtém um serviço do container com cast de tipo
func (c *Container) GetAs(name string, target interface{}) error {
	service, err := c.Get(name)
	if err != nil {
		return err
	}

	serviceValue := reflect.ValueOf(service)
	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	targetElem := targetValue.Elem()
	if !serviceValue.Type().AssignableTo(targetElem.Type()) {
		return fmt.Errorf("cannot assign %s to %s", serviceValue.Type(), targetElem.Type())
	}

	targetElem.Set(serviceValue)
	return nil
}

type singleton struct {
	factory  func() interface{}
	instance interface{}
	once     *sync.Once
}
