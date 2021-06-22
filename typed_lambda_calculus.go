package main

import "fmt"

type Expression interface {
	pretty() string
	prettyType() string
	reduce() Expression
	getType() Type
}

type TypeFlag int

const (
	TypeIllTyped TypeFlag = 0
	TypeInt      TypeFlag = 1
	TypeBool     TypeFlag = 2
	TypeFunction TypeFlag = 3
)

type Type interface {
	pretty() string
	//getFlag() TypeFlag
}

type IllTyped int
type TyInt int
type TyBool bool
type TyFunction [2]Type

func (tyIllType IllTyped) pretty() string {
	return "IllTyped"
}

func (tyInt TyInt) pretty() string {
	return "Int"
}

func (tyBool TyBool) pretty() string {
	return "Bool"
}

func (tyFunction TyFunction) pretty() string {
	return "(" + tyFunction[0].pretty() + " -> " + tyFunction[1].pretty() + ")"
}

// type ValueKind int

// const (
// 	ValueInt  ValueKind = 0
// 	ValueBool ValueKind = 1
// 	Undefined ValueKind = 2
// )

// func createBool(x bool) Value {
// 	return Value{flag: ValueBool, valueBool: x}
// }

// func createInt(x int) Value {
// 	return Value{flag: ValueInt, valueInt: x}
// }

// func createUndefinde() Value {
// 	return Value{flag: Undefined}
// }

func getPrettyType(t Type) string {
	switch t.(type) {
	case TyFunction:
		new_function := t.(TyFunction)
		return getPrettyType(new_function[1])
	}

	return t.pretty()
}

// type Value struct {
// 	flag      ValueKind
// 	valueInt  int
// 	valueBool bool
// }

type Variable struct {
	name         string
	variableType Type
}

type Lambda struct {
	body     Expression
	argument Variable
}

type Application struct {
	function Expression
	argument Expression
}

func (variable Variable) pretty() string {
	return variable.name
}

func (variable Variable) prettyType() string {
	return variable.variableType.pretty()
}

func (variable Variable) reduce() Expression {
	return variable
}

// func (variable Variable) evaluate() Value {
// 	return variable.value
// }

func (variable Variable) getType() Type {
	return variable.variableType
}

func (application Application) pretty() string {
	return "(" + application.function.pretty() + " " + application.argument.pretty() + ")"
}

func (application Application) prettyType() string {
	return getPrettyType(application.function.getType())
}

func (x Application) reduce() Expression {

	x.function = x.function.reduce()
	x.argument = x.argument.reduce()

	switch x.function.(type) {
	case Lambda:
		new_abstraction := x.function.(Lambda)
		substitution := substitute(new_abstraction.body, new_abstraction.argument, x.argument)
		//fmt.Printf(substitution.pretty() + "\n")
		return substitution.reduce()
	}

	return Application{x.function, x.argument}
}

func substitute(function_body Expression, variable Variable, argument Expression) Expression {
	switch function_body.(type) {
	case Variable:
		v := function_body.(Variable)
		if v.name == variable.name {
			return argument
		} else {
			return function_body
		}

	case Application:
		app := function_body.(Application)
		return Application{substitute(app.function, variable, argument),
			substitute(app.argument, variable, argument)}

	default:
		new_abstraction := function_body.(Lambda)
		new_variable := (Variable)(new_abstraction.argument)
		new_body := substitute(new_abstraction.body, new_abstraction.argument, new_variable)
		return Lambda{substitute(new_body, variable, argument), new_variable}
	}
}

// func (application Application) evaluate() Value {
// 	return application.function.evaluate()
// }

func (application Application) getType() Type {

	return application.function.getType()
}

func (lambda Lambda) pretty() string {
	return "(" + "\\" + lambda.argument.pretty() + " -> " + lambda.body.pretty() + ")"
}

func (lambda Lambda) prettyType() string {
	return "(" + lambda.argument.prettyType() + " -> " + lambda.body.prettyType() + ")"
}

func (x Lambda) reduce() Expression {
	return Lambda{x.body.reduce(), x.argument}
}

// func (lambda Lambda) evaluate() Value {
// 	return lambda.body.evaluate()
// }

func (lambda Lambda) getType() Type {
	return lambda.body.getType()
}

func variable(name string, variableType Type) Variable {
	return Variable{name, variableType}
}

func lambda(name Variable, body Expression) Lambda {
	return Lambda{body, name}
}

func application(function Expression, variable Expression) Application {
	return Application{function, variable}
}

func tyBool() TyBool {
	return (TyBool)(true)
}

func tyInt() TyInt {
	return (TyInt)(0)
}

func illType() IllTyped {
	return (IllTyped)(0)
}

func tyFunction(x Type, y Type) TyFunction {
	return TyFunction{x, y}
}

func main() {
	x := application(lambda(variable("x", tyBool()), lambda(variable("y", tyInt()), variable("y", tyInt()))), variable("z", tyInt()))
	fmt.Printf(x.reduce().prettyType() + "\n")

	// y := x.reduce()
	// fmt.Printf(y.pretty() + "\n\n")

	// lambda := application(abstraction(variable("f"),
	// 	abstraction(variable("x"), application(variable("f"), application(variable("f"), variable("x"))))), application(abstraction(variable("x"), variable("x")), abstraction(variable("x"), variable("x"))))
	// fmt.Printf(lambda.pretty() + " => " + lambda.reduce().pretty() + "\n\n")

	// terminate := application(abstraction(variable("x"), application(variable("x"), variable("x"))), abstraction(variable("x"), application(variable("x"), variable("x"))))
	// fmt.Printf(terminate.pretty() + "\n")
	// reduced_terminate := terminate.reduce()
	// fmt.Printf(reduced_terminate.pretty() + "\n")

	// lambda2 := application(application(abstraction(variable("x"), application(variable("x"), variable("y"))), abstraction(variable("x"), variable("x"))), abstraction(variable("x"), variable("x")))
	// fmt.Printf(lambda2.pretty() + " => " + lambda2.reduce().pretty() + "\n\n")

	// lambda3 := application(application(abstraction(variable("u"), Variable("u")), variable("z")), variable("x"))
	// fmt.Printf(lambda3.pretty() + " => " + lambda3.reduce().pretty() + "\n\n")

	// lambda4 := application(abstraction(variable("x"), application(variable("x"), variable("y"))), application(abstraction(variable("x"), variable("x")), abstraction(variable("x"), variable("x"))))
	// fmt.Printf(lambda4.pretty() + " => " + lambda4.reduce().pretty() + "\n\n")

	typed_lambda := lambda(variable("y", tyInt()), lambda(variable("x", tyFunction(tyInt(), tyBool())), application(variable("x", tyFunction(tyInt(), tyBool())), variable("y", tyInt()))))
	fmt.Printf(typed_lambda.prettyType() + "\n")

}
