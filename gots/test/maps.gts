// Test Map feature

// Create an empty map
let m: Map<string, int> = {}

// Set values
m.set("one", 1)
m.set("two", 2)
m.set("three", 3)

// Get values
let val: int = m.get("one")
println(val)

// Check if key exists
if (m.has("two")) {
    println("has two")
}

// Delete key
m.delete("one")

// Get keys and values
let keys: string[] = m.keys()
let values: int[] = m.values()

println("Keys:")
for (let k of keys) {
    println(k)
}

println("Values:")
for (let v of values) {
    println(v)
}
