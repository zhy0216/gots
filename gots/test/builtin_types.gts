// Test: Builtin function typing

// Test 1: len returns int
let arr: int[] = [1, 2, 3]
let length: int = len(arr)
println(tostring(length))

// Test 2: push - verify it works
push(arr, 4)
println(tostring(len(arr)))

// Test 3: pop returns element type
let last: int = pop(arr)
println(tostring(last))

// Test 4: len on string
let str: string = "hello"
let strLen: int = len(str)
println(tostring(strLen))

// Test 5: Assign pop result to typed variable
let nums: float[] = [1.5, 2.5, 3.5]
let lastNum: float = pop(nums)
println(tostring(lastNum))

println("Done!")
