package generator

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/object"
	"github.com/pkg/errors"
)

type callContext struct {
	object.CallContext
}

func (c *callContext) Generate(node ast.Node, env object.Environment) (object.RubyObject, error) {
	return Generate(node, env)
}

type rubyObjects []object.RubyObject

func (r rubyObjects) Inspect() string {
	toS := make([]string, len(r))
	for i, e := range r {
		toS[i] = e.Inspect()
	}
	return strings.Join(toS, ", ")
}
func (r rubyObjects) Type() object.Type       { return "" }
func (r rubyObjects) Class() object.RubyClass { return nil }

func expandToArrayIfNeeded(obj object.RubyObject) object.RubyObject {
	fmt.Println("ExpandingToArray:", obj)
	arr, ok := obj.(rubyObjects)
	if !ok {
		return obj
	}

	return object.NewArray(arr...)
}

var file, err = os.Create("compiled.go")
var variableCount = 0
var existingArrays = make(map[string]string)

// Eval evaluates the given node and traverses recursive over its children
func Generate(node ast.Node, env object.Environment) (object.RubyObject, error) {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		fmt.Println("Creating file at ast.Program")
		createMainFunc()
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		fmt.Println("Expression statement", node, node.Expression)
		return Generate(node.Expression, env)
	case *ast.ReturnStatement:
		fmt.Println("Return statement", node)
		val, err := Generate(node.ReturnValue, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval of return statement")
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.BlockStatement:
		fmt.Println("/////////// BlockStatment:", node.Statements, node, node.End(), node.Pos(), node.EndToken)
		return evalBlockStatement(node, env)
	case (*ast.IntegerLiteral):
		// Literals
		fmt.Println("----- IntegerLiteral:", node.Value)
		return outputInteger(node.Value), nil
		//return object.NewInteger(node.Value), nil
	case (*ast.Boolean):
		return nativeBoolToBooleanObject(node.Value), nil
	case (*ast.Nil):
		return object.NIL, nil
	case (*ast.Self):
		self, _ := env.Get("self")
		return self, nil
	case (*ast.Keyword__FILE__):
		return &object.String{Value: node.Filename}, nil
	case (*ast.InstanceVariable):
		fmt.Println("InstanceVariable", node)
		self, _ := env.Get("self")
		selfObj := self.(*object.Self)
		selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(
					fmt.Errorf("instance variable not allowed for %s", selfObj.Name),
				),
			)
		}

		val, ok := selfAsEnv.Get(node.String())
		fmt.Println("INstanceVariable value: ", val)
		if !ok {
			return object.NIL, nil
		}
		return val, nil
	case *ast.Identifier:
		fmt.Println("Identifier:", node)
		outputGetIdentifier(node)
		return evalIdentifier(node, env)
	case *ast.Global:
		fmt.Println("Global", node)
		val, ok := env.Get(node.Value)
		if !ok {
			return object.NIL, nil
		}
		return val, nil
	case *ast.StringLiteral:
		fmt.Println("StringLiteral:", node)
		return outputString(node.Value), nil
		//return &object.String{Value: node.Value}, nil
	case *ast.SymbolLiteral:
		fmt.Println("SymbolLiteral:", node)
		switch value := node.Value.(type) {
		case *ast.Identifier:
			fmt.Println("SymBolLiteralIdentifier:", node)
			return &object.Symbol{Value: value.Value}, nil
		case *ast.StringLiteral:
			fmt.Println("SymBolLiteralStringLiteral:", node)
			str, err := Generate(value, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval symbol literal string")
			}
			if str, ok := str.(*object.String); ok {
				return &object.Symbol{Value: str.Value}, nil
			}
			panic(errors.WithStack(
				fmt.Errorf("error while parsing SymbolLiteral: expected *object.String, got %T", str),
			))
		default:
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("malformed symbol AST: %T", value)),
			)
		}
	case *ast.FunctionLiteral:
		fmt.Println("FunctionLiteral:", node, node.Receiver)
		context, _ := env.Get("self")
		_, inClassOrModule := context.(*object.Self).RubyObject.(object.Environment)
		if node.Receiver != nil {
			fmt.Println("FunctionLiteral Reciever: ", node.Receiver)
			rec, err := Generate(node.Receiver, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function receiver")
			}
			context = rec
			_, recIsEnv := context.(object.Environment)
			if recIsEnv || inClassOrModule {
				inClassOrModule = true
				context = context.Class().(object.RubyClassObject)
			}
		}
		params := make([]*object.FunctionParameter, len(node.Parameters))
		for i, param := range node.Parameters {
			def, err := Generate(param.Default, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function literal param")
			}
			params[i] = &object.FunctionParameter{Name: param.Name.Value, Default: def}
			fmt.Println("function params: ", params[i])
		}
		body := node.Body
		function := &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}
		extended := object.AddMethod(context, node.Name.Value, function)
		if node.Receiver != nil && !inClassOrModule {
			envInfo, _ := object.EnvStat(env, context)
			envInfo.Env().Set(node.Receiver.Value, extended)
		}
		return &object.Symbol{Value: node.Name.Value}, nil
	case *ast.BlockExpression:
		fmt.Println("BlockExpression:", node)
		params := node.Parameters
		body := node.Body
		block := &object.Proc{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
		return block, nil
	case *ast.ArrayLiteral:
		fmt.Println("Arrayliteral:", node, node.Elements, node.String())
		//elements, err := evalExpressions(node.Elements, env)
		variableName, err := outputArray(node.Elements, env, node.String())
		if err != nil {
			return nil, errors.WithMessage(err, "Output array is nil")
		}
		return &object.String{Value: variableName}, nil
	case *ast.HashLiteral:
		fmt.Println("HashLiteral:")
		var hash object.Hash
		for k, v := range node.Map {
			key, err := Generate(k, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash key")
			}
			value, err := Generate(v, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash value")
			}
			hash.Set(key, value)
		}
		return &hash, nil
	case ast.ExpressionList:
		fmt.Println("ExpressionList:", node)
		var objects []object.RubyObject
		for _, e := range node {
			obj, err := Generate(e, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval expression list")
			}
			objects = append(objects, obj)
		}
		return rubyObjects(objects), nil

	// Expressions
	case *ast.Assignment:
		fmt.Println("Assignment here:", node)
		right, err := Generate(node.Right, env)
		fmt.Println("After generating assignment:", right)
		if err != nil {
			return nil, errors.WithMessage(err, "eval right hand Assignment side")
		}

		fmt.Println("Assignment Right:", node.Right, right)
		fmt.Println("Assignment Left:", node.Left)

		switch left := node.Left.(type) {
		case *ast.IndexExpression:
			fmt.Println("IndexExpression:")
			indexLeft, err := Generate(left.Left, env)
			fmt.Println("IndexExpressionLeft:", indexLeft)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
			}
			index, err := Generate(left.Index, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
			}
			return evalIndexExpressionAssignment(indexLeft, index, expandToArrayIfNeeded(right))
		case *ast.InstanceVariable:
			fmt.Println("InstanceVariable:")
			self, _ := env.Get("self")
			selfObj := self.(*object.Self)
			selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
			if !ok {
				return nil, errors.Wrap(
					object.NewSyntaxError(fmt.Errorf("instance variable not allowed for %s", selfObj.Name)),
					"eval left hand Assignment side",
				)
			}

			right = expandToArrayIfNeeded(right)
			fmt.Println("InstanceVariable:", left, right)
			selfAsEnv.Set(left.String(), right)
			return right, nil
		case *ast.Identifier:
			fmt.Println("Here setting Identifier:", left, right)
			right = expandToArrayIfNeeded(right)
			fmt.Println("LEFT RIGHT", left, right)
			outputValue(left.Value, right)
			//appendToFile("\tenv.Set(\"" + left.Value + "\", right)")
			env.Set(left.Value, right)
			fmt.Println("Identifier:", left, right)
			return right, nil
		case *ast.Global:
			fmt.Println("Global:")
			right = expandToArrayIfNeeded(right)
			env.SetGlobal(left.Value, right)
			return right, nil
		case ast.ExpressionList:
			fmt.Println("ExpressionList")
			values := []object.RubyObject{right}
			if list, ok := right.(rubyObjects); ok {
				values = list
			}
			if len(left) > len(values) {
				// enlarge slice
				for len(values) <= len(left) {
					values = append(values, object.NIL)
				}
			}
			for i, exp := range left {
				if _, ok := exp.(*ast.InstanceVariable); ok {
					self, _ := env.Get("self")
					selfObj := self.(*object.Self)
					selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
					if !ok {
						return nil, errors.Wrap(
							object.NewSyntaxError(fmt.Errorf("instance variable not allowed for %s", selfObj.Name)),
							"eval left hand Assignment side",
						)
					}

					selfAsEnv.Set(exp.String(), values[i])
					continue
				}
				if indexExp, ok := exp.(*ast.IndexExpression); ok {
					indexLeft, err := Generate(indexExp.Left, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
					}
					index, err := Generate(indexExp.Index, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
					}
					evalIndexExpressionAssignment(indexLeft, index, values[i])
					continue
				}
				env.Set(exp.String(), values[i])
			}
			return expandToArrayIfNeeded(right), nil
		default:
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("Assignment not supported to %T", node.Left)),
			)
		}
	case *ast.ModuleExpression:
		fmt.Println("ModuleExpression:")
		module, ok := env.Get(node.Name.Value)
		if !ok {
			module = object.NewModule(node.Name.Value, env)
		}
		moduleEnv := module.(object.Environment)
		moduleEnv.Set("self", &object.Self{RubyObject: module, Name: node.Name.Value})
		bodyReturn, err := Generate(node.Body, moduleEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval Module body")
		}
		selfObject, _ := moduleEnv.Get("self")
		self := selfObject.(*object.Self)
		env.Set(node.Name.Value, self.RubyObject)
		return bodyReturn, nil
	case *ast.ClassExpression:
		fmt.Println("ClassExpression:", node.Name.Value)
		outputClassExpression(node, env)
		superClassName := "Object"
		if node.SuperClass != nil {
			superClassName = node.SuperClass.Value
		}
		superClass, ok := env.Get(superClassName)
		if !ok {
			return nil, errors.Wrap(
				object.NewUninitializedConstantNameError(superClassName),
				"eval class superclass",
			)
		}
		class, ok := env.Get(node.Name.Value)
		if !ok {
			class = object.NewClass(node.Name.Value, superClass.(object.RubyClassObject), env)
		}
		classEnv := class.(object.Environment)
		classEnv.Set("self", &object.Self{RubyObject: class, Name: node.Name.Value})
		bodyReturn, err := Generate(node.Body, classEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval class body")
		}
		selfObject, _ := classEnv.Get("self")
		self := selfObject.(*object.Self)
		env.Set(node.Name.Value, self.RubyObject)
		return bodyReturn, nil
	case *ast.ContextCallExpression:
		fmt.Println("ContextCallExpression:", node)
		//TODO: Support methods

		outputContextCalls(node)
		context, err := Generate(node.Context, env)
		fmt.Println("Context: ", context)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call receiver")
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call arguments")
		}
		if node.Block != nil {
			fmt.Println("Block")
			block, err := Generate(node.Block, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval method call block")
			}
			args = append(args, block)
		}
		callContext := &callContext{object.NewCallContext(env, context)}
		fmt.Println("CallContext:", context)
		return object.Send(callContext, node.Function.Value, args...)
	case *ast.YieldExpression:
		fmt.Println("YieldExpression:", node)
		selfObject, _ := env.Get("self")
		self := selfObject.(*object.Self)
		if self.Block == nil {
			return nil, errors.WithStack(object.NewNoBlockGivenLocalJumpError())
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval yield arguments")
		}
		callContext := &callContext{object.NewCallContext(env, self)}
		return self.Block.Call(callContext, args...)
	case *ast.IndexExpression:
		fmt.Println("IndexExpression:", node)
		left, err := Generate(node.Left, env)
		fmt.Println("IndexExpression:", left)
		if err != nil {
			return nil, errors.WithMessage(err, "eval IndexExpression left side")
		}
		index, err := Generate(node.Index, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval IndexExpression index")
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		fmt.Println("PrefixExpression:", node)
		right, err := Generate(node.Right, env)
		fmt.Println("PrefixExpression right", right)
		if err != nil {
			return nil, errors.WithMessage(err, "eval prefix right side")
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		fmt.Println(">>>>>>> InfixExpression:", node)
		left, err := Generate(node.Left, env)
		fmt.Println("InfixExpression generated:", left)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator left side")
		}

		if node.IsControlExpression() && !node.MustEvaluateRight() && isTruthy(left) {
			return left, nil
		}

		fmt.Println("Generating right >>>>>", node.Operator)

		right, err := Generate(node.Right, env)

		fmt.Println("RIGHT INFIX", node.Right, node.Left)
		if node.Operator == ">" || node.Operator == "<" {
			generatedCondition := outputInfixIdentifers(node.Right, node.Left, node.Operator)

			if len(generatedCondition) > 0 {
				fmt.Println("Returning loop", generatedCondition)
				return &object.String{Value: generatedCondition}, nil
			}

		}
		fmt.Println("Right infix generated:", right)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator right side")
		}

		if node.IsControlExpression() {
			fmt.Println("ControlExpression")
			return right, nil
		}

		if node.Operator == "+" {
			identifier := outputInfixString(left, node.Left, right)
			return &object.String{Value: identifier}, nil
		}
		context := &callContext{object.NewCallContext(env, left)}
		fmt.Println("CONTEXT ________", context)
		return &object.String{Value: right.Inspect()}, nil
		//return object.Send(context, node.Operator, right)
	case *ast.ConditionalExpression:
		fmt.Println("ConditionalExpression:", node)
		return evalConditionalExpression(node, env)
	case *ast.ScopedIdentifier:
		fmt.Println("ScopedIdentifier:", node)
		self, _ := env.Get("self")
		outer, ok := env.Get(node.Outer.Value)
		if !ok {
			return nil, errors.Wrap(
				object.NewUndefinedLocalVariableOrMethodNameError(self, node.Outer.Value),
				"eval scope outer",
			)
		}
		outerEnv, ok := outer.(object.Environment)
		if !ok {
			return nil, errors.Wrap(
				object.NewUndefinedLocalVariableOrMethodNameError(self, node.Outer.Value),
				"eval scope outer",
			)
		}
		inner, err := Generate(node.Inner, outerEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval scope inner")
		}
		return inner, nil
	case *ast.ExceptionHandlingBlock:
		fmt.Println("ExceptionHandlingBlock", node)
		bodyReturn, err := Generate(node.TryBody, env)
		if err == nil {
			return bodyReturn, nil
		}
		return handleException(err, node.Rescues, env)

	case *ast.Comment:
		// ignore comments
		return nil, nil

	case *ast.LoopExpression:
		fmt.Println("While loop", node.Condition, node.Block)
		outputForLoop(node.Condition, env)
		ret, err := Generate(node.Block, env)
		fmt.Println("__________________")
		outputLoopEnd()
		if err == nil {
			return ret, nil
		}
		return nil, nil
	case nil:
		return nil, nil
	default:
		err := object.NewException("Unknown AST: %T", node)
		return nil, errors.WithStack(err)
	}

}

