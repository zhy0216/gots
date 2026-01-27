// Test: const declarations
// Based on test262 const declaration tests

// Basic const
const PI: float = 3.14159
if (PI == 3.14159) { println("PASS: const float") } else { println("FAIL: const float") }

// Const int
const MAX: int = 100
if (MAX == 100) { println("PASS: const int") } else { println("FAIL: const int") }

// Const string
const NAME: string = "goTS"
if (NAME == "goTS") { println("PASS: const string") } else { println("FAIL: const string") }

// Const boolean
const ENABLED: boolean = true
if (ENABLED == true) { println("PASS: const boolean") } else { println("FAIL: const boolean") }

// Multiple consts
const A: int = 1
const B: int = 2
const C: int = 3
if (A + B + C == 6) { println("PASS: multiple consts") } else { println("FAIL: multiple consts") }

// Const in expressions
let x: float = PI * 2.0
if (x == 6.28318) { println("PASS: const in expression") } else { println("FAIL: const in expression") }

// Const with inference
const INFERRED = 42
if (INFERRED == 42) { println("PASS: const with inference") } else { println("FAIL: const with inference") }
