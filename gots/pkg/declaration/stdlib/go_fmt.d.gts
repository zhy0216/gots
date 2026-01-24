// Declaration file for Go's fmt package
declare module "go:fmt" {
    function Println(...args: any[]): void
    function Print(...args: any[]): void
    function Printf(format: string, ...args: any[]): void
    function Sprintf(format: string, ...args: any[]): string
    function Errorf(format: string, ...args: any[]): Error
    function Sscanf(str: string, format: string, ...args: any[]): int
    function Fscanf(r: any, format: string, ...args: any[]): int
}
