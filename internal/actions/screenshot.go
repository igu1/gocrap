package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("screenshot", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := step.Validate(step, []string{"Filename"})
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("screenshot validation failed: %v", err))
			return false
		}

		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(step.Filename),
		})
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("screenshot failed: %v", err))
			return false
		}

		fmt.Println("Saved screenshot:", step.Filename)
		return true
	})
}
