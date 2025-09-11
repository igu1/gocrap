package actions

import (
	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("extract_multi", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		isvalid, err := step.Validate(step, []string{"Selector", "StoreAs"})
		if !isvalid {
			flow.ErrorLog("Error extracting:" + err.Error())
			return false
		}
		elementHandles, err := page.QuerySelectorAll(step.Selector)
		if err != nil {
			flow.ErrorLog("Error navigating:" + err.Error())
			return false
		}
		if elementHandles == nil {
			flow.ErrorLog("No element found for selector: " + step.Selector)
			return false
		}
		var datas []string
		for _, e := range elementHandles {
			var data string
			if step.Attribute != "" {
				data, _ = e.GetAttribute(step.Attribute)
			} else {
				data, _ = e.TextContent()
			}
			if data != "" {
				datas = append(datas, data)
			}
		}
		flow.Mem[step.StoreAs] = datas
		return true
	})
}
