// Test: let variable declarations
// Based on test262 let declaration tests

// Basic let with type
let x: int = 42
if (x == 42) { println("PASS: let with int type") } else { println("FAIL: let with int type") }

// Let with inference
let y = 100
if (y == 100) { println("PASS: let with inference") } else { println("FAIL: let with inference") }

// Let reassignment
let z: int = 1
z = 2
if (z == 2) { println("PASS: let reassignment") } else { println("FAIL: let reassignment") }

// Multiple lets
let a: int = 1
let b: int = 2
let c: int = 3
if (a + b + c == 6) { println("PASS: multiple lets") } else { println("FAIL: multiple lets") }

// Let in block scope
let outer: int = 10
if (true) {
    let inner: int = 20
    if (inner == 20) { println("PASS: let in block") } else { println("FAIL: let in block") }
    if (outer == 10) { println("PASS: outer visible in block") } else { println("FAIL: outer visible in block") }
}

// Let with different types
let str: string = "test"
let bool: boolean = true
let num: number = 3.14

if (str == "test") { println("PASS: let string") } else { println("FAIL: let string") }
if (bool == true) { println("PASS: let boolean") } else { println("FAIL: let boolean") }
if (num == 3.14) { println("PASS: let number") } else { println("FAIL: let number") }
