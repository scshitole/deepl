# deepl

# GPT-3.5 Turbo API Example

This is a simple example program that demonstrates how to use the OpenAI GPT-3.5 Turbo API to generate text based on a given prompt. The program sends a request to the OpenAI API with a prompt text and receives a generated text in response.


## Prerequisites

Before running the program, make sure you have the following:

- Go programming language installed
- OpenAI API key (Sign up at https://openai.com to get an API key)

## How it works?
- When you run the code as ```run go c1.go``` it will prompt you for ChatGPT prompt ( background it will connect to GPT-3.5)	
- Once you enter the prompt it will display the completions or reponse on the screen
- It will record the ***prompt*** and ***completion*** and put the details under directory ```pc```
- It will ask you if you want to upload the details github repo

## Getting Started

1. Clone the repository:

```bash
git clone <repository-url>

Change to the project directory:
cd deepl

Set the OpenAI API key as an environment variable:
export OPENAI_API_KEY=<your-api-key>

Create a file named prompt.txt and enter your desired prompt text.

Run the program:

go run c1.go
