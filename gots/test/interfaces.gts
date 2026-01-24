// Test interface declaration and structural typing

interface Drawable {
    draw(): void
    getArea(): float
}

class Circle {
    radius: float

    constructor(radius: float) {
        this.radius = radius
    }

    draw(): void {
        println("Drawing circle with radius " + tostring(this.radius))
    }

    getArea(): float {
        return 3.14159 * this.radius * this.radius
    }
}

class Rectangle {
    width: float
    height: float

    constructor(width: float, height: float) {
        this.width = width
        this.height = height
    }

    draw(): void {
        println("Drawing rectangle " + tostring(this.width) + "x" + tostring(this.height))
    }

    getArea(): float {
        return this.width * this.height
    }
}

// Test with class types
let circle: Circle = new Circle(5.0)
circle.draw()
println("Circle area: " + tostring(circle.getArea()))

let rect: Rectangle = new Rectangle(4.0, 3.0)
rect.draw()
println("Rectangle area: " + tostring(rect.getArea()))

// Test interface type assignment (structural typing)
let shape1: Drawable = circle
let shape2: Drawable = rect

shape1.draw()
shape2.draw()
println("Shape 1 area via interface: " + tostring(shape1.getArea()))
println("Shape 2 area via interface: " + tostring(shape2.getArea()))
