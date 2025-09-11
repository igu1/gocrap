package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("forward", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := page.GoForward()
		if err != nil {
			flow.ErrorLog(fmt.Sprintf("forward failed: %v", err))
			return false
		}
		return true
	})
}
