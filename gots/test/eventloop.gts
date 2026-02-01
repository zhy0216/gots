// Test file for JavaScript-compatible event loop

// Test callback functions
function onTimeout(): void {
    println("setTimeout fired")
}

function onMicrotask(): void {
    println("microtask ran")
}

// Test async function
async function testAsync(): Promise<void> {
    println("async start")
    await Promise.resolve(null)
    println("async after await")
}

// Main test execution
println("Event loop test started")
println("1. sync code")

// Queue a macrotask (setTimeout)
setTimeout(onTimeout, 10)

// Queue a microtask
queueMicrotask(onMicrotask)

println("2. still sync")

// Start async function
testAsync()

println("3. sync end")
