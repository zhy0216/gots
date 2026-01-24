// Test: try/catch/throw error handling

function mayFail(shouldFail: boolean): int {
    if (shouldFail) {
        throw "Something went wrong!"
    }
    return 42
}

// Test 1: No error case
try {
    let result: int = mayFail(false)
    println(result)
} catch (e) {
    println("Caught an error!")
}

// Test 2: Error case
try {
    let result: int = mayFail(true)
    println(result)
} catch (e) {
    println("Caught error in test 2")
}

// Test 3: Throw in main
try {
    throw "Custom error"
} catch (e) {
    println("Caught custom error")
}

println("Done!")
