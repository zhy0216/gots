// Package main implements the GoTS CLI.
package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketlang/gots/pkg/bytecode"
	"github.com/pocketlang/gots/pkg/compiler"
	"github.com/pocketlang/gots/pkg/lexer"
	"github.com/pocketlang/gots/pkg/parser"
	"github.com/pocketlang/gots/pkg/vm"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots run <file.gts>")
			os.Exit(1)
		}
		runFile(os.Args[2])

	case "compile":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots compile <file.gts> [output.gtsb]")
			os.Exit(1)
		}
		input := os.Args[2]
		output := ""
		if len(os.Args) > 3 {
			output = os.Args[3]
		} else {
			output = strings.TrimSuffix(input, filepath.Ext(input)) + ".gtsb"
		}
		compileFile(input, output)

	case "exec":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots exec <file.gtsb>")
			os.Exit(1)
		}
		execFile(os.Args[2])

	case "disasm":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots disasm <file.gtsb|file.gts>")
			os.Exit(1)
		}
		disasmFile(os.Args[2])

	case "repl":
		runRepl()

	case "version":
		fmt.Printf("gots version %s\n", version)

	case "help", "--help", "-h":
		printUsage()

	default:
		if fileExists(cmd) {
			runFile(cmd)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("GoTS - A TypeScript-like language compiler and VM")
	fmt.Println()
	fmt.Println("Usage: gots <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.gts>           Compile and run a source file")
	fmt.Println("  compile <file.gts>       Compile to bytecode (.gtsb)")
	fmt.Println("  exec <file.gtsb>         Execute compiled bytecode")
	fmt.Println("  disasm <file>            Disassemble bytecode")
	fmt.Println("  repl                     Start interactive mode")
	fmt.Println("  version                  Show version")
	fmt.Println("  help                     Show this help")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// runFile compiles and runs a .gts source file
func runFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	chunk, err := compileSource(string(source), path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := executeChunk(chunk); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// compileFile compiles a .gts file to .gtsb bytecode
func compileFile(input, output string) {
	source, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	chunk, err := compileSource(string(source), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Write to file
	f, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var w io.Writer = f
	// Use gzip compression if the output ends with .gz
	if strings.HasSuffix(output, ".gz") {
		gzw := gzip.NewWriter(f)
		defer gzw.Close()
		w = gzw
	}

	if err := bytecode.WriteBinary(w, chunk); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing bytecode: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compiled %s -> %s\n", input, output)
}

// execFile executes a .gtsb bytecode file
func execFile(path string) {
	chunk, err := loadBytecode(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bytecode: %v\n", err)
		os.Exit(1)
	}

	if err := executeChunk(chunk); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// disasmFile disassembles a bytecode or source file
func disasmFile(path string) {
	var chunk *bytecode.Chunk
	var err error

	if strings.HasSuffix(path, ".gtsb") || strings.HasSuffix(path, ".gtsb.gz") {
		chunk, err = loadBytecode(path)
	} else {
		source, readErr := os.ReadFile(path)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", readErr)
			os.Exit(1)
		}
		chunk, err = compileSource(string(source), path)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Print(bytecode.Disassemble(chunk, filepath.Base(path)))
}

// runRepl starts the REPL
func runRepl() {
	fmt.Println("GoTS REPL v" + version)
	fmt.Println("Type 'exit' or press Ctrl+D to quit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	globals := make(map[string]vm.Value)

	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}

		// Handle multi-line input
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "(") {
			for {
				fmt.Print("... ")
				moreLine, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				line += "\n" + moreLine
				moreLine = strings.TrimSpace(moreLine)
				if moreLine == "" || moreLine == "}" || moreLine == ")" {
					break
				}
			}
		}

		// Try to evaluate the input
		result, newGlobals, err := evalRepl(line, globals)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else if result != "" {
			fmt.Println(result)
		}
		globals = newGlobals
	}
}

// compileSource compiles source code to bytecode.
func compileSource(source, filename string) (*bytecode.Chunk, error) {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		return nil, fmt.Errorf("Parser errors:\n  %s", strings.Join(errs, "\n  "))
	}

	c := compiler.New()
	chunk, err := c.Compile(program)
	if err != nil {
		return nil, fmt.Errorf("Compiler error: %v", err)
	}

	return chunk, nil
}

// loadBytecode loads a bytecode file.
func loadBytecode(path string) (*bytecode.Chunk, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gzr.Close()
		r = gzr
	}

	return bytecode.ReadBinary(r)
}

// executeChunk executes a bytecode chunk.
func executeChunk(chunk *bytecode.Chunk) error {
	v := vm.New(chunk)
	return v.Run()
}

// evalRepl evaluates a line in the REPL context.
func evalRepl(source string, globals map[string]vm.Value) (string, map[string]vm.Value, error) {
	chunk, err := compileSource(source, "<repl>")
	if err != nil {
		return "", globals, err
	}

	v := vm.New(chunk)
	for name, val := range globals {
		v.SetGlobal(name, val)
	}

	if err := v.Run(); err != nil {
		return "", globals, err
	}

	newGlobals := v.GetGlobals()
	lastVal := v.LastPopped()
	if !lastVal.IsNull() {
		return lastVal.String(), newGlobals, nil
	}

	return "", newGlobals, nil
}
