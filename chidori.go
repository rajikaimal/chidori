package main

import (
	"fmt"
	"github.com/goruby/goruby/compiler"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	//input := `
	//names = ["FOO", "BAR"]
	//`
	command := os.Args[1]
	rubyFilePath := os.Args[2]

	if command == "compile" {
		if len(rubyFilePath) == 0 {
			panic("Please provide Ruby file path")
		}

		bytes, err := ioutil.ReadFile(rubyFilePath)
		if err != nil {
			panic("Error when reading file")
		}

		input := string(bytes)
		c := compiler.New()

		_, errCompile := c.Compile("", input)

		if errCompile != nil {
			fmt.Println("Compile errors: ", err)
		}

		cmd := exec.Command("bash", "-c", "go fmt compiled.go")
		cmd.Run()
		fmt.Println("Ruby file compiled")
	}
	if command == "run" {
		output, err := exec.Command("go", "run", "compiled.go").Output()

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		fmt.Println(string(output))
	}
}
