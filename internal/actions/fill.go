package actions

import (
	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("fill", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		ok, err := step.Validate(step, []string{"Target", "Value"})
		if !ok {
			flow.ErrorLog(err.Error())
			return false
		}
		loc, err := page.WaitForSelector(step.Target)
		if err != nil {
			flow.ErrorLog(err.Error())
			return false
		}
		err = loc.Fill(step.Value)
		if err != nil {
			flow.ErrorLog(err.Error())
			return false
		}
		return true
	})
}
