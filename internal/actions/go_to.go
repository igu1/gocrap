package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("go_to", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		isvalid, err := step.Validate(step, []string{"Target"})
		if !isvalid {
			flow.ErrorLog("Error navigating:" + err.Error())
			return false
		}
		fmt.Println("Go to", step.Target)
		_, err = page.Goto(flow.Url + step.Target)
		if err != nil {
			flow.ErrorLog("Error navigating:" + err.Error())
			return false
		}
		flow.Mem["curr_page"] = flow.Url + step.Target
		return true
	})
}
