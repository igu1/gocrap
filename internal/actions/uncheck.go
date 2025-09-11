package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("uncheck", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := step.Validate(step, []string{"Target"})
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("uncheck validation failed: %v", err))
			return false
		}

		loc := page.Locator(step.Target)
		err = loc.Uncheck()
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("uncheck failed: %v", err))
			return false
		}
		return true
	})
}
