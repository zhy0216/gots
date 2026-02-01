// Generics Test File

// Generic identity function
function identity<T>(x: T): T {
    return x
}

// Generic function with multiple type parameters
function pair<A, B>(a: A, b: B): A {
    println(a)
    println(b)
    return a
}

// Test generic functions with inference
let num: int = identity(42)
println(num)

let str: string = identity("hello")
println(str)

let result: int = pair(1, "test")
println(result)

// Test explicit type arguments on function calls
let x: number = identity<number>(99)
println(x)

let y: string = identity<string>("explicit")
println(y)

// Generic class
class Box<T> {
    value: T

    constructor(v: T) {
        this.value = v
    }

    get(): T {
        return this.value
    }

    set(v: T): void {
        this.value = v
    }
}

// Use generic class with int
let intBox: Box<int> = new Box(10)
println(intBox.get())
intBox.set(20)
println(intBox.get())

// Use generic class with string
let strBox: Box<string> = new Box("world")
println(strBox.get())

// Default type parameters
function makeValue<T = number>(x: T): T {
    return x
}

let defaultVal: int = makeValue(42)
println(defaultVal)

println("Generics test passed!")
