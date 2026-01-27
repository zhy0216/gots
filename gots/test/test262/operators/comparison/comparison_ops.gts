// Test: Comparison operators
// Based on test262 comparison tests

let a: int = 5
let b: int = 3
let c: int = 5

// Equality
if (a == c) { println("PASS: 5 == 5") } else { println("FAIL: 5 == 5") }
if (!(a == b)) { println("PASS: 5 != 3") } else { println("FAIL: 5 != 3") }

// Not equal
if (a != b) { println("PASS: 5 != 3") } else { println("FAIL: 5 != 3") }
if (!(a != c)) { println("PASS: !(5 != 5)") } else { println("FAIL: !(5 != 5)") }

// Less than
if (b < a) { println("PASS: 3 < 5") } else { println("FAIL: 3 < 5") }
if (!(a < b)) { println("PASS: !(5 < 3)") } else { println("FAIL: !(5 < 3)") }
if (!(a < c)) { println("PASS: !(5 < 5)") } else { println("FAIL: !(5 < 5)") }

// Less than or equal
if (b <= a) { println("PASS: 3 <= 5") } else { println("FAIL: 3 <= 5") }
if (a <= c) { println("PASS: 5 <= 5") } else { println("FAIL: 5 <= 5") }
if (!(a <= b)) { println("PASS: !(5 <= 3)") } else { println("FAIL: !(5 <= 3)") }

// Greater than
if (a > b) { println("PASS: 5 > 3") } else { println("FAIL: 5 > 3") }
if (!(b > a)) { println("PASS: !(3 > 5)") } else { println("FAIL: !(3 > 5)") }
if (!(a > c)) { println("PASS: !(5 > 5)") } else { println("FAIL: !(5 > 5)") }

// Greater than or equal
if (a >= b) { println("PASS: 5 >= 3") } else { println("FAIL: 5 >= 3") }
if (a >= c) { println("PASS: 5 >= 5") } else { println("FAIL: 5 >= 5") }
if (!(b >= a)) { println("PASS: !(3 >= 5)") } else { println("FAIL: !(3 >= 5)") }

// String comparison
let s1: string = "abc"
let s2: string = "abc"
let s3: string = "def"
if (s1 == s2) { println("PASS: abc == abc") } else { println("FAIL: abc == abc") }
if (s1 != s3) { println("PASS: abc != def") } else { println("FAIL: abc != def") }
