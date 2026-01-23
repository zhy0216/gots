// Higher-Order Functions and Currying Test for GoTS
// Demonstrates functional programming patterns

// Currying: transforming a function with multiple arguments
// into a sequence of functions each with a single argument
function curry_add(a: number): Function {
    return function(b: number): number {
        return a + b
    }
}

let add5: Function = curry_add(5)
let add10: Function = curry_add(10)

println("Currying:")
println(add5(3))   // 8
println(add10(3))  // 13

// Function composition: (f . g)(x) = f(g(x))
function compose(f: Function, g: Function): Function {
    return function(x: number): number {
        return f(g(x))
    }
}

function double(x: number): number {
    return x * 2
}

function increment(x: number): number {
    return x + 1
}

function square(x: number): number {
    return x * x
}

let doubleAndIncrement: Function = compose(increment, double)
let incrementAndDouble: Function = compose(double, increment)

println("Function composition:")
println(doubleAndIncrement(5))  // (5 * 2) + 1 = 11
println(incrementAndDouble(5))  // (5 + 1) * 2 = 12

// Triple composition
let tripleCompose: Function = compose(compose(double, increment), square)
println(tripleCompose(3))  // ((3^2) + 1) * 2 = 20

// Partial application using closures
function makeMultiplier(factor: number): Function {
    return function(x: number): number {
        return x * factor
    }
}

let triple: Function = makeMultiplier(3)
let quadruple: Function = makeMultiplier(4)

println("Partial application:")
println(triple(7))     // 21
println(quadruple(7))  // 28

// Higher-order array operations (manual implementation)
function forEach(arr: number[], f: Function): void {
    let i: number = 0
    while (i < len(arr)) {
        f(arr[i])
        i = i + 1
    }
}

function map(arr: number[], f: Function): number[] {
    let result: number[] = []
    let i: number = 0
    while (i < len(arr)) {
        push(result, f(arr[i]))
        i = i + 1
    }
    return result
}

function filter(arr: number[], predicate: Function): number[] {
    let result: number[] = []
    let i: number = 0
    while (i < len(arr)) {
        if (predicate(arr[i])) {
            push(result, arr[i])
        }
        i = i + 1
    }
    return result
}

function reduce(arr: number[], f: Function, initial: number): number {
    let acc: number = initial
    let i: number = 0
    while (i < len(arr)) {
        acc = f(acc, arr[i])
        i = i + 1
    }
    return acc
}

let numbers: number[] = [1, 2, 3, 4, 5]

println("Map (double each element):")
let doubled: number[] = map(numbers, double)
forEach(doubled, function(x: number): void {
    println(x)
})

println("Filter (keep even numbers):")
let evens: number[] = filter(numbers, function(x: number): boolean {
    return x % 2 == 0
})
forEach(evens, function(x: number): void {
    println(x)
})

println("Reduce (sum all elements):")
let sum: number = reduce(numbers, function(acc: number, x: number): number {
    return acc + x
}, 0)
println(sum)  // 15

println("Reduce (product of all elements):")
let product: number = reduce(numbers, function(acc: number, x: number): number {
    return acc * x
}, 1)
println(product)  // 120

// Chaining: sum of squares of even numbers
let data: number[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evenSquareSum: number = reduce(
    map(
        filter(data, function(x: number): boolean {
            return x % 2 == 0
        }),
        square
    ),
    function(acc: number, x: number): number {
        return acc + x
    },
    0
)
println("Sum of squares of even numbers 1-10:")
println(evenSquareSum)  // 4+16+36+64+100 = 220

println("Higher-order functions test completed!")
