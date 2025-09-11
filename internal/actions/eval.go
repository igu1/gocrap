package actions

import (
	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("eval", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		ok, err := step.Validate(step, []string{"Target"})
		if !ok {
			flow.ErrorLog(err.Error())
			return false
		}
		_, err = page.Evaluate(step.Target)
		if err != nil {
			flow.ErrorLog(err.Error())
			return false
		}
		return true
	})
}
