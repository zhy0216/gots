package main

import (
	"fmt"
	"reflect"
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

type Person struct {
	Name string
}

func NewPerson(name string) *Person {
	this := &Person{}
	this.Name = name
	return this
}

func main() {
	var name *string = func() *string { v := "Alice"; return &v }()
	var result1 string = func() string {
		if name != nil {
			return *name
		}
		return "Default"
	}()
	fmt.Println(result1)
	var nullName *string = nil
	var result2 string = func() string {
		if nullName != nil {
			return *nullName
		}
		return "Default"
	}()
	fmt.Println(result2)
	var person *Person = NewPerson("Bob")
	var result3 *Person = func() *Person {
		if person != nil {
			return person
		}
		return NewPerson("Default")
	}()
	fmt.Println(result3.Name)
	var nullPerson *Person = nil
	var result4 *Person = func() *Person {
		if nullPerson != nil {
			return nullPerson
		}
		return NewPerson("Default")
	}()
	fmt.Println(result4.Name)
	var a *string = nil
	var b *string = nil
	var c *string = func() *string { v := "Charlie"; return &v }()
	var result5 string = func() string {
		if func() *string {
			if func() *string {
				if a != nil {
					return *a
				}
				return b
			}() != nil {
				return *func() *string {
					if a != nil {
						return *a
					}
					return b
				}()
			}
			return c
		}() != nil {
			return *func() *string {
				if func() *string {
					if a != nil {
						return *a
					}
					return b
				}() != nil {
					return *func() *string {
						if a != nil {
							return *a
						}
						return b
					}()
				}
				return c
			}()
		}
		return "None"
	}()
	fmt.Println(result5)
	fmt.Println("Done!")
}
