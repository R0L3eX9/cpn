package main

import (
    "os"
    "fmt"
)

func main() {
    args := os.Args
    if len(args) < 2 || len(args) > 3 {
        fmt.Println("Invalid number of arguments")
        fmt.Println("Use: cpn help")
        os.Exit(1)
    }

    file := "main.cpp"
    command := "help"

    if len(args) == 3 {
        file = args[2]
    }
    command = args[1]

    cli, err := NewCli(file, command)
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
    cli.Execute()
}
