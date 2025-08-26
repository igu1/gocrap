package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

const (
	GOTO             = "go_to"
	EXTRACT_TEXT     = "extract_text"
	EXTRACT_ELEMENTS = "extract_els"
	SAVE             = "save"
	VISIT_EACH       = "visit_each"
)

type Step struct {
	Action      string `json:"action"`
	Target      string `json:"target,omitempty"`
	Description string `json:"description,omitempty"`
	Selector    string `json:"selector,omitempty"`
	StoreAs     string `json:"store_as,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Step        []Step `json:"step,omitempty"`
}

func (step Step) run(flow *Flow, page playwright.Page) bool {
	switch step.Action {
	case GOTO:
		fmt.Println("Go to", step.Target)
		_, err := page.Goto(flow.Url + step.Target)
		if err != nil {
			fmt.Println("Error navigating:", err)
			return false
		}
		flow.Mem["curr_page"] = flow.Url + step.Target

	case EXTRACT_TEXT:
		fmt.Printf("Extracting from: %s, selector: %s\n", flow.Mem["curr_page"], step.Selector)

		_, err := page.WaitForSelector(step.Selector)
		if err != nil {
			fmt.Println("Error waiting for selector:", err)
			return false
		}

		elements, err := page.QuerySelectorAll(step.Selector)
		if err != nil {
			fmt.Println("Error selecting elements:", err)
			return false
		}

		var results []string
		for _, el := range elements {
			text, _ := el.TextContent()
			if text != "" {
				fmt.Println("Extracted:", strings.TrimSpace(text))
				results = append(results, strings.TrimSpace(text))
			}
		}
		flow.Mem[step.StoreAs] = strings.Join(results, "; ")

	case EXTRACT_ELEMENTS:
		elements, err := page.QuerySelectorAll(step.Selector)
		if err != nil {
			fmt.Println("Error selecting elements:", err)
			return false
		}
		var results []playwright.ElementHandle
		for _, el := range elements {
			results = append(results, el)
		}
		flow.Mem[step.StoreAs] = results

	case VISIT_EACH:
		//ITS TAKING FROM MEMORY, SO ONLY PREV NODE CAN PASS THE URL
		urls_str, ok := flow.Mem[step.Target].(string)
		urls := strings.Split(urls_str, ",")
		fmt.Println("Visiting urls:", urls)
		if !ok {
			fmt.Println("No URLs found in memory for key:", step.Target)
			return false
		}
		for _, url := range urls {
			fmt.Println("Visiting:", url)
			flow.Mem["curr_page"] = url
			_, err := page.Goto(url)
			if err != nil {
				fmt.Println("Error navigating:", err)
				return false
			}
			for _, subStep := range step.Step {
				if !subStep.run(flow, page) {
					return false
				}
			}
		}

	case SAVE:
		filename := step.Filename
		fmt.Println("Saving to", filename)
		data, err := json.MarshalIndent(flow.Mem, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			return false
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
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

func validateJson(flow Flow) error {
	return nil
}

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatal(err)
	}
	flow := Flow{
		Mem: make(map[string]interface{}),
	}
	file, err := ioutil.ReadFile("flow.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(file, &flow)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateJson(flow)
	if err != nil {
		fmt.Println(err)
		return
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	for _, step := range flow.Path {
		step.run(&flow, page)
	}
	defer page.Close()
	err = browser.Close()
	if err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
}
