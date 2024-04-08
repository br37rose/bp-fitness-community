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
	CreateFitnessPlan(ctx context.Context, prompt, availableExcercises string) (openai.Run, error)
	GetRunner(ctx context.Context, threadID, runnerID string) (openai.Run, error)
	SubmitToolOutputs(ctx context.Context, threadID, runnerID string, request openai.SubmitToolOutputsRequest) (response openai.Run, err error)
}

type openAIConnector struct {
	Logger             *slog.Logger
	Client             *openai.Client
	FitnessAssistantID string
}

func NewOpenAIConnector(cfg *c.Conf, logger *slog.Logger, dbClient *mongo_client.Client) OpenAIConnector {
	logger.Debug("openai initializing...")
	client := openai.NewClient(cfg.AI.APIKey)
	logger.Debug("openai initialized with mongodb as backend")
	return &openAIConnector{
		Logger:             logger,
		Client:             client,
		FitnessAssistantID: cfg.AI.FitnessPlanAssistantID,
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

func (impl *openAIConnector) CreateFitnessPlan(ctx context.Context, prompt, availableExcercises string) (openai.Run, error) {

	runnerResp, err := impl.Client.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{Role: openai.ThreadMessageRoleUser, Content: prompt},
				{Role: openai.ThreadMessageRoleUser, Content: "Available exercises: " + availableExcercises},
			},
		},
		RunRequest: openai.RunRequest{
			AssistantID: impl.FitnessAssistantID,
		},
	})

	return runnerResp, err
}

func (impl *openAIConnector) GetRunner(ctx context.Context, threadID, runnerID string) (openai.Run, error) {
	return impl.Client.RetrieveRun(ctx, threadID, runnerID)
}

func (impl *openAIConnector) SubmitToolOutputs(ctx context.Context, threadID, runnerID string, request openai.SubmitToolOutputsRequest) (response openai.Run, err error) {
	return impl.Client.SubmitToolOutputs(ctx, threadID, runnerID, request)
}
