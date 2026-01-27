// Test: Classes
// Based on test262 class tests

// Basic class
class Point {
    x: int
    y: int

    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }

    getX(): int {
        return this.x
    }

    getY(): int {
        return this.y
    }
}

let p: Point = new Point(3, 4)
if (p.x == 3) { println("PASS: class field access x") } else { println("FAIL: class field access x") }
if (p.y == 4) { println("PASS: class field access y") } else { println("FAIL: class field access y") }
if (p.getX() == 3) { println("PASS: class method getX") } else { println("FAIL: class method getX") }
if (p.getY() == 4) { println("PASS: class method getY") } else { println("FAIL: class method getY") }

// Method that uses this
class Counter {
    count: int

    constructor() {
        this.count = 0
    }

    increment(): void {
        this.count = this.count + 1
    }

    get(): int {
        return this.count
    }
}

let ctr: Counter = new Counter()
ctr.increment()
ctr.increment()
ctr.increment()
if (ctr.get() == 3) { println("PASS: method modifies this") } else { println("FAIL: method modifies this") }

// Class with computed method
class Rectangle {
    width: int
    height: int

    constructor(w: int, h: int) {
        this.width = w
        this.height = h
    }

    area(): int {
        return this.width * this.height
    }

    perimeter(): int {
        return 2 * (this.width + this.height)
    }
}

let rect: Rectangle = new Rectangle(5, 3)
if (rect.area() == 15) { println("PASS: area method") } else { println("FAIL: area method") }
if (rect.perimeter() == 16) { println("PASS: perimeter method") } else { println("FAIL: perimeter method") }

// Class inheritance
class Animal {
    name: string

    constructor(name: string) {
        this.name = name
    }

    speak(): string {
        return this.name + " makes a sound"
    }
}

class Dog extends Animal {
    constructor(name: string) {
        super(name)
    }

    speak(): string {
        return this.name + " barks"
    }
}

let animal: Animal = new Animal("Generic")
let dog: Dog = new Dog("Rex")

if (animal.speak() == "Generic makes a sound") { println("PASS: base class method") } else { println("FAIL: base class method") }
if (dog.speak() == "Rex barks") { println("PASS: overridden method") } else { println("FAIL: overridden method") }
if (dog.name == "Rex") { println("PASS: inherited field") } else { println("FAIL: inherited field") }
