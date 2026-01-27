// Test: Boolean negation operator

// !true == false
if (!true == false) {
    println("PASS: !true == false")
} else {
    println("FAIL: !true should equal false")
}

// !false == true
if (!false == true) {
    println("PASS: !false == true")
} else {
    println("FAIL: !false should equal true")
}

// Double negation
if (!!true == true) {
    println("PASS: !!true == true")
} else {
    println("FAIL: !!true should equal true")
}

if (!!false == false) {
    println("PASS: !!false == false")
} else {
    println("FAIL: !!false should equal false")
}

// Negation with variables
let t: boolean = true
let f: boolean = false

if (!t == false) {
    println("PASS: !t == false (t is true)")
} else {
    println("FAIL: !t should be false")
}

if (!f == true) {
    println("PASS: !f == true (f is false)")
} else {
    println("FAIL: !f should be true")
}
