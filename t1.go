package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// API endpoint URL
	url := "https://api.openai.com/v1/engines/text-davinci-003/completions"

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable.")
		return
	}

	// Read prompt from external file
	prompt, err := ioutil.ReadFile("prompt.txt")
	if err != nil {
		fmt.Println("Error reading prompt file:", err)
		return
	}

	// Request payload
	payload := map[string]interface{}{
		"prompt":     string(prompt),
		"max_tokens": 3000,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Check if the response is an error
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response:", resp.Status)
		fmt.Println(string(respBody))
		return
	}

	// Parse response JSON
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Extract generated text from response
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		fmt.Println("Error extracting generated text from response")
		return
	}
	generatedText, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		fmt.Println("Error extracting generated text from response")
		return
	}

	fmt.Println("Generated text:", generatedText)
}

