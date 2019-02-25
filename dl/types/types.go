package types

type MalType interface {
	IsMalType() string
}

type MalList []MalType

func (l MalList) IsMalType() string { return "List" }

type MalInt int64

func (n MalInt) IsMalType() string { return "Int" }

type MalSymbol string

func (s MalSymbol) IsMalType() string { return "Symbol" }

type MalFunction func(args ...MalType) (MalType, error)

func (f MalFunction) IsMalType() string { return "Function" }

type MalBool bool

func (b MalBool) IsMalType() string { return "Bool" }

type MalString string

func (s MalString) IsMalType() string { return "String" }
