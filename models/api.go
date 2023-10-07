package models

type ChatGPTRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatGPTRequestMessage `json:"messages"`
	Stream   bool                    `json:"stream"`
}

type ChatGPTRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatGPTRequestChoice `json:"choices"`
}

type ChatGPTResponseDelta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequestChoice struct {
	Index        int                  `json:"index"`
	Delta        ChatGPTResponseDelta `json:"delta"`
	FinishReason string               `json:"finish_reason"`
}
