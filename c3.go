package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const pcFilename = "pc.json" // Name of the output JSON file

type PromptCompletion struct {
	Prompt      string `json:"prompt"`
	Completions string `json:"completions"`
}

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

	// Create prompt-completion object
	pc := PromptCompletion{
		Prompt:      prompt,
		Completions: generatedText,
	}

	// Convert prompt-completion object to JSON with indentation
	pcBytes, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling prompt-completion:", err)
		return
	}

	// Create pc.json file if it doesn't exist
	fileExists, err := fileExists(pcFilename)
	if err != nil {
		fmt.Println("Error checking if pc.json file exists:", err)
		return
	}
	if !fileExists {
		err = ioutil.WriteFile(pcFilename, pcBytes, 0644)
		if err != nil {
			fmt.Println("Error creating pc.json file:", err)
			return
		}
	} else {
		// Append prompt-completion JSON to pc.json file with spacing
		err = appendToFile(pcFilename, ",\n\n"+string(pcBytes))
		if err != nil {
			fmt.Println("Error appending to pc.json file:", err)
			return
		}
	}

	// Print prompt and completions side by side
	fmt.Println("Prompt:      ", prompt)
	fmt.Println("Completions: ", generatedText)
}

// Check if a file exists
func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// Append data to a file
func appendToFile(filename, data string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		return err
	}
	return nil
}

