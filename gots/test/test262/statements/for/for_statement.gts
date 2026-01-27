// Test: For loops
// Based on test262 for tests

// Basic for loop
let sum: int = 0
for (let i: int = 0; i < 5; i = i + 1) {
    sum = sum + i
}
if (sum == 10) { println("PASS: basic for loop") } else { println("FAIL: basic for loop, got " + tostring(sum)) }

// For loop with different step
sum = 0
for (let i: int = 0; i < 10; i = i + 2) {
    sum = sum + i
}
if (sum == 20) { println("PASS: for with step 2") } else { println("FAIL: for with step 2, got " + tostring(sum)) }

// For count down
sum = 0
for (let i: int = 5; i > 0; i = i - 1) {
    sum = sum + i
}
if (sum == 15) { println("PASS: for count down") } else { println("FAIL: for count down, got " + tostring(sum)) }

// Nested for
let count: int = 0
for (let i: int = 0; i < 3; i = i + 1) {
    for (let j: int = 0; j < 3; j = j + 1) {
        count = count + 1
    }
}
if (count == 9) { println("PASS: nested for") } else { println("FAIL: nested for, got " + tostring(count)) }

// For with compound assignment
sum = 0
for (let i: int = 1; i <= 5; i += 1) {
    sum += i
}
if (sum == 15) { println("PASS: for with +=") } else { println("FAIL: for with +=, got " + tostring(sum)) }
