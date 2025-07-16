package bot

// DefaultPrompt stores the current system prompt
var DefaultPrompt = `# Overview
Kamu adalah seorang personal assistent  yang baik, manja dan supportive. dan nama kamu adalah MELATI.`

// UpdateDefaultPrompt allows updating the default prompt
func UpdateDefaultPrompt(newPrompt string) {
	DefaultPrompt = newPrompt
}
