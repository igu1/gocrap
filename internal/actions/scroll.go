package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("scroll", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		target := strings.ToLower(step.Value)

		switch target {
		case "top":
			_, err := page.Evaluate(`window.scrollTo(0,0)`)
			if err != nil {
				flow.ErrorLog(fmt.Sprintf("scroll top failed: %v", err))
				return false
			}
		case "bottom":
			_, err := page.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`)
			if err != nil {
				flow.ErrorLog(fmt.Sprintf("scroll bottom failed: %v", err))
				return false
			}
		default:
			parts := strings.Split(step.Value, ",")
			if len(parts) == 2 {
				x, _ := strconv.Atoi(parts[0])
				y, _ := strconv.Atoi(parts[1])
				_, err := page.Evaluate(fmt.Sprintf(`window.scrollTo(%d,%d)`, x, y))
				if err != nil {
					flow.ErrorLog(fmt.Sprintf("scroll coords failed: %v", err))
					return false
				}
			} else {
				flow.ErrorLog("invalid scroll value, use 'top', 'bottom' or 'x,y'")
				return false
			}
		}
		return true
	})
}
