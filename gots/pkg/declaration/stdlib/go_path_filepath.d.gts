// Declaration file for Go's path/filepath package
declare module "go:path/filepath" {
    // Path separator
    const Separator: string
    const ListSeparator: string

    // Functions
    function Join(...elem: string[]): string
    function Split(path: string): (string, string)
    function Dir(path: string): string
    function Base(path: string): string
    function Ext(path: string): string
    function Clean(path: string): string
    function Abs(path: string): string
    function Rel(basepath: string, targpath: string): string
    function IsAbs(path: string): boolean
    function Match(pattern: string, name: string): boolean
    function Glob(pattern: string): string[]
    function EvalSymlinks(path: string): string
    function FromSlash(path: string): string
    function ToSlash(path: string): string
    function VolumeName(path: string): string
    function SplitList(path: string): string[]
    function HasPrefix(p: string, prefix: string): boolean
}
