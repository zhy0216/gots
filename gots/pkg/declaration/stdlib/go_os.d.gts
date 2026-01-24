// Declaration file for Go's os package
declare module "go:os" {
    // Exit terminates the program with given status code
    function Exit(code: int): void

    // Environment
    function Getenv(key: string): string
    function Setenv(key: string, value: string): Error | null
    function Unsetenv(key: string): Error | null
    function Environ(): string[]
    function LookupEnv(key: string): (string, boolean)
    function ExpandEnv(s: string): string
    function Clearenv(): void

    // Working directory
    function Getwd(): string
    function Chdir(dir: string): Error | null

    // File operations
    function Mkdir(name: string, perm: int): Error | null
    function MkdirAll(path: string, perm: int): Error | null
    function Remove(name: string): Error | null
    function RemoveAll(path: string): Error | null
    function Rename(oldpath: string, newpath: string): Error | null
    function Chmod(name: string, mode: int): Error | null

    // File reading/writing
    function ReadFile(name: string): byte[]
    function WriteFile(name: string, data: byte[], perm: int): Error | null

    // File info
    function Stat(name: string): FileInfo
    function Lstat(name: string): FileInfo
    function IsNotExist(err: Error): boolean
    function IsExist(err: Error): boolean
    function IsPermission(err: Error): boolean

    // File interface
    interface File {
        Name(): string
        Read(b: byte[]): int
        Write(b: byte[]): int
        WriteString(s: string): int
        Close(): Error | null
        Sync(): Error | null
    }

    // FileInfo interface
    interface FileInfo {
        Name(): string
        Size(): int
        Mode(): int
        IsDir(): boolean
    }

    // File creation/opening
    function Create(name: string): File
    function Open(name: string): File
    function OpenFile(name: string, flag: int, perm: int): File

    // Standard file descriptors
    const Stdin: File
    const Stdout: File
    const Stderr: File

    // File open flags
    const O_RDONLY: int
    const O_WRONLY: int
    const O_RDWR: int
    const O_APPEND: int
    const O_CREATE: int
    const O_EXCL: int
    const O_SYNC: int
    const O_TRUNC: int

    // Args contains command-line arguments
    const Args: string[]

    // Hostname
    function Hostname(): string

    // User info
    function Getuid(): int
    function Geteuid(): int
    function Getgid(): int
    function Getegid(): int
    function Getpid(): int
    function Getppid(): int

    // Temp directory
    function TempDir(): string
    function UserHomeDir(): string
    function UserCacheDir(): string
    function UserConfigDir(): string
}
