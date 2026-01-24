// Declaration file for Go's time package
declare module "go:time" {
    // Duration type (nanoseconds)
    type Duration = int

    // Duration constants
    const Nanosecond: int
    const Microsecond: int
    const Millisecond: int
    const Second: int
    const Minute: int
    const Hour: int

    // Time represents an instant in time
    interface Time {
        Unix(): int
        UnixMilli(): int
        UnixNano(): int
        Year(): int
        Month(): int
        Day(): int
        Hour(): int
        Minute(): int
        Second(): int
        Weekday(): int
        YearDay(): int
        Format(layout: string): string
        String(): string
        Add(d: int): Time
        Sub(u: Time): int
        Before(u: Time): boolean
        After(u: Time): boolean
        Equal(u: Time): boolean
        IsZero(): boolean
    }

    // Functions
    function Now(): Time
    function Sleep(d: int): void
    function Since(t: Time): int
    function Until(t: Time): int
    function Parse(layout: string, value: string): Time
    function ParseDuration(s: string): int
    function Date(year: int, month: int, day: int, hour: int, min: int, sec: int, nsec: int, loc: any): Time
    function Unix(sec: int, nsec: int): Time
    function UnixMilli(msec: int): Time
}
