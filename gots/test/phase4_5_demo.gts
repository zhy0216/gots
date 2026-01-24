// Example demonstrating Phase 4 & 5 features

// String functions
let message: string = "Hello, World!"
let words: string[] = split(message, ", ")
println("Split:", words)

let joined: string = join(words, " | ")
println("Joined:", joined)

let replaced: string = replace(message, "World", "goTS")
println("Replaced:", replaced)

let trimmed: string = trim("  spaces  ")
println("Trimmed:", trimmed)

if (startsWith(message, "Hello")) {
    println("Message starts with Hello")
}

if (endsWith(message, "World!")) {
    println("Message ends with World!")
}

if (includes(message, "World")) {
    println("Message includes World")
}

// Array functions
let numbers: int[] = [1, 2, 3, 4, 5]

// map
let doubled: int[] = map(numbers, (x: int): int => x * 2)
println("Doubled:", doubled)

// filter
let evens: int[] = filter(numbers, (x: int): boolean => x % 2 == 0)
println("Evens:", evens)

// reduce
let sum: int = reduce(numbers, 0, (acc: int, x: int): int => acc + x)
println("Sum:", sum)

// find
let found: int = find(numbers, (x: int): boolean => x > 3)
println("Found (>3):", found)

// findIndex
let idx: int = findIndex(numbers, (x: int): boolean => x > 3)
println("Index of first >3:", idx)

// some
let hasEven: boolean = some(numbers, (x: int): boolean => x % 2 == 0)
println("Has even number:", hasEven)

// every
let allPositive: boolean = every(numbers, (x: int): boolean => x > 0)
println("All positive:", allPositive)

// Complex example: processing data
let scores: int[] = [85, 92, 78, 95, 88]

let passed: int[] = filter(scores, (s: int): boolean => s >= 80)
let total: int = reduce(passed, 0, (acc: int, s: int): int => acc + s)
let average: float = tofloat(total) / tofloat(len(passed))

println("Passing scores:", passed)
println("Average of passing scores:", average)

// String manipulation chain
let text: string = "  the quick brown fox  "
text = trim(text)
let parts: string[] = split(text, " ")
let capitalized: string[] = map(parts, (w: string): string => w)
let result: string = join(capitalized, " ")
println("Processed text:", result)
