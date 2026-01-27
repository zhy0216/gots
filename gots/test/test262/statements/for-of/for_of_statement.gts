// Test: For-of loops
// Based on test262 for-of tests

// Basic for-of with array
let arr: int[] = [1, 2, 3, 4, 5]
let sum: int = 0
for (let x of arr) {
    sum = sum + x
}
if (sum == 15) { println("PASS: basic for-of") } else { println("FAIL: basic for-of, got " + tostring(sum)) }

// For-of with empty array
let empty: int[] = []
let count: int = 0
for (let y of empty) {
    count = count + y
}
if (count == 0) { println("PASS: for-of empty array") } else { println("FAIL: for-of empty array") }

// For-of with string array
let names: string[] = ["Alice", "Bob", "Charlie"]
let result: string = ""
for (let name of names) {
    result = result + name + ","
}
if (result == "Alice,Bob,Charlie,") { println("PASS: for-of string array") } else { println("FAIL: for-of string array") }

// For-of single element
let singleVal: int = 42
let single: int[] = [singleVal]
sum = 0
for (let z of single) {
    sum = sum + z
}
if (sum == 42) { println("PASS: for-of single element") } else { println("FAIL: for-of single element") }

// Sum test (simpler nested test)
let nums1: int[] = [1, 2, 3]
let nums2: int[] = [4, 5, 6]
sum = 0
for (let a of nums1) {
    sum = sum + a
}
for (let b of nums2) {
    sum = sum + b
}
if (sum == 21) { println("PASS: multiple for-of") } else { println("FAIL: multiple for-of, got " + tostring(sum)) }
