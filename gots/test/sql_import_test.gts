import { connect } from "sql"

interface User {
    id: int
    name: string
}

const db = connect("./test.db")

db.sql`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL)`

db.sql`INSERT INTO users (name) VALUES (${"Alice"})`

const users = db.sql<User[]>`SELECT id, name FROM users`
println("Users: " + tostring(len(users)))

db.close()
