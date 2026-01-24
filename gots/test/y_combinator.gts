// Y Combinator Test for goTS
// The Y combinator enables anonymous recursion without explicit self-reference
// We use the Z combinator variant which works with strict (call-by-value) evaluation

// Since goTS doesn't support recursive types, we use a class wrapper
// to enable self-application: x(x) becomes x.call(x)
class FuncWrapper {
    call: Function
    constructor(f: Function) {
        this.call = f
    }
}

// The Z combinator (strict Y combinator)
// Mathematical form: Z = λf. (λx. f(λv. x(x)(v))) (λx. f(λv. x(x)(v)))
function Y(f: Function): Function {
    let wrapper: FuncWrapper = new FuncWrapper(function(x: FuncWrapper): Function {
        return f(function(v: int): int {
            return x.call(x)(v)
        })
    })
    return wrapper.call(wrapper)
}

// Factorial using Y combinator - no explicit recursion!
// The function doesn't call itself by name; recursion is achieved through Y
let factorial: Function = Y(function(rec: Function): Function {
    return function(n: int): int {
        if (n <= 1) {
            return 1
        }
        return n * rec(n - 1)
    }
})

println("Factorial using Y combinator:")
println(factorial(0))  // 1
println(factorial(1))  // 1
println(factorial(5))  // 120
println(factorial(10)) // 3628800

// Fibonacci using Y combinator
let fib: Function = Y(function(rec: Function): Function {
    return function(n: int): int {
        if (n < 2) {
            return n
        }
        return rec(n - 1) + rec(n - 2)
    }
})

println("Fibonacci using Y combinator:")
println(fib(0))  // 0
println(fib(1))  // 1
println(fib(10)) // 55
println(fib(15)) // 610

println("Y combinator tests completed!")
