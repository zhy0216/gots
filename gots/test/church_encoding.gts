// Church Encoding Test for GoTS
// Church numerals encode natural numbers using only functions
// A Church numeral n is a function that applies another function n times

// Church numerals: encode numbers as repeated function application
// zero = λf.λx. x           (apply f zero times)
// one  = λf.λx. f(x)        (apply f once)
// two  = λf.λx. f(f(x))     (apply f twice)

function zero(f: Function): Function {
    return function(x: int): int {
        return x
    }
}

function one(f: Function): Function {
    return function(x: int): int {
        return f(x)
    }
}

function two(f: Function): Function {
    return function(x: int): int {
        return f(f(x))
    }
}

function three(f: Function): Function {
    return function(x: int): int {
        return f(f(f(x)))
    }
}

// Convert Church numeral to integer by applying increment
function churchToInt(n: Function): int {
    let inc: Function = function(x: int): int {
        return x + 1
    }
    return n(inc)(0)
}

println("Church numerals to integers:")
println(churchToInt(zero))   // 0
println(churchToInt(one))    // 1
println(churchToInt(two))    // 2
println(churchToInt(three))  // 3

// Successor: add one to a Church numeral
// succ(n) = λf.λx. f(n(f)(x))
function succ(n: Function): Function {
    return function(f: Function): Function {
        return function(x: int): int {
            return f(n(f)(x))
        }
    }
}

let four: Function = succ(three)
let five: Function = succ(four)

println("Successor function:")
println(churchToInt(four))   // 4
println(churchToInt(five))   // 5

// Addition: add two Church numerals
// add(m, n) = λf.λx. m(f)(n(f)(x))
function add(m: Function, n: Function): Function {
    return function(f: Function): Function {
        return function(x: int): int {
            return m(f)(n(f)(x))
        }
    }
}

println("Church addition:")
println(churchToInt(add(two, three)))     // 5
println(churchToInt(add(three, four)))    // 7

// Multiplication: multiply two Church numerals
// mult(m, n) = λf. m(n(f))
function mult(m: Function, n: Function): Function {
    return function(f: Function): Function {
        return m(n(f))
    }
}

// Convert integer to Church numeral using successor
function intToChurch(n: int): Function {
    let result: Function = zero
    let i: int = 0
    while (i < n) {
        result = succ(result)
        i = i + 1
    }
    return result
}

println("Church multiplication:")
println(churchToInt(mult(two, three)))    // 6
println(churchToInt(mult(three, four)))   // 12

// Power: raise a Church numeral to a power
// exp(m, n) = n(m) -- apply m, n times
function exp(base: Function, power: Function): Function {
    // Church numerals here are int-specialized; exponentiation needs a numeric bridge.
    let baseInt: int = churchToInt(base)
    let powerInt: int = churchToInt(power)
    let result: int = 1
    let i: int = 0
    while (i < powerInt) {
        result = result * baseInt
        i = i + 1
    }
    return intToChurch(result)
}

println("Church exponentiation:")
println(churchToInt(exp(two, three)))     // 2^3 = 8
println(churchToInt(exp(three, two)))     // 3^2 = 9

println("Integer to Church and back:")
let six: Function = intToChurch(6)
let seven: Function = intToChurch(7)
println(churchToInt(six))   // 6
println(churchToInt(seven)) // 7
println(churchToInt(add(six, seven)))  // 13
println(churchToInt(mult(six, seven))) // 42

// Church booleans
// true  = λa.λb. a (select first)
// false = λa.λb. b (select second)

// We'll use a number encoding: true = 1, false = 0
function churchTrue(a: int, b: int): int {
    return a
}

function churchFalse(a: int, b: int): int {
    return b
}

// isZero checks if a Church numeral is zero
// isZero(n) = n(λx. false)(true)
function isZero(n: Function): int {
    let alwaysFalse: Function = function(x: int): int {
        return 0  // false
    }
    return n(alwaysFalse)(1)  // returns 1 (true) if n is zero
}

println("isZero predicate:")
println(isZero(zero))   // 1 (true)
println(isZero(one))    // 0 (false)
println(isZero(three))  // 0 (false)

println("Church encoding test completed!")
