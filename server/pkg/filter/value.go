package filter

import "fmt"

type Value interface {
	// Type returns the type for the value.
	Type() Type
	// Returns true if the value is NULL.
	IsNull() bool
	// Value returns the value with the correct Go type as an interface{}.
	Value() interface{}
	// Formats the value as a string.
	String() string
}

func String(s string) StringValue {
	return StringValue(s)
}

type StringValue string

func (s StringValue) Type() Type {
	return StringType
}

func (s StringValue) IsNull() bool {
	return false
}

func (s StringValue) Value() interface{} {
	return string(s)
}

func (s StringValue) String() string {
	return fmt.Sprintf(`"%s"`, string(s))
}

func Bool(v bool) BoolValue {
	return BoolValue(v)
}

type BoolValue bool

func (b BoolValue) Type() Type {
	return BoolType
}

func (b BoolValue) IsNull() bool {
	return false
}

func (b BoolValue) Value() interface{} {
	return bool(b)
}

func (b BoolValue) String() string {
	return fmt.Sprint(bool(b))
}

func Null() NullValue {
	return NullValue{}
}

type NullValue struct{}

func (b NullValue) Type() Type {
	return NullType
}

func (b NullValue) IsNull() bool {
	return true
}

func (b NullValue) Value() interface{} {
	return nil
}

func (b NullValue) String() string {
	return "null"
}

func Int(v int64) IntValue {
	return IntValue(v)
}

type IntValue int64

func (s IntValue) Type() Type {
	return IntType
}

func (s IntValue) IsNull() bool {
	return false
}

func (s IntValue) Value() interface{} {
	return int64(s)
}

func (s IntValue) String() string {
	return fmt.Sprint(int64(s))
}

func Float(v float64) FloatValue {
	return FloatValue(v)
}

type FloatValue float64

func (s FloatValue) Type() Type {
	return FloatType
}

func (s FloatValue) IsNull() bool {
	return false
}

func (s FloatValue) Value() interface{} {
	return float64(s)
}

func (s FloatValue) String() string {
	return fmt.Sprint(float64(s))
}
