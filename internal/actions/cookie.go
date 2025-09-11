package actions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("cookies", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		context := page.Context()

		if step.Value == "save" {
			cookies, err := context.Cookies()
			if err != nil {
				flow.ErrorLog(fmt.Sprintf("get cookies failed: %v", err))
				return false
			}
			data, _ := json.MarshalIndent(cookies, "", "  ")
			os.WriteFile(step.Filename, data, 0644)
			fmt.Println("Cookies saved to", step.Filename)
		} else if step.Value == "load" {
			data, err := os.ReadFile(step.Filename)
			if err != nil {
				flow.ErrorLog(fmt.Sprintf("read cookies failed: %v", err))
				return false
			}
			var cookies []playwright.OptionalCookie
			json.Unmarshal(data, &cookies)
			err = context.AddCookies(cookies)
			if err != nil {
				flow.ErrorLog(fmt.Sprintf("load cookies failed: %v", err))
				return false
			}
			fmt.Println("Cookies loaded from", step.Filename)
		}
		return true
	})
}
