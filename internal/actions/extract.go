package actions

import (
	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("extract", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		isvalid, err := step.Validate(step, []string{"Selector", "StoreAs"})
		if !isvalid {
			flow.ErrorLog("Error extracting:" + err.Error())
			return false
		}
		el, err := page.WaitForSelector(step.Selector)
		if err != nil {
			flow.ErrorLog("Error navigating:" + err.Error())
			return false
		}
		if el == nil {
			flow.ErrorLog("No element found for selector: " + step.Selector)
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
		return true
	})
}
