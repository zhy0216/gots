// Test262-style tests for additional Array methods
// Tests array methods not yet implemented in goTS

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// ============================================
// Array.at (negative index support)
// ============================================

let arr: int[] = [10, 20, 30, 40, 50]

printResult("at(0) = 10", arr.at(0) == 10)
printResult("at(2) = 30", arr.at(2) == 30)
printResult("at(-1) = 50 (last)", arr.at(-1) == 50)
printResult("at(-2) = 40", arr.at(-2) == 40)

// ============================================
// Array.lastIndexOf
// ============================================

let arr2: int[] = [1, 2, 3, 2, 1]
printResult("lastIndexOf(2) = 3", arr2.lastIndexOf(2) == 3)
printResult("lastIndexOf(1) = 4", arr2.lastIndexOf(1) == 4)
printResult("lastIndexOf(5) = -1 (not found)", arr2.lastIndexOf(5) == -1)

// ============================================
// Array.fill
// ============================================

let fillArr: int[] = [1, 2, 3, 4, 5]
fillArr.fill(0)
printResult("fill(0) fills entire array", fillArr[0] == 0 && fillArr[4] == 0)

let fillArr2: int[] = [1, 2, 3, 4, 5]
fillArr2.fill(9, 1, 4)
printResult("fill(9, 1, 4) fills partial range", fillArr2[0] == 1 && fillArr2[1] == 9 && fillArr2[3] == 9 && fillArr2[4] == 5)

// ============================================
// Array.copyWithin
// ============================================

let copyArr: int[] = [1, 2, 3, 4, 5]
copyArr.copyWithin(0, 3)
printResult("copyWithin(0, 3) copies from index 3 to start", copyArr[0] == 4 && copyArr[1] == 5)

// ============================================
// Array.flat (simple 1-level flatten)
// ============================================

// Note: goTS arrays are typed, so flat only works with same-type nested arrays
// This is a simplified test for conceptual support

// ============================================
// Array.isArray (static method)
// ============================================

let testArr: int[] = [1, 2, 3]
printResult("Array.isArray([1,2,3]) = true", Array.isArray(testArr) == true)

let notArr: string = "hello"
printResult("Array.isArray('hello') = false", Array.isArray(notArr) == false)

println("")
println("========== Array Additional Methods Tests Complete ==========")
