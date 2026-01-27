// Test: Logical operators
// Based on test262 logical operator tests

// Logical AND
if (true && true) { println("PASS: true && true") } else { println("FAIL: true && true") }
if (!(true && false)) { println("PASS: !(true && false)") } else { println("FAIL: !(true && false)") }
if (!(false && true)) { println("PASS: !(false && true)") } else { println("FAIL: !(false && true)") }
if (!(false && false)) { println("PASS: !(false && false)") } else { println("FAIL: !(false && false)") }

// Logical OR
if (true || true) { println("PASS: true || true") } else { println("FAIL: true || true") }
if (true || false) { println("PASS: true || false") } else { println("FAIL: true || false") }
if (false || true) { println("PASS: false || true") } else { println("FAIL: false || true") }
if (!(false || false)) { println("PASS: !(false || false)") } else { println("FAIL: !(false || false)") }

// Logical NOT
if (!false) { println("PASS: !false") } else { println("FAIL: !false") }
if (!(!true)) { println("PASS: !(!true)") } else { println("FAIL: !(!true)") }

// Combined
if ((true && true) || false) { println("PASS: (true && true) || false") } else { println("FAIL: (true && true) || false") }
if (true && (false || true)) { println("PASS: true && (false || true)") } else { println("FAIL: true && (false || true)") }
if (!(false && true) || false) { println("PASS: !(false && true) || false") } else { println("FAIL: !(false && true) || false") }

// With comparison
let a: int = 5
let b: int = 3
if (a > 0 && b > 0) { println("PASS: a > 0 && b > 0") } else { println("FAIL: a > 0 && b > 0") }
if (a > 10 || b > 0) { println("PASS: a > 10 || b > 0") } else { println("FAIL: a > 10 || b > 0") }
if (!(a < 0) && b < 10) { println("PASS: !(a < 0) && b < 10") } else { println("FAIL: !(a < 0) && b < 10") }
