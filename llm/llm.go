package llm

import (
	"os"
	"sync"
)

type LLMProvider struct {
	Url      string
	LLM      string
	ApiToken string
}

var (
	mu         sync.Mutex
	CurrentLLM = &LLMProvider{
		Url:      os.Getenv("OPENAI_API_URL"),
		LLM:      os.Getenv("GPT_TURBO"),
		ApiToken: os.Getenv("OPENAI_API_TOKEN"),
	}
)

func (p *LLMProvider) SetModel(url string, name string, apiToken string) {
	mu.Lock()
	defer mu.Unlock()
	p.Url = url
	p.LLM = name
	p.ApiToken = apiToken
}

func GetCurrentLLM() *LLMProvider {
	mu.Lock()
	defer mu.Unlock()
	return CurrentLLM
}
