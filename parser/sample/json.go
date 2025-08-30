package sample

import (
	"fmt"
	"iter"
	"reflect"
	"strconv"
	"strings"

	"github.com/fuwjax/gopase/parser"
)

const jsonGrammar = `
Value = WS (String / Object / Array / Number / Literal) WS
Object = "{" (WS String WS ":" Value ("," WS String WS ":" Value)* / WS) "}"
Array = "[" (Value ("," Value)* / WS) "]"
String = '"' ("\\u" Hex / "\\" Escape / Plain)* '"'
Number = "-"? ("0" / [1-9][0-9]*) ("." [0-9]+)? ([eE][+-]?[0-9]+)?
Literal = "true" / "false" / "null"
Plain = [^\\"]+
Escape = [/\\"bfnrt]
Hex = [0-9a-fA-F] [0-9a-fA-F] [0-9a-fA-F] [0-9a-fA-F]
WS = [ \r\n\t]*
`

var JsonGrammar = parser.Preserve2(parser.Bootstrap(jsonGrammar))

func ParseJson(input string) (any, error) {
	return ParseJsonFrom("Value", input)
}

func ParseJsonFrom(root, input string) (any, error) {
	grammar, err := JsonGrammar()
	if err != nil {
		return nil, err
	}
	result, err := parser.Parse(root, grammar, parser.WrapHandler(jsonHandler{}), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ConvertJson[T any](data any) (T, error) {
	value, err := ConvertJsonValue(data, reflect.TypeFor[T]())
	var t T
	if err != nil || value == nil {
		return t, err
	}
	return value.(T), err
}

func ConvertJsonValue(data any, target reflect.Type) (any, error) {
	source := reflect.TypeOf(data)
	if source == nil {
		return reflect.Zero(target), nil
	}
	value := reflect.ValueOf(data)
	if value.CanConvert(target) {
		return value.Convert(target).Interface(), nil
	}
	switch source.Kind() {
	case reflect.Map:
		return ConvertJsonObject(data.(map[string]any), target)
	case reflect.Slice:
		return ConvertJsonArray(data.([]any), target)
		// strings, numbers, and bools should have converted
	}
	return nil, fmt.Errorf("cannot convert json %v to type %v", source, target)
}

func ConvertJsonObject(data map[string]any, target reflect.Type) (any, error) {
	switch target.Kind() {
	case reflect.Map:
		dest := reflect.MakeMap(target)
		for key, value := range data {
			elem, err := ConvertJsonValue(value, target.Elem())
			if err != nil {
				return nil, err
			}
			dest.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(reflect.ValueOf(elem)))
		}
		return dest.Interface(), nil
	case reflect.Struct:
		dest := reflect.New(target).Elem()
		for i := range target.NumField() {
			field := target.Field(i)
			jsonName := field.Tag.Get("jsonName")
			if len(jsonName) == 0 {
				jsonName = field.Name
			}
			value, ok := data[jsonName]
			if ok {
				elem, err := ConvertJsonValue(value, field.Type)
				if err != nil {
					return nil, err
				}
				dest.Field(i).Set(reflect.ValueOf(elem))
			}
		}
		return dest.Interface(), nil
	}
	return nil, fmt.Errorf("cannot convert json object to %v", target)
}

func ConvertJsonArray(data []any, target reflect.Type) (any, error) {
	switch target.Kind() {
	case reflect.Slice:
		dest := reflect.MakeSlice(target, len(data), cap(data))
		for i, value := range data {
			elem, err := ConvertJsonValue(value, target.Elem())
			if err != nil {
				return nil, err
			}
			dest.Index(i).Set(reflect.ValueOf(elem))
		}
		return dest.Interface(), nil
	}
	return nil, fmt.Errorf("cannot convert json array to %v", target)
}

type jsonHandler struct{}

func (h jsonHandler) Value(results iter.Seq2[string, any]) (any, error) {
	name, value := parser.FirstOf(results, "String", "Object", "Array", "Number", "Literal")
	switch name {
	case "Number":
		return strconv.ParseFloat(value.(string), 64)
	case "Literal":
		switch value.(string) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "null":
			return nil, nil
		}
	}
	return value, nil
}

func (h jsonHandler) Object(results iter.Seq2[string, any]) (any, error) {
	var key string
	obj := make(map[string]any)
	for name, result := range results {
		switch name {
		case "String":
			key = result.(string)
		case "Value":
			obj[key] = result
		}
	}
	return obj, nil
}

func (h jsonHandler) Array(results iter.Seq2[string, any]) (any, error) {
	return parser.ListOf(results, "Value"), nil
}

func (h jsonHandler) String(results iter.Seq2[string, any]) (any, error) {
	var sb strings.Builder
	for name, result := range results {
		switch name {
		case "Plain":
			sb.WriteString(result.(string))
		case "Hex":
			i, err := strconv.ParseInt(result.(string), 16, 32)
			if err != nil {
				return nil, err
			}
			sb.WriteRune(rune(i))
		case "Escape":
			switch result.(string) {
			case "b":
				sb.WriteString("\b")
			case "f":
				sb.WriteString("\f")
			case "n":
				sb.WriteString("\n")
			case "r":
				sb.WriteString("\r")
			case "t":
				sb.WriteString("\t")
			default:
				sb.WriteString(result.(string))
			}
		}
	}
	return sb.String(), nil
}
