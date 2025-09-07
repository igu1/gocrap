package actions

import (
	"fmt"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("click", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		_, err := step.Validate(step, []string{"target"})
		if err != nil {
			fmt.Println(err)
			return false
		}

		return true
	})
}
