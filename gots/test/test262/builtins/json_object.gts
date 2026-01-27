// Test262-style tests for the JSON object
// Tests JSON.parse and JSON.stringify

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// ============================================
// JSON.stringify
// ============================================

// Stringify primitives
let strNum: string = JSON.stringify(42)
printResult("JSON.stringify(42) = '42'", strNum == "42")

let strBool: string = JSON.stringify(true)
printResult("JSON.stringify(true) = 'true'", strBool == "true")

let strStr: string = JSON.stringify("hello")
printResult("JSON.stringify('hello') = '\"hello\"'", strStr == "\"hello\"")

// Stringify arrays
let arr: int[] = [1, 2, 3]
let strArr: string = JSON.stringify(arr)
printResult("JSON.stringify([1,2,3]) = '[1,2,3]'", strArr == "[1,2,3]")

// ============================================
// JSON.parse
// ============================================

// Parse number
let parsedNum: number = JSON.parse("42")
printResult("JSON.parse('42') = 42", parsedNum == 42)

let parsedFloat: number = JSON.parse("3.14")
printResult("JSON.parse('3.14') = 3.14", parsedFloat == 3.14)

// Parse boolean
let parsedTrue: boolean = JSON.parse("true")
printResult("JSON.parse('true') = true", parsedTrue == true)

let parsedFalse: boolean = JSON.parse("false")
printResult("JSON.parse('false') = false", parsedFalse == false)

// Parse string
let parsedStr: string = JSON.parse("\"hello\"")
printResult("JSON.parse('\"hello\"') = 'hello'", parsedStr == "hello")

// Parse null
let parsedNull: null = JSON.parse("null")
printResult("JSON.parse('null') = null", parsedNull == null)

println("")
println("========== JSON Object Tests Complete ==========")
