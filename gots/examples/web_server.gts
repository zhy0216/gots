// Flask-style parameterized decorator example
// Demonstrates @obj.method(args) decorator syntax

class App {
    routes: string[]

    constructor() {
        this.routes = []
    }

    get(path: string): Function {
        let routes: string[] = this.routes
        return function(handler: Function): Function {
            push(routes, "GET " + path)
            println("Registered GET " + path)
            return handler
        }
    }

    post(path: string): Function {
        let routes: string[] = this.routes
        return function(handler: Function): Function {
            push(routes, "POST " + path)
            println("Registered POST " + path)
            return handler
        }
    }

    listen(port: int): void {
        println(`Server has ${len(this.routes)} routes`)
        println(`Listening on :${port}`)
    }
}

const app = new App()

@app.get("/hello")
function hello(name: string): string {
    return `hello, ${name}`
}

@app.get("/goodbye")
function goodbye(name: string): string {
    return `goodbye, ${name}`
}

app.listen(8080)
