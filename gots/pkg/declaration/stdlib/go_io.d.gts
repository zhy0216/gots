// Declaration file for Go's io package
declare module "go:io" {
    // EOF is the error returned by Read when no more input is available
    const EOF: Error

    // Reader interface
    interface Reader {
        Read(p: byte[]): int
    }

    // Writer interface
    interface Writer {
        Write(p: byte[]): int
    }

    // Closer interface
    interface Closer {
        Close(): Error | null
    }

    // ReadCloser combines Reader and Closer
    interface ReadCloser {
        Read(p: byte[]): int
        Close(): Error | null
    }

    // WriteCloser combines Writer and Closer
    interface WriteCloser {
        Write(p: byte[]): int
        Close(): Error | null
    }

    // ReadWriter combines Reader and Writer
    interface ReadWriter {
        Read(p: byte[]): int
        Write(p: byte[]): int
    }

    // ReadWriteCloser combines Reader, Writer, and Closer
    interface ReadWriteCloser {
        Read(p: byte[]): int
        Write(p: byte[]): int
        Close(): Error | null
    }

    // Functions
    function Copy(dst: Writer, src: Reader): int
    function CopyN(dst: Writer, src: Reader, n: int): int
    function CopyBuffer(dst: Writer, src: Reader, buf: byte[]): int
    function ReadAll(r: Reader): byte[]
    function ReadFull(r: Reader, buf: byte[]): int
    function WriteString(w: Writer, s: string): int
    function Pipe(): (Reader, Writer)
    function LimitReader(r: Reader, n: int): Reader
    function MultiReader(...readers: Reader[]): Reader
    function MultiWriter(...writers: Writer[]): Writer
    function TeeReader(r: Reader, w: Writer): Reader
    function NopCloser(r: Reader): ReadCloser

    // Discard is a Writer on which all Write calls succeed without doing anything
    const Discard: Writer
}
