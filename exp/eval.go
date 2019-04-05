package exp

import (
	"fmt"
	"math"
	"strconv"
)

type Expr interface{
	Eval(env Env) float64
	String() string
	GetGraphy() Graphy
}

//--------------------------------------------------------------------------------------
type Var string

func (v Var)Eval(env Env) float64{
	return env[v]
}

func (v Var)String() string {
	return string(v)
}

func (v Var)GetGraphy() Graphy {
	return &VarGraphy{Name:string(v)}
}

//---------------------------------------------------------------------------------------
type Literal string

func (l Literal)Eval(env Env) float64 {
	f, err := strconv.ParseFloat(string(l), 64)
	if err != nil {
		panic(fmt.Sprintf("literal eval, conv fail: %s", l))
	}
	return f
}

func (l Literal) String() string {
	return string(l)
}

func (l Literal)GetGraphy() Graphy {
	return &VarGraphy{Name:string(l)}
}

//--------------------------------------------------------------------------------------
type Unary struct {
	op rune // one of '+', '-'
	x  Expr
}

func (u Unary)Eval(env Env) float64 {
	switch u.op {
	case '+':
		return u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported unary op: %q", u.op))
	}
}

func (u Unary)String() string {
	switch u.op {
	case '+':
		return u.x.String()
	case '-':
		switch x := u.x.(type) {
		case Binary:
			if x.op == '+' || x.op == '-' {
				return "-(" + x.String() + ")"
			} else {
				return "-" + x.String()
			}
		default:
			return "-" + x.String()
		}
	default:
		panic(fmt.Sprintf("unsupported unary op: %q", u.op))
	}
}

func (u Unary)GetGraphy() Graphy {
	return &UnaryGraphy{Op:string(u.op), X:u.x.GetGraphy()}
}
//-----------------------------------------------------------------------------------------

type Binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

func (b Binary)Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported binary op: %q", b.op))
	}
}

func (b Binary)String() string {
	switch b.op {
	case '+','-':
		return b.x.String() + string(b.op) + b.y.String()
	case '*','/':
		var xStr, yStr string
		switch x := b.x.(type) {
		case Binary:
			if x.op == '+' || x.op == '-' {
				xStr = x.wrap()
			} else {
				xStr = x.String()
			}
		default:
			xStr = x.String()
		}

		switch y := b.y.(type) {
		case Binary:
			if y.op == '+' || y.op == '-' {
				yStr = y.wrap()
			} else {
				yStr = y.String()
			}
		default:
			yStr = y.String()
		}
		return xStr + string(b.op) + yStr
	default:
		panic(fmt.Sprintf("unsupported binary op: %q", b.op))
	}
}

func (b Binary)wrap() string {
	return "(" + b.String() + ")"
}

func (b Binary)GetGraphy() Graphy {
	if b.op == '/' {
		return &DivideGraphy{X:b.x.GetGraphy(),Y:b.y.GetGraphy()}
	}
	return &BinaryGraphy{Op:string(b.op),X:b.x.GetGraphy(),Y:b.y.GetGraphy()}
}

//-------------------------------------------------------------------------------------------

type Call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

func (c Call)Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	default:
		panic(fmt.Sprintf("unsupported call: %q", c.fn))
	}
}

func (c Call)String() string {
	switch c.fn {
	case "pow":
		return c.fn + "(" + c.args[0].String() + "," + c.args[1].String() + ")"
	default:
		return c.fn + "(" + c.args[0].String() + ")"
	}
}

func (c Call)GetGraphy() Graphy {
	switch c.fn {
	case "pow":
		return &PowGraphy{X:c.args[0].GetGraphy(),Y:c.args[1].GetGraphy()}
	case "sqrt":
		return &SqrtGraphy{X:c.args[0].GetGraphy()}
	default:
		return &TriangleGraphy{Op:c.fn, X:c.args[0].GetGraphy()}
	}
}

//---------------------------------------------------------------------------------------------
type Env map[Var]float64