// Test regex literals
let re: RegExp = /hello/i
println(re.test("Hello World"))  // true
println(re.test("goodbye"))  // false

// Test regex with multiple flags
let patternRegex: RegExp = /world/im
println(patternRegex.test("Hello\nWorld"))  // true

// Test exec method
let digitRegex: RegExp = /\d+/
let result: string[] | null = digitRegex.exec("abc123def")
if (result != null) {
    println(result[0])  // 123
}

// Test regex without flags
let simpleRegex: RegExp = /test/
println(simpleRegex.test("testing"))  // true
