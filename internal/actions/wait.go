package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("wait", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		fmt.Println("Waiting for ", step.Duration, "ms...")
		page.WaitForTimeout(float64(step.Duration))
		return true
	})
}
