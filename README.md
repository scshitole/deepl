# deepl
# GPT-3.5 Turbo API Example

This is a simple example program that demonstrates how to use the OpenAI GPT-3.5 Turbo API to generate text based on a given prompt.

## Prerequisites

Before running the program, make sure you have the following:

- Go programming language installed
- OpenAI API key (Sign up at https://openai.com to get an API key)

## Installation

-  Clone the repository:

```bash
git clone <repository-url>

-  Change to the project directory:
cd deepl

-  Set the OpenAI API key as an environment variable:
export OPENAI_API_KEY=<your-api-key>

- Create a file named prompt.txt and enter your desired prompt text.

- Run the program:

go run main.go


# GPT-3.5 Turbo API Example

This is a simple example program that demonstrates how to use the OpenAI GPT-3.5 Turbo API to generate text based on a given prompt. The program sends a request to the OpenAI API with a prompt text and receives a generated text in response.

## How It Works

1. The program reads the prompt text from an external file named `prompt.txt`. You can provide your desired prompt in the `prompt.txt` file.

2. It sends a request to the OpenAI API using the GPT-3.5 Turbo model, providing the prompt text and a maximum number of tokens for the generated text.

3. The API returns a response with the generated text based on the provided prompt.

4. The program extracts the generated text from the response and prints it to the console.

## Prerequisites

Before running the program, make sure you have the following:

- Go programming language installed
- OpenAI API key (Sign up at https://openai.com to get an API key)

## Getting Started

1. Clone the repository:

```bash
git clone <repository-url>
