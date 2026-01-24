// Test: For-of over strings

// Test 1: Basic string iteration
let greeting: string = "Hello"
for (let char of greeting) {
    println(char)
}

// Test 2: String iteration with concatenation
let result: string = ""
for (let c of "abc") {
    result = result + c + "-"
}
println(result)

// Test 3: Empty string (should not enter loop)
let emptyCount: int = 0
for (let c of "") {
    emptyCount = emptyCount + len(c)
}
println("Empty count: " + tostring(emptyCount))

println("Done!")
