/**

This code is for quadratic calculator
using Golang
**/

package main

import (
	"fmt"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head>
<title>Quadratic Equation Solver</title><body>
<h3>Quadratic Equation Solver</h3><p>Solves equations of the form
a<i>x</i>² + b<i>x</i> + c</p>`
	form = `<form action="/" method="POST">
<input type="text" name="a" size="1"><label for="a"><i>x</i>²</label> +
<input type="text" name="b" size="1"><label for="b"><i>x</i></label> +
<input type="text" name="c" size="1"><label for="c"> →</label>
<input type="submit" name="calculate" value="Calculate">
</form>`
	pageBottom = "</body></html>"
	error      = `<p class="error">%s</p>`
	solution   = "<p>%s → %s</p>"
)

func main() {
	http.HandleFunc("/", mainPage)

	if err := http.ListenAndServe(":9020", nil); err != nil {
		log.Fatal("Enable to connect", err)
	}

}

func mainPage(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()

	fmt.Fprint(writer, pageTop, form)

	if err != nil {
		fmt.Print("Error")

	} else {
		if values, message, ok := processInfo(request); ok && message != "zeros" {
			roots := rootFinder(values)
			question := formatQuestion(request.Form)
			answer := formatSolution(roots[0], roots[1])

			fmt.Fprintf(writer, solution, question, answer)
		}
	}
}

func EqualComplex(root1 complex128, root2 complex128) bool {
	if cmplx.Tanh(root1) == cmplx.Tanh(root2) {
		return true
	}
	return false
}

func formatQuestion(form map[string][]string) string {

	return fmt.Sprintf("%s<i>x</i>² + %s<i>x</i> + %s", form["a"][0],
		form["b"][0], form["c"][0])
}
func formatSolution(x1, x2 complex128) string {
	// Ugly formatting since it always shows a complex even if the imag
	// part is 0i.
	if EqualComplex(x1, x2) {
		return fmt.Sprintf("<i>x</i>=%f", x1)
	}
	return fmt.Sprintf("<i>x</i>=%f or <i>x</i>=%f", x1, x2)
}

func processInfo(request *http.Request) ([3]float64, string, bool) {

	var values [3]float64
	for index, value := range []string{"a", "b", "c"} {
		if item, found := request.Form[value]; found {
			if value, err := strconv.ParseFloat(item[0], 64); err != nil {
				return values, item[0] + "is invalid input", false
			} else {
				values[index] = value
			}
		}
	}
	if values[0] == values[1] && values[2] == values[1] && values[2] == float64(0) {
		return values, "zeros", true
	}
	return values, "", true
}

func rootFinder(values [3]float64) [2]complex128 {
	var roots [2]complex128

	a, b, c := complex(values[0], 0), complex(values[1], 0),
		complex(values[2], 0)
	root := cmplx.Sqrt(cmplx.Pow(b, 2) - (4 * a * c))
	roots[1] = (-b + root) / (2 * a)
	roots[0] = (-b - root) / (2 * a)

	return roots

}
