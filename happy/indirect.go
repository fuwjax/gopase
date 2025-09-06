package happy

import (
	"fmt"
	"iter"
	"reflect"
	"strconv"

	"github.com/fuwjax/gopase/funki"
)

var zero reflect.Value

/*
Attempts a direct cast and standard get. Currently supports map[string]any and []any, which are both
returned by sample.ParseJson.
*/
func getDirect(data any, key string) (any, bool) {
	mapping, ok := data.(map[string]any)
	if ok {
		result, ok := mapping[key]
		return result, ok
	}
	slice, ok := data.([]any)
	if ok {
		index, err := strconv.ParseInt(key, 0, 0)
		if err != nil {
			return nil, false
		}
		i := int(index)
		if i < 0 || i >= len(slice) {
			return nil, false
		}
		return slice[i], true
	}
	return nil, false
}

func getMap(value reflect.Value, name string) reflect.Value {
	key := reflect.ValueOf(name)
	return value.MapIndex(key)
}

func getStruct(value reflect.Value, name string) reflect.Value {
	result := value.FieldByName(name)
	if result.IsValid() {
		return result
	}
	return value.MethodByName(name)
}

func getPointerOrInterface(value reflect.Value, name string) reflect.Value {
	result := value.MethodByName(name)
	if result.IsValid() {
		return result
	}
	return getIndirect(value.Elem(), name)
}

func getSliceOrArray(value reflect.Value, name string) reflect.Value {
	index, err := strconv.ParseInt(name, 0, 0)
	if err != nil {
		return zero
	}
	i := int(index)
	if i < 0 || i >= value.Len() {
		return zero
	}
	return value.Index(i)
}

func getIndirect(value reflect.Value, key string) reflect.Value {
	if !value.IsValid() {
		return zero
	}
	switch value.Kind() {
	case reflect.Map:
		return getMap(value, key)
	case reflect.Struct:
		return getStruct(value, key)
	case reflect.Interface, reflect.Pointer:
		return getPointerOrInterface(value, key)
	case reflect.Slice, reflect.Array:
		return getSliceOrArray(value, key)
	// case reflect.Chan, reflect.Func: must go through Call not Get
	default:
		return zero
	}
}

/*
Gets an element/property/method/value from data by key. Scalars (bools/numbers/strings) are explicitly not
supported. No default methods or operators are implemented. Get is undefined for channels and functions.
*/
func Get(data any, key string) (any, bool) {
	elem, ok := getDirect(data, key)
	if ok {
		return elem, ok
	}
	return devalue(getIndirect(reflect.ValueOf(data), key))
}

func Truthy(data any) bool {
	if data == nil {
		return false
	}
	return truthyIndirect(reflect.ValueOf(data))
}

func truthyIndirect(value reflect.Value) bool {
	if !value.IsValid() { // trust that IsZero catches nil?
		return false
	}
	switch value.Kind() {
	case reflect.Pointer, reflect.Interface:
		return !value.IsNil() && truthyIndirect(value.Elem())
	case reflect.Array, reflect.Map, reflect.Slice:
		return value.Len() > 0
		//case reflect.Func, reflect.Chan: I guess we have to just see where it goes?
	case reflect.Struct:
		return true
	default:
		return !value.IsZero()
	}
}

func devalue(result reflect.Value) (any, bool) {
	if !result.IsValid() {
		return nil, false
	}
	if result.CanInterface() {
		return result.Interface(), true
	}
	if result.CanAddr() {
		return devalue(reflect.NewAt(result.Type(), result.Addr().UnsafePointer()).Elem())
	}
	return nil, false
}

func callResults(results []reflect.Value) (any, bool) {
	if len(results) == 1 {
		return devalue(results[0])
	}
	if len(results) == 2 {
		if results[1].Type().Implements(reflect.TypeFor[error]()) && results[1].IsNil() {
			return devalue(results[0])
		}
		if results[1].Kind() == reflect.Bool && results[1].Bool() {
			return devalue(results[0])
		}
	}
	return nil, false
}

/*
Attempts to "call" data with the args. For methods and functions, this is the normal meaning of call.
For anything else, a "call" with no args is treated as an identity, and returns the data itself. Otherwise,
an attempt is made to Get from data with the first arg, and then Call the resulting value with the remaining
args. This recursion will continue. It is not currently possible to use a subset of the args for a normal
call. Channels only support receive, and only with zero args. Functions must return a single value,
a value-error tuple, or a value-bool tuple. The result's second element must be nil or true, or the Call
will return a nil-false.
*/
func Call(data any, args []any) (any, bool) {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Func {
		//check value.Type().NumIn()/.In() so we don't panic
		results := value.Call(funki.Apply(args, reflect.ValueOf))
		return callResults(results)
	}
	if len(args) == 0 {
		if value.Kind() == reflect.Chan {
			result, ok := value.TryRecv()
			if ok {
				return devalue(result)
			}
			return nil, false
		}
		return data, true
	}
	name, ok := String(args[0])
	if !ok || name == "" {
		return nil, false
	}
	result, ok := Get(data, name)
	if !ok || result == nil {
		return nil, false
	}
	return Call(result, args[1:])
}

func iterDirect(data any) (iter.Seq2[any, any], bool) {
	mapping, ok := data.(map[string]any)
	if ok {
		return func(yield func(any, any) bool) {
			for key, data := range mapping {
				if !yield(key, data) {
					return
				}
			}
		}, true
	}
	slice, ok := data.([]any)
	if ok {
		return func(yield func(any, any) bool) {
			for index, data := range slice {
				if !yield(index, data) {
					return
				}
			}
		}, true
	}
	return nil, false
}

func Iter(data any) (iter.Seq2[any, any], bool) {
	iter, ok := iterDirect(data)
	if ok {
		return iter, ok
	}
	return iterIndirect(reflect.ValueOf(data))
}

func iterIndirect(value reflect.Value) (iter.Seq2[any, any], bool) {
	switch value.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Func: //.Func handles iter.Seq2
		seq := value.Seq2()
		return func(yield func(any, any) bool) {
			seq(func(v1, v2 reflect.Value) bool {
				r1, ok1 := devalue(v1)
				r2, ok2 := devalue(v2)
				return ok1 && ok2 && yield(r1, r2)
			})
		}, true
		/*
		   // we are explicitly not iterating over struct fields. This would break the expected behavior of (^*^)
		   	case reflect.Struct:
		   		t := value.Type()
		   		return func(yield func(any, any) bool) {
		   			for i := 0; i < t.NumField(); i++ {
		   				field := t.Field(i)
		   				result, ok := devalue(value.Field(i))
		   				if ok {
		   					if !yield(field.Name, result) {
		   						return
		   					}
		   				}
		   			}
		   		}, true
		*/
	case reflect.Interface, reflect.Pointer:
		return iterIndirect(value.Elem())
	// really not sure about reflect.Chan
	default:
		return nil, false
	}
}

func String(data any) (string, bool) {
	switch d := data.(type) {
	case nil:
		return "", false
	case string:
		return d, true
	default:
		return fmt.Sprint(data), true
	}
}
