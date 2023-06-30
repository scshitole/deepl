package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

	// Read prompt from keyboard input
	fmt.Println("Enter the prompt:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	prompt := scanner.Text()

	// Trim any trailing whitespace and % characters
	prompt = strings.TrimSpace(prompt)
	prompt = strings.TrimRight(prompt, "%")

	// Save prompt to external file
	err := ioutil.WriteFile("prompt.txt", []byte(prompt), 0644)
	if err != nil {
		fmt.Println("Error saving prompt to file:", err)
		return
	}

	// Request payload
	payload := map[string]interface{}{
		"prompt":     prompt,
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

	// Create the desired JSON structure
	declaration := map[string]interface{}{
		"class":           "ADC",
		"schemaVersion":   "3.14.0",
		"id":              "adcDeclaration_12345",
		"remarks":         "Example F5 BIG-IP Application Services 3 (AS3) declaration for an HTTP application.",
		"label":           "httpApp",
		"templateVersion": "3.14.0",
		"template": map[string]interface{}{
			"class":         "Template",
			"httpProfiles": []map[string]interface{}{
				{
					"class":       "HTTPProfile",
					"name":        "http_profile_12345",
					"idleTimeout": 3600,
				},
			},
			"httpServices": []map[string]interface{}{
				{
					"class":             "HTTPService",
					"virtualAddresses":  []string{"1.1.1.1"},
					"virtualPort":       80,
					"profiles": []map[string]interface{}{
						{
							"class":   "ProfileContext",
							"profile": "http_profile_12345",
						},
					},
					"name": "http_service_12345",
					"persistenceMethods": []map[string]interface{}{
						{
							"class":        "ClientSourceAddress",
							"matchForward": "enabled",
						},
					},
				},
			},
		},
	}

	// Create the final response JSON structure
	output := map[string]interface{}{
		"class":       "AS3",
		"action":      "deploy",
		"declaration": declaration,
		"persist":     true,
	}

	// Convert output to JSON
	outputBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling output:", err)
		return
	}

	// Remove trailing "%" symbol from the generated text
	generatedText := string(outputBytes)
	generatedText = strings.TrimRight(generatedText, "%")

	// Save output to response.json file
	err = ioutil.WriteFile("response.json", []byte(generatedText), 0644)
	if err != nil {
		fmt.Println("Error saving output to file:", err)
		return
	}

	// Print the generated text
	fmt.Println("Generated text:\n", generatedText)
}

