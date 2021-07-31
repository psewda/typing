package ioc

import (
	"errors"
	"reflect"

	"github.com/psewda/typing/internal/utils"
)

// ActivatorFunc creates a new instance of specified type.
type ActivatorFunc func(params ...interface{}) (interface{}, error)

// InstanceType represents the instance type enum.
type InstanceType uint16

// Container is the ioc container. It manages the life cycle of instance creation.
type Container interface {
	// Add inserts the instance creation activator with the type value.
	Add(t InstanceType, activator ActivatorFunc)

	// GetInstance creates new instance associated with the type value.
	GetInstance(t InstanceType, p ...interface{}) (interface{}, error)
}

// DefaultContainer implements the dependency injection container.
type DefaultContainer struct {
	activators map[InstanceType]ActivatorFunc
}

// Add inserts the instance creation activator in the container.
func (c *DefaultContainer) Add(t InstanceType, activator ActivatorFunc) {
	c.activators[t] = activator
}

// GetInstance creates new instance of the specified type value.
func (c *DefaultContainer) GetInstance(t InstanceType, params ...interface{}) (interface{}, error) {
	activator, ok := c.activators[t]
	if !ok {
		return nil, errors.New("activator not found")
	}

	in := make([]reflect.Value, 0)
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	returns := reflect.ValueOf(activator).Call(in)
	instance := returns[0].Interface()
	err := returns[1].Interface()
	if err != nil {
		msg := "instance creation failed, check the activator"
		return nil, utils.Error(msg, err.(error))
	}

	return instance, nil
}

// New create a new instance of default container.
func New() *DefaultContainer {
	return &DefaultContainer{
		activators: make(map[InstanceType]ActivatorFunc),
	}
}
