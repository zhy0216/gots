// Test: Continue statement
// Based on test262 continue tests

// Continue in while
let i: int = 0
let sum: int = 0
while (i < 10) {
    i = i + 1
    if (i % 2 == 0) {
        continue
    }
    sum = sum + i
}
if (sum == 25) { println("PASS: continue in while") } else { println("FAIL: continue in while, got " + tostring(sum)) }

// Continue in for
sum = 0
for (let j: int = 1; j <= 10; j = j + 1) {
    if (j % 2 == 0) {
        continue
    }
    sum = sum + j
}
if (sum == 25) { println("PASS: continue in for") } else { println("FAIL: continue in for, got " + tostring(sum)) }

// Continue in for-of
let arr: int[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
sum = 0
for (let x of arr) {
    if (x % 2 == 0) {
        continue
    }
    sum = sum + x
}
if (sum == 25) { println("PASS: continue in for-of") } else { println("FAIL: continue in for-of, got " + tostring(sum)) }

// Continue skips rest of iteration
let result: string = ""
for (let n: int = 1; n <= 5; n = n + 1) {
    if (n == 3) {
        continue
    }
    result = result + tostring(n)
}
if (result == "1245") { println("PASS: continue skips correctly") } else { println("FAIL: continue skips correctly, got " + result) }
