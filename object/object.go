package object

import (
	"bytes"
	"compiler/ast"
	"strconv"
)

const (
	INT        = "INT"
	FUNCTION   = "FUNCTION"
	RETURN_OBJ = "RETURN_OBJ"
	BOOLEAN    = "BOOLEAN"
	ARRAY      = "ARRAY"
	STRING     = "STRING"
	ERROR      = "ERROR"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	String() string
}

type Integer struct {
	Value int64
}

func (intObj *Integer) Type() ObjectType { return INT }
func (intObj *Integer) String() string   { return strconv.FormatInt(intObj.Value, 10) }

type Function struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (funcObj *Function) Type() ObjectType { return FUNCTION }
func (funcObj *Function) String() string   { return "" }

type Return struct {
	ReturnValue Object
}

func (returnObj *Return) Type() ObjectType { return RETURN_OBJ }
func (returnObj *Return) String() string   { return "" }

type Boolean struct {
	Value bool
}

func (boolObj *Boolean) Type() ObjectType { return BOOLEAN }
func (boolObj *Boolean) String() string   { return strconv.FormatBool(boolObj.Value) }

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return "ARRAY" }
func (arr *Array) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, e := range arr.Elements {
		out.WriteString(e.String())
	}
	out.WriteString(")")

	return out.String()
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

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType {
	return ERROR
}
func (err *Error) String() string {
	return err.Message
}
