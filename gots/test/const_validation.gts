// Test: const variables work correctly
const PI: float = 3.14159
const MAX_SIZE: int = 100
const NAME: string = "GoTS"

// This should compile - using const values
println(PI)
println(MAX_SIZE)
println(NAME)

// let variables can be reassigned
let counter: int = 0
counter = counter + 1
counter += 5
counter++
println(counter)
