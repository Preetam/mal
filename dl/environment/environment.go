package environment

import (
	"errors"

	"github.com/Preetam/mal/dl/types"
)

var (
	ErrNotFound = errors.New("environment: not found")
)

type Env struct {
	outer  *Env
	values map[types.MalSymbol]types.MalType
}

func New(outer *Env) *Env {
	return &Env{
		outer:  outer,
		values: map[types.MalSymbol]types.MalType{},
	}
}

func (e *Env) Set(key types.MalSymbol, value types.MalType) {
	e.values[key] = value
}

func (e *Env) Find(key types.MalSymbol) *Env {
	if _, ok := e.values[key]; ok {
		return e
	}
	if e.outer != nil {
		return e.outer.Find(key)
	}
	return nil
}

func (e *Env) Get(key types.MalSymbol) (types.MalType, error) {
	env := e.Find(key)
	if env == nil {
		return nil, ErrNotFound
	}
	return env.values[key], nil
}
