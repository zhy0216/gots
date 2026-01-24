// Example demonstrating String and Array Methods

// String methods
let message: string = "Hello, World!"
let words: string[] = message.split(", ")
println("Split:", words)

let joined: string = words.join(" | ")
println("Joined:", joined)

let replaced: string = message.replace("World", "goTS")
println("Replaced:", replaced)

let trimmed: string = "  spaces  ".trim()
println("Trimmed:", trimmed)

if (message.startsWith("Hello")) {
    println("Message starts with Hello")
}

if (message.endsWith("World!")) {
    println("Message ends with World!")
}

if (message.includes("World")) {
    println("Message includes World")
}

let lower: string = message.toLowerCase()
println("Lowercase:", lower)

let upper: string = message.toUpperCase()
println("Uppercase:", upper)

let idx: int = message.indexOf("World")
println("Index of 'World':", idx)

// Array methods
let numbers: int[] = [1, 2, 3, 4, 5]

// map
let doubled: int[] = numbers.map((x: int): int => x * 2)
println("Doubled:", doubled)

// filter
let evens: int[] = numbers.filter((x: int): boolean => x % 2 == 0)
println("Evens:", evens)

// reduce
let sum: int = numbers.reduce((acc: int, x: int): int => acc + x, 0)
println("Sum:", sum)

// find
let found: int | null = numbers.find((x: int): boolean => x > 3)
println("Found (>3):", found)

// findIndex
let foundIdx: int = numbers.findIndex((x: int): boolean => x > 3)
println("Index of first >3:", foundIdx)

// some
let hasEven: boolean = numbers.some((x: int): boolean => x % 2 == 0)
println("Has even number:", hasEven)

// every
let allPositive: boolean = numbers.every((x: int): boolean => x > 0)
println("All positive:", allPositive)

// Complex example: processing data
let scores: int[] = [85, 92, 78, 95, 88]

let passed: int[] = scores.filter((s: int): boolean => s >= 80)
let total: int = passed.reduce((acc: int, s: int): int => acc + s, 0)
let average: float = tofloat(total) / tofloat(len(passed))

println("Passing scores:", passed)
println("Average of passing scores:", average)

// String manipulation chain
let text: string = "  the quick brown fox  "
text = text.trim()
let parts: string[] = text.split(" ")
let result: string = parts.join(" ")
println("Processed text:", result)

// Array method chaining
let data: int[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let processedData: int[] = data.filter((n: int): boolean => n % 2 == 0).map((n: int): int => n * n)
println("Squared evens:", processedData)
