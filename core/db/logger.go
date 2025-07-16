// core/db/logger.go
package db

func LogInteraction(prompt, response string) {
	_, _ = DB.Exec("INSERT INTO interactions (user_input, llm_response) VALUES (?, ?)", prompt, response)
}
