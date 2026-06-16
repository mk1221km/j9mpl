package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "--print" {
		fmt.Fprintln(os.Stderr, "Usage: llm --print <prompt>")
		os.Exit(1)
	}
	prompt := os.Args[2]

	apiKey := resolveAPIKey()
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] No DEEPSEEK_API_KEY found in env or ~/.env")
		os.Exit(1)
	}

	content, err := fetchCompletion(apiKey, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}
	fmt.Print(content)
}

func resolveAPIKey() string {
	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		return key
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	envFile := filepath.Join(home, ".env")
	f, err := os.Open(envFile)
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "DEEPSEEK_API_KEY=") {
			val := strings.TrimPrefix(line, "DEEPSEEK_API_KEY=")
			val = strings.Trim(val, `"'`)
			return val
		}
	}
	return ""
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "pom.xml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

func logLLMCall(prompt, response string, attempt, maxRetries int, latency time.Duration, status string) {
	projectRoot := findProjectRoot()
	logPath := filepath.Join(projectRoot, ".context", "llm_usage.log")
	
	// Ensure parent directory exists
	os.MkdirAll(filepath.Dir(logPath), 0755)
	
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return // ignore logging errors to avoid breaking the tool
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	var logBlock strings.Builder
	logBlock.WriteString("================================================================================\n")
	logBlock.WriteString(fmt.Sprintf("TIMESTAMP: %s\n", timestamp))
	logBlock.WriteString(fmt.Sprintf("LATENCY: %.2fs\n", latency.Seconds()))
	logBlock.WriteString(fmt.Sprintf("ATTEMPT: %d/%d\n", attempt, maxRetries))
	logBlock.WriteString("MODEL: deepseek-v4-flash\n")
	logBlock.WriteString(fmt.Sprintf("STATUS: %s\n", status))
	logBlock.WriteString("PROMPT:\n")
	logBlock.WriteString("--------------------------------------------------------------------------------\n")
	logBlock.WriteString(prompt)
	logBlock.WriteString("\n--------------------------------------------------------------------------------\n")
	if status == "SUCCESS" {
		logBlock.WriteString("RESPONSE:\n")
		logBlock.WriteString("--------------------------------------------------------------------------------\n")
		logBlock.WriteString(response)
		logBlock.WriteString("\n--------------------------------------------------------------------------------\n")
	}
	logBlock.WriteString("================================================================================\n\n")

	f.WriteString(logBlock.String())
}

func fetchCompletion(apiKey, prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: "deepseek-v4-flash",
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.1,
		MaxTokens:   4096,
	}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	maxRetries := 5
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		startTime := time.Now()
		req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
		if err != nil {
			logLLMCall(prompt, "", attempt, maxRetries, time.Since(startTime), fmt.Sprintf("FAILED (Create request error: %v)", err))
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Set client timeout to 90 seconds for large outputs
		client := &http.Client{Timeout: 90 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			latency := time.Since(startTime)
			logLLMCall(prompt, "", attempt, maxRetries, latency, fmt.Sprintf("FAILED (API request failed: %v)", err))
			fmt.Fprintf(os.Stderr, "[WARNING] API request failed (attempt %d/%d): %v\n", attempt, maxRetries, err)
			if attempt < maxRetries {
				sleepDuration := time.Duration(attempt*5) * time.Second
				time.Sleep(sleepDuration)
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			latency := time.Since(startTime)
			logLLMCall(prompt, "", attempt, maxRetries, latency, fmt.Sprintf("FAILED (Reading response body failed: %v)", err))
			fmt.Fprintf(os.Stderr, "[WARNING] Reading API response failed (attempt %d/%d): %v\n", attempt, maxRetries, err)
			if attempt < maxRetries {
				sleepDuration := time.Duration(attempt*5) * time.Second
				time.Sleep(sleepDuration)
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
			latency := time.Since(startTime)
			logLLMCall(prompt, "", attempt, maxRetries, latency, fmt.Sprintf("FAILED (API error HTTP %d: %s)", resp.StatusCode, string(body)))
			fmt.Fprintf(os.Stderr, "[WARNING] API returned HTTP status %d (attempt %d/%d)\n", resp.StatusCode, attempt, maxRetries)
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusBadRequest {
				return "", lastErr
			}
			if attempt < maxRetries {
				sleepDuration := time.Duration(attempt*5) * time.Second
				time.Sleep(sleepDuration)
			}
			continue
		}

		var chatResp ChatResponse
		if err := json.Unmarshal(body, &chatResp); err != nil {
			latency := time.Since(startTime)
			logLLMCall(prompt, "", attempt, maxRetries, latency, fmt.Sprintf("FAILED (JSON unmarshal failed: %v)", err))
			return "", err
		}
		if len(chatResp.Choices) == 0 {
			latency := time.Since(startTime)
			logLLMCall(prompt, "", attempt, maxRetries, latency, "FAILED (Empty choices in response)")
			return "", fmt.Errorf("empty response from API")
		}

		latency := time.Since(startTime)
		result := chatResp.Choices[0].Message.Content
		logLLMCall(prompt, result, attempt, maxRetries, latency, "SUCCESS")
		return result, nil
	}

	return "", fmt.Errorf("failed after %d attempts: %v", maxRetries, lastErr)
}
