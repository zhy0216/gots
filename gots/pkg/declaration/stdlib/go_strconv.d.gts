// Declaration file for Go's strconv package
declare module "go:strconv" {
    function Itoa(i: int): string
    function Atoi(s: string): int
    function ParseInt(s: string, base: int, bitSize: int): int
    function ParseFloat(s: string, bitSize: int): float
    function FormatInt(i: int, base: int): string
    function FormatFloat(f: float, fmt: int, prec: int, bitSize: int): string
    function ParseBool(str: string): boolean
    function FormatBool(b: boolean): string
    function Quote(s: string): string
    function Unquote(s: string): string
}
