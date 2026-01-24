package main

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// Runtime helpers

func gts_len(v interface{}) int {
	switch x := v.(type) {
	case string:
		return len(x)
	case []interface{}:
		return len(x)
	case []int:
		return len(x)
	case []float64:
		return len(x)
	case []string:
		return len(x)
	case []bool:
		return len(x)
	default:
		return 0
	}
}

func gts_typeof(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "object"
	}
}

func gts_tostring(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func gts_toint(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		var n int
		fmt.Sscanf(x, "%d", &n)
		return n
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func gts_tofloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case string:
		var n float64
		fmt.Sscanf(x, "%f", &n)
		return n
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func gts_call(fn interface{}, args ...interface{}) interface{} {
	v := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	fnType := v.Type()
	for i, arg := range args {
		if i < fnType.NumIn() {
			// Convert argument to expected type
			expectedType := fnType.In(i)
			argVal := reflect.ValueOf(arg)
			if argVal.Type().ConvertibleTo(expectedType) {
				in[i] = argVal.Convert(expectedType)
			} else {
				in[i] = argVal
			}
		} else {
			in[i] = reflect.ValueOf(arg)
		}
	}
	out := v.Call(in)
	if len(out) > 0 {
		return out[0].Interface()
	}
	return nil
}

func gts_tobool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}

func gts_toarr_float(v []interface{}) []float64 {
	result := make([]float64, len(v))
	for i, x := range v {
		result[i] = gts_tofloat(x)
	}
	return result
}

func gts_toarr_int(v []interface{}) []int {
	result := make([]int, len(v))
	for i, x := range v {
		result[i] = gts_toint(x)
	}
	return result
}

func main() {
	var upper string = strings.ToUpper("hello world")
	var lower string = strings.ToLower("HELLO WORLD")
	fmt.Println(("Upper: " + upper))
	fmt.Println(("Lower: " + lower))
	var hasHello bool = strings.Contains("Hello World", "Hello")
	fmt.Println(("Contains Hello: " + gts_tostring(hasHello)))
	var parts []string = strings.Split("a,b,c,d", ",")
	fmt.Println("Split result:")
	for i := 0; i < gts_len(parts); i = (i + 1) {
		fmt.Println(("  " + parts[int(i)]))
	}
	var joined string = strings.Join(parts, "-")
	fmt.Println(("Joined: " + joined))
	var sqrtVal float64 = math.Sqrt(16)
	fmt.Println(("Sqrt(16) = " + gts_tostring(sqrtVal)))
	var powVal float64 = math.Pow(2, 10)
	fmt.Println(("Pow(2, 10) = " + gts_tostring(powVal)))
	var maxVal float64 = math.Max(3.14, 2.71)
	var minVal float64 = math.Min(3.14, 2.71)
	fmt.Println(("Max(3.14, 2.71) = " + gts_tostring(maxVal)))
	fmt.Println(("Min(3.14, 2.71) = " + gts_tostring(minVal)))
}
