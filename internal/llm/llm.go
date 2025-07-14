package llm

var baseUrl, apiKey, modelName string

func SetLLMInfo(llmBaseUrl, llmApiKey, llmModel string) {
	baseUrl, apiKey, modelName = llmBaseUrl, llmApiKey, llmModel
}
