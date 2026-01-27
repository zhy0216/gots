// Test: Boolean literal false
// Based on test262/test/language/literals/boolean/S7.8.2_A1_T2.js

// Check that false literal has correct value
let f: boolean = false
if (f == false) {
    println("PASS: false literal equals false")
} else {
    println("FAIL: false literal does not equal false")
}

// Check false in condition
if (false) {
    println("FAIL: false should not be truthy")
} else {
    println("PASS: false is falsy")
}

// Check typeof
if (typeof(false) == "boolean") {
    println("PASS: typeof false is boolean")
} else {
    println("FAIL: typeof false should be boolean")
}
