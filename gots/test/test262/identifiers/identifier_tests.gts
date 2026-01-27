// Test: Identifiers
// Based on test262 identifier tests

// Basic identifiers
let a: int = 1
let abc: int = 2
let myVar: int = 3
let MY_CONST: int = 4
let _private: int = 5
let __double: int = 6
let var123: int = 7
let camelCase: int = 8
let PascalCase: int = 9

if (a == 1) { println("PASS: single char identifier") } else { println("FAIL: single char identifier") }
if (abc == 2) { println("PASS: multi char identifier") } else { println("FAIL: multi char identifier") }
if (myVar == 3) { println("PASS: camelCase identifier") } else { println("FAIL: camelCase identifier") }
if (MY_CONST == 4) { println("PASS: upper snake case") } else { println("FAIL: upper snake case") }
if (_private == 5) { println("PASS: underscore prefix") } else { println("FAIL: underscore prefix") }
if (__double == 6) { println("PASS: double underscore") } else { println("FAIL: double underscore") }
if (var123 == 7) { println("PASS: identifier with digits") } else { println("FAIL: identifier with digits") }
if (camelCase == 8) { println("PASS: camelCase") } else { println("FAIL: camelCase") }
if (PascalCase == 9) { println("PASS: PascalCase") } else { println("FAIL: PascalCase") }

// Identifiers as function names
function myFunc(): int {
    return 42
}
if (myFunc() == 42) { println("PASS: function identifier") } else { println("FAIL: function identifier") }

// Identifiers as class names
class MyClass {
    value: int
    constructor(v: int) {
        this.value = v
    }
}
let obj: MyClass = new MyClass(100)
if (obj.value == 100) { println("PASS: class identifier") } else { println("FAIL: class identifier") }

// Long identifiers
let thisIsAVeryLongIdentifierNameThatShouldStillWork: int = 42
if (thisIsAVeryLongIdentifierNameThatShouldStillWork == 42) { println("PASS: long identifier") } else { println("FAIL: long identifier") }
