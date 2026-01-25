// Test async/await and Promise functionality

// Simple async function
async function fetchNumber(): Promise<int> {
    return 42
}

// Async function with computation
async function computeSum(a: int, b: int): Promise<int> {
    return a + b
}

// Using await inside another async function
async function testAwait(): Promise<int> {
    let num: int = await fetchNumber()
    return num * 2
}

// Async function to test everything
async function runTests(): Promise<void> {
    // Test basic async
    let result: int = await testAwait()
    println(result)  // Should print 84

    // Test direct async call
    let sum: int = await computeSum(10, 20)
    println(sum)  // Should print 30

    // Test sequential awaits
    let a: int = await fetchNumber()
    let b: int = await fetchNumber()
    println(a + b)  // Should print 84
}

// Simple sync test
function simpleTest(): void {
    println("Testing simple async...")
    let num: int = 10
    println(num)
}

// Run tests
simpleTest()
println("Async tests will run in background...")
runTests()  // The Promise executes asynchronously
