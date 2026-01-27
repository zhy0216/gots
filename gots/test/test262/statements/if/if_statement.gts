// Test: If statements
// Based on test262 if statement tests

// Basic if true
let x: int = 0
if (true) {
    x = 1
}
if (x == 1) { println("PASS: if true executes") } else { println("FAIL: if true executes") }

// Basic if false
x = 0
if (false) {
    x = 1
}
if (x == 0) { println("PASS: if false skips") } else { println("FAIL: if false skips") }

// If-else true branch
x = 0
if (true) {
    x = 1
} else {
    x = 2
}
if (x == 1) { println("PASS: if-else true branch") } else { println("FAIL: if-else true branch") }

// If-else false branch
x = 0
if (false) {
    x = 1
} else {
    x = 2
}
if (x == 2) { println("PASS: if-else false branch") } else { println("FAIL: if-else false branch") }

// If-else if-else
x = 0
let n: int = 2
if (n == 1) {
    x = 1
} else if (n == 2) {
    x = 2
} else {
    x = 3
}
if (x == 2) { println("PASS: if-else if-else") } else { println("FAIL: if-else if-else") }

// Nested if
x = 0
if (true) {
    if (true) {
        x = 1
    }
}
if (x == 1) { println("PASS: nested if") } else { println("FAIL: nested if") }

// Condition with comparison
let a: int = 5
if (a > 3) {
    x = 1
} else {
    x = 0
}
if (x == 1) { println("PASS: if with comparison") } else { println("FAIL: if with comparison") }

// Condition with logical
if (a > 0 && a < 10) {
    x = 1
} else {
    x = 0
}
if (x == 1) { println("PASS: if with logical") } else { println("FAIL: if with logical") }
