package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("check", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := step.Validate(step, []string{"Target"})
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("check validation failed: %v", err))
			return false
		}

		loc := page.Locator(step.Target)
		err = loc.Check()
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("check failed: %v", err))
			return false
		}
		return true
	})
}
