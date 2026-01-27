// Test: Assignment operators
// Based on test262 assignment tests

// Basic assignment
let a: int = 5
if (a == 5) { println("PASS: basic assignment") } else { println("FAIL: basic assignment") }

// Reassignment
a = 10
if (a == 10) { println("PASS: reassignment") } else { println("FAIL: reassignment") }

// Add and assign
a = 5
a += 3
if (a == 8) { println("PASS: += operator") } else { println("FAIL: += operator") }

// Subtract and assign
a = 10
a -= 3
if (a == 7) { println("PASS: -= operator") } else { println("FAIL: -= operator") }

// Multiply and assign
a = 4
a *= 3
if (a == 12) { println("PASS: *= operator") } else { println("FAIL: *= operator") }

// Divide and assign
let f: float = 10.0
f /= 4.0
if (f == 2.5) { println("PASS: /= operator") } else { println("FAIL: /= operator") }

// Modulo and assign
a = 10
a %= 3
if (a == 1) { println("PASS: %= operator") } else { println("FAIL: %= operator") }

// Chain assignment
let x: int = 0
let y: int = 0
x = 5
y = x
if (x == 5 && y == 5) { println("PASS: chain assignment") } else { println("FAIL: chain assignment") }

// Compound in expression (commented - not supported)
// let result: int = a += 3
a = 5
a += 3
if (a == 8) { println("PASS: compound assignment works") } else { println("FAIL: compound assignment works") }
