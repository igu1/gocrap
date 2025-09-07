package core

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type StepHandler func(step Step, flow *Flow, page playwright.Page) bool

var handlers = map[string]StepHandler{}

func RegisterHandler(name string, fn StepHandler) {
	fmt.Println("Registering Handler: " + name)
	handlers[name] = fn
}
