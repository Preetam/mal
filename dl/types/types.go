package types

type MalType interface {
	IsMalType()
}

type MalList []MalType

func (l MalList) IsMalType() {}

type MalInt int64

func (n MalInt) IsMalType() {}

type MalSymbol string

func (s MalSymbol) IsMalType() {}

type MalFunction func(args ...MalType) (MalType, error)

func (f MalFunction) IsMalType() {}

type MalBool bool

func (b MalBool) IsMalType() {}
