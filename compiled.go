package main
 
import "github.com/goruby/goruby/object"

func main() { 
	env := object.NewMainEnvironment()

	right := object.NewInteger(1)
	env.Set("x", right)
}
package main
 
import "github.com/goruby/goruby/object"

func main() { 
	env := object.NewMainEnvironment()

	right := object.NewInteger(1)
	env.Set("x", right)
}
