// Test: Arithmetic operators
// Based on test262 arithmetic tests

// Addition
let a: int = 5
let b: int = 3
if (a + b == 8) { println("PASS: 5 + 3 == 8") } else { println("FAIL: 5 + 3 == 8") }
if (0 + 0 == 0) { println("PASS: 0 + 0 == 0") } else { println("FAIL: 0 + 0 == 0") }
if (-1 + 1 == 0) { println("PASS: -1 + 1 == 0") } else { println("FAIL: -1 + 1 == 0") }

// Subtraction
if (a - b == 2) { println("PASS: 5 - 3 == 2") } else { println("FAIL: 5 - 3 == 2") }
if (0 - 5 == -5) { println("PASS: 0 - 5 == -5") } else { println("FAIL: 0 - 5 == -5") }
if (b - a == -2) { println("PASS: 3 - 5 == -2") } else { println("FAIL: 3 - 5 == -2") }

// Multiplication
if (a * b == 15) { println("PASS: 5 * 3 == 15") } else { println("FAIL: 5 * 3 == 15") }
if (a * 0 == 0) { println("PASS: 5 * 0 == 0") } else { println("FAIL: 5 * 0 == 0") }
if (a * 1 == 5) { println("PASS: 5 * 1 == 5") } else { println("FAIL: 5 * 1 == 5") }
if (-2 * 3 == -6) { println("PASS: -2 * 3 == -6") } else { println("FAIL: -2 * 3 == -6") }
if (-2 * -3 == 6) { println("PASS: -2 * -3 == 6") } else { println("FAIL: -2 * -3 == 6") }

// Division
let x: float = 10.0
let y: float = 4.0
if (x / y == 2.5) { println("PASS: 10.0 / 4.0 == 2.5") } else { println("FAIL: 10.0 / 4.0 == 2.5") }
if (x / 2.0 == 5.0) { println("PASS: 10.0 / 2.0 == 5.0") } else { println("FAIL: 10.0 / 2.0 == 5.0") }

// Modulo (int only)
if (10 % 3 == 1) { println("PASS: 10 % 3 == 1") } else { println("FAIL: 10 % 3 == 1") }
if (15 % 5 == 0) { println("PASS: 15 % 5 == 0") } else { println("FAIL: 15 % 5 == 0") }
if (7 % 2 == 1) { println("PASS: 7 % 2 == 1") } else { println("FAIL: 7 % 2 == 1") }

// Unary minus
let n: int = 5
if (-n == -5) { println("PASS: -n == -5") } else { println("FAIL: -n == -5") }
if (-(-n) == 5) { println("PASS: -(-n) == 5") } else { println("FAIL: -(-n) == 5") }

// Operator precedence
if (2 + 3 * 4 == 14) { println("PASS: 2 + 3 * 4 == 14") } else { println("FAIL: 2 + 3 * 4 == 14") }
if ((2 + 3) * 4 == 20) { println("PASS: (2 + 3) * 4 == 20") } else { println("FAIL: (2 + 3) * 4 == 20") }
