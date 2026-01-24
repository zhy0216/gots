// Declaration file for Go's strings package
declare module "go:strings" {
    function Join(a: string[], sep: string): string
    function Split(s: string, sep: string): string[]
    function Contains(s: string, substr: string): boolean
    function HasPrefix(s: string, prefix: string): boolean
    function HasSuffix(s: string, suffix: string): boolean
    function ToUpper(s: string): string
    function ToLower(s: string): string
    function TrimSpace(s: string): string
    function Replace(s: string, old: string, new: string, n: int): string
    function ReplaceAll(s: string, old: string, new: string): string
    function Index(s: string, substr: string): int
    function Count(s: string, substr: string): int
    function Trim(s: string, cutset: string): string
    function TrimPrefix(s: string, prefix: string): string
    function TrimSuffix(s: string, suffix: string): string
    function Repeat(s: string, count: int): string
    function Fields(s: string): string[]
    function SplitN(s: string, sep: string, n: int): string[]
}
