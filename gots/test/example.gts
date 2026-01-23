// Test program for GoTS

let x: number = 10
let y: number = 20
println(x + y)

function factorial(n: number): number {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))

// Test built-in functions
let arr: number[] = [1, 2, 3]
push(arr, 4)
println(len(arr))

println(sqrt(16))
println(abs(-5))
println(floor(3.7))
println(ceil(3.2))

// Test class
class Counter {
    count: number

    constructor() {
        this.count = 0
    }

    increment(): void {
        this.count = this.count + 1
    }

    getCount(): number {
        return this.count
    }
}

let counter: Counter = new Counter()
counter.increment()
counter.increment()
counter.increment()
println(counter.getCount())

println("All tests passed!")
