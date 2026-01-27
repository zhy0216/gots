// Test262-style tests for the Date object
// Tests Date constructor and methods

function printResult(name: string, condition: boolean): void {
    if (condition) {
        println("PASS: " + name)
    } else {
        println("FAIL: " + name)
    }
}

// ============================================
// Date.now() static method
// ============================================

let timestamp: number = Date.now()
printResult("Date.now() returns a number", timestamp > 0)

// ============================================
// Date constructor with timestamp
// ============================================

// Create date from timestamp (Jan 1, 2024 00:00:00 UTC = 1704067200000)
let date1: Date = new Date(1704067200000)
printResult("new Date(timestamp) creates date", date1.getTime() == 1704067200000)

// ============================================
// Date getter methods
// ============================================

// Use a known date: Jan 15, 2024 14:30:45 UTC (1705328445000)
let knownDate: Date = new Date(1705328445000)

// getTime
printResult("getTime() returns timestamp", knownDate.getTime() == 1705328445000)

// getFullYear
let year: int = knownDate.getFullYear()
printResult("getFullYear() returns 2024", year == 2024)

// getMonth (0-indexed: January = 0)
let month: int = knownDate.getMonth()
printResult("getMonth() returns 0 for January", month == 0)

// getDate (day of month)
let day: int = knownDate.getDate()
printResult("getDate() returns 15", day == 15)

// getDay (day of week: Monday = 1)
let dayOfWeek: int = knownDate.getDay()
printResult("getDay() returns day of week", dayOfWeek >= 0 && dayOfWeek <= 6)

// getHours (UTC)
let hours: int = knownDate.getHours()
printResult("getHours() returns hours", hours >= 0 && hours <= 23)

// getMinutes
let minutes: int = knownDate.getMinutes()
printResult("getMinutes() returns valid minutes", minutes >= 0 && minutes <= 59)

// getSeconds
let seconds: int = knownDate.getSeconds()
printResult("getSeconds() returns valid seconds", seconds >= 0 && seconds <= 59)

// getMilliseconds
let ms: int = knownDate.getMilliseconds()
printResult("getMilliseconds() returns 0", ms == 0)

// ============================================
// Date setter methods
// ============================================

let mutableDate: Date = new Date(1704067200000)

// setFullYear
mutableDate.setFullYear(2025)
printResult("setFullYear(2025) updates year", mutableDate.getFullYear() == 2025)

// setMonth
mutableDate.setMonth(5)  // June (0-indexed)
printResult("setMonth(5) sets to June", mutableDate.getMonth() == 5)

// setDate
mutableDate.setDate(20)
printResult("setDate(20) sets day of month", mutableDate.getDate() == 20)

// setHours
mutableDate.setHours(10)
printResult("setHours(10) sets hours", mutableDate.getHours() == 10)

// setMinutes
mutableDate.setMinutes(30)
printResult("setMinutes(30) sets minutes", mutableDate.getMinutes() == 30)

// setSeconds
mutableDate.setSeconds(15)
printResult("setSeconds(15) sets seconds", mutableDate.getSeconds() == 15)

// ============================================
// Date string methods
// ============================================

let dateForStr: Date = new Date(1704067200000)
let isoStr: string = dateForStr.toISOString()
printResult("toISOString() returns ISO format", isoStr.includes("2024"))

println("")
println("========== Date Object Tests Complete ==========")
