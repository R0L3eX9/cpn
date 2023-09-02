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
    CLI_VERSION = 1.10

	Help CliCommand = iota
	Run
	Compile
	Debug
	CompileRun
	DebugRun
	Parse
	StressTest
	StressFiles
    TestCases
	Template
	Version
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
	case "stress-files", "sf":
		return StressFiles
    case "test-cases", "tc":
        return TestCases
	case "template", "temp":
		return Template
	case "version", "v":
		return Version
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
		help()
	case Run:
		err := run(c.File, c.FileName)
		if err != nil {
			fmt.Println(err)
		}
	case Compile:
		err := compile(c.File, c.FileName, false)
		if err != nil {
			fmt.Println(err)
		}
	case Debug:
		err := compile(c.File, c.FileName, true)
		if err != nil {
			fmt.Println(err)
		}
	case CompileRun:
		err := compileRun(c.File, c.FileName)
		if err != nil {
			fmt.Println(err)
		}
	case DebugRun:
		err := debugRun(c.File, c.FileName)
		if err != nil {
			fmt.Println(err)
		}
	case Parse:
		err := parse()
		if err != nil {
			fmt.Println(err)
		}
	case StressTest:
		err := stressTest(c.File, c.FileName)
		if err != nil {
			fmt.Println(err)
		}
	case StressFiles:
		err := stressFiles()
		if err != nil {
			fmt.Println(err)
		}
    case TestCases:
        err := testCases(c.File, c.FileName)
        if err != nil {
            fmt.Println(err)
        }
	case Template:
		err := template(c.File)
		if err != nil {
			fmt.Println(err)
		}
    case Version:
        version()
	}
}

func help() {
	fmt.Println("Usage: cpn [command] [file]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help, h: Show this help")
	fmt.Println("  run, r: Run the program")
	fmt.Println("  compile, c: Compile the program")
	fmt.Println("  debug, d: Compile the program in debug mode")
	fmt.Println("  compile-run, cr: Compile and run the program")
	fmt.Println("  debug-run, dr: Compile and run the program in debug mode")
	fmt.Println("  parse: Parse the problem using Competitive Companion")
	fmt.Println("  stress-test, st: Run the program in stress test mode")
    fmt.Println("  stress-files, sf: Generate files needed for stress testing")
	fmt.Println("  test-cases, tc: Test your program against parsed test-cases")
	fmt.Println("  template: Generate a file from personal template")
	fmt.Println()
}

func run(file, fileName string) error {
	fmt.Println("Running", file)
	cmd := exec.Command("./" + fileName)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
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

func compile(file, fileName string, debug bool) error {
	fmt.Printf("Compiling %s\n", file)
	var cmd *exec.Cmd
	if debug {
		flags := DEBUG_FLAGS + " " + fileName + " " + file
		arg := strings.Split(flags, " ")
		cmd = exec.Command("g++", arg...)
	} else {
		flags := COMPILATION_FLAGS + " " + fileName + " " + file
		arg := strings.Split(flags, " ")
		cmd = exec.Command("g++", arg...)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	fmt.Println("Finished compiling", file)
	return nil
}

func compileRun(file, fileName string) error {
	err := compile(file, fileName, false)
	if err != nil {
		return err
	}
	err = run(file, fileName)
	if err != nil {
		return err
	}
	return nil
}

func debugRun(file, fileName string) error {
	err := compile(file, fileName, true)
	if err != nil {
		return err
	}
	err = run(file, fileName)
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

func parse() error {
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
	err = cmd.Run()

	if err != nil {
		return err
	}

	err = os.WriteFile(out, stdout.Bytes(), 0666)
	if err != nil {
		return err
	}

	return nil
}

func stressTest(file, fileName string) error {
	fmt.Printf("Running in stress test mode\n")
    err := compile(file, fileName, false)
    if err != nil {
        return err
    }
    err = compile("brute.cpp", "brute", false)
    if err != nil {
        return err
    }
    err = compile("gen.cpp", "gen", false)
    if err != nil {
        return err
    }
	testCase := 0
	for {
		fmt.Printf("Running test %d\n", testCase)
		cmd := exec.Command("./gen")

		var input bytes.Buffer
		cmd.Stdout = &input
		err := cmd.Run()

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

		err = runStress(fileName, "out2")
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
			fmt.Println("Expected:")
            fmt.Println(out1)
			fmt.Println("Got:")
            fmt.Println(out2)
            break
		}
		testCase++
	}
    return nil
}

// TODO: read from config file
const MAIN_TEMPLATE_PATH = "/home/razvan/Templates/template.cpp"
const GENERATOR_TEMPLATE_PATH = "/home/razvan/Templates/gen.cpp"

func template(file string) error {
	err := FromTemplate(MAIN_TEMPLATE_PATH, file)
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

func stressFiles() error {
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

func testCases(file, fileName string) error {
    err := compile(file, fileName, false)
    if err != nil {
        return err
    }
	fmt.Printf("Running Test Cases\n")
    testIdx := '0'
    for {
        input, err := os.ReadFile(fmt.Sprintf("./test-cases/%s-%c.in", fileName, testIdx))
        if os.IsNotExist(err) {
            break
        }
        expectedOutput, err := os.ReadFile(fmt.Sprintf("./test-cases/%s-%c.out", fileName, testIdx))
        if os.IsNotExist(err) {
            break
        }

        cmd := exec.Command(fmt.Sprintf("./%s", fileName))
        stdin := bytes.NewBuffer(input)
        var stdout bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdin = stdin
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr
        err = cmd.Run()
        if err != nil {
            log.Println(stderr)
        }

        userOutput := stdout.String()
        if userOutput != string(expectedOutput) {
            fmt.Printf("Failed on test %c\n", testIdx)
            fmt.Println("Expected:")
            fmt.Println(string(expectedOutput))
            fmt.Println("Got:")
            fmt.Println(userOutput)
            break
        }
        fmt.Printf("Test %c passed!\n", testIdx)
        testIdx++
    }
    fmt.Printf("Finished all %c Test Cases\n", testIdx)
    return nil
}

func version() {
    fmt.Printf("cpn version %.2f\n", CLI_VERSION)
}
