// Test262-style tests for the Object static methods
// Tests Object.keys, Object.values, Object.entries

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// Note: Object static methods work with Map<string, T> types in goTS
// since plain object literals compile to Go structs

// ============================================
// Object.keys
// ============================================

let map1: Map<string, int> = new Map<string, int>()
map1.set("a", 1)
map1.set("b", 2)
map1.set("c", 3)

let keys: string[] = Object.keys(map1)
printResult("Object.keys returns array", len(keys) == 3)
printResult("Object.keys contains 'a'", keys.includes("a"))
printResult("Object.keys contains 'b'", keys.includes("b"))
printResult("Object.keys contains 'c'", keys.includes("c"))

// ============================================
// Object.values
// ============================================

let values: int[] = Object.values(map1)
printResult("Object.values returns array", len(values) == 3)
printResult("Object.values contains 1", values.includes(1))
printResult("Object.values contains 2", values.includes(2))
printResult("Object.values contains 3", values.includes(3))

// ============================================
// Object.entries
// ============================================

// Object.entries returns array of [key, value] pairs
// In goTS, this maps to [][]interface{} or similar

// ============================================
// Object.assign (creates new map by copying)
// ============================================

let target: Map<string, int> = new Map<string, int>()
target.set("x", 10)

let source: Map<string, int> = new Map<string, int>()
source.set("y", 20)

// After assign, target should have both x and y
let result: Map<string, int> = Object.assign(target, source)
printResult("Object.assign merges maps", result.has("x") && result.has("y"))
printResult("Object.assign preserves original value", result.get("x") == 10)
printResult("Object.assign adds new value", result.get("y") == 20)

// ============================================
// Object.hasOwn
// ============================================

let map2: Map<string, string> = new Map<string, string>()
map2.set("name", "Alice")

printResult("Object.hasOwn returns true for existing key", Object.hasOwn(map2, "name") == true)
printResult("Object.hasOwn returns false for non-existing key", Object.hasOwn(map2, "age") == false)

println("")
println("========== Object Static Methods Tests Complete ==========")
