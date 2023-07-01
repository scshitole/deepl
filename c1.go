package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	// Create prompt directory if it doesn't exist
	err := os.MkdirAll("prompt", 0755)
	if err != nil {
		fmt.Println("Error creating prompt directory:", err)
		return
	}

	// Find the next available index for the prompt file
	files, err := ioutil.ReadDir("prompt")
	if err != nil {
		fmt.Println("Error reading prompt directory:", err)
		return
	}
	index := len(files) + 1
	promptFilename := filepath.Join("prompt", "prompt"+strconv.Itoa(index)+".txt")

	// Save prompt to file
	err = ioutil.WriteFile(promptFilename, []byte(prompt), 0644)
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

	// Create pc directory if it doesn't exist
	err = os.MkdirAll("pc", 0755)
	if err != nil {
		fmt.Println("Error creating pc directory:", err)
		return
	}

	// Find the next available index for the pc file
	files, err = ioutil.ReadDir("pc")
	if err != nil {
		fmt.Println("Error reading pc directory:", err)
		return
	}
	index = len(files) + 1
	pcFilename := filepath.Join("pc", "pc"+strconv.Itoa(index)+".json")

	// Remove whitespaces from pcFilename
	pcFilename = strings.ReplaceAll(pcFilename, " ", "")

	// Save prompt and completions to file
	pcFile, err := os.Create(pcFilename)
	if err != nil {
		fmt.Println("Error creating pc file:", err)
		return
	}
	defer pcFile.Close()

	// Write prompt and completions to the pc file
	pcFile.WriteString("Prompt: " + prompt + "\n\n")
	pcFile.WriteString("Completions: " + generatedText)

	// Create compl directory if it doesn't exist
	err = os.MkdirAll("compl", 0755)
	if err != nil {
		fmt.Println("Error creating compl directory:", err)
		return
	}

	// Find the next available index for the compl file
	files, err = ioutil.ReadDir("compl")
	if err != nil {
		fmt.Println("Error reading compl directory:", err)
		return
	}
	index = len(files) + 1
	complFilename := filepath.Join("compl", "compl"+strconv.Itoa(index)+".json")

	// Remove whitespaces from complFilename
	complFilename = strings.ReplaceAll(complFilename, " ", "")

	// Save completions to file
	err = ioutil.WriteFile(complFilename, []byte(generatedText), 0644)
	if err != nil {
		fmt.Println("Error saving completions to file:", err)
		return
	}

	// Print prompt and completions side by side
	fmt.Println("Prompt:      ", prompt)
	fmt.Println("Completions: ", generatedText)

	// Ask user if they want to proceed with updating the GitHub repository
	fmt.Print("Do you want to proceed with updating the GitHub repository? (yes/no): ")
	scanner.Scan()
	confirmation := scanner.Text()

	if confirmation == "yes" {
		// Run git commands: git add ., git commit, and git push
		runGitCommand("add", ".")
		runGitCommand("commit", "-m", "Added prompt and completion files")
		runGitCommand("push")
	}
}

// Helper function to run git commands
func runGitCommand(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running git command:", err)
	}
}

