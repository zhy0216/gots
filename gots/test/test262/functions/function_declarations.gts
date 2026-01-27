// Test: Function declarations
// Based on test262 function tests

// Basic function
function add(a: int, b: int): int {
    return a + b
}
if (add(2, 3) == 5) { println("PASS: basic function") } else { println("FAIL: basic function") }

// Function with no return value (using closure pattern)
function makeResultSetter(): Function {
    let result: int = 0
    return function(x: int): int {
        result = x
        return result
    }
}
let setter: Function = makeResultSetter()
let setterResult: int = setter(42)
if (setterResult == 42) { println("PASS: void-like function") } else { println("FAIL: void-like function") }

// Function with multiple parameters
function sum3(a: int, b: int, c: int): int {
    return a + b + c
}
if (sum3(1, 2, 3) == 6) { println("PASS: function with 3 params") } else { println("FAIL: function with 3 params") }

// Function with no parameters
function getFortyTwo(): int {
    return 42
}
if (getFortyTwo() == 42) { println("PASS: function no params") } else { println("FAIL: function no params") }

// Recursive function
function factorial(n: int): int {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}
if (factorial(5) == 120) { println("PASS: recursive function") } else { println("FAIL: recursive function") }

// Function with different types
function greet(name: string): string {
    return "Hello, " + name
}
if (greet("World") == "Hello, World") { println("PASS: string function") } else { println("FAIL: string function") }

// Function returning boolean
function isPositive(n: int): boolean {
    return n > 0
}
if (isPositive(5) == true) { println("PASS: boolean return true") } else { println("FAIL: boolean return true") }
if (isPositive(-5) == false) { println("PASS: boolean return false") } else { println("FAIL: boolean return false") }

// Function with local variables
function compute(x: int): int {
    let temp: int = x * 2
    let result: int = temp + 1
    return result
}
if (compute(5) == 11) { println("PASS: function with locals") } else { println("FAIL: function with locals") }
