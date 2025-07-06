package sylva

import "fmt"

type ProgramLocation struct {
	Line   int
	Column int
}

type CommandDebugData struct {
	FileLocation        string
	FunctionDeclaration struct {
		FunctionName       string
		FunctionLineNumber int
	}
	Start ProgramLocation
	End   ProgramLocation
}

func (l ProgramLocation) StringRepresentation() string {
	return fmt.Sprintf("(%v, %v)", l.Line, l.Column)
}

func (d CommandDebugData) StringRepresentation() string {
	return fmt.Sprintf(
		"[\"%v\", (\"%v\", %v), %v, %v]",
		d.FileLocation,
		d.FunctionDeclaration.FunctionName,
		d.FunctionDeclaration.FunctionLineNumber,
		d.Start.StringRepresentation(),
		d.End.StringRepresentation(),
	)
}

type Command interface {
	GetDebugData() CommandDebugData
	StringRepresentation() string
	GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error)
}

type LoadCommand struct {
	DebugData CommandDebugData
	Register  string
	Value     string
}

func (c *LoadCommand) GetDebugData() CommandDebugData {
	return c.DebugData
}

func (c *LoadCommand) StringRepresentation() string {
	return fmt.Sprintf(
		"load %v, %v %v",
		c.Register,
		c.Value,
		c.DebugData.StringRepresentation(),
	)
}

const (
	NOOP        = 0
	LOAD_REG    = 1
	LOAD_INT    = 2
	LOAD_FLOAT  = 3
	LOAD_STRING = 4
	LOAD_TRUE   = 5
	LOAD_FALSE  = 6
	LOAD_NIL    = 7
	LIST        = 8

	FREE = 16

	ADD    = 32
	SUB    = 33
	MUL    = 34
	DIV    = 35
	MOD    = 36
	CONCAT = 37

	UMN = 42
)

func (c *LoadCommand) GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error) {
	return func(runtime *SylvaRuntime) ([]uint64, error) {
		reg := runtime.GetRegisterID(c.Register)
		data, err := runtime.GetByteCodeForValue(c.Value)
		if err != nil {
			return nil, err
		}

		var command uint64
		switch data.Type {
		case BYTECODE_VALUE_TYPE_NIL:
			command = LOAD_NIL
		case BYTECODE_VALUE_TYPE_REGISTER:
			command = LOAD_REG
		case BYTECODE_VALUE_TYPE_STRING:
			command = LOAD_STRING
		case BYTECODE_VALUE_TYPE_INT:
			command = LOAD_INT
		case BYTECODE_VALUE_TYPE_FLOAT:
			command = LOAD_FLOAT
		case BYTECODE_VALUE_TYPE_FALSE:
			command = LOAD_FALSE
		case BYTECODE_VALUE_TYPE_TRUE:
			command = LOAD_TRUE
		default:
			command = NOOP
		}

		byteCode := []uint64{command, reg}
		byteCode = append(byteCode, data.Code...)
		return byteCode, nil
	}
}

type FreeCommand struct {
	DebugData CommandDebugData
	Register  string
}

func (c *FreeCommand) GetDebugData() CommandDebugData {
	return c.DebugData
}

func (c *FreeCommand) StringRepresentation() string {
	return fmt.Sprintf(
		"free %v %v",
		c.Register,
		c.DebugData.StringRepresentation(),
	)
}

func (c *FreeCommand) GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error) {
	return func(runtime *SylvaRuntime) ([]uint64, error) {
		reg := runtime.GetRegisterID(c.Register)
		return []uint64{FREE, reg}, nil
	}
}

const (
	OP_ADD    = "add"
	OP_SUB    = "sub"
	OP_MUL    = "mul"
	OP_DIV    = "div"
	OP_MOD    = "mod"
	OP_CONCAT = "concat"
)

type BinOpCommand struct {
	DebugData CommandDebugData
	Register  string
	A         string
	B         string
	Operation string
}

func (c *BinOpCommand) GetDebugData() CommandDebugData {
	return c.DebugData
}

func (c *BinOpCommand) StringRepresentation() string {
	return fmt.Sprintf(
		"%v %v, %v, %v %v",
		c.Operation,
		c.Register,
		c.A,
		c.B,
		c.DebugData.StringRepresentation(),
	)
}

func (c *BinOpCommand) GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error) {
	return func(runtime *SylvaRuntime) ([]uint64, error) {
		reg := runtime.GetRegisterID(c.Register)
		a := runtime.GetRegisterID(c.A)
		b := runtime.GetRegisterID(c.B)
		var command uint64
		switch c.Operation {
		case "add":
			command = ADD
		case "sub":
			command = SUB
		case "mul":
			command = MUL
		case "div":
			command = DIV
		case "mod":
			command = MOD
		case "concat":
			command = CONCAT
		default:
			return nil, fmt.Errorf("unknown binop command: %v", c.Operation)
		}
		return []uint64{command, reg, a, b}, nil
	}
}

type ListCommand struct {
	DebugData CommandDebugData
	Register  string
}

func (c *ListCommand) GetDebugData() CommandDebugData {
	return c.DebugData
}

func (c *ListCommand) StringRepresentation() string {
	return fmt.Sprintf(
		"list %v %v",
		c.Register,
		c.DebugData.StringRepresentation(),
	)
}

func (c *ListCommand) GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error) {
	return func(runtime *SylvaRuntime) ([]uint64, error) {
		reg := runtime.GetRegisterID(c.Register)
		return []uint64{LIST, reg}, nil
	}
}

type UmnCommand struct {
	DebugData CommandDebugData
	Register  string
	Value     string
}

func (c *UmnCommand) GetDebugData() CommandDebugData {
	return c.DebugData
}

func (c *UmnCommand) StringRepresentation() string {
	return fmt.Sprintf(
		"umn %v, %v %v",
		c.Register,
		c.Value,
		c.DebugData.StringRepresentation(),
	)
}

func (c *UmnCommand) GetByteCode() func(runtime *SylvaRuntime) ([]uint64, error) {
	return func(runtime *SylvaRuntime) ([]uint64, error) {
		reg := runtime.GetRegisterID(c.Register)
		val := runtime.GetRegisterID(c.Value)
		return []uint64{UMN, reg, val}, nil
	}
}
