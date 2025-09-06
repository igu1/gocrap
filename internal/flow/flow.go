package flow

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/playwright-community/playwright-go"
)

const (
	GOTO          = "go_to"
	EXTRACT       = "extract"
	EXTRACT_MULTI = "extract_multi"
	SAVE          = "save"
)

type Step struct {
	Action        string `json:"action"`
	Target        string `json:"target,omitempty"`
	Description   string `json:"description,omitempty"`
	Selector      string `json:"selector,omitempty"`
	ChildSelector string `json:"child_selector,omitempty"`
	Attribute     string `json:"attribute,omitempty"`
	StoreAs       string `json:"store_as,omitempty"`
	Filename      string `json:"filename,omitempty"`
	Step          []Step `json:"step,omitempty"`
}

func (step Step) validate(s Step, mands []string) (bool, error) {
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
	switch step.Action {
	case GOTO:
		isvalid, err := step.validate(step, []string{"Target"})
		if !isvalid {
			flow.errorLog("Error navigating:" + err.Error())
			return false
		}
		fmt.Println("Go to", step.Target)
		_, err = page.Goto(flow.Url + step.Target)
		if err != nil {
			flow.errorLog("Error navigating:" + err.Error())
			return false
		}
		flow.Mem["curr_page"] = flow.Url + step.Target

	case EXTRACT:
		isvalid, err := step.validate(step, []string{"Selector", "StoreAs"})
		if !isvalid {
			flow.errorLog("Error extracting:" + err.Error())
			return false
		}
		el, err := page.WaitForSelector(step.Selector)
		if err != nil {
			flow.errorLog("Error navigating:" + err.Error())
			return false
		}
		if el == nil {
			flow.errorLog("No element found for selector: " + step.Selector)
			return false
		}
		var data string
		if step.Attribute != "" {
			data, _ = el.GetAttribute(step.Attribute)
		} else {
			data, _ = el.TextContent()
		}
		if data != "" {
			if old, ok := flow.Mem[step.StoreAs].(string); ok {
				flow.Mem[step.StoreAs] = []string{old, data}
			} else {
				flow.Mem[step.StoreAs] = data
			}
		}
	case EXTRACT_MULTI:
		isvalid, err := step.validate(step, []string{"Selector", "StoreAs"})
		if !isvalid {
			flow.errorLog("Error muti extracting:" + err.Error())
			return false
		}
		elements, err := page.QuerySelectorAll(step.Selector)
		if err != nil {
			flow.errorLog("Error selecting child elements:" + err.Error())
			return false
		}
		extractedData := []string{}
		for _, el := range elements {
			if el == nil {
				continue
			}
			var data string
			if step.Attribute != "" {
				data, _ = el.GetAttribute(step.Attribute)
			} else {
				data, _ = el.TextContent()
			}
			if data != "" {
				extractedData = append(extractedData, data)
			}
		}
		if old, ok := flow.Mem[step.StoreAs]; ok {
			if slice, ok := old.([]string); ok {
				flow.Mem[step.StoreAs] = append(slice, extractedData...)
			} else {
				flow.Mem[step.StoreAs] = extractedData
			}
		} else {
			flow.Mem[step.StoreAs] = extractedData
		}

	case SAVE:
		isvalid, err := step.validate(step, []string{"Filename"})
		if !isvalid {
			flow.errorLog("Error saving:" + err.Error())
			return false
		}
		filename := step.Filename
		fmt.Println("Saving to", filename)
		data, err := json.MarshalIndent(flow.Mem, "", "  ")
		if err != nil {
			flow.errorLog("Error marshaling data:" + err.Error())
			return false
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			flow.errorLog("Error writing file:" + err.Error())
			return false
		}
	}
	return true
}

type Flow struct {
	Title string                 `json:"title"`
	Url   string                 `json:"url"`
	Path  []Step                 `json:"path"`
	Mem   map[string]interface{} `json:"mem"`
}

func (f *Flow) errorLog(err string) {
	errs, _ := f.Mem["errors"].([]string)
	errs = append(errs, err)
	f.Mem["errors"] = errs
	fmt.Println(err)
}

func (f Flow) Run(flow Flow) Flow {
	err := validateJson(flow)
	if err != nil {
		f.errorLog(err.Error())
		return f
	}
	pw, err := playwright.Run()
	if err != nil {
		f.errorLog(err.Error())
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Timeout:  playwright.Float(50000),
	})
	if err != nil {
		f.errorLog(err.Error())
	}
	page, err := browser.NewPage()
	if err != nil {
		f.errorLog(err.Error())
	}
	for _, step := range flow.Path {
		step.run(&flow, page)
	}
	defer page.Close()
	err = browser.Close()
	if err != nil {
		f.errorLog(err.Error())
	}
	return flow
}

func validateJson(flow Flow) error { // TODO: implement schema checks if needed
	return nil
}
