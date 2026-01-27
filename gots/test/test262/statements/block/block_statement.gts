// Test: Block statement
// Based on test262 block tests

// Basic block
let x: int = 1
{
    let y: int = 2
    x = y
}
if (x == 2) { println("PASS: basic block") } else { println("FAIL: basic block") }

// Nested blocks
x = 0
{
    x = x + 1
    {
        x = x + 1
        {
            x = x + 1
        }
    }
}
if (x == 3) { println("PASS: nested blocks") } else { println("FAIL: nested blocks") }

// Block scope shadowing
let a: int = 10
{
    let a: int = 20
    if (a == 20) { println("PASS: inner block shadow") } else { println("FAIL: inner block shadow") }
}
if (a == 10) { println("PASS: outer preserved after block") } else { println("FAIL: outer preserved after block") }

// Block with control flow
x = 0
if (true) {
    {
        x = 5
    }
}
if (x == 5) { println("PASS: block in if") } else { println("FAIL: block in if") }

// Multiple statements in block
let sum: int = 0
{
    let a: int = 1
    let b: int = 2
    let c: int = 3
    sum = a + b + c
}
if (sum == 6) { println("PASS: multiple statements in block") } else { println("FAIL: multiple statements in block") }
