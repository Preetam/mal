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

func New(outer *Env, binds, exprs types.MalType) (*Env, error) {
	env := &Env{
		outer:  outer,
		values: map[types.MalSymbol]types.MalType{},
	}

	if binds != nil {
		bindsList, ok := binds.(types.MalList)
		if !ok {
			return nil, errors.New("binds is not a list")
		}
		exprsList, ok := exprs.(types.MalList)
		if !ok {
			return nil, errors.New("exprs is not a list")
		}
		if len(bindsList) != len(exprsList) {
			return nil, errors.New("binds list length does not equal exprs list length")
		}

		for i, bind := range bindsList {
			env.Set(bind.(types.MalSymbol), exprsList[i])
		}
	}

	return env, nil
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
