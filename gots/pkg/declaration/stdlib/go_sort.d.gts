// Declaration file for Go's sort package
declare module "go:sort" {
    // Sorting functions for basic types
    function Ints(x: int[]): void
    function Float64s(x: float[]): void
    function Strings(x: string[]): void

    // Check if sorted
    function IntsAreSorted(x: int[]): boolean
    function Float64sAreSorted(x: float[]): boolean
    function StringsAreSorted(x: string[]): boolean

    // Search functions (binary search)
    function SearchInts(a: int[], x: int): int
    function SearchFloat64s(a: float[], x: float): int
    function SearchStrings(a: string[], x: string): int

    // Reverse returns the reverse order
    function Reverse(data: any): any

    // Generic sort with comparison function
    function Slice(x: any[], less: (i: int, j: int) => boolean): void
    function SliceStable(x: any[], less: (i: int, j: int) => boolean): void
    function SliceIsSorted(x: any[], less: (i: int, j: int) => boolean): boolean

    // Sort interface
    interface Interface {
        Len(): int
        Less(i: int, j: int): boolean
        Swap(i: int, j: int): void
    }

    // General sort function
    function Sort(data: Interface): void
    function Stable(data: Interface): void
    function IsSorted(data: Interface): boolean
}
