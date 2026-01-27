// Test: Regular expressions
// Based on test262 regexp tests adapted for goTS

// Basic regex literal
let re: RegExp = /hello/

// test() method
if (re.test("hello world")) { println("PASS: test matches") } else { println("FAIL: test matches") }
if (!re.test("goodbye")) { println("PASS: test no match") } else { println("FAIL: test no match") }

// Case insensitive flag
let reI: RegExp = /hello/i
if (reI.test("HELLO")) { println("PASS: case insensitive") } else { println("FAIL: case insensitive") }
if (reI.test("Hello")) { println("PASS: mixed case") } else { println("FAIL: mixed case") }

// Digit pattern
let digits: RegExp = /\d+/
if (digits.test("abc123")) { println("PASS: digits match") } else { println("FAIL: digits match") }
if (!digits.test("abcdef")) { println("PASS: no digits") } else { println("FAIL: no digits") }

// Word boundary
let word: RegExp = /\btest\b/
if (word.test("this is a test")) { println("PASS: word boundary match") } else { println("FAIL: word boundary match") }
if (!word.test("testing")) { println("PASS: word boundary no match") } else { println("FAIL: word boundary no match") }

// exec() method
let numRe: RegExp = /\d+/
let execResult: string[] | null = numRe.exec("abc123def")
if (execResult != null) {
    if (execResult[0] == "123") { println("PASS: exec returns match") } else { println("FAIL: exec returns match, got " + execResult[0]) }
} else {
    println("FAIL: exec returned null")
}

// exec() no match
let noMatch: string[] | null = numRe.exec("abcdef")
if (noMatch == null) { println("PASS: exec null on no match") } else { println("FAIL: exec null on no match") }

// Email-like pattern
let email: RegExp = /[a-z]+@[a-z]+\.[a-z]+/i
if (email.test("test@example.com")) { println("PASS: email pattern") } else { println("FAIL: email pattern") }
if (!email.test("invalid")) { println("PASS: invalid email") } else { println("FAIL: invalid email") }

// Start and end anchors
let startEnd: RegExp = /^hello$/
if (startEnd.test("hello")) { println("PASS: exact match") } else { println("FAIL: exact match") }
if (!startEnd.test("hello world")) { println("PASS: not exact match") } else { println("FAIL: not exact match") }
