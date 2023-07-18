package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CliCommand int

const (
    Help CliCommand = iota
    Run
    Compile
    Debug
    CompileRun
    DebugRun
    Parse
    StressTest
    Template
)

func toCliCommand(command string) CliCommand {
    switch command {
        case "help", "h":
            return Help
        case "run", "r":
            return Run
        case "compile", "c":
            return Compile
        case "debug", "d":
            return Debug
        case "compile-run", "cr":
            return CompileRun
        case "debug-run", "dr":
            return DebugRun
        case "parse":
            return Parse
        case "stress-test", "st":
            return StressTest
        case "template", "temp":
            return Template
        default:
            return Help
    }
}

type Cli struct {
    File            string
    FileExtension   string
    FileName        string
    Command         CliCommand
}


func NewCli(file, command string) (*Cli, error) {
    dotIdx := strings.Index(file, ".")
    if dotIdx == -1 {
        return nil, errors.New("Not a valid file format")
    }

    fileName := file[:dotIdx]
    fileExtension := file[dotIdx:]

    return &Cli{
        File: file,
        FileName: fileName,
        FileExtension: fileExtension,
        Command: toCliCommand(command),
    }, nil
}

func (c *Cli) Execute() {
    switch c.Command {
        case Help:
            c.help()
        case Run:
            c.run()
        case Compile:
            err := c.compile(false)
            if err != nil {
                fmt.Println(err)
            }
        case Debug:
            err := c.compile(true)
            if err != nil {
                fmt.Println(err)
            }
        case CompileRun:
            c.compileRun()
        case DebugRun:
            c.debugRun()
        case Parse:
            c.parse()
        case StressTest:
            c.stressTest()
        case Template:
            c.template()
    }
}

func (c *Cli) help() {
    fmt.Println("Usage: cpn [command] [file]")
    fmt.Println()
    fmt.Println("Commands:")
    fmt.Println("  help, h: Show this help")
    fmt.Println("  run, r: Run the program")
    fmt.Println("  compile, c: Compile the program")
    fmt.Println("  debug, d: Run the program in debug mode")
    fmt.Println("  compile-run, cr: Compile and run the program")
    fmt.Println("  debug-run, dr: Compile and run the program in debug mode")
    fmt.Println("  parse: Parse the program")
    fmt.Println("  stress-test: Run the program in stress test mode")
    fmt.Println("  template: Generate a template")
    fmt.Println()
}

func (c *Cli) run() error {
    fmt.Println("Running", c.File)
    cmd := exec.Command("./" + c.FileName)
    fmt.Println("Executing:", cmd)

    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    err := cmd.Run()
    if err != nil {
        return errors.New("Couldn't run the program")
    }
    return nil
}

const (
    COMPILATION_FLAGS = "-fsanitize=address -std=c++17 -Wno-unused-result" +
                        " -Wall -Wextra -Wshadow -DONPC -O2 -o"

    DEBUG_FLAGS = "-pedantic -DONPC -Wformat=2 -Wfloat-equal -Wconversion" +
                  " -Wlogical-op -Wshift-overflow=2 -Wduplicated-cond" +
                  " -Wcast-qual -Wcast-align -Wno-sign-conversion" +
                  " -fsanitize=undefined -fsanitize=float-divide-by-zero" +
                  " -fsanitize=float-cast-overflow -fno-sanitize-recover=all" +
                  " -fstack-protector-all -D_FORTIFY_SOURCE=2 -D_GLIBCXX_DEBUG" +
                  " -D_GLIBCXX_DEBUG_PEDANTIC -o"
)


func (c *Cli) compile(debug bool) error {
    fmt.Printf("Compiling %s\n", c.File)
    var cmd *exec.Cmd
    if debug == true {
        cmd = exec.Command("g++", "-o", c.FileName, DEBUG_FLAGS, c.File, )
    } else {
        flags := COMPILATION_FLAGS + " " + c.FileName + " " + c.File
        arg := strings.Split(flags, " ")
        cmd = exec.Command("g++", arg...)
    }
    err := cmd.Run()
    if err != nil {
        return errors.New("CompilationError")
    }
    fmt.Println("Finished compiling", c.File)
    return nil
}

func (c *Cli) compileRun() {
    fmt.Printf("Compiling and running %s\n", c.File)
}

func (c *Cli) debugRun() {
    fmt.Printf("Debugging and running %s\n", c.File)
}

func (c *Cli) parse() {
    fmt.Printf("Parsing %s\n", c.File)
}

func (c *Cli) stressTest() {
    fmt.Printf("Running in stress test mode\n")
}

func (c *Cli) template() {
    fmt.Printf("Generating a template\n")
}
