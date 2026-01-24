// Math module - export some math functions
export function add(a: int, b: int): int {
    return a + b
}

export function multiply(a: int, b: int): int {
    return a * b
}

export class Vector {
    x: int
    y: int

    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }

    add(other: Vector): Vector {
        return new Vector(this.x + other.x, this.y + other.y)
    }

    toString(): string {
        return "(" + tostring(this.x) + ", " + tostring(this.y) + ")"
    }
}
