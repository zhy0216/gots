// Package main implements the goTS CLI.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/codegen"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

const version = "0.2.0"

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

	case "build":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots build <file.gts> [-o output]")
			os.Exit(1)
		}
		input := os.Args[2]
		output := ""
		emitGo := false

		for i := 3; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-o":
				if i+1 < len(os.Args) {
					output = os.Args[i+1]
					i++
				}
			case "--emit-go":
				emitGo = true
			}
		}

		if output == "" {
			output = strings.TrimSuffix(input, filepath.Ext(input))
			if emitGo {
				output += ".go"
			}
		}

		buildFile(input, output, emitGo)

	case "emit-go":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gots emit-go <file.gts> [output.go]")
			os.Exit(1)
		}
		input := os.Args[2]
		output := ""
		if len(os.Args) > 3 {
			output = os.Args[3]
		} else {
			output = strings.TrimSuffix(input, filepath.Ext(input)) + ".go"
		}
		emitGoFile(input, output)

	case "repl":
		runRepl()

	case "version":
		fmt.Printf("gots version %s (Go transpiler)\n", version)

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
	fmt.Println("goTS - A TypeScript-like language that compiles to Go")
	fmt.Println()
	fmt.Println("Usage: gots <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.gts>           Compile and run a source file")
	fmt.Println("  build <file.gts>         Compile to native binary")
	fmt.Println("    -o <output>            Specify output file name")
	fmt.Println("    --emit-go              Output Go source instead of binary")
	fmt.Println("  emit-go <file.gts>       Generate Go source code")
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

	goCode, err := compileToGo(string(source), path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Create a temporary directory for the Go code
	tmpDir, err := os.MkdirTemp("", "gots-run-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Write the Go code to a file
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, goCode, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing Go file: %v\n", err)
		os.Exit(1)
	}

	// If code uses external dependencies (like modernc.org/sqlite), init a go module
	usesModules := strings.Contains(string(goCode), "modernc.org/sqlite")
	if usesModules {
		initCmd := exec.Command("go", "mod", "init", "gts_temp")
		initCmd.Dir = tmpDir
		if out, err := initCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing module: %v\n%s\n", err, out)
			os.Exit(1)
		}
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = tmpDir
		if out, err := tidyCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching dependencies: %v\n%s\n", err, out)
			os.Exit(1)
		}
	}

	// Run the Go code
	var cmd *exec.Cmd
	if usesModules {
		cmd = exec.Command("go", "run", ".")
		cmd.Dir = tmpDir
	} else {
		cmd = exec.Command("go", "run", goFile)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running: %v\n", err)
		os.Exit(1)
	}
}

// buildFile compiles a .gts file to a native binary or Go source
func buildFile(input, output string, emitGo bool) {
	source, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	goCode, err := compileToGo(string(source), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if emitGo {
		// Just write the Go source
		if err := os.WriteFile(output, goCode, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing Go file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", output)
		return
	}

	// Create a temporary directory for the Go code
	tmpDir, err := os.MkdirTemp("", "gots-build-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Write the Go code to a file
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, goCode, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing Go file: %v\n", err)
		os.Exit(1)
	}

	// If code uses external dependencies (like modernc.org/sqlite), init a go module
	usesModules := strings.Contains(string(goCode), "modernc.org/sqlite")
	if usesModules {
		initCmd := exec.Command("go", "mod", "init", "gts_temp")
		initCmd.Dir = tmpDir
		if out, err := initCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing module: %v\n%s\n", err, out)
			os.Exit(1)
		}
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = tmpDir
		if out, err := tidyCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching dependencies: %v\n%s\n", err, out)
			os.Exit(1)
		}
	}

	// Build the binary
	absOutput, _ := filepath.Abs(output)
	var cmd *exec.Cmd
	if usesModules {
		cmd = exec.Command("go", "build", "-o", absOutput, ".")
		cmd.Dir = tmpDir
	} else {
		cmd = exec.Command("go", "build", "-o", absOutput, goFile)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Build error: %v\n%s\n", err, out)
		os.Exit(1)
	}

	fmt.Printf("Built %s -> %s\n", input, output)
}

// emitGoFile generates Go source code from a .gts file
func emitGoFile(input, output string) {
	source, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	goCode, err := compileToGo(string(source), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(output, goCode, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing Go file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s -> %s\n", input, output)
}

// compileToGo compiles source code to Go.
func compileToGo(source, filename string) ([]byte, error) {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		return nil, fmt.Errorf("Parser errors:\n  %s", strings.Join(errs, "\n  "))
	}

	// Build typed AST
	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		var errMsgs []string
		for _, e := range builder.Errors() {
			errMsgs = append(errMsgs, e.String())
		}
		return nil, fmt.Errorf("Type errors:\n  %s", strings.Join(errMsgs, "\n  "))
	}

	// Generate Go code
	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		return nil, fmt.Errorf("Codegen error: %v", err)
	}

	return goCode, nil
}

// runRepl starts the REPL
func runRepl() {
	fmt.Println("goTS REPL v" + version + " (Go transpiler)")
	fmt.Println("Type 'exit' or press Ctrl+D to quit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	var history []string

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

		// Add previous declarations to make REPL stateful
		fullSource := strings.Join(history, "\n") + "\n" + line

		// Try to compile and run
		goCode, err := compileToGo(fullSource, "<repl>")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		// Create a temporary directory for the Go code
		tmpDir, err := os.MkdirTemp("", "gots-repl-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		goFile := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(goFile, goCode, 0644); err != nil {
			os.RemoveAll(tmpDir)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		cmd := exec.Command("go", "run", goFile)
		output, err := cmd.CombinedOutput()
		os.RemoveAll(tmpDir)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n%s", err, output)
			continue
		}

		if len(output) > 0 {
			fmt.Print(string(output))
		}

		// If successful, add to history (for declarations)
		if strings.HasPrefix(strings.TrimSpace(line), "let ") ||
			strings.HasPrefix(strings.TrimSpace(line), "const ") ||
			strings.HasPrefix(strings.TrimSpace(line), "function ") ||
			strings.HasPrefix(strings.TrimSpace(line), "class ") {
			history = append(history, line)
		}
	}
}
