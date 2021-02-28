package filter

type Type string

const (
	NullType    Type = "null"
	IntType     Type = "int"
	FloatType   Type = "float"
	StringType  Type = "string"
	BoolType    Type = "bool"
)

func (t Type) Kind() Type {
	return t
}

func (t Type) Equal(t2 Type) bool {
	return t == t2
}

func (t Type) String() string {
	return string(t)
}
