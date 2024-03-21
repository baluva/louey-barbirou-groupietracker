package main

import (
	"projet/routeur"
	"projet/templates"
)

func main() {
	templates.InitTemplate()
	routeur.Initserv()
}
