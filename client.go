// Package ultravox provides a client for interacting with the Ultravox AI voice API.
package ultravox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Constants for default configuration values
const (
	DefaultAPIBaseURL       = "https://api.ultravox.ai/api"
	DefaultModel            = "fixie-ai/ultravox"
	DefaultVoice            = "Mark"
	DefaultInputSampleRate  = 8000
	DefaultOutputSampleRate = 8000
	DefaultTimeout          = 15 * time.Second
	DefaultSystemPrompt     = "You are a helpful AI assistant that provides clear and concise information."
)

// OutputMediumType defines the type of output medium
type OutputMediumType string

// Predefined output medium constants
const (
	OutputMediumVoice OutputMediumType = "MESSAGE_MEDIUM_VOICE"
	OutputMediumText  OutputMediumType = "MESSAGE_MEDIUM_TEXT"
)

// FirstSpeakerType defines who speaks first in a conversation
type FirstSpeakerType string

// Predefined first speaker constants
const (
	FirstSpeakerAgent FirstSpeakerType = "FIRST_SPEAKER_AGENT"
	FirstSpeakerUser  FirstSpeakerType = "FIRST_SPEAKER_USER"
)

// EndBehaviorType defines behaviors after a message is spoken
type EndBehaviorType string

// Predefined end behavior constants
const (
	EndBehaviorDefault    EndBehaviorType = "END_BEHAVIOR_UNSPECIFIED"
	EndBehaviorHangUpSoft EndBehaviorType = "END_BEHAVIOR_HANG_UP_SOFT"
	EndBehaviorHangUpHard EndBehaviorType = "END_BEHAVIOR_HANG_UP_STRICT"
)

// Config holds the client configuration
type Config struct {
	CallRequest
	APIKey      string
	APIBaseURL  string
	HTTPTimeout time.Duration
}

// Option is a function that modifies the client configuration
type Option func(*Config)

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithAPIBaseURL sets the base URL for API requests
func WithAPIBaseURL(url string) Option {
	return func(c *Config) {
		c.APIBaseURL = url
	}
}

// WithSystemPrompt sets the system prompt for the agent
func WithSystemPrompt(prompt string) Option {
	return func(c *Config) {
		c.SystemPrompt = prompt
	}
}

// WithTemperature sets the temperature for model generation
func WithTemperature(temperature float64) Option {
	return func(c *Config) {
		c.Temperature = temperature
	}
}

// WithModel sets the model for the agent
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

// WithVoice sets the voice for the agent
func WithVoice(voice string) Option {
	return func(c *Config) {
		c.Voice = voice
	}
}

// WithExternalVoice sets an external voice provider for the agent
func WithExternalVoice(voice *ExternalVoice) Option {
	return func(c *Config) {
		c.ExternalVoice = voice
	}
}

// WithLanguageHint sets a language hint to guide speech recognition
func WithLanguageHint(languageHint string) Option {
	return func(c *Config) {
		c.LanguageHint = languageHint
	}
}

// WithFirstSpeaker sets who speaks first in the conversation
func WithFirstSpeaker(speaker FirstSpeakerType) Option {
	return func(c *Config) {
		c.FirstSpeaker = speaker
	}
}

// WithFirstSpeakerSettings sets detailed configuration for who speaks first
func WithFirstSpeakerSettings(settings *FirstSpeakerSettings) Option {
	return func(c *Config) {
		c.FirstSpeakerSettings = settings
	}
}

// WithHTTPTimeout sets the timeout for HTTP requests
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.HTTPTimeout = timeout
	}
}

// WithJoinTimeout sets the join timeout for the client configuration
func WithJoinTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.JoinTimeout = UltravoxDuration(timeout)
	}
}

// WithMaxDuration sets the maximum duration for the client configuration
func WithMaxDuration(duration time.Duration) Option {
	return func(c *Config) {
		c.MaxDuration = UltravoxDuration(duration)
	}
}

// WithInitialOutputMedium sets the initial output medium (voice or text)
func WithInitialOutputMedium(medium OutputMediumType) Option {
	return func(c *Config) {
		c.InitialOutputMedium = medium
	}
}

// WithVadSettings sets voice activity detection settings
func WithVadSettings(settings *VadSettings) Option {
	return func(c *Config) {
		c.VadSettings = settings
	}
}

// WithDataConnection sets data connection configuration
func WithDataConnection(config *DataConnectionConfig) Option {
	return func(c *Config) {
		c.DataConnection = config
	}
}

// WithRecordingEnabled sets whether call recording is enabled
func WithRecordingEnabled(enabled bool) Option {
	return func(c *Config) {
		c.RecordingEnabled = enabled
	}
}

