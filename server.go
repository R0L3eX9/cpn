package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Parser struct {
	ServerType string
	ServerHost string
	ServerPort string
}

func NewParser(serverType, serverHost, serverPort string) *Parser {
	return &Parser{
		ServerType: serverType,
		ServerHost: serverHost,
		ServerPort: serverPort,
	}
}

func (p *Parser) Listen() error {
	fmt.Println("Listening for problems!")
	server, err := net.Listen(p.ServerType, p.ServerHost+":"+p.ServerPort)

	if err != nil {
		return err
	}

	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			return err
		}

		errs := make(chan error, 1)
		go func() {
			err := processClient(connection)
			errs <- err
			close(errs)
		}()

		if <-errs != nil {
			return <-errs
		}
	}
}

func processClient(client net.Conn) error {
	defer client.Close()
	buffer := make([]byte, 1024)
	_, err := client.Read(buffer)
	if err != nil {
		return err
	}
	problem, err := parseData(buffer)
	if err != nil {
		return err
	}
	err = problem.create()
	if err != nil {
		return err
	}
	return nil
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Problem struct {
	Name        string     `json:"name"`
	Group       string     `json:"group"`
	Url         string     `json:"url"`
	Interactive bool       `json:"interactive"`
	MemoryLimit int        `json:"memoryLimit"`
	TimeLimit   int        `json:"timeLimit"`
	Tests       []TestCase `json:"tests"`
}

func parseData(b []byte) (*Problem, error) {
	data := string(b)

	problemStart := strings.Index(data, "{")
	if problemStart == -1 {
		return &Problem{}, errors.New("Invalid data format")
	}

	problemEnd := strings.Index(data, "]")
	if problemEnd == -1 {
		return &Problem{}, errors.New("Invalid data format")
	}

	processedData := string(data[problemStart : problemEnd+1])
	processedData += "}"
	fmt.Println("Marshaled data:", processedData)

	var problem Problem
	err := json.Unmarshal([]byte(processedData), &problem)
	if err != nil {
		return &Problem{}, err
	}
	return &problem, nil
}

func (p *Problem) create() error {
	problemName := p.Name

	codeforcesSpecifier := strings.Index(p.Name, ".")
	atcoderSpecifier := strings.Index(p.Name, "-")

	if codeforcesSpecifier != -1 {
		problemName = p.Name[:codeforcesSpecifier]
	} else if atcoderSpecifier != -1 {
		problemName = p.Name[:atcoderSpecifier-1]
	}
	err := FromTemplate(MAIN_TEMPLATE_PATH, problemName)
	if err != nil {
		return err
	}
	err = p.testCases()
	if err != nil {
		return err
	}
	return nil
}

func (p *Problem) testCases() error {
	fmt.Println("Generating test-cases")
	err := os.Mkdir("test-cases", 0777)
	if !os.IsExist(err) {
		return err
	}

	for idx, test := range p.Tests {
		fmt.Println("Creating test:", idx)
		testName := p.Name + "-" + fmt.Sprint(idx)
		err = test.create(testName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TestCase) create(testName string) error {
	err := os.WriteFile("test-cases/"+testName+".in", []byte(t.Input), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile("test-cases/"+testName+".out", []byte(t.Output), 0644)
	if err != nil {
		return err
	}
	return nil
}
