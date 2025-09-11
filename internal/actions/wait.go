package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("wait", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		ok, err := step.Validate(step, []string{"Duration"})
		if !ok || err != nil {
			flow.ErrorLog(fmt.Sprintf("%s: %v", err, step))
			return false
		}
		fmt.Println("Waiting for ", step.Duration, "ms...")
		page.WaitForTimeout(float64(step.Duration))
		return true
	})
}
