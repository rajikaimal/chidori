package compiler

import (
	"go/token"
	"log"
	"os"

	"github.com/goruby/goruby/generator"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

// Compiler defines the methods of an compiler
type Compiler interface {
	Compile(filename string, input interface{}) (object.RubyObject, error)
}

// New returns an Compiler ready to use and with the environment set to
// object.NewMainEnvironment()
func New() Compiler {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get working directory: %s\n", err)
	}
	env := object.NewMainEnvironment()
	loadPath, _ := env.Get("$:")
	loadPathArr := loadPath.(*object.Array)
	loadPathArr.Elements = append(loadPathArr.Elements, &object.String{Value: cwd})
	env.SetGlobal("$:", loadPathArr)
	return &compiler{environment: env}
}

type compiler struct {
	environment object.Environment
}

func (c *compiler) Compile(filename string, input interface{}) (object.RubyObject, error) {
	node, err := parser.ParseFile(token.NewFileSet(), filename, input, 0)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return generator.Eval(node, c.environment)
}
