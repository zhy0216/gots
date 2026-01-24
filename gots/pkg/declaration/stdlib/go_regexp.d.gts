// Declaration file for Go's regexp package
declare module "go:regexp" {
    // MatchString checks if pattern matches string
    function MatchString(pattern: string, s: string): boolean

    // Compile parses a regular expression
    function Compile(expr: string): Regexp | null

    // MustCompile parses a regular expression and panics if invalid
    function MustCompile(expr: string): Regexp

    // QuoteMeta escapes all regex metacharacters
    function QuoteMeta(s: string): string

    // Regexp represents a compiled regular expression
    interface Regexp {
        MatchString(s: string): boolean
        Match(b: byte[]): boolean
        FindString(s: string): string
        FindStringIndex(s: string): int[] | null
        FindStringSubmatch(s: string): string[]
        FindAllString(s: string, n: int): string[]
        FindAllStringIndex(s: string, n: int): int[][]
        FindAllStringSubmatch(s: string, n: int): string[][]
        ReplaceAllString(src: string, repl: string): string
        ReplaceAllLiteralString(src: string, repl: string): string
        ReplaceAllStringFunc(src: string, repl: (s: string) => string): string
        Split(s: string, n: int): string[]
        SubexpNames(): string[]
        NumSubexp(): int
        String(): string
    }
}
