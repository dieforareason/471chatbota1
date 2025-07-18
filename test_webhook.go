package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Test payload
	payload := map[string]interface{}{
		"event": "message_received",
		"data": map[string]interface{}{
			"key": map[string]interface{}{
				"remoteJid": "yournumber@s.whatsapp.net", //your number with @s.whatsapp.net
			},
			"message": map[string]interface{}{
				"conversation": "Hello bot",
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}

	// Send request to webhook
	resp, err := http.Post("http://localhost:8080/webhook/wa", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body))
}
