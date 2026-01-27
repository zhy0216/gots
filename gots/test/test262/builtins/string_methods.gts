// Test262-style tests for additional String methods
// Tests string methods not yet implemented in goTS

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// ============================================
// String.charCodeAt
// ============================================

let str: string = "Hello"
printResult("charCodeAt(0) = 72 (H)", str.charCodeAt(0) == 72)
printResult("charCodeAt(1) = 101 (e)", str.charCodeAt(1) == 101)
printResult("charCodeAt(4) = 111 (o)", str.charCodeAt(4) == 111)

// ============================================
// String.at
// ============================================

printResult("at(0) = 'H'", str.at(0) == "H")
printResult("at(-1) = 'o' (last char)", str.at(-1) == "o")
printResult("at(-2) = 'l'", str.at(-2) == "l")

// ============================================
// String.slice
// ============================================

printResult("slice(0, 2) = 'He'", str.slice(0, 2) == "He")
printResult("slice(1, 4) = 'ell'", str.slice(1, 4) == "ell")
printResult("slice(-2) = 'lo'", str.slice(-2) == "lo")

// ============================================
// String.repeat
// ============================================

let ab: string = "ab"
printResult("repeat(3) = 'ababab'", ab.repeat(3) == "ababab")
printResult("repeat(1) = 'ab'", ab.repeat(1) == "ab")
printResult("repeat(0) = ''", ab.repeat(0) == "")

// ============================================
// String.padStart
// ============================================

let num: string = "5"
printResult("padStart(3, '0') = '005'", num.padStart(3, "0") == "005")
printResult("padStart(5, '0') = '00005'", num.padStart(5, "0") == "00005")

// ============================================
// String.padEnd
// ============================================

printResult("padEnd(3, '0') = '500'", num.padEnd(3, "0") == "500")
printResult("padEnd(5, '0') = '50000'", num.padEnd(5, "0") == "50000")

// ============================================
// String.trimStart / trimEnd
// ============================================

let padded: string = "  hello  "
printResult("trimStart() = 'hello  '", padded.trimStart() == "hello  ")
printResult("trimEnd() = '  hello'", padded.trimEnd() == "  hello")

// ============================================
// String.replaceAll
// ============================================

let text: string = "foo bar foo baz foo"
printResult("replaceAll('foo', 'x') replaces all", text.replaceAll("foo", "x") == "x bar x baz x")

println("")
println("========== String Additional Methods Tests Complete ==========")
