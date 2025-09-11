package core

import (
	"fmt"
	"reflect"

	"github.com/playwright-community/playwright-go"
)

type Step struct {
	Action        string `json:"action"`
	Target        string `json:"target,omitempty"`
	Description   string `json:"description,omitempty"`
	Duration      int    `json:"duration,omitempty"`
	Selector      string `json:"selector,omitempty"`
	ChildSelector string `json:"child_selector,omitempty"`
	Attribute     string `json:"attribute,omitempty"`
	StoreAs       string `json:"store_as,omitempty"`
	Filename      string `json:"filename,omitempty"`
	Value         string `json:"value,omitempty"`
	Step          []Step `json:"step,omitempty"`
}

func (step Step) Validate(s Step, mands []string) (bool, error) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	for _, mand := range mands {
		field, ok := typ.FieldByName(mand)
		if !ok {
			return false, fmt.Errorf("unknown field: %s", mand)
		}
		f := val.FieldByName(field.Name)
		if f.Kind() == reflect.String && f.String() == "" {
			return false, fmt.Errorf("missing required field: %s", mand)
		}
	}
	return true, nil
}

func (step Step) run(flow *Flow, page playwright.Page) bool {
	if handler, ok := handlers[step.Action]; ok {
		return handler(step, flow, page)
	}
	fmt.Printf("Unknown action: %s\n", step.Action)
	return false
}
