package modelconfig

type ModelConfig struct {
	OllamaHost   string `yaml:"ollama_host"`
	Name         string `yaml:"name"`
	SystemPrompt string `yaml:"system_prompt"`
	UserPrompt   string `yaml:"user_prompt"`
}
