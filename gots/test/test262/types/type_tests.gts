// Test: Type system
// Based on test262 type tests adapted for goTS

// Type inference
let inferred_int = 42
let inferred_float = 3.14
let inferred_string = "hello"
let inferred_bool = true

if (typeof(inferred_int) == "number") { println("PASS: inferred int type") } else { println("FAIL: inferred int type, got " + typeof(inferred_int)) }
if (typeof(inferred_float) == "number") { println("PASS: inferred float type") } else { println("FAIL: inferred float type, got " + typeof(inferred_float)) }
if (typeof(inferred_string) == "string") { println("PASS: inferred string type") } else { println("FAIL: inferred string type") }
if (typeof(inferred_bool) == "boolean") { println("PASS: inferred bool type") } else { println("FAIL: inferred bool type") }

// Explicit types
let explicit_int: int = 42
let explicit_float: float = 3.14
let explicit_number: number = 100
let explicit_string: string = "test"
let explicit_bool: boolean = false

// Note: typeof returns "number" for both int and float types in goTS
if (typeof(explicit_int) == "number") { println("PASS: explicit int") } else { println("FAIL: explicit int, got " + typeof(explicit_int)) }
if (typeof(explicit_float) == "number") { println("PASS: explicit float") } else { println("FAIL: explicit float") }
if (typeof(explicit_number) == "number") { println("PASS: explicit number") } else { println("FAIL: explicit number") }
if (typeof(explicit_string) == "string") { println("PASS: explicit string") } else { println("FAIL: explicit string") }
if (typeof(explicit_bool) == "boolean") { println("PASS: explicit bool") } else { println("FAIL: explicit bool") }

// Array types
let int_arr: int[] = [1, 2, 3]
let str_arr: string[] = ["a", "b"]
let empty_arr: int[] = []

if (len(int_arr) == 3) { println("PASS: int array") } else { println("FAIL: int array") }
if (len(str_arr) == 2) { println("PASS: string array") } else { println("FAIL: string array") }
if (len(empty_arr) == 0) { println("PASS: empty array") } else { println("FAIL: empty array") }

// Type conversions
let n: int = 42
let s: string = tostring(n)
if (s == "42") { println("PASS: tostring int") } else { println("FAIL: tostring int") }

let f: float = 3.14
let sf: string = tostring(f)
if (sf == "3.14") { println("PASS: tostring float") } else { println("FAIL: tostring float") }

let si: string = "123"
let i: int = toint(si)
if (i == 123) { println("PASS: toint") } else { println("FAIL: toint") }

let sfloat: string = "3.14"
let ff: float = tofloat(sfloat)
if (ff == 3.14) { println("PASS: tofloat") } else { println("FAIL: tofloat") }
