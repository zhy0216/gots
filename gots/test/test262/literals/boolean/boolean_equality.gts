// Test: Boolean equality comparisons

// true != false
if (true != false) {
    println("PASS: true != false")
} else {
    println("FAIL: true should not equal false")
}

// true == true
if (true == true) {
    println("PASS: true == true")
} else {
    println("FAIL: true should equal true")
}

// false == false
if (false == false) {
    println("PASS: false == false")
} else {
    println("FAIL: false should equal false")
}

// Boolean variable comparison
let a: boolean = true
let b: boolean = false
let c: boolean = true

if (a == c) {
    println("PASS: two true variables are equal")
} else {
    println("FAIL: two true variables should be equal")
}

if (a != b) {
    println("PASS: true variable != false variable")
} else {
    println("FAIL: true variable should not equal false variable")
}
