// Package bot provides the main bot functionality
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func init() {
	// Set up logging with timestamps
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

// WhatsAppWebhookPayload represents the incoming webhook payload structure
type WhatsAppWebhookPayload struct {
	ChatID  string `json:"chat_id"`
	From    string `json:"from"`
	Message struct {
		Text          string `json:"text"`
		ID            string `json:"id"`
		RepliedID     string `json:"replied_id"`
		QuotedMessage string `json:"quoted_message"`
	} `json:"message"`
	PushName  string `json:"pushname"`
	SenderID  string `json:"sender_id"`
	Timestamp string `json:"timestamp"`
}

// WhatsAppResponse represents the response format for WhatsApp API
type WhatsAppResponse struct {
	Phone       string `json:"phone"`        // Phone number with @s.whatsapp.net
	Message     string `json:"message"`      // Text message to send
	IsForwarded bool   `json:"is_forwarded"` // Whether message is forwarded
}

// WhatsAppAPIResponse represents the API response structure
type WhatsAppAPIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Results struct {
		MessageID string `json:"message_id"`
		Status    string `json:"status"`
	} `json:"results"`
}

// HandleWhatsAppWebhook processes incoming WhatsApp webhook requests
func (b *Bot) HandleWhatsAppWebhook(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("Received webhook request from %s", r.RemoteAddr)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Raw webhook payload: %s", string(body))

	// Parse webhook payload
	var payload WhatsAppWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		log.Printf("Invalid JSON: %s", string(body))
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("Processing message from %s (%s): %s",
		payload.From, payload.PushName, payload.Message.Text)

	// Process message using existing bot logic
	response, err := b.HandleMessage(r.Context(), payload.Message.Text)
	if err != nil {
		log.Printf("Error processing message: %v", err)
		http.Error(w, "Error processing message", http.StatusInternalServerError)
		return
	}

	log.Printf("Got LLM response: %s", response)

	// Send response back to WhatsApp
	if err := b.sendDirectWhatsAppResponse(r.Context(), payload.From, response); err != nil {
		log.Printf("Error sending WhatsApp response: %v", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// sendDirectWhatsAppResponse sends the response directly to WhatsApp
func (b *Bot) sendDirectWhatsAppResponse(ctx context.Context, jid string, message string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"phone":        jid,
		"message":      message,
		"is_forwarded": false,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	log.Printf("Sending to WhatsApp API: %s", string(jsonData))

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:3000/send/message", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Referer", "http://localhost:3000/")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	log.Printf("WhatsApp API response: %s", string(respBody))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// RegisterWebhook registers our webhook URL with the WhatsApp API
func (b *Bot) RegisterWebhook(webhookURL string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare webhook registration payload
	payload := map[string]interface{}{
		"webhook_url": webhookURL,
		"events": []string{
			"message_received",
			"message_sent",
			"message_updated",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling webhook payload: %w", err)
	}

	log.Printf("Registering webhook with WhatsApp API: %s", string(jsonData))

	// Send registration request
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/webhook/register", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error registering webhook: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	log.Printf("Webhook registration response: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook registration failed: %s", string(respBody))
	}

	return nil
}

// StartWebhookServer starts the webhook server
func (b *Bot) StartWebhookServer(port int) error {
	mux := http.NewServeMux()

	// Add handlers with CORS support
	mux.HandleFunc("/webhook/wa", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("Received webhook request from %s", r.RemoteAddr)
		b.HandleWhatsAppWebhook(w, r)
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Webhook server is running"))
	})

	// Create listener on all interfaces
	listener, err := net.Listen("tcp", "0.0.0.0:4444")
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	server := &http.Server{
		Handler: mux,
	}

	log.Printf("Starting webhook server on 0.0.0.0:4444")
	log.Printf("If running WhatsApp API in WSL, configure webhook URL as:")
	log.Printf("1. Try: http://172.0.0.1:4444/webhook/wa")
	log.Printf("2. Or: http://host.docker.internal:4444/webhook/wa")
	log.Printf("3. Or use your Windows IP address: http://<windows-ip>:4444/webhook/wa")

	return server.Serve(listener)
}
