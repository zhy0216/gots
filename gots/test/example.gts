// Test program for GoTS

let x: int = 10
let y: int = 20
println(x + y)

function factorial(n: int): int {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))

// Test built-in functions
let arr: int[] = [1, 2, 3]
push(arr, 4)
println(len(arr))

println(sqrt(16.0))
println(abs(-5.0))
println(floor(3.7))
println(ceil(3.2))

// Test class
class Counter {
    count: int

    constructor() {
        this.count = 0
    }

    increment(): void {
        this.count = this.count + 1
    }

    getCount(): int {
        return this.count
    }
}

let counter: Counter = new Counter()
counter.increment()
counter.increment()
counter.increment()
println(counter.getCount())

println("All tests passed!")
