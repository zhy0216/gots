// Test: Class inheritance with super() call

class Animal {
    name: string
    age: int
    constructor(name: string, age: int) {
        this.name = name
        this.age = age
    }
    speak(): void {
        println(this.name + " makes a sound")
    }
}

class Dog extends Animal {
    breed: string
    constructor(name: string, age: int, breed: string) {
        super(name, age)
        this.breed = breed
    }
    speak(): void {
        println(this.name + " barks")
    }
    info(): void {
        println(this.name + " is a " + tostring(this.age) + " year old " + this.breed)
    }
}

// Test 1: Create dog and access inherited fields
let dog: Dog = new Dog("Buddy", 5, "Golden Retriever")
println(dog.name)
println(tostring(dog.age))
println(dog.breed)

// Test 2: Call overridden method
dog.speak()

// Test 3: Call child-specific method
dog.info()

// Test 4: Access inherited fields directly
let animal: Animal = new Animal("Generic", 10)
animal.speak()

println("Done!")
