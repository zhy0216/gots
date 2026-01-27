// Test262-style tests for the Math object
// Tests Math methods and constants

// Helper for float comparison (due to floating point precision)
function approxEqual(a: number, b: number, epsilon: number): boolean {
    let diff: number = a - b
    if (diff < 0) {
        diff = -diff
    }
    return diff < epsilon
}

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// ============================================
// Math Constants
// ============================================

// Math.PI
let pi: number = Math.PI
printResult("Math.PI is approximately 3.14159", approxEqual(pi, 3.14159, 0.00001))

// Math.E
let e: number = Math.E
printResult("Math.E is approximately 2.71828", approxEqual(e, 2.71828, 0.00001))

// ============================================
// Rounding Methods
// ============================================

// Math.round
printResult("Math.round(4.5) = 5", Math.round(4.5) == 5)
printResult("Math.round(4.4) = 4", Math.round(4.4) == 4)
// Note: Go rounds -4.5 to -5 (away from zero), JS rounds to -4 (toward +infinity)
printResult("Math.round(-4.5) = -5", Math.round(-4.5) == -5)

// Math.floor
printResult("Math.floor(4.7) = 4", Math.floor(4.7) == 4)
printResult("Math.floor(-4.7) = -5", Math.floor(-4.7) == -5)

// Math.ceil
printResult("Math.ceil(4.3) = 5", Math.ceil(4.3) == 5)
printResult("Math.ceil(-4.3) = -4", Math.ceil(-4.3) == -4)

// Math.trunc
printResult("Math.trunc(4.7) = 4", Math.trunc(4.7) == 4)
printResult("Math.trunc(-4.7) = -4", Math.trunc(-4.7) == -4)

// ============================================
// Power and Root Methods
// ============================================

// Math.pow
printResult("Math.pow(2, 3) = 8", Math.pow(2, 3) == 8)
printResult("Math.pow(2, 0.5) approx sqrt(2)", approxEqual(Math.pow(2, 0.5), 1.41421, 0.0001))

// Math.sqrt
printResult("Math.sqrt(4) = 2", Math.sqrt(4) == 2)
printResult("Math.sqrt(2) approx 1.41421", approxEqual(Math.sqrt(2), 1.41421, 0.0001))

// Math.cbrt (cube root)
printResult("Math.cbrt(8) = 2", Math.cbrt(8) == 2)
printResult("Math.cbrt(27) = 3", Math.cbrt(27) == 3)

// Math.exp
printResult("Math.exp(0) = 1", Math.exp(0) == 1)
printResult("Math.exp(1) approx E", approxEqual(Math.exp(1), Math.E, 0.0001))

// Math.log (natural log)
printResult("Math.log(1) = 0", Math.log(1) == 0)
printResult("Math.log(E) approx 1", approxEqual(Math.log(Math.E), 1, 0.0001))

// Math.log10
printResult("Math.log10(100) = 2", Math.log10(100) == 2)
printResult("Math.log10(1000) = 3", Math.log10(1000) == 3)

// Math.log2
printResult("Math.log2(8) = 3", Math.log2(8) == 3)
printResult("Math.log2(16) = 4", Math.log2(16) == 4)

// ============================================
// Absolute Value and Sign
// ============================================

// Math.abs
printResult("Math.abs(-5) = 5", Math.abs(-5) == 5)
printResult("Math.abs(5) = 5", Math.abs(5) == 5)

// Math.sign
printResult("Math.sign(5) = 1", Math.sign(5) == 1)
printResult("Math.sign(-5) = -1", Math.sign(-5) == -1)
printResult("Math.sign(0) = 0", Math.sign(0) == 0)

// ============================================
// Min/Max
// ============================================

// Math.min
printResult("Math.min(1, 2) = 1", Math.min(1, 2) == 1)
printResult("Math.min(5, 3, 8) = 3", Math.min(5, 3, 8) == 3)
printResult("Math.min(-1, -5) = -5", Math.min(-1, -5) == -5)

// Math.max
printResult("Math.max(1, 2) = 2", Math.max(1, 2) == 2)
printResult("Math.max(5, 3, 8) = 8", Math.max(5, 3, 8) == 8)
printResult("Math.max(-1, -5) = -1", Math.max(-1, -5) == -1)

// ============================================
// Trigonometric Methods
// ============================================

// Math.sin
printResult("Math.sin(0) = 0", Math.sin(0) == 0)
printResult("Math.sin(PI/2) approx 1", approxEqual(Math.sin(Math.PI / 2), 1, 0.0001))

// Math.cos
printResult("Math.cos(0) = 1", Math.cos(0) == 1)
printResult("Math.cos(PI) approx -1", approxEqual(Math.cos(Math.PI), -1, 0.0001))

// Math.tan
printResult("Math.tan(0) = 0", Math.tan(0) == 0)
printResult("Math.tan(PI/4) approx 1", approxEqual(Math.tan(Math.PI / 4), 1, 0.0001))

// Math.asin
printResult("Math.asin(0) = 0", Math.asin(0) == 0)
printResult("Math.asin(1) approx PI/2", approxEqual(Math.asin(1), Math.PI / 2, 0.0001))

// Math.acos
printResult("Math.acos(1) = 0", Math.acos(1) == 0)
printResult("Math.acos(0) approx PI/2", approxEqual(Math.acos(0), Math.PI / 2, 0.0001))

// Math.atan
printResult("Math.atan(0) = 0", Math.atan(0) == 0)
printResult("Math.atan(1) approx PI/4", approxEqual(Math.atan(1), Math.PI / 4, 0.0001))

// Math.atan2
printResult("Math.atan2(1, 1) approx PI/4", approxEqual(Math.atan2(1, 1), Math.PI / 4, 0.0001))
printResult("Math.atan2(0, 1) = 0", Math.atan2(0, 1) == 0)

// ============================================
// Random
// ============================================

// Math.random returns a value between 0 (inclusive) and 1 (exclusive)
let r1: number = Math.random()
printResult("Math.random() >= 0", r1 >= 0)
printResult("Math.random() < 1", r1 < 1)

println("")
println("========== Math Object Tests Complete ==========")
