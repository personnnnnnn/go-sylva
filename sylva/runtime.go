package sylva

import (
	"math"
	"strconv"
	"sylva/util"
)

func EncodeFloat64ToUInt64(f float64) uint64 {
	return uint64(math.Float64bits(f))
}

func DecodeUInt64ToFloat64(i uint64) float64 {
	return math.Float64frombits(i)
}

type SylvaRuntime struct {
	Registers          map[uint64]Value
	RegisterNames      map[string]uint64
	Bytecode           []uint64
	IP                 int
	topRegisterID      uint64
	BytecodeToCommands map[int]int
	CommandDebugData   map[int]CommandDebugData
}

func CreateSylvaRuntime() *SylvaRuntime {
	return &SylvaRuntime{
		IP:                 0,
		Registers:          map[uint64]Value{},
		RegisterNames:      map[string]uint64{},
		Bytecode:           []uint64{},
		topRegisterID:      0,
		BytecodeToCommands: map[int]int{},
		CommandDebugData:   map[int]CommandDebugData{},
	}
}

func (runtime *SylvaRuntime) GetRegisterID(name string) uint64 {
	if id, ok := runtime.RegisterNames[name]; ok {
		return id
	} else {
		id := runtime.topRegisterID
		runtime.topRegisterID++
		runtime.RegisterNames[name] = id
		return id
	}
}

type BytecodeValue struct {
	Code []uint64
	Type int
}

const (
	BYTECODE_VALUE_TYPE_REGISTER = 0
	BYTECODE_VALUE_TYPE_FALSE    = 1
	BYTECODE_VALUE_TYPE_TRUE     = 2
	BYTECODE_VALUE_TYPE_NIL      = 3
	BYTECODE_VALUE_TYPE_INT      = 4
	BYTECODE_VALUE_TYPE_FLOAT    = 5
	BYTECODE_VALUE_TYPE_STRING   = 6
)

func IsAlphaNumeric(s string) bool {
	for _, v := range s {
		if !(v == '1' ||
			v == '2' ||
			v == '3' ||
			v == '4' ||
			v == '5' ||
			v == '6' ||
			v == '7' ||
			v == '8' ||
			v == '9' ||
			v == '0' ||
			v == '.') {
			return false
		}
	}
	return true
}

func ContainsDot(s string) bool {
	for _, v := range s {
		if v == '.' {
			return true
		}
	}
	return false
}

func (runtime *SylvaRuntime) GetByteCodeForValue(value string) (BytecodeValue, error) {
	if IsAlphaNumeric(value) {
		if ContainsDot(value) { // float
			f, _ := strconv.ParseFloat(value, 64) // will never fail since there are only digits and dots
			encoding := EncodeFloat64ToUInt64(f)
			return BytecodeValue{
				Type: BYTECODE_VALUE_TYPE_FLOAT,
				Code: []uint64{encoding},
			}, nil
		} else { // int
			n, _ := strconv.Atoi(value) // will never fail since there are only digits
			return BytecodeValue{
				Type: BYTECODE_VALUE_TYPE_INT,
				Code: []uint64{uint64(n)},
			}, nil
		}
	}

	switch value[0] {
	case '$': // register
		otherReg := runtime.GetRegisterID(value)
		return BytecodeValue{Code: []uint64{otherReg}, Type: BYTECODE_VALUE_TYPE_REGISTER}, nil
	case '"': // string
		// remove quotes
		proccesedString := value[1:]
		proccesedString = proccesedString[:len(proccesedString)-1]

		proccesedString, err := util.ParseString(proccesedString)
		if err != nil {
			return BytecodeValue{}, err
		}

		proccesedStringLen := uint64(len(proccesedString))
		byteCode := []uint64{proccesedStringLen}

		for _, char := range proccesedString {
			charInt := uint64(char)
			byteCode = append(byteCode, charInt)
		}
		return BytecodeValue{Code: byteCode, Type: BYTECODE_VALUE_TYPE_STRING}, nil
	case 't': // true boolean
		return BytecodeValue{
			Code: []uint64{},
			Type: BYTECODE_VALUE_TYPE_TRUE,
		}, nil
	case 'f': // false boolean
		return BytecodeValue{
			Code: []uint64{},
			Type: BYTECODE_VALUE_TYPE_FALSE,
		}, nil
	}

	return BytecodeValue{ // nil
		Code: []uint64{},
		Type: BYTECODE_VALUE_TYPE_NIL,
	}, nil
}

