package exp

import (
	"testing"
	"fmt"
)

func TestParse(t *testing.T) {
	scenes := []struct{
		input string
		want string
	}{
		{"x","x"},
		{"3.14","3.14"},
		{"-(x)","-x"},
		{"-(x+1)","-(x+1)"},
		{"x+y","x+y"},
		{"(x+y)*(x-y)","(x+y)*(x-y)"},
		{"x+y*x-y","x+y*x-y"},
		{"x+-(y*y+1)*x","x+-(y*y+1)*x"},
		{"sin(x+y*x)","sin(x+y*x)"},
		{"sin(x/pow(y,2))","sin(x/pow(y,2))"},
		{"sin(-x)*pow(1.5,-sqrt(x*x+y*y))","sin(-x)*pow(1.5,-sqrt(x*x+y*y))"},
		{"pow(2,sin(y))*pow(2,sin(x))/12","pow(2,sin(y))*pow(2,sin(x))/12"},
		{"sin(x*y/10)/10","sin(x*y/10)/10"},
	}
	for _, scene := range scenes {
		expr, err := Parse(scene.input)
		if err != nil {
			t.Errorf("parse err: %s", err)
		}
		str := expr.String()
		if str != scene.want {
			t.Errorf("toString: %q, want: %q", str, scene.want)
		}
	}
}

func TestConv(t *testing.T) {
	var x = 1.6789
	var y = FloatToFixed(x)
	var z = FixedToFloat(y)
	fmt.Printf("%v\n%b\n%v\n%v", x,y,y,z)

	fmt.Println("\n--------------------")
	x = -x
	y = FloatToFixed(x)
	z = FixedToFloat(y)
	fmt.Printf("%v\n%b\n%v\n%v", x,y,y,z)
}