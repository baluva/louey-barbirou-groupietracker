package templates

import (
	"fmt"
	"html/template"
	"os"
)

var Temp *template.Template

func InitTemplate() {

	temp, errTemp := template.ParseGlob("./templates/*.html")
	if errTemp != nil {
		fmt.Println("err")
		os.Exit(1)
	}
	Temp = temp
}
