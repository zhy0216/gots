// Test: Boolean literal true
// Based on test262/test/language/literals/boolean/S7.8.2_A1_T1.js

// Check that true literal has correct value
let t: boolean = true
if (t == true) {
    println("PASS: true literal equals true")
} else {
    println("FAIL: true literal does not equal true")
}

// Check true in condition
if (true) {
    println("PASS: true is truthy")
} else {
    println("FAIL: true should be truthy")
}

// Check typeof
if (typeof(true) == "boolean") {
    println("PASS: typeof true is boolean")
} else {
    println("FAIL: typeof true should be boolean")
}
