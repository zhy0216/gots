// Test: Closures
// Based on test262 closure/scope tests

// Basic closure
function makeCounter(): Function {
    let count: int = 0
    return function(): int {
        count = count + 1
        return count
    }
}

let counter: Function = makeCounter()
let c1: int = counter()
let c2: int = counter()
let c3: int = counter()
if (c1 == 1 && c2 == 2 && c3 == 3) { println("PASS: basic closure") } else { println("FAIL: basic closure") }

// Independent closures
let counter1: Function = makeCounter()
let counter2: Function = makeCounter()
let a1: int = counter1()
let a2: int = counter1()
let b1: int = counter2()
if (a1 == 1 && a2 == 2 && b1 == 1) { println("PASS: independent closures") } else { println("FAIL: independent closures") }

// Closure capturing parameter
function makeAdder(x: int): Function {
    return function(y: int): int {
        return x + y
    }
}

let add5: Function = makeAdder(5)
let add10: Function = makeAdder(10)
if (add5(3) == 8) { println("PASS: closure captures param 1") } else { println("FAIL: closure captures param 1") }
if (add10(3) == 13) { println("PASS: closure captures param 2") } else { println("FAIL: closure captures param 2") }

// Nested closure
function outer(a: int): Function {
    return function(b: int): Function {
        return function(c: int): int {
            return a + b + c
        }
    }
}

let step1: Function = outer(1)
let step2: Function = step1(2)
let final: int = step2(3)
if (final == 6) { println("PASS: nested closure") } else { println("FAIL: nested closure") }

// Closure modifying outer variable
function createAccumulator(): Function {
    let total: int = 0
    return function(n: int): int {
        total = total + n
        return total
    }
}

let acc: Function = createAccumulator()
acc(5)
acc(10)
let accResult: int = acc(3)
if (accResult == 18) { println("PASS: closure modifies outer") } else { println("FAIL: closure modifies outer, got " + tostring(accResult)) }
