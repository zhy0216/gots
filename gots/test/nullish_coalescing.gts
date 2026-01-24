// Test: Nullish coalescing operator (??)

class Person {
    name: string
    constructor(name: string) {
        this.name = name
    }
}

// Test 1: Nullish with non-null value
let name: string | null = "Alice"
let result1: string = name ?? "Default"
println(result1)

// Test 2: Nullish with null value
let nullName: string | null = null
let result2: string = nullName ?? "Default"
println(result2)

// Test 3: Nullish with object
let person: Person | null = new Person("Bob")
let result3: Person = person ?? new Person("Default")
println(result3.name)

// Test 4: Nullish with null object
let nullPerson: Person | null = null
let result4: Person = nullPerson ?? new Person("Default")
println(result4.name)

// Test 5: Chained nullish
let a: string | null = null
let b: string | null = null
let c: string | null = "Charlie"
let result5: string = a ?? b ?? c ?? "None"
println(result5)

println("Done!")