func evalProgram(stmts []ast.Statement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range stmts {
		if _, ok := statement.(*ast.Comment); ok {
			continue
		}

		result, err = Generate(statement, env)

		if err != nil {
			return nil, errors.WithMessage(err, "eval program statement")
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}

	}

	endMainFunc()
	return result, nil
}

func evalExpressions(exps []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	fmt.Println("EvalExpression: ", exps)
	for _, e := range exps {
		fmt.Println("Inside loop")
		evaluated, err := Generate(e, env)
		fmt.Println("Evaluated in eval exp:", e)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	fmt.Println("eval result:", result)
	return result, nil
}

func evalPrefixExpression(operator string, right object.RubyObject) (object.RubyObject, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: %s%s", operator, right.Type()))
	}
}

func evalBangOperatorExpression(right object.RubyObject) object.RubyObject {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NIL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.RubyObject) (object.RubyObject, error) {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}, nil
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: -%s", right.Type()))
	}
}

func evalConditionalExpression(ce *ast.ConditionalExpression, env object.Environment) (object.RubyObject, error) {
	fmt.Println("Eval condition expr:", ce.Condition)
	condition, err := Generate(ce.Condition, env)
	if err != nil {
		return nil, err
	}
	evaluateConsequence := isTruthy(condition)
	if ce.IsNegated() {
		evaluateConsequence = !evaluateConsequence
	}
	if evaluateConsequence {
		return Generate(ce.Consequence, env)
	} else if ce.Alternative != nil {
		return Generate(ce.Alternative, env)
	} else {
		return object.NIL, nil
	}
}

