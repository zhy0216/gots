// Test: Built-in functions
// Based on goTS built-in functions

// println - tested implicitly throughout

// len
let arr: int[] = [1, 2, 3, 4, 5]
if (len(arr) == 5) { println("PASS: len array") } else { println("FAIL: len array") }

let str: string = "hello"
if (len(str) == 5) { println("PASS: len string") } else { println("FAIL: len string") }

let empty: int[] = []
if (len(empty) == 0) { println("PASS: len empty") } else { println("FAIL: len empty") }

// push and pop
let stack: int[] = []
push(stack, 1)
push(stack, 2)
push(stack, 3)
if (len(stack) == 3) { println("PASS: push increases len") } else { println("FAIL: push increases len") }

let p: int = pop(stack)
if (p == 3) { println("PASS: pop returns last") } else { println("FAIL: pop returns last") }
if (len(stack) == 2) { println("PASS: pop decreases len") } else { println("FAIL: pop decreases len") }

// sqrt
let s: float = sqrt(16.0)
if (s == 4.0) { println("PASS: sqrt(16) = 4") } else { println("FAIL: sqrt(16) = 4, got " + tostring(s)) }

let s2: float = sqrt(2.0)
if (s2 > 1.41 && s2 < 1.42) { println("PASS: sqrt(2) ~ 1.414") } else { println("FAIL: sqrt(2) ~ 1.414") }

// floor
let fl: float = floor(3.7)
if (fl == 3.0) { println("PASS: floor(3.7) = 3") } else { println("FAIL: floor(3.7) = 3") }

let fl2: float = floor(-2.3)
if (fl2 == -3.0) { println("PASS: floor(-2.3) = -3") } else { println("FAIL: floor(-2.3) = -3, got " + tostring(fl2)) }

// ceil
let ce: float = ceil(3.2)
if (ce == 4.0) { println("PASS: ceil(3.2) = 4") } else { println("FAIL: ceil(3.2) = 4") }

let ce2: float = ceil(-2.7)
if (ce2 == -2.0) { println("PASS: ceil(-2.7) = -2") } else { println("FAIL: ceil(-2.7) = -2, got " + tostring(ce2)) }

// abs
let ab: float = abs(-5.0)
if (ab == 5.0) { println("PASS: abs(-5) = 5") } else { println("FAIL: abs(-5) = 5") }

let ab2: float = abs(3.0)
if (ab2 == 3.0) { println("PASS: abs(3) = 3") } else { println("FAIL: abs(3) = 3") }

// typeof - Note: goTS may return "int" for integer literals
let typeofVal: number = 42.0
if (typeof(typeofVal) == "number") { println("PASS: typeof number") } else { println("FAIL: typeof number, got " + typeof(typeofVal)) }
if (typeof("hi") == "string") { println("PASS: typeof string") } else { println("FAIL: typeof string") }
if (typeof(true) == "boolean") { println("PASS: typeof boolean") } else { println("FAIL: typeof boolean") }
if (typeof(arr) == "object") { println("PASS: typeof array") } else { println("FAIL: typeof array, got " + typeof(arr)) }
