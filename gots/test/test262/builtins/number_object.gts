// Test262-style tests for the Number object
// Tests Number static methods and constants

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// Helper for float comparison
function approxEqual(a: number, b: number, epsilon: number): boolean {
    let diff: number = a - b
    if (diff < 0) {
        diff = -diff
    }
    return diff < epsilon
}

// ============================================
// Number Static Methods
// ============================================

// Number.isFinite
printResult("Number.isFinite(42) = true", Number.isFinite(42) == true)
printResult("Number.isFinite(3.14) = true", Number.isFinite(3.14) == true)
printResult("Number.isFinite(0) = true", Number.isFinite(0) == true)

// Number.isNaN
printResult("Number.isNaN(42) = false", Number.isNaN(42) == false)
printResult("Number.isNaN(3.14) = false", Number.isNaN(3.14) == false)

// Number.isInteger
printResult("Number.isInteger(42) = true", Number.isInteger(42) == true)
printResult("Number.isInteger(42.0) = true", Number.isInteger(42.0) == true)
printResult("Number.isInteger(3.14) = false", Number.isInteger(3.14) == false)

// Number.parseFloat
printResult("Number.parseFloat('3.14') = 3.14", approxEqual(Number.parseFloat("3.14"), 3.14, 0.001))
printResult("Number.parseFloat('42') = 42", Number.parseFloat("42") == 42)
printResult("Number.parseFloat('-5.5') = -5.5", approxEqual(Number.parseFloat("-5.5"), -5.5, 0.001))

// Number.parseInt
printResult("Number.parseInt('42') = 42", Number.parseInt("42") == 42)
printResult("Number.parseInt('-10') = -10", Number.parseInt("-10") == -10)
printResult("Number.parseInt('ff', 16) = 255", Number.parseInt("ff", 16) == 255)
printResult("Number.parseInt('1010', 2) = 10", Number.parseInt("1010", 2) == 10)

// ============================================
// Number Constants
// ============================================

// Number.MAX_SAFE_INTEGER
printResult("Number.MAX_SAFE_INTEGER is 9007199254740991", Number.MAX_SAFE_INTEGER == 9007199254740991)

// Number.MIN_SAFE_INTEGER
printResult("Number.MIN_SAFE_INTEGER is -9007199254740991", Number.MIN_SAFE_INTEGER == -9007199254740991)

// Global number functions can be tested with Number.NaN constant
// once NaN constant support is added

// isNaN (global)
printResult("isNaN(42) = false", isNaN(42) == false)
printResult("isNaN(3.14) = false", isNaN(3.14) == false)

// isFinite (global)
printResult("isFinite(42) = true", isFinite(42) == true)
printResult("isFinite(3.14) = true", isFinite(3.14) == true)

// parseFloat (global) - already exists as tofloat
printResult("parseFloat('3.14') = 3.14", approxEqual(parseFloat("3.14"), 3.14, 0.001))

println("")
println("========== Number Object Tests Complete ==========")
