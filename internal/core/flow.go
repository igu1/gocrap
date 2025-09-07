package core

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
)

type Flow struct {
	Title string                 `json:"title"`
	Url   string                 `json:"url"`
	Path  []Step                 `json:"path"`
	Mem   map[string]interface{} `json:"mem"`
}

func (f *Flow) ErrorLog(err string) {
	errs, _ := f.Mem["errors"].([]string)
	errs = append(errs, err)
	f.Mem["errors"] = errs
	fmt.Println(err)
}

func (f Flow) Run(flow Flow) Flow {
	err := validateJson(flow)
	if err != nil {
		f.ErrorLog(err.Error())
		return f
	}
	pw, err := playwright.Run()
	if err != nil {
		f.ErrorLog(err.Error())
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Timeout:  playwright.Float(50000),
	})
	if err != nil {
		f.ErrorLog(err.Error())
	}
	page, err := browser.NewPage()
	if err != nil {
		f.ErrorLog(err.Error())
	}
	for _, step := range flow.Path {
		step.run(&flow, page)
	}
	defer page.Close()
	err = browser.Close()
	if err != nil {
		f.ErrorLog(err.Error())
	}
	return flow
}

func validateJson(flow Flow) error { // TODO: implement schema checks if needed
	return nil
}
