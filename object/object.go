package object

import (
	"compiler/ast"
	"strconv"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	String() string
}

type IntegerObject struct {
	Value int64
}

func (intObj *IntegerObject) Type() ObjectType { return "INT" }
func (intObj *IntegerObject) String() string   { return strconv.FormatInt(intObj.Value, 10) }

type FunctionObject struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (funcObj *FunctionObject) Type() ObjectType { return "FUNCTION" }
func (funcObj *FunctionObject) String() string   { return "" }

type ReturnObject struct {
	ReturnValue Object
}

func (returnObj *ReturnObject) Type() ObjectType { return "RETURN_OBJ" }
func (returnObj *ReturnObject) String() string   { return "" }

type BooleanObject struct {
	Value bool
}

func (boolObj *BooleanObject) Type() ObjectType { return "BOOLEAN" }
func (boolObj *BooleanObject) String() string   { return strconv.FormatBool(boolObj.Value) }
