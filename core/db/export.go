// core/db/export.go
package db

import (
	"encoding/json"
	"fmt"
	"os"
)

type Interaction struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

func ExportAsJSONL(filename string) error {
	rows, err := DB.Query("SELECT user_input, llm_response FROM interactions")
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for rows.Next() {
		var entry Interaction
		if err := rows.Scan(&entry.Prompt, &entry.Completion); err != nil {
			return err
		}
		line, _ := json.Marshal(entry)
		fmt.Fprintln(file, string(line))
	}
	return nil
}
