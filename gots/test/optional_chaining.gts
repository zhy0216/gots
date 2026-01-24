// Test: Optional chaining codegen

class Person {
    name: string
    age: int
    constructor(name: string, age: int) {
        this.name = name
        this.age = age
    }
}

// Test 1: Non-null case - direct property access on non-nullable
let person: Person = new Person("Alice", 30)
println(person.name)

// Test 2: Nullable but with value - using optional chaining
let maybePerson: Person | null = new Person("Bob", 25)
if (maybePerson != null) {
    println(maybePerson.name)
}

// Test 3: Null case
let nullPerson: Person | null = null
if (nullPerson?.name != null) {
    println("Has name")
} else {
    println("No name (null)")
}

println("Done!")
