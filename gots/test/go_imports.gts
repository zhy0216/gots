// Test Go package imports
import { ToUpper, ToLower, Contains, Split, Join } from "go:strings"
import { Sqrt, Pow, Max, Min } from "go:math"

// Test strings functions
let upper: string = ToUpper("hello world")
let lower: string = ToLower("HELLO WORLD")
println("Upper: " + upper)
println("Lower: " + lower)

let hasHello: boolean = Contains("Hello World", "Hello")
println("Contains Hello: " + tostring(hasHello))

let parts: string[] = Split("a,b,c,d", ",")
println("Split result:")
for (let i: int = 0; i < len(parts); i = i + 1) {
    println("  " + parts[i])
}

let joined: string = Join(parts, "-")
println("Joined: " + joined)

// Test math functions
let sqrtVal: float = Sqrt(16.0)
println("Sqrt(16) = " + tostring(sqrtVal))

let powVal: float = Pow(2.0, 10.0)
println("Pow(2, 10) = " + tostring(powVal))

let maxVal: float = Max(3.14, 2.71)
let minVal: float = Min(3.14, 2.71)
println("Max(3.14, 2.71) = " + tostring(maxVal))
println("Min(3.14, 2.71) = " + tostring(minVal))