func (runtime *SylvaRuntime) IsDone() bool {
	return runtime.IP >= len(runtime.Bytecode)
}

func (runtime *SylvaRuntime) ReadBytecode() uint64 {
	runtime.Inc()
	return runtime.Bytecode[runtime.IP-1]
}

func (runtime *SylvaRuntime) Inc() {
	runtime.IP++
}
func (runtime *SylvaRuntime) Step() error {
	if runtime.IsDone() {
		return nil
	}

	command := runtime.ReadBytecode()
	switch command {
	case NOOP:
	case LOAD_NIL:
		reg := runtime.ReadBytecode()
		runtime.Registers[reg] = nil
	case LOAD_TRUE:
		reg := runtime.ReadBytecode()
		runtime.Registers[reg] = true
	case LOAD_FALSE:
		reg := runtime.ReadBytecode()
		runtime.Registers[reg] = false
	case LOAD_INT:
		reg := runtime.ReadBytecode()
		val := runtime.ReadBytecode()
		i := int64(val)
		runtime.Registers[reg] = i
	case LOAD_FLOAT:
		reg := runtime.ReadBytecode()
		val := runtime.ReadBytecode()
		f := DecodeUInt64ToFloat64(val)
		runtime.Registers[reg] = f
	case LOAD_REG:
		regA := runtime.ReadBytecode()
		regB := runtime.ReadBytecode()
		runtime.Registers[regA] = runtime.Registers[regB]
	case LOAD_STRING:
		reg := runtime.ReadBytecode()
		strLen := runtime.ReadBytecode()
		str := make([]rune, 0, strLen)
		for range strLen {
			char := runtime.ReadBytecode()
			str = append(str, rune(char))
		}
		runtime.Registers[reg] = string(str)
	case FREE:
		reg := runtime.ReadBytecode()
		delete(runtime.Registers, reg)
	case ADD:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res, err := Add(a, b)
		if err != nil {
			return err
		}
		runtime.Registers[reg] = res
	case SUB:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res, err := Sub(a, b)
		if err != nil {
			return err
		}
		runtime.Registers[reg] = res
	case MUL:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res, err := Mul(a, b)
		if err != nil {
			return err
		}
		runtime.Registers[reg] = res
	case DIV:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res, err := Div(a, b)
		if err != nil {
			return err
		}
		runtime.Registers[reg] = res
	case MOD:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res, err := Mod(a, b)
		if err != nil {
			return err
		}
		runtime.Registers[reg] = res
	case CONCAT:
		reg := runtime.ReadBytecode()
		aReg := runtime.ReadBytecode()
		bReg := runtime.ReadBytecode()
		a := runtime.Registers[aReg]
		b := runtime.Registers[bReg]
		res := Concat(a, b)
		runtime.Registers[reg] = res
	}

	return nil
}

func (runtime *SylvaRuntime) ExecuteUntilDone() error {
	runtime.IP = 0
	for !runtime.IsDone() {
		// oldIP := runtime.IP
		err := runtime.Step()
		if err == nil {
			continue
		}
		// TODO: use the given debug data to give a better error messages
		return err
	}
	return nil
}

func (runtime *SylvaRuntime) ConvertToBytecode(commands []Command) error {
	for i, command := range commands {
		debug := command.GetDebugData()
		byteCode, err := command.GetByteCode()(runtime)
		if err != nil {
			return err
		}
		byteCodeIndex := len(runtime.Bytecode)
		for j := range len(byteCode) {
			runtime.BytecodeToCommands[byteCodeIndex+j] = i
		}
		runtime.CommandDebugData[i] = debug
		runtime.Bytecode = append(runtime.Bytecode, byteCode...)
	}
	return nil
}
