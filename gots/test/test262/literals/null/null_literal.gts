// Test: Null literal
// Based on test262 null literal tests

// Null assignment
let n: string | null = null

if (n == null) { println("PASS: null equals null") } else { println("FAIL: null equals null") }

// Null comparison
let x: int | null = null
let y: int | null = null

if (x == y) { println("PASS: null == null") } else { println("FAIL: null == null") }
if (x == null) { println("PASS: x == null") } else { println("FAIL: x == null") }

// Non-null then null
let s: string | null = "hello"
if (s != null) { println("PASS: s != null when assigned") } else { println("FAIL: s != null when assigned") }

s = null
if (s == null) { println("PASS: s == null after reassignment") } else { println("FAIL: s == null after reassignment") }

// typeof null
if (typeof(null) == "null") { println("PASS: typeof(null) == null") } else { println("FAIL: typeof(null) == null, got: " + typeof(null)) }
