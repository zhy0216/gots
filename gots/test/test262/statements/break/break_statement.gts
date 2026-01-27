// Test: Break statement
// Based on test262 break tests

// Break in while
let i: int = 0
while (true) {
    if (i >= 5) {
        break
    }
    i = i + 1
}
if (i == 5) { println("PASS: break in while") } else { println("FAIL: break in while") }

// Break in for
let sum: int = 0
for (let j: int = 0; j < 100; j = j + 1) {
    if (j >= 5) {
        break
    }
    sum = sum + j
}
if (sum == 10) { println("PASS: break in for") } else { println("FAIL: break in for, got " + tostring(sum)) }

// Break in for-of
let arr: int[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
sum = 0
for (let x of arr) {
    if (x > 5) {
        break
    }
    sum = sum + x
}
if (sum == 15) { println("PASS: break in for-of") } else { println("FAIL: break in for-of, got " + tostring(sum)) }

// Break only breaks inner loop
let count: int = 0
for (let a: int = 0; a < 3; a = a + 1) {
    for (let b: int = 0; b < 10; b = b + 1) {
        if (b >= 2) {
            break
        }
        count = count + 1
    }
}
if (count == 6) { println("PASS: break inner loop only") } else { println("FAIL: break inner loop only, got " + tostring(count)) }
