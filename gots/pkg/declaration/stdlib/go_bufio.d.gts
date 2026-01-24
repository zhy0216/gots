// Declaration file for Go's bufio package
declare module "go:bufio" {
    // Reader provides buffered reading
    interface Reader {
        Read(p: byte[]): int
        ReadByte(): byte
        UnreadByte(): Error | null
        ReadRune(): (int, int)
        UnreadRune(): Error | null
        ReadSlice(delim: byte): byte[]
        ReadLine(): (byte[], boolean)
        ReadString(delim: byte): string
        ReadBytes(delim: byte): byte[]
        Buffered(): int
        Peek(n: int): byte[]
        Discard(n: int): int
        Reset(r: any): void
        Size(): int
    }

    // Writer provides buffered writing
    interface Writer {
        Write(p: byte[]): int
        WriteByte(c: byte): Error | null
        WriteRune(r: int): int
        WriteString(s: string): int
        Flush(): Error | null
        Buffered(): int
        Available(): int
        AvailableBuffer(): byte[]
        Reset(w: any): void
        Size(): int
    }

    // Scanner provides line/word/byte scanning
    interface Scanner {
        Scan(): boolean
        Text(): string
        Bytes(): byte[]
        Err(): Error | null
        Buffer(buf: byte[], max: int): void
        Split(split: any): void
    }

    // Functions
    function NewReader(rd: any): Reader
    function NewReaderSize(rd: any, size: int): Reader
    function NewWriter(w: any): Writer
    function NewWriterSize(w: any, size: int): Writer
    function NewScanner(r: any): Scanner
    function NewReadWriter(r: Reader, w: Writer): any

    // Split functions for Scanner
    const ScanLines: any
    const ScanWords: any
    const ScanRunes: any
    const ScanBytes: any
}
