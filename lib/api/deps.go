package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Deps interface {
	Get(name string) any
	Set(name string, dep any)
	GetContext(name string) Deps
	SetContext(name string) Deps
	String() string
}

func NewDeps(data map[string]any) Deps {
	id := uuid.NewString()
	return &DefaultDeps{
		id:      id,
		data:    data,
		parent:  nil,
		RWMutex: sync.RWMutex{},
	}
}

type DefaultDeps struct {
	id     string
	data   map[string]any
	parent *DefaultDeps
	sync.RWMutex
}

func (deps *DefaultDeps) resolve() map[string]any {
	result := make(map[string]any)
	for k, v := range deps.data {
		if deps, ok := v.(*DefaultDeps); ok {
			v = deps.resolve()
		} else {
			v = fmt.Sprintf("%+v", v)
		}
		result[k] = v
	}
	return result
}

func (deps *DefaultDeps) String() string {
	b, err := json.MarshalIndent(deps.resolve(), "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (deps *DefaultDeps) Get(name string) any {
	deps.RWMutex.RLock()
	defer deps.RWMutex.RUnlock()
	d := deps
	for {
		if value, ok := d.data[name]; ok {
			return value
		}
		if d.parent != nil {
			d = d.parent
		} else {
			return nil
		}
	}

}
func (deps *DefaultDeps) Set(name string, dep any) {
	deps.RWMutex.Lock()
	defer deps.RWMutex.Unlock()
	deps.data[name] = dep
}

func (deps *DefaultDeps) GetContext(name string) Deps {
	deps.RWMutex.RLock()
	defer deps.RWMutex.RUnlock()
	if ctx, ok := deps.data[name]; ok {
		if ctx, ok := ctx.(*DefaultDeps); ok {
			return ctx
		}
	}
	return deps
}

func (deps *DefaultDeps) SetContext(name string) Deps {
	deps.RWMutex.Lock()
	defer deps.RWMutex.Unlock()
	if ctx, ok := deps.data[name]; ok {
		if ctx, ok := ctx.(*DefaultDeps); ok {
			return ctx
		}
	}
	ctx := &DefaultDeps{
		id:      uuid.NewString(),
		parent:  deps,
		data:    make(map[string]any),
		RWMutex: sync.RWMutex{},
	}
	deps.data[name] = ctx
	return ctx
}
