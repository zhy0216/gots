// Test: String literals
// Based on test262 string literal tests

// Basic string literals with double quotes
let s1: string = "hello"
let s2: string = "world"
let s3: string = ""
let s4: string = "a"

if (s1 == "hello") { println("PASS: hello string") } else { println("FAIL: hello string") }
if (s2 == "world") { println("PASS: world string") } else { println("FAIL: world string") }
if (s3 == "") { println("PASS: empty string") } else { println("FAIL: empty string") }
if (s4 == "a") { println("PASS: single char string") } else { println("FAIL: single char string") }

// String with spaces
let s5: string = "hello world"
if (s5 == "hello world") { println("PASS: string with space") } else { println("FAIL: string with space") }

// String concatenation
let concat: string = "hello" + " " + "world"
if (concat == "hello world") { println("PASS: string concatenation") } else { println("FAIL: string concatenation") }

// String with numbers
let s6: string = "test123"
if (s6 == "test123") { println("PASS: string with digits") } else { println("FAIL: string with digits") }

// Special characters (basic)
let s7: string = "line1\nline2"
let s8: string = "col1\tcol2"

if (len(s7) > 0) { println("PASS: newline string") } else { println("FAIL: newline string") }
if (len(s8) > 0) { println("PASS: tab string") } else { println("FAIL: tab string") }

// Length check
if (len(s1) == 5) { println("PASS: len(hello) == 5") } else { println("FAIL: len(hello) == 5") }
if (len(s3) == 0) { println("PASS: len(empty) == 0") } else { println("FAIL: len(empty) == 0") }
