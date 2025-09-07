package actions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/igu1/gocrap/internal/core"
	"github.com/playwright-community/playwright-go"
)

func init() {
	core.RegisterHandler("save", func(step core.Step, flow *core.Flow, page playwright.Page) bool {
		isvalid, err := step.Validate(step, []string{"Filename"})
		if !isvalid {
			flow.ErrorLog("Error saving:" + err.Error())
			return false
		}
		filename := step.Filename
		fmt.Println("Saving to", filename)
		data, err := json.MarshalIndent(flow.Mem, "", "  ")
		if err != nil {
			flow.ErrorLog("Error marshaling data:" + err.Error())
			return false
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			flow.ErrorLog("Error writing file:" + err.Error())
			return false
		}
		return true
	})
}
