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

type Drawable interface {
	Draw()
	GetArea() float64
}

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) *Circle {
	this := &Circle{}
	this.Radius = radius
	return this
}

func (this *Circle) Draw() {
	fmt.Println(("Drawing circle with radius " + gts_tostring(this.Radius)))
}

func (this *Circle) GetArea() float64 {
	return ((3.14159 * this.Radius) * this.Radius)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func NewRectangle(width float64, height float64) *Rectangle {
	this := &Rectangle{}
	this.Width = width
	this.Height = height
	return this
}

func (this *Rectangle) Draw() {
	fmt.Println(((("Drawing rectangle " + gts_tostring(this.Width)) + "x") + gts_tostring(this.Height)))
}

func (this *Rectangle) GetArea() float64 {
	return (this.Width * this.Height)
}

func main() {
	var circle *Circle = NewCircle(5)
	circle.Draw()
	fmt.Println(("Circle area: " + gts_tostring(circle.GetArea())))
	var rect *Rectangle = NewRectangle(4, 3)
	rect.Draw()
	fmt.Println(("Rectangle area: " + gts_tostring(rect.GetArea())))
	var shape1 Drawable = circle
	var shape2 Drawable = rect
	shape1.Draw()
	shape2.Draw()
	fmt.Println(("Shape 1 area via interface: " + gts_tostring(shape1.GetArea())))
	fmt.Println(("Shape 2 area via interface: " + gts_tostring(shape2.GetArea())))
}
