// Test export functionality
export function greet(name: string): string {
    return "Hello, " + name + "!"
}

export let PI: float = 3.14159

println(greet("World"))
println("PI = " + tostring(PI))
