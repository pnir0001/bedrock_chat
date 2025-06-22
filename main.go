package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	region  = "us-east-1"
	modelId = "anthropic.claude-3-haiku-20240307-v1:0" // 利用モデルIDは適宜変更
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ClaudeMessage struct {
	Role    string               `json:"role"`
	Content []ClaudeContentBlock `json:"content"`
}

type ClaudeContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	Messages         []ClaudeMessage `json:"messages"`
	MaxTokens        int             `json:"max_tokens"`
	Temperature      float64         `json:"temperature"`
	TopP             float64         `json:"top_p"`
	StopSequences    []string        `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}
	client := bedrockruntime.NewFromConfig(cfg)

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Claude 3形式でリクエストを作成
		input := ClaudeRequest{
			AnthropicVersion: "bedrock-2023-05-31",
			Messages: []ClaudeMessage{
				{
					Role: "user",
					Content: []ClaudeContentBlock{
						{Type: "text", Text: req.Message},
					},
				},
			},
			MaxTokens:   1024,
			Temperature: 0.7,
			TopP:        0.9,
		}

		payload, err := json.Marshal(input)
		if err != nil {
			http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
			return
		}

		out, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(modelId),
			ContentType: aws.String("application/json"),
			Body:        payload,
		})
		if err != nil {
			log.Println("Bedrock API error:", err)
			http.Error(w, "Bedrock API error", http.StatusInternalServerError)
			return
		}

		// Claude 3のレスポンス形式に応じてパース
		var resp ClaudeResponse
		if err := json.Unmarshal(out.Body, &resp); err != nil {
			http.Error(w, "Failed to parse Bedrock response", http.StatusInternalServerError)
			return
		}

		// 応答テキストを返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"response": resp.Content[0].Text,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
