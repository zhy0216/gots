// Declaration file for Go's bytes package
declare module "go:bytes" {
    // Buffer is a variable-sized buffer of bytes
    interface Buffer {
        Bytes(): byte[]
        String(): string
        Len(): int
        Cap(): int
        Truncate(n: int): void
        Reset(): void
        Grow(n: int): void
        Write(p: byte[]): int
        WriteString(s: string): int
        WriteByte(c: byte): Error | null
        WriteRune(r: int): int
        Read(p: byte[]): int
        ReadByte(): byte
        ReadRune(): (int, int)
        ReadBytes(delim: byte): byte[]
        ReadString(delim: byte): string
        Next(n: int): byte[]
        UnreadByte(): Error | null
        UnreadRune(): Error | null
    }

    // Reader reads from a byte slice
    interface Reader {
        Len(): int
        Size(): int
        Read(b: byte[]): int
        ReadAt(b: byte[], off: int): int
        ReadByte(): byte
        UnreadByte(): Error | null
        ReadRune(): (int, int)
        UnreadRune(): Error | null
        Seek(offset: int, whence: int): int
        Reset(b: byte[]): void
    }

    // Functions
    function NewBuffer(buf: byte[]): Buffer
    function NewBufferString(s: string): Buffer
    function NewReader(b: byte[]): Reader

    // String operations on byte slices
    function Contains(b: byte[], subslice: byte[]): boolean
    function ContainsAny(b: byte[], chars: string): boolean
    function ContainsRune(b: byte[], r: int): boolean
    function Count(s: byte[], sep: byte[]): int
    function Equal(a: byte[], b: byte[]): boolean
    function EqualFold(s: byte[], t: byte[]): boolean
    function HasPrefix(s: byte[], prefix: byte[]): boolean
    function HasSuffix(s: byte[], suffix: byte[]): boolean
    function Index(s: byte[], sep: byte[]): int
    function IndexAny(s: byte[], chars: string): int
    function IndexByte(b: byte[], c: byte): int
    function IndexRune(s: byte[], r: int): int
    function Join(s: byte[][], sep: byte[]): byte[]
    function LastIndex(s: byte[], sep: byte[]): int
    function LastIndexAny(s: byte[], chars: string): int
    function LastIndexByte(s: byte[], c: byte): int
    function Repeat(b: byte[], count: int): byte[]
    function Replace(s: byte[], old: byte[], new: byte[], n: int): byte[]
    function ReplaceAll(s: byte[], old: byte[], new: byte[]): byte[]
    function Split(s: byte[], sep: byte[]): byte[][]
    function SplitN(s: byte[], sep: byte[], n: int): byte[][]
    function Fields(s: byte[]): byte[][]
    function Title(s: byte[]): byte[]
    function ToLower(s: byte[]): byte[]
    function ToUpper(s: byte[]): byte[]
    function ToTitle(s: byte[]): byte[]
    function Trim(s: byte[], cutset: string): byte[]
    function TrimSpace(s: byte[]): byte[]
    function TrimPrefix(s: byte[], prefix: byte[]): byte[]
    function TrimSuffix(s: byte[], suffix: byte[]): byte[]
    function TrimLeft(s: byte[], cutset: string): byte[]
    function TrimRight(s: byte[], cutset: string): byte[]
    function Compare(a: byte[], b: byte[]): int
}
