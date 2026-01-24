// Test new array methods and Set type

// Array .length property
let arr: int[] = [1, 2, 3, 4, 5]
println(arr.length)

// Array methods
let doubled: int[] = arr.map(function(v: int, i: int): int { return v * 2 })
println(doubled)

let evens: int[] = arr.filter(function(v: int, i: int): boolean { return v % 2 == 0 })
println(evens)

let sum: int = arr.reduce(function(acc: int, v: int, i: int): int { return acc + v }, 0)
println(sum)

// indexOf, includes
println(arr.indexOf(3))
println(arr.includes(3))

// Set type
let s: Set<int> = new Set<int>()
s.add(1)
s.add(2)
s.add(3)
println(s.size)
println(s.has(2))
println(s.has(10))

// Map methods
let m: Map<string, int> = new Map<string, int>()
m.set("a", 1)
m.set("b", 2)
println(m.size)
println(m.get("a"))
println(m.has("b"))
