package object

import (
	"bytes"
	"compiler/ast"
	"compiler/code"
	"fmt"
	"strconv"
)

const (
	INT               = "INT"
	FUNCTION          = "FUNCTION"
	RETURN_OBJ        = "RETURN_OBJ"
	BOOLEAN           = "BOOLEAN"
	ARRAY             = "ARRAY"
	MAP               = "MAP"
	STRING            = "STRING"
	ERROR             = "ERROR"
	NULL              = "NULL"
	BUILTIN           = "BUILTIN"
	COMPILED_FUNCTION = "COMPILED_FUNCTION"
	CLOSURE           = "CLOSURE"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	String() string
}

type Hashable interface {
	Object
	Hash() string
}

type Integer struct {
	Value int64
}

func (intObj *Integer) Type() ObjectType { return INT }
func (intObj *Integer) String() string   { return strconv.FormatInt(intObj.Value, 10) }
func (intObj *Integer) Hash() string {
	return fmt.Sprintf("%s: %d", intObj.Type(), intObj.Value)
}

type Function struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (funcObj *Function) Type() ObjectType { return FUNCTION }
func (funcObj *Function) String() string   { return "NULL" }

type Return struct {
	ReturnValue Object
}

func (returnObj *Return) Type() ObjectType { return RETURN_OBJ }
func (returnObj *Return) String() string   { return "NULL" }

type Boolean struct {
	Value bool
}

func (boolObj *Boolean) Type() ObjectType { return BOOLEAN }
func (boolObj *Boolean) String() string   { return strconv.FormatBool(boolObj.Value) }
func (boolObj *Boolean) Hash() string {
	return fmt.Sprintf("%s: %t", boolObj.Type(), boolObj.Value)
}

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return ARRAY }
func (arr *Array) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, e := range arr.Elements {
		out.WriteString(e.String())
	}
	out.WriteString(")")

	return out.String()
}

type Map struct {
	Entries map[string]Object
}

func (mapObj *Map) Type() ObjectType { return MAP }
func (mapObj *Map) String() string {
	return "map"
}

type String struct {
	Value string
}

func (str *String) Type() ObjectType {
	return STRING
}
func (str *String) String() string {
	return str.Value
}
func (str *String) Hash() string {
	return fmt.Sprintf("%s: %s", str.Type(), str.Value)
}

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType {
	return ERROR
}
func (err *Error) String() string {
	return err.Message
}

type Null struct {
}

func (null *Null) Type() ObjectType {
	return NULL
}
func (null *Null) String() string {
	return "NULL"
}

type Builtin struct {
	Name string
	Fn   func(args ...Object) interface{}
}

func (builtin *Builtin) Type() ObjectType {
	return BUILTIN
}
func (builtin *Builtin) String() string {
	return ""
}

type CompiledFunction struct {
	Instructions code.Instructions

	NumParams int
	NumLocals int
}

func (compiledFn *CompiledFunction) Type() ObjectType {
	return COMPILED_FUNCTION
}
func (compiledFn *CompiledFunction) String() string {
	return ""
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType {
	return CLOSURE
}
func (c *Closure) String() string {
	return fmt.Sprintf("Closure[%p]", c)
}
