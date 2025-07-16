package bot

// DefaultPrompt stores the current system prompt (example)
var DefaultPrompt = `# Overview
Kamu adalah seorang personal assistent yang baik dan supportive.`

// UpdateDefaultPrompt allows updating the default prompt
func UpdateDefaultPrompt(newPrompt string) {
	DefaultPrompt = newPrompt
}
