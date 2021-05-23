package di

import (
	"errors"
	"reflect"

	"github.com/psewda/typing/internal/utils"
)

// NewInstanceFunc is a function type which creates a new instance.
type NewInstanceFunc func(params ...interface{}) (interface{}, error)

// InstanceType represents the instance type enum.
type InstanceType uint16

const (
	// InstanceTypeAuth is the enum member of type auth.
	InstanceTypeAuth InstanceType = iota

	// InstanceTypeUserinfo is the enum member of type userinfo.
	InstanceTypeUserinfo

	// InstanceTypeNotestore is the enum member of type notestore.
	InstanceTypeNotestore

	// InstanceTypeSectionstore is the enum member of type sectionstore.
	InstanceTypeSectionstore
)

// Container is the dependency injection container. It manages
// the life cycle of instance creation.
type Container interface {
	// Add inserts the instance creation handler with the type value.
	Add(it InstanceType, handler NewInstanceFunc)

	// GetInstance creates new instance associated with the type value.
	GetInstance(it InstanceType, p ...interface{}) (interface{}, error)
}

// DefaultContainer implements the dependency injection container.
type DefaultContainer struct {
	handlers map[InstanceType]NewInstanceFunc
}

// Add inserts the instance creation handler in the container.
func (c *DefaultContainer) Add(t InstanceType, handler NewInstanceFunc) {
	c.handlers[t] = handler
}

// GetInstance creates new instance of the specified type value.
func (c *DefaultContainer) GetInstance(t InstanceType, params ...interface{}) (interface{}, error) {
	handler, ok := c.handlers[t]
	if !ok {
		return nil, errors.New("handler not found")
	}

	in := make([]reflect.Value, 0)
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	returns := reflect.ValueOf(handler).Call(in)
	instance := returns[0].Interface()
	err := returns[1].Interface()
	if err != nil {
		msg := "instance creation failed, check the handler"
		return nil, utils.Error(msg, err.(error))
	}

	return instance, nil
}

// New create a new instance of default container.
func New() *DefaultContainer {
	return &DefaultContainer{
		handlers: make(map[InstanceType]NewInstanceFunc),
	}
}
