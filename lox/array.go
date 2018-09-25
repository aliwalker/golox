package lox

import (
	"strconv"
)

// LoxArray is the runtime object for lox array.
var LoxArray *LoxClass

// NOTE: this is a hack.
type _arrayInsType struct {
	*LoxInstance
}

func newArraryInsType(o *LoxInstance) *_arrayInsType {
	return &_arrayInsType{o}
}

// we gonna reimplement the general method.
func (o *_arrayInsType) String() string {
	stringified := "["
	list, _ := o.props["list"].([]interface{})
	for _, item := range list {
		var itemStr string
		switch val := item.(type) {
		case string:
			itemStr = val
		case int:
			itemStr = strconv.Itoa(val)
		case float64:
			itemStr = strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			if val == true {
				itemStr = "true"
			} else {
				itemStr = "false"
			}
		}
		stringified += itemStr + ", "
	}
	if len(list) > 0 {
		stringified = stringified[:len(stringified)-2]
	}
	//stringified = stringified[:len(stringified)-1]
	stringified += "]"
	return stringified
}

// init Array class. This function will be called when an Interpreter is instantiated.
func initArray() {
	// Array static methods
	var statics = map[string]Callable{
		"isArray": NewBuiltinFunc(1, func(i *LoxInstance, args ...interface{}) interface{} {
			obj, ok := args[0].(*LoxInstance)
			if ok != true {
				return false
			}
			if obj.class.Name != "Array" {
				return false
			}
			return true
		}),
	}

	// instance methods
	var methods = map[string]Callable{
		// We mark the arity to be -1, means we accept inifinite args.
		"init": NewBuiltinFunc(-1, func(i *LoxInstance, args ...interface{}) interface{} {
			argsLen := len(args)
			list := []interface{}{}
			if argsLen != 0 {
				for _, obj := range args {
					list = append(list, obj)
				}
			}
			i.props["list"] = list
			return newArraryInsType(i)
		}),
		"append": NewBuiltinFunc(1, func(i *LoxInstance, args ...interface{}) interface{} {
			list, _ := i.props["list"].([]interface{})
			list = append(list, args[0])
			i.props["list"] = list
			return args[0]
		}),
		"pop": NewBuiltinFunc(0, func(i *LoxInstance, args ...interface{}) interface{} {
			list, _ := i.props["list"].([]interface{})
			returned := list[len(list)-1]
			list = list[:len(list)]
			i.props["list"] = list
			return returned
		}),
	}

	// instance getters
	var getters = map[string]Callable{
		"length": NewBuiltinFunc(0, func(i *LoxInstance, args ...interface{}) interface{} {
			list, _ := i.props["list"].([]interface{})
			return len(list)
		}),
	}

	LoxArray = NewLoxClass("Array", nil, statics, methods, getters, nil)
}
