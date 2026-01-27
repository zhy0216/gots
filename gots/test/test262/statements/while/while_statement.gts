// Test: While loops
// Based on test262 while tests

// Basic while
let i: int = 0
let sum: int = 0
while (i < 5) {
    sum = sum + i
    i = i + 1
}
if (sum == 10) { println("PASS: basic while loop") } else { println("FAIL: basic while loop, got " + tostring(sum)) }

// While with false condition (never executes)
let x: int = 0
while (false) {
    x = 1
}
if (x == 0) { println("PASS: while false never executes") } else { println("FAIL: while false never executes") }

// While count down
i = 5
let count: int = 0
while (i > 0) {
    count = count + 1
    i = i - 1
}
if (count == 5) { println("PASS: while count down") } else { println("FAIL: while count down") }

// Nested while
i = 0
let j: int = 0
count = 0
while (i < 3) {
    j = 0
    while (j < 3) {
        count = count + 1
        j = j + 1
    }
    i = i + 1
}
if (count == 9) { println("PASS: nested while") } else { println("FAIL: nested while, got " + tostring(count)) }

// While with complex condition
i = 0
sum = 0
while (i < 10 && sum < 20) {
    sum = sum + i
    i = i + 1
}
if (sum == 21) { println("PASS: while with complex condition") } else { println("FAIL: while with complex condition, got " + tostring(sum)) }
