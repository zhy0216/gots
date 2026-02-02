interface User {
    id: int
    name: string
}

const db = connect("./test.db")

db.sql`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL)`

db.begin(function(tx: Transaction): void {
    tx.sql`INSERT INTO users (name) VALUES (${"Alice"})`
    tx.sql`INSERT INTO users (name) VALUES (${"Bob"})`
})

const users = db.sql<User[]>`SELECT id, name FROM users ORDER BY id`
for (let i: int = 0; i < len(users); i = i + 1) {
    println(users[i].name)
}

db.close()
