// Test: Numeric literals - integers
// Based on test262 numeric literal tests

// Basic integer literals
let n0: number = 0
let n1: number = 1
let n2: number = 2
let n9: number = 9
let n42: number = 42
let n100: number = 100
let n999: number = 999

// Verify values
if (n0 == 0) { println("PASS: 0 literal") } else { println("FAIL: 0 literal") }
if (n1 == 1) { println("PASS: 1 literal") } else { println("FAIL: 1 literal") }
if (n2 == 2) { println("PASS: 2 literal") } else { println("FAIL: 2 literal") }
if (n9 == 9) { println("PASS: 9 literal") } else { println("FAIL: 9 literal") }
if (n42 == 42) { println("PASS: 42 literal") } else { println("FAIL: 42 literal") }
if (n100 == 100) { println("PASS: 100 literal") } else { println("FAIL: 100 literal") }
if (n999 == 999) { println("PASS: 999 literal") } else { println("FAIL: 999 literal") }

// Int type
let i0: int = 0
let i42: int = 42
let i1000: int = 1000

if (i0 == 0) { println("PASS: int 0") } else { println("FAIL: int 0") }
if (i42 == 42) { println("PASS: int 42") } else { println("FAIL: int 42") }
if (i1000 == 1000) { println("PASS: int 1000") } else { println("FAIL: int 1000") }

// Negative integers
let neg1: int = -1
let neg100: int = -100

if (neg1 == -1) { println("PASS: -1 literal") } else { println("FAIL: -1 literal") }
if (neg100 == -100) { println("PASS: -100 literal") } else { println("FAIL: -100 literal") }
