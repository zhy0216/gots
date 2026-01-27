// Test: Numeric literals - floats
// Based on test262 decimal literal tests

// Float literals with decimal point
let f0: float = 0.0
let f1: float = 1.0
let f05: float = 0.5
let f314: float = 3.14
let f999: float = 999.999

if (f999 > 999.0) { println("PASS: 999.999 literal") } else { println("FAIL: 999.999 literal") }
if (f0 == 0.0) { println("PASS: 0.0 literal") } else { println("FAIL: 0.0 literal") }
if (f1 == 1.0) { println("PASS: 1.0 literal") } else { println("FAIL: 1.0 literal") }
if (f05 == 0.5) { println("PASS: 0.5 literal") } else { println("FAIL: 0.5 literal") }
if (f314 == 3.14) { println("PASS: 3.14 literal") } else { println("FAIL: 3.14 literal") }

// Number type (default numeric)
let n0: number = 0.0
let n314: number = 3.14159

if (n0 == 0.0) { println("PASS: number 0.0") } else { println("FAIL: number 0.0") }
if (n314 == 3.14159) { println("PASS: number 3.14159") } else { println("FAIL: number 3.14159") }

// Negative floats
let negf1: float = -1.5
let negf314: float = -3.14

if (negf1 == -1.5) { println("PASS: -1.5 literal") } else { println("FAIL: -1.5 literal") }
if (negf314 == -3.14) { println("PASS: -3.14 literal") } else { println("FAIL: -3.14 literal") }

// Small decimals
let small1: float = 0.001
let small2: float = 0.0001

if (small1 == 0.001) { println("PASS: 0.001 literal") } else { println("FAIL: 0.001 literal") }
if (small2 == 0.0001) { println("PASS: 0.0001 literal") } else { println("FAIL: 0.0001 literal") }
