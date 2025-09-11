package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("press", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := step.Validate(step, []string{"Target", "Value"})
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("press validation failed: %v", err))
			return false
		}

		loc := page.Locator(step.Target)
		err = loc.Press(step.Value) // e.g. "Enter"
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("press failed: %v", err))
			return false
		}
		return true
	})
}
