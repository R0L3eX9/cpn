package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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
	StressFiles
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
	case "stress-files", "sf":
		return StressFiles
	default:
		return Help
	}
}

type Cli struct {
	File          string
	FileExtension string
	FileName      string
	Command       CliCommand
}

func NewCli(file, command string) (*Cli, error) {
	dotIdx := strings.Index(file, ".")
	if dotIdx == -1 {
		return nil, errors.New("Not a valid file format")
	}

	fileName := file[:dotIdx]
	fileExtension := file[dotIdx:]

	return &Cli{
		File:          file,
		FileName:      fileName,
		FileExtension: fileExtension,
		Command:       toCliCommand(command),
	}, nil
}

func (c *Cli) Execute() {
	switch c.Command {
	case Help:
		c.help()
	case Run:
		err := c.run()
		if err != nil {
			fmt.Println(err)
		}
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
		err := c.compileRun()
		if err != nil {
			fmt.Println(err)
		}
	case DebugRun:
		err := c.debugRun()
		if err != nil {
			fmt.Println(err)
		}
	case Parse:
		err := c.parse()
		if err != nil {
			fmt.Println(err)
		}
	case StressTest:
        err := c.stressTest()
        if err != nil {
            fmt.Println(err)
        }
	case StressFiles:
		err := c.stressFiles()
		if err != nil {
			fmt.Println(err)
		}
	case Template:
		err := c.template()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// TODO: Complete with the rest of the commands
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
	if debug {
		flags := DEBUG_FLAGS + " " + c.FileName + " " + c.File
		arg := strings.Split(flags, " ")
		cmd = exec.Command("g++", arg...)
	} else {
		flags := COMPILATION_FLAGS + " " + c.FileName + " " + c.File
		arg := strings.Split(flags, " ")
		cmd = exec.Command("g++", arg...)
	}
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	fmt.Println("Finished compiling", c.File)
	return nil
}

func (c *Cli) compileRun() error {
	err := c.compile(false)
	if err != nil {
		return err
	}
	err = c.run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cli) debugRun() error {
	err := c.compile(true)
	if err != nil {
		return err
	}
	err = c.run()
	if err != nil {
		return err
	}
	return nil
}

const (
	SERVER_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
)

func (c *Cli) parse() error {
	parser := NewParser(SERVER_TYPE, SERVER_HOST, SERVER_PORT)
	err := parser.Listen()
	if err != nil {
		return err
	}
	log.Println("Ended parsing")
	return nil
}

func runStress(file string, out string) error {
    input, err := os.ReadFile("int")
    if err != nil {
        return err
    }
    cmd := exec.Command("./" + file)
    var stdout bytes.Buffer
    cmd.Stdin = bytes.NewBuffer(input)
    cmd.Stdout = &stdout
    err = cmd.Run();

    if err != nil {
        return err
    }

    err = os.WriteFile(out, stdout.Bytes(), 0666)
    if err != nil {
        return err
    }

    return nil
}

func (c *Cli) stressTest() error {
	fmt.Printf("Running in stress test mode\n")
    testCase := 0
    for {
        fmt.Printf("Running test %d\n", testCase)
        cmd := exec.Command("./gen")

        var input bytes.Buffer
        cmd.Stdout = &input
        err := cmd.Run();

        if err != nil {
            return err
        }

        err = os.WriteFile("int", input.Bytes(), 0666)
        if err != nil {
            return err
        }

        err = runStress("brute", "out1")
        if err != nil {
            return err
        }

        err = runStress(c.FileName, "out2")
        if err != nil {
            return err
        }

        out1, err := os.ReadFile("out1")
        if err != nil {
            return err
        }

        out2, err := os.ReadFile("out2")
        if err != nil {
            return err
        }
        if string(out1) != string(out2) {
            fmt.Println("Test case failed:")
            fmt.Println(input.String())
            fmt.Printf("Expected:\n %s\n", out1)
            fmt.Printf("Got:\n %s\n", out2)
            os.Exit(1)
        }
        testCase++
    }
}

// TODO: read from config file
const MAIN_TEMPLATE_PATH = "/home/razvan/Templates/template.cpp"
const GENERATOR_TEMPLATE_PATH = "/home/razvan/Templates/gen.cpp"

func (c *Cli) template() error {
	err := FromTemplate(MAIN_TEMPLATE_PATH, c.File)
	if err != nil {
		return err
	}
	return nil
}

func FromTemplate(templatePath string, file string) error {
    fmt.Printf("Generating %s\n", file)
	cmd := exec.Command("cp", templatePath, file)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cli) stressFiles() error {
	err := FromTemplate(MAIN_TEMPLATE_PATH, "brute.cpp")
	if err != nil {
		return err
	}
	err = FromTemplate(GENERATOR_TEMPLATE_PATH, "gen.cpp")
	if err != nil {
		return err
	}
	println("Stress test files created!")
	return nil
}
