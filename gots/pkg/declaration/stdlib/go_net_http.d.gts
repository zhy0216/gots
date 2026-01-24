// Declaration file for Go's net/http package
declare module "go:net/http" {
    // HTTP method constants
    const MethodGet: string
    const MethodHead: string
    const MethodPost: string
    const MethodPut: string
    const MethodPatch: string
    const MethodDelete: string
    const MethodConnect: string
    const MethodOptions: string
    const MethodTrace: string

    // Status code constants
    const StatusOK: int
    const StatusCreated: int
    const StatusAccepted: int
    const StatusNoContent: int
    const StatusMovedPermanently: int
    const StatusFound: int
    const StatusNotModified: int
    const StatusBadRequest: int
    const StatusUnauthorized: int
    const StatusForbidden: int
    const StatusNotFound: int
    const StatusMethodNotAllowed: int
    const StatusInternalServerError: int
    const StatusBadGateway: int
    const StatusServiceUnavailable: int

    // Header represents HTTP headers
    interface Header {
        Add(key: string, value: string): void
        Set(key: string, value: string): void
        Get(key: string): string
        Del(key: string): void
        Values(key: string): string[]
    }

    // Request represents an HTTP request
    interface Request {
        Method: string
        URL: URL
        Header: Header
        Body: any
        ContentLength: int
        Host: string
        RemoteAddr: string
        RequestURI: string
        FormValue(key: string): string
        FormFile(key: string): any
        ParseForm(): Error | null
        ParseMultipartForm(maxMemory: int): Error | null
        Cookie(name: string): Cookie
        Cookies(): Cookie[]
        AddCookie(c: Cookie): void
        Context(): any
        WithContext(ctx: any): Request
    }

    // Response represents an HTTP response from a server
    interface Response {
        Status: string
        StatusCode: int
        Header: Header
        Body: any
        ContentLength: int
    }

    // ResponseWriter interface for writing HTTP responses
    interface ResponseWriter {
        Header(): Header
        Write(data: byte[]): int
        WriteHeader(statusCode: int): void
    }

    // Cookie represents an HTTP cookie
    interface Cookie {
        Name: string
        Value: string
        Path: string
        Domain: string
        Expires: any
        MaxAge: int
        Secure: boolean
        HttpOnly: boolean
        SameSite: int
        String(): string
    }

    // URL represents a parsed URL
    interface URL {
        Scheme: string
        Host: string
        Path: string
        RawPath: string
        RawQuery: string
        Fragment: string
        String(): string
        Query(): any
        Hostname(): string
        Port(): string
    }

    // Handler interface
    type HandlerFunc = (w: ResponseWriter, r: Request) => void

    // Client for making HTTP requests
    interface Client {
        Get(url: string): Response
        Post(url: string, contentType: string, body: any): Response
        Do(req: Request): Response
        Head(url: string): Response
        PostForm(url: string, data: any): Response
    }

    // Server represents an HTTP server
    interface Server {
        Addr: string
        Handler: any
        ListenAndServe(): Error | null
        ListenAndServeTLS(certFile: string, keyFile: string): Error | null
        Shutdown(ctx: any): Error | null
        Close(): Error | null
    }

    // Functions
    function Get(url: string): Response
    function Post(url: string, contentType: string, body: any): Response
    function Head(url: string): Response
    function PostForm(url: string, data: any): Response
    function NewRequest(method: string, url: string, body: any): Request
    function ListenAndServe(addr: string, handler: any): Error | null
    function ListenAndServeTLS(addr: string, certFile: string, keyFile: string, handler: any): Error | null
    function Handle(pattern: string, handler: any): void
    function HandleFunc(pattern: string, handler: HandlerFunc): void
    function Redirect(w: ResponseWriter, r: Request, url: string, code: int): void
    function Error(w: ResponseWriter, error: string, code: int): void
    function NotFound(w: ResponseWriter, r: Request): void
    function ServeFile(w: ResponseWriter, r: Request, name: string): void
    function SetCookie(w: ResponseWriter, cookie: Cookie): void
    function StatusText(code: int): string

    // Default client
    const DefaultClient: Client
}
