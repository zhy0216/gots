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

type Animal struct {
	Name string
	Age  int
}

func NewAnimal(name string, age int) *Animal {
	this := &Animal{}
	this.Name = name
	this.Age = age
	return this
}

func (this *Animal) Speak() {
	fmt.Println((this.Name + " makes a sound"))
}

type Dog struct {
	Animal
	Breed string
}

func NewDog(name string, age int, breed string) *Dog {
	this := &Dog{}
	this.Animal = *NewAnimal(name, age)
	this.Breed = breed
	return this
}

func (this *Dog) Speak() {
	fmt.Println((this.Name + " barks"))
}

func (this *Dog) Info() {
	fmt.Println(((((this.Name + " is a ") + gts_tostring(this.Age)) + " year old ") + this.Breed))
}

func main() {
	var dog *Dog = NewDog("Buddy", 5, "Golden Retriever")
	fmt.Println(dog.Name)
	fmt.Println(gts_tostring(dog.Age))
	fmt.Println(dog.Breed)
	dog.Speak()
	dog.Info()
	var animal *Animal = NewAnimal("Generic", 10)
	animal.Speak()
	fmt.Println("Done!")
}