// HTTPClient defines the interface for making HTTP requests
// This makes testing easier by allowing mock implementations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client handles communication with the Ultravox API
type Client struct {
	config Config
	http   HTTPClient
}

// NewClient creates a new Ultravox client with the provided options
func NewClient(opts ...Option) *Client {
	// Set default configuration
	config := Config{
		HTTPTimeout: DefaultTimeout,
		APIBaseURL:  DefaultAPIBaseURL,
		APIKey:      os.Getenv("ULTRAVOX_API_KEY"),
		CallRequest: CallRequest{
			Model:               DefaultModel,
			Voice:               DefaultVoice,
			FirstSpeaker:        FirstSpeakerAgent,
			SystemPrompt:        DefaultSystemPrompt,
			JoinTimeout:         UltravoxDuration(30 * time.Second),
			MaxDuration:         UltravoxDuration(10 * time.Minute),
			Temperature:         0.0,
			InitialOutputMedium: OutputMediumVoice,
			RecordingEnabled:    false,
			Medium: &CallMedium{
				ServerWebSocket: &WebSocketMedium{
					InputSampleRate:  DefaultInputSampleRate,
					OutputSampleRate: DefaultOutputSampleRate,
				},
			},
		},
	}

	// Apply provided options
	for _, opt := range opts {
		opt(&config)
	}

	return &Client{
		config: config,
		http:   &http.Client{Timeout: config.HTTPTimeout},
	}
}

// WithHTTPClient sets a custom HTTP client
func (c *Client) WithHTTPClient(httpClient HTTPClient) {
	c.http = httpClient
}

// Call initiates a new call with the Ultravox API
// Optional CallOption parameters can be provided to override default configuration for this specific call
func (c *Client) Call(ctx context.Context, opts ...CallOption) (*Call, error) {
	// Start with default configuration from client
	request := CallRequest{
		SystemPrompt:         c.config.SystemPrompt,
		Temperature:          c.config.Temperature,
		Model:                c.config.Model,
		Voice:                c.config.Voice,
		ExternalVoice:        c.config.ExternalVoice,
		LanguageHint:         c.config.LanguageHint,
		MaxDuration:          c.config.MaxDuration,
		JoinTimeout:          c.config.JoinTimeout,
		FirstSpeaker:         c.config.FirstSpeaker,
		FirstSpeakerSettings: c.config.FirstSpeakerSettings,
		InitialOutputMedium:  c.config.InitialOutputMedium,
		VadSettings:          c.config.VadSettings,
		RecordingEnabled:     c.config.RecordingEnabled,
		DataConnection:       c.config.DataConnection,
		Medium:               c.config.Medium,
		TemplateContext:      c.config.TemplateContext,
	}

	// Apply any call-specific options
	for _, opt := range opts {
		opt(&request)
	}

	// Validate required configuration
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Build the URL with query parameters if needed
	url := c.buildCallURL(&request)
	// api/agents/${AGENT_ID}/calls
	// Add query parameters if specified
	hasParams := false
	if request.EnableGreetingPrompt {
		if !hasParams {
			url += "?enableGreetingPrompt=true"
			hasParams = true
		} else {
			url += "&enableGreetingPrompt=true"
		}
	}

	if request.PriorCallId != "" {
		if !hasParams {
			url += "?priorCallId=" + request.PriorCallId
			hasParams = true
		} else {
			url += "&priorCallId=" + request.PriorCallId
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned non-success status: %d", resp.StatusCode)
	}

	var callResp Call
	if err := json.NewDecoder(resp.Body).Decode(&callResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if callResp.JoinURL == "" {
		return nil, fmt.Errorf("API did not return a valid join URL")
	}

	return &callResp, nil
}

// CallAgent initiates a call to a specific agent using the Ultravox API.
// This method is designed to interact with a specific agent endpoint, allowing
// for customized interactions based on the agent's configuration and context.
func (c *Client) CallAgent(ctx context.Context, agentID string, opts ...CallOption) (*Call, error) {
	opts = append(opts, WithCallAgentID(agentID))
	return c.Call(ctx, opts...)
}

// buildCallURL returns the appropriate API endpoint for creating a call.
// If the request includes an AgentID, it targets the agent-scoped endpoint:
//
//	/api/agents/{agentId}/calls
//
// Otherwise, it uses the default endpoint:
//
//	/api/calls
func (c *Client) buildCallURL(req *CallRequest) string {
	if req.AgentID != "" {
		return fmt.Sprintf("%s/agents/%s/calls", c.config.APIBaseURL, req.AgentID)
	}
	return fmt.Sprintf("%s/calls", c.config.APIBaseURL)
}
