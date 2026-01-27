// Test: Switch statement
// Based on test262 switch tests

// Basic switch
let x: int = 2
let result: string = ""

switch (x) {
    case 1:
        result = "one"
        break
    case 2:
        result = "two"
        break
    case 3:
        result = "three"
        break
    default:
        result = "other"
}
if (result == "two") { println("PASS: basic switch") } else { println("FAIL: basic switch, got " + result) }

// Switch with default
x = 10
switch (x) {
    case 1:
        result = "one"
        break
    default:
        result = "default"
}
if (result == "default") { println("PASS: switch default") } else { println("FAIL: switch default") }

// Switch without break (fallthrough) - Note: Go doesn't auto-fallthrough
// In goTS, each case block is independent (Go semantics)
x = 1
let count: int = 0
switch (x) {
    case 1:
        count = count + 1
        // fallthrough not supported in goTS
        break
    case 2:
        count = count + 1
        break
    case 3:
        count = count + 1
}
if (count == 1) { println("PASS: switch no fallthrough") } else { println("FAIL: switch no fallthrough, got " + tostring(count)) }

// Switch with string
let s: string = "hello"
switch (s) {
    case "hello":
        result = "greeting"
        break
    case "bye":
        result = "farewell"
        break
    default:
        result = "unknown"
}
if (result == "greeting") { println("PASS: switch string") } else { println("FAIL: switch string") }

// Nested switch
x = 1
let y: int = 2
let nested: string = ""
switch (x) {
    case 1:
        switch (y) {
            case 2:
                nested = "x1y2"
                break
            default:
                nested = "x1other"
        }
        break
    default:
        nested = "other"
}
if (nested == "x1y2") { println("PASS: nested switch") } else { println("FAIL: nested switch") }
