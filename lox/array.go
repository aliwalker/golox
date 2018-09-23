package lox

// LoxArray is the runtime object for lox array.
var LoxArray *LoxClass

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
			return i
		}),
		"append": NewBuiltinFunc(1, func(i *LoxInstance, args ...interface{}) interface{} {
			list, _ := i.props["list"].([]interface{})
			list = append(list, args[0])
			return args[0]
		}),
		"pop": NewBuiltinFunc(0, func(i *LoxInstance, args ...interface{}) interface{} {
			list, _ := i.props["list"].([]interface{})
			returned := list[len(list)-1]
			list = list[:len(list)]
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
