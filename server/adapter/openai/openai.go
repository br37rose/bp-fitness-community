package openai

import (
	"context"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	mongo_client "go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type OpenAIConnector interface {
	Shutdown()
	CreateChatCompletion(prompt string) (string, error)
}

type openAIConnector struct {
	Logger *slog.Logger
	Client *openai.Client
}

func NewOpenAIConnector(cfg *c.Conf, logger *slog.Logger, dbClient *mongo_client.Client) OpenAIConnector {
	logger.Debug("openai initializing...")
	client := openai.NewClient(cfg.AI.APIKey)
	logger.Debug("openai initialized with mongodb as backend")
	return &openAIConnector{
		Logger: logger,
		Client: client,
	}
}

func (impl *openAIConnector) Shutdown() {
	// Do nothing...
}

func (impl *openAIConnector) CreateChatCompletion(prompt string) (string, error) {
	resp, err := impl.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature:      1,
			MaxTokens:        2048,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		},
	)
	if err != nil {
		impl.Logger.Error("completion error", slog.Any("error", err))
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}
