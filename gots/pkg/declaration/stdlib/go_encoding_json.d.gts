// Declaration file for Go's encoding/json package
declare module "go:encoding/json" {
    // Marshal returns the JSON encoding of v
    function Marshal(v: any): byte[]

    // Unmarshal parses JSON-encoded data
    function Unmarshal(data: byte[], v: any): Error | null

    // MarshalIndent is like Marshal but with indentation
    function MarshalIndent(v: any, prefix: string, indent: string): byte[]

    // Valid reports whether data is a valid JSON encoding
    function Valid(data: byte[]): boolean

    // Compact appends compacted JSON to dst
    function Compact(dst: any, src: byte[]): Error | null

    // Indent appends indented JSON to dst
    function Indent(dst: any, src: byte[], prefix: string, indent: string): Error | null

    // HTMLEscape appends HTML-escaped JSON to dst
    function HTMLEscape(dst: any, src: byte[]): void

    // Encoder writes JSON to an output stream
    interface Encoder {
        Encode(v: any): Error | null
        SetIndent(prefix: string, indent: string): void
        SetEscapeHTML(on: boolean): void
    }

    // Decoder reads JSON from an input stream
    interface Decoder {
        Decode(v: any): Error | null
        Buffered(): any
        More(): boolean
        Token(): any
        UseNumber(): void
        DisallowUnknownFields(): void
    }

    // NewEncoder returns a new encoder that writes to w
    function NewEncoder(w: any): Encoder

    // NewDecoder returns a new decoder that reads from r
    function NewDecoder(r: any): Decoder

    // RawMessage is a raw encoded JSON value
    type RawMessage = byte[]
}