func evalIndexExpressionAssignment(left, index, right object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		integer, ok := index.(*object.Integer)
		if !ok {
			return nil, errors.Wrap(
				object.NewImplicitConversionTypeError(integer, index),
				"eval array index",
			)
		}
		idx := int(integer.Value)
		if idx >= len(target.Elements) {
			// enlarge slice
			for len(target.Elements) <= idx {
				target.Elements = append(target.Elements, object.NIL)
			}
		}
		target.Elements[idx] = right
		return right, nil
	case *object.Hash:
		target.Set(index, right)
		return right, nil
	default:
		return nil, errors.Wrap(
			object.NewException("assignment target not supported: %s", left.Type()),
			"eval IndexExpression Assignment",
		)
	}
}

func evalIndexExpression(left, index object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		return evalArrayIndexExpression(target, index), nil
	case *object.Hash:
		return evalHashIndexExpression(target, index), nil
	default:
		return nil, errors.WithStack(object.NewException("index operator not supported: %s", left.Type()))
	}
}

func evalArrayIndexExpression(arrayObject *object.Array, index object.RubyObject) object.RubyObject {
	idx := index.(*object.Integer).Value
	maxNegative := -int64(len(arrayObject.Elements))
	maxPositive := maxNegative*-1 - 1
	if maxPositive < 0 {
		return object.NIL
	}

	if idx > 0 && idx > maxPositive {
		return object.NIL
	}
	if idx < 0 && idx < maxNegative {
		return object.NIL
	}
	if idx < 0 {
		return arrayObject.Elements[len(arrayObject.Elements)+int(idx)]
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash *object.Hash, index object.RubyObject) object.RubyObject {
	result, ok := hash.Get(index)
	if !ok {
		return object.NIL
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env object.Environment) (object.RubyObject, error) {
	fmt.Println(block)
	var result object.RubyObject
	var err error
	for _, statement := range block.Statements {
		result, err = Generate(statement, env)
		fmt.Println("BlockSt:", statement, statement.String())
		if err != nil {
			return nil, err
		}
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ {
				return result, nil
			}

		}
	}
	if result == nil {
		return object.NIL, nil
	}
	return result, nil
}

func evalIdentifier(node *ast.Identifier, env object.Environment) (object.RubyObject, error) {
	fmt.Println(".ValueEvalIdentifer:", node)
	val, ok := env.Get(node.Value)
	fmt.Println("EvalIdentifer val:", val, ok)
	if ok {
		return val, nil
	}

	if node.IsConstant() {
		return nil, errors.Wrap(
			object.NewUninitializedConstantNameError(node.Value),
			"eval identifier",
		)
	}

	self, _ := env.Get("self")
	context := &callContext{object.NewCallContext(env, self)}
	val, err := object.Send(context, node.Value)
	if err != nil {
		return nil, errors.Wrap(
			object.NewUndefinedLocalVariableOrMethodNameError(self, node.Value),
			"eval ident as method call",
		)
	}
	fmt.Println("EvalIdentifer value:", val)
	return val, nil
}

func unwrapReturnValue(obj object.RubyObject) object.RubyObject {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func handleException(err error, rescues []*ast.RescueBlock, env object.Environment) (object.RubyObject, error) {
	if err != nil && len(rescues) == 0 {
		return nil, err
	}
	errorObject := err.(object.RubyObject)
	errClass := errorObject.Class().Name()
	rescueEnv := object.WithScopedLocalVariables(env)

	var catchAll *ast.RescueBlock
	for _, r := range rescues {
		if len(r.ExceptionClasses) == 0 {
			catchAll = r
			continue
		}
		if r.Exception != nil {
			rescueEnv.Set(r.Exception.Value, errorObject)
		}
		for _, cl := range r.ExceptionClasses {
			if cl.Value == errClass {
				rescueRet, err := Generate(r.Body, rescueEnv)
				return rescueRet, err
			}
		}
	}

	if catchAll != nil {
		ancestors := getAncestors(errorObject)
		sort.Strings(ancestors)
		if sort.SearchStrings(ancestors, "StandardError") >= len(ancestors) {
			return nil, err
		}

		if catchAll.Exception != nil {
			rescueEnv.Set(catchAll.Exception.Value, errorObject)
		}
		rescueRet, err := Generate(catchAll.Body, rescueEnv)
		return rescueRet, err
	}

	return nil, err
}

func getAncestors(obj object.RubyObject) []string {
	class := obj.Class()
	if c, ok := obj.(object.RubyClass); ok {
		class = c
	}
	var ancestors []string
	ancestors = append(ancestors, class.Name())

	superClass := class.SuperClass()
	if superClass != nil {
		superAncestors := getAncestors(superClass.(object.RubyClassObject))
		ancestors = append(ancestors, superAncestors...)
	}
	return ancestors
}

func isTruthy(obj object.RubyObject) bool {
	switch obj {
	case object.NIL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

// IsError returns true if the given RubyObject is an object.Error or an
// object.Exception (or any subclass of object.Exception)
func IsError(obj object.RubyObject) bool {
	if obj != nil {
		return obj.Type() == object.EXCEPTION_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) object.RubyObject {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func createMainFunc() {
	src := `package main
 
import "github.com/goruby/goruby/object"

func main() { 
	env := object.NewMainEnvironment()`
	appendToFile(src)
}

func endMainFunc() {
	src := `}`
	appendToFile(src)
}

func outputValue(value string, right object.RubyObject) {
	appendToFile("\tenv.Set(\"" + value + "\", " + right.Inspect() + ")")
}

func outputInteger(value int64) object.RubyObject {
	newVar := getNewVariableName()
	fmt.Println("********** Writing integer", newVar)
	src := `
	` + newVar + ` := object.NewInteger(` + strconv.FormatInt(int64(value), 10) + `)`
	appendToFile(src)
	stringObj := &object.String{Value: newVar}
	return stringObj
}

func outputString(value string) object.RubyObject {
	fmt.Println("********** Writing string")
	newVar := getNewVariableName()
	src := `
	` + newVar + ` := &object.String{Value:"` + value + `"}`
	appendToFile(src)
	stringObj := &object.String{Value: newVar}
	return stringObj
}

func outputGetIdentifier(node ast.Node) {
	//TODO: Don't disregard ok
	src := "\t" + node.String() + ", _ := env.Get(\"" + node.String() + "\") \n"
	//src = src + "\tfmt.Println(" + node.String() + ")"
	appendToFile(src)
}

func outputArray(exps []ast.Expression, env object.Environment, identifier string) (string, error) {
	exists := existingArrays[identifier]
	if len(exists) == 0 {
		appendToFile(`
	var result []object.RubyObject`)
	}
	for _, e := range exps {
		evaluated, err := Generate(e, env)
		//evalString := fmt.Sprintf("%v", e)
		fmt.Println("Evaluated within write ---:", evaluated, reflect.TypeOf(evaluated), reflect.ValueOf(e))

		if err != nil {
			return "", err
		}

		appendToFile(`
	result = append(result, ` + evaluated.Inspect() + `)`)

		//appendToFile(`
		//result = append(result, &object.String{Value:"` + evalString + `"})`)
	}
	newVar := getNewVariableName()
	src := `
	` + newVar + ` := &object.Array{Elements: result}`
	appendToFile(src)
	return newVar, nil
}

func outputForLoop(condition ast.Expression, env object.Environment) error {
	generatedCondition, _ := Generate(condition, env)
	src := `
	for ` + generatedCondition.Inspect() + ` {
	`
	appendToFile(src)
	return nil
}

func outputInfixIdentifers(right ast.Node, left ast.Node, operator string) string {
	appendToFile(`
	` + right.String() + `Val, _ := ` + right.String() + `.(*object.Integer)
	`)
	appendToFile(`
	` + left.String() + `Val, _ := ` + left.String() + `.(*object.Integer)
	`)

	return left.String() + "Val.Value" + operator + right.String() + "Val.Value"
}

func outputLoopEnd() {
	src := "}"

	appendToFile(src)
}

func outputInfixString(left object.RubyObject, nLeft ast.Expression, right object.RubyObject) string {
	identifier := getNewVariableName()
	fmt.Println("Output Infix op:", left.Type(), right.Inspect())
	leftVal := left.(*object.String)
	isInt := leftVal

	fmt.Println("Is it Int:", isInt, left.Inspect())

	if true {
		appendToFile("\t" + nLeft.String() + "Val_ := " + nLeft.String() + ".(*object.Integer)")
		appendToFile("\t" + nLeft.String() + "Val = " + nLeft.String() + "Val_")
		appendToFile("\t" + identifier + " := " + nLeft.String() + "Val_.Value +" + right.Inspect() + ".Value")

		return "&object.Integer{ Value:" + identifier + "}"
	} else {
		appendToFile("\t" + identifier + " := " + left.Inspect() + ".Inspect() +" + right.Inspect() + ".Value")
		return "&object.String{ Value:" + identifier + "}"
	}
}

func outputClassExpression(nodeVal ast.Node, env object.Environment) error {
	node := nodeVal.(*ast.ClassExpression)
	fmt.Println("ClassExpressionFunc:", node.Name.Value)
	src := `
	superClassName := "Object"
	`
	if node.SuperClass != nil {
		src = src + `superClassName = ` + node.SuperClass.Value
	}
	src = src + `
	superClass, ok := env.Get(superClassName)
	class, ok := env.Get("` + node.Name.Value + `")
	if !ok {
		class = object.NewClass("` + node.Name.Value + `", superClass.(object.RubyClassObject), env)
	}
	classEnv := class.(object.Environment)
	classEnv.Set("self", &object.Self{RubyObject: class, Name: "` + node.Name.Value + `"})
	`
	superClassName := "Object"
	if node.SuperClass != nil {
		superClassName = node.SuperClass.Value
	}
	superClass, ok := env.Get(superClassName)
	if !ok {
		return errors.Wrap(
			object.NewUninitializedConstantNameError(superClassName),
			"eval class superclass",
		)
	}
	class, ok := env.Get(node.Name.Value)
	if !ok {
		class = object.NewClass(node.Name.Value, superClass.(object.RubyClassObject), env)
	}
	classEnv := class.(object.Environment)
	classEnv.Set("self", &object.Self{RubyObject: class, Name: node.Name.Value})
	_, err := Generate(node.Body, classEnv)

	if err != nil {
		return errors.WithMessage(err, "eval class body")
	}
	src = src + `
	selfObject, _ := classEnv.Get("self")
	self := selfObject.(*object.Self)
	env.Set("` + node.Name.Value + `", self.RubyObject)
	`
	appendToFile(src)
	return nil
}

func outputContextCalls(node ast.Node) {
	nodeVal := node.String()
	fmt.Println(nodeVal)

	if strings.Contains(nodeVal, "puts") {
		appendToFile("fmt.Println(" + nodeVal + ")")
	}
}

func getNewVariableName() string {
	variableCount = variableCount + 1

	return "v" + strconv.Itoa(variableCount)
}

func appendToFile(content string) {
	f, err := os.OpenFile("compiled.go", os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(content + "\n"); err != nil {
		panic(err)
	}
}
