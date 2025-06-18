package ultravox

import "time"

type TemplateContext struct {
	UserFirstname      string `json:"userFirstname,omitempty"`
	LastCallTranscript string `json:"lastCallTranscript,omitempty"`
}

// CallRequest represents the request structure for initiating a call
type CallRequest struct {
	// Basic properties
	SystemPrompt        string           `json:"systemPrompt,omitempty"`
	Temperature         float64          `json:"temperature,omitempty"`
	Model               string           `json:"model,omitempty"`
	Voice               string           `json:"voice,omitempty"`
	ExternalVoice       *ExternalVoice   `json:"externalVoice,omitempty"`
	LanguageHint        string           `json:"languageHint,omitempty"`
	InitialMessages     []Message        `json:"initialMessages,omitempty"`
	JoinTimeout         UltravoxDuration `json:"joinTimeout,omitempty"`
	MaxDuration         UltravoxDuration `json:"maxDuration,omitempty"`
	TimeExceededMessage string           `json:"timeExceededMessage,omitempty"`
	InactivityMessages  []TimedMessage   `json:"inactivityMessages,omitempty"`

	// Tool configuration
	SelectedTools []SelectedTool `json:"selectedTools,omitempty"`

	// Medium configuration
	Medium           *CallMedium `json:"medium,omitempty"`
	RecordingEnabled bool        `json:"recordingEnabled,omitempty"`

	// First speaker configuration
	FirstSpeaker         FirstSpeakerType      `json:"firstSpeaker,omitempty"` // Deprecated
	InitialOutputMedium  OutputMediumType      `json:"initialOutputMedium,omitempty"`
	FirstSpeakerSettings *FirstSpeakerSettings `json:"firstSpeakerSettings,omitempty"`

	// Advanced settings
	VadSettings          *VadSettings          `json:"vadSettings,omitempty"`
	ExperimentalSettings interface{}           `json:"experimentalSettings,omitempty"`
	Metadata             map[string]string     `json:"metadata,omitempty"`
	InitialState         interface{}           `json:"initialState,omitempty"`
	DataConnection       *DataConnectionConfig `json:"dataConnection,omitempty"`

	// For creating a call from a prior call
	PriorCallId          string `json:"priorCallId,omitempty"`
	EnableGreetingPrompt bool   `json:"enableGreetingPrompt,omitempty"`

	// For Agent Calls
	AgentID         string           `json:"-"`
	TemplateContext *TemplateContext `json:"templateContext,omitempty"`
}

// Call contains the response from a call creation request
type Call struct {
	CallID               string                `json:"callId"`
	ClientVersion        string                `json:"clientVersion,omitempty"`
	JoinURL              string                `json:"joinUrl"`
	Created              string                `json:"created"`
	Joined               string                `json:"joined,omitempty"`
	Ended                string                `json:"ended,omitempty"`
	EndReason            string                `json:"endReason,omitempty"`
	MaxDuration          UltravoxDuration      `json:"maxDuration"`
	JoinTimeout          UltravoxDuration      `json:"joinTimeout"`
	FirstSpeaker         FirstSpeakerType      `json:"firstSpeaker,omitempty"`
	FirstSpeakerSettings *FirstSpeakerSettings `json:"firstSpeakerSettings,omitempty"`
	InitialOutputMedium  OutputMediumType      `json:"initialOutputMedium,omitempty"`
	Medium               *CallMedium           `json:"medium,omitempty"`
	RecordingEnabled     bool                  `json:"recordingEnabled"`
	ErrorCount           int                   `json:"errorCount"`
	ShortSummary         string                `json:"shortSummary,omitempty"`
	Summary              string                `json:"summary,omitempty"`
}

// CallOption defines a function that modifies a call request
type CallOption func(*CallRequest)

// WithCallJoinTimeout overrides the join timeout for a specific call
func WithCallJoinTimeout(timeout time.Duration) CallOption {
	return func(r *CallRequest) {
		r.JoinTimeout = UltravoxDuration(timeout)
	}
}

// WithCallMaxDuration overrides the maximum duration for a specific call
func WithCallMaxDuration(duration time.Duration) CallOption {
	return func(r *CallRequest) {
		r.MaxDuration = UltravoxDuration(duration)
	}
}

// WithCallSystemPrompt overrides the system prompt for a specific call
func WithCallSystemPrompt(prompt string) CallOption {
	return func(r *CallRequest) {
		r.SystemPrompt = prompt
	}
}

// WithCallTemperature overrides the temperature for a specific call
func WithCallTemperature(temperature float64) CallOption {
	return func(r *CallRequest) {
		r.Temperature = temperature
	}
}

// WithCallModel overrides the model for a specific call
func WithCallModel(model string) CallOption {
	return func(r *CallRequest) {
		r.Model = model
	}
}

// WithCallVoice overrides the voice for a specific call
func WithCallVoice(voice string) CallOption {
	return func(r *CallRequest) {
		r.Voice = voice
	}
}

// WithCallExternalVoice overrides the external voice for a specific call
func WithCallExternalVoice(voice *ExternalVoice) CallOption {
	return func(r *CallRequest) {
		r.ExternalVoice = voice
	}
}

// WithCallFirstSpeaker overrides who speaks first for a specific call
// Deprecated: Use WithCallFirstSpeakerSettings instead
func WithCallFirstSpeaker(speaker FirstSpeakerType) CallOption {
	return func(r *CallRequest) {
		r.FirstSpeaker = speaker
	}
}

// WithCallFirstSpeakerSettings sets detailed configuration for who speaks first
func WithCallFirstSpeakerSettings(settings *FirstSpeakerSettings) CallOption {
	return func(r *CallRequest) {
		r.FirstSpeakerSettings = settings
	}
}

// WithCallMedium overrides the medium configuration for a specific call
func WithCallMedium(medium *CallMedium) CallOption {
	return func(r *CallRequest) {
		r.Medium = medium
	}
}

// WithCallWebSocketMedium configures the call to use WebSocket with specified sample rates
func WithCallWebSocketMedium(inputRate, outputRate int) CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			ServerWebSocket: &WebSocketMedium{
				InputSampleRate:  inputRate,
				OutputSampleRate: outputRate,
			},
		}
	}
}

// WithCallWebRTCMedium configures the call to use WebRTC
func WithCallWebRTCMedium() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			WebRTC: &WebRTCMedium{},
		}
	}
}

// WithCallTwilioMedium configures the call to use Twilio
func WithCallTwilioMedium() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			Twilio: &TwilioMedium{},
		}
	}
}

// WithCallTelnyxMedium configures the call to use Telnyx
func WithCallTelnyxMedium() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			Telnyx: &TelnyxMedium{},
		}
	}
}

// WithCallPlivoMedium configures the call to use Plivo
func WithCallPlivoMedium() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			Plivo: &PlivoMedium{},
		}
	}
}

// WithCallExotelMedium configures the call to use Exotel
func WithCallExotelMedium() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			Exotel: &ExotelMedium{},
		}
	}
}

// WithCallSIPOutgoing configures the call to use outgoing SIP
func WithCallSIPOutgoing(to, from, username, password string) CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			SIP: &SIPMedium{
				Outgoing: &SIPOutgoing{
					To:       to,
					From:     from,
					Username: username,
					Password: password,
				},
			},
		}
	}
}

// WithCallSIPIncoming configures the call to use incoming SIP
func WithCallSIPIncoming() CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			SIP: &SIPMedium{
				Incoming: &SIPIncoming{},
			},
		}
	}
}

// WithCallLanguageHint sets a language hint for a specific call
func WithCallLanguageHint(languageHint string) CallOption {
	return func(r *CallRequest) {
		r.LanguageHint = languageHint
	}
}

// WithCallInitialMessages sets initial messages for a specific call
func WithCallInitialMessages(messages []Message) CallOption {
	return func(r *CallRequest) {
		r.InitialMessages = messages
	}
}

// WithCallTimeExceededMessage sets a message to be spoken when time is exceeded
func WithCallTimeExceededMessage(message string) CallOption {
	return func(r *CallRequest) {
		r.TimeExceededMessage = message
	}
}

// WithCallInactivityMessages sets messages to be spoken during inactivity
func WithCallInactivityMessages(messages []TimedMessage) CallOption {
	return func(r *CallRequest) {
		r.InactivityMessages = messages
	}
}

// WithCallRecordingEnabled sets whether recording is enabled for a specific call
func WithCallRecordingEnabled(enabled bool) CallOption {
	return func(r *CallRequest) {
		r.RecordingEnabled = enabled
	}
}

// WithCallInitialOutputMedium sets the initial output medium for a specific call
func WithCallInitialOutputMedium(medium OutputMediumType) CallOption {
	return func(r *CallRequest) {
		r.InitialOutputMedium = medium
	}
}

// WithCallVadSettings sets voice activity detection settings for a specific call
func WithCallVadSettings(settings *VadSettings) CallOption {
	return func(r *CallRequest) {
		r.VadSettings = settings
	}
}

// WithCallExperimentalSettings sets experimental settings for a specific call
func WithCallExperimentalSettings(settings interface{}) CallOption {
	return func(r *CallRequest) {
		r.ExperimentalSettings = settings
	}
}

// WithCallMetadata sets metadata for a specific call
func WithCallMetadata(metadata map[string]string) CallOption {
	return func(r *CallRequest) {
		r.Metadata = metadata
	}
}

// WithCallInitialState sets the initial state for a specific call
func WithCallInitialState(state interface{}) CallOption {
	return func(r *CallRequest) {
		r.InitialState = state
	}
}

// WithCallDataConnection sets the data connection for a specific call
func WithCallDataConnection(config *DataConnectionConfig) CallOption {
	return func(r *CallRequest) {
		r.DataConnection = config
	}
}

// WithCallPriorCallId sets the prior call ID for a specific call
func WithCallPriorCallId(callId string) CallOption {
	return func(r *CallRequest) {
		r.PriorCallId = callId
	}
}

// WithCallEnableGreetingPrompt sets whether to enable the greeting prompt
func WithCallEnableGreetingPrompt(enable bool) CallOption {
	return func(r *CallRequest) {
		r.EnableGreetingPrompt = enable
	}
}

// WithTemplateContext sets the entire TemplateContext for the call
func WithTemplateContext(ctx *TemplateContext) CallOption {
	return func(r *CallRequest) {
		r.TemplateContext = ctx
	}
}

// WithTemplateUserFirstname sets the UserFirstname in the TemplateContext
func WithTemplateUserFirstname(firstname string) CallOption {
	return func(r *CallRequest) {
		if r.TemplateContext == nil {
			r.TemplateContext = &TemplateContext{}
		}
		r.TemplateContext.UserFirstname = firstname
	}
}

// WithTemplateLastCallTranscript sets the LastCallTranscript in the TemplateContext
func WithTemplateLastCallTranscript(transcript string) CallOption {
	return func(r *CallRequest) {
		if r.TemplateContext == nil {
			r.TemplateContext = &TemplateContext{}
		}
		r.TemplateContext.LastCallTranscript = transcript
	}
}

// WithCallAgentID sets the AgentID for a specific call
func WithCallAgentID(agentID string) CallOption {
	return func(r *CallRequest) {
		r.AgentID = agentID
	}
}

// Tool-related call options
func WithCallToolByID(toolID string) CallOption {
	return func(r *CallRequest) {
		if r.SelectedTools == nil {
			r.SelectedTools = []SelectedTool{}
		}
		r.SelectedTools = append(r.SelectedTools, SelectedTool{
			ToolID: toolID,
		})
	}
}

func WithCallToolByName(toolName string) CallOption {
	return func(r *CallRequest) {
		if r.SelectedTools == nil {
			r.SelectedTools = []SelectedTool{}
		}
		r.SelectedTools = append(r.SelectedTools, SelectedTool{
			ToolName: toolName,
		})
	}
}

func WithCallTemporaryTool(tool *BaseToolDefinition) CallOption {
	return func(r *CallRequest) {
		if r.SelectedTools == nil {
			r.SelectedTools = []SelectedTool{}
		}
		r.SelectedTools = append(r.SelectedTools, SelectedTool{
			TemporaryTool: tool,
		})
	}
}

// Medium-specific call options with additional configuration
func WithCallWebSocketMediumBuffered(inputRate, outputRate, bufferSizeMs int) CallOption {
	return func(r *CallRequest) {
		r.Medium = &CallMedium{
			ServerWebSocket: &WebSocketMedium{
				InputSampleRate:    inputRate,
				OutputSampleRate:   outputRate,
				ClientBufferSizeMs: bufferSizeMs,
			},
		}
	}
}

// Voice configuration options with advanced settings
func WithCallElevenLabsVoice(voiceID string, options *ElevenLabsVoiceOptions) CallOption {
	return func(r *CallRequest) {
		voice := &ElevenLabsVoice{
			VoiceID: voiceID,
		}
		if options != nil {
			voice.Model = options.Model
			voice.Speed = options.Speed
			voice.UseSpeakerBoost = options.UseSpeakerBoost
			voice.Style = options.Style
			voice.SimilarityBoost = options.SimilarityBoost
			voice.Stability = options.Stability
			voice.OptimizeStreamingLatency = options.OptimizeStreamingLatency
			voice.MaxSampleRate = options.MaxSampleRate
		}
		r.ExternalVoice = &ExternalVoice{ElevenLabs: voice}
	}
}

func WithCallCartesiaVoice(voiceID string, options *CartesiaVoiceOptions) CallOption {
	return func(r *CallRequest) {
		voice := &CartesiaVoice{
			VoiceID: voiceID,
		}
		if options != nil {
			voice.Model = options.Model
			voice.Speed = options.Speed
			voice.Emotion = options.Emotion
			voice.Emotions = options.Emotions
		}
		r.ExternalVoice = &ExternalVoice{Cartesia: voice}
	}
}

func WithCallPlayHtVoice(userID, voiceID string, options *PlayHtVoiceOptions) CallOption {
	return func(r *CallRequest) {
		voice := &PlayHtVoice{
			UserID:  userID,
			VoiceID: voiceID,
		}
		if options != nil {
			voice.Model = options.Model
			voice.Speed = options.Speed
			voice.Quality = options.Quality
			voice.Temperature = options.Temperature
			voice.Emotion = options.Emotion
			voice.VoiceGuidance = options.VoiceGuidance
			voice.StyleGuidance = options.StyleGuidance
			voice.TextGuidance = options.TextGuidance
			voice.VoiceConditioningSeconds = options.VoiceConditioningSeconds
		}
		r.ExternalVoice = &ExternalVoice{PlayHt: voice}
	}
}

func WithCallLmntVoice(voiceID string, options *LmntVoiceOptions) CallOption {
	return func(r *CallRequest) {
		voice := &LmntVoice{
			VoiceID: voiceID,
		}
		if options != nil {
			voice.Model = options.Model
			voice.Speed = options.Speed
			voice.Conversational = options.Conversational
		}
		r.ExternalVoice = &ExternalVoice{Lmnt: voice}
	}
}

// Voice options structures for advanced configuration
type ElevenLabsVoiceOptions struct {
	Model                    string  `json:"model,omitempty"`
	Speed                    float64 `json:"speed,omitempty"`
	UseSpeakerBoost          bool    `json:"useSpeakerBoost,omitempty"`
	Style                    float64 `json:"style,omitempty"`
	SimilarityBoost          float64 `json:"similarityBoost,omitempty"`
	Stability                float64 `json:"stability,omitempty"`
	OptimizeStreamingLatency int     `json:"optimizeStreamingLatency,omitempty"`
	MaxSampleRate            int     `json:"maxSampleRate,omitempty"`
}

type CartesiaVoiceOptions struct {
	Model    string   `json:"model,omitempty"`
	Speed    float64  `json:"speed,omitempty"`
	Emotion  string   `json:"emotion,omitempty"`
	Emotions []string `json:"emotions,omitempty"`
}

type PlayHtVoiceOptions struct {
	Model                    string  `json:"model,omitempty"`
	Speed                    float64 `json:"speed,omitempty"`
	Quality                  string  `json:"quality,omitempty"`
	Temperature              float64 `json:"temperature,omitempty"`
	Emotion                  float64 `json:"emotion,omitempty"`
	VoiceGuidance            float64 `json:"voiceGuidance,omitempty"`
	StyleGuidance            float64 `json:"styleGuidance,omitempty"`
	TextGuidance             float64 `json:"textGuidance,omitempty"`
	VoiceConditioningSeconds float64 `json:"voiceConditioningSeconds,omitempty"`
}

type LmntVoiceOptions struct {
	Model          string  `json:"model,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
	Conversational bool    `json:"conversational,omitempty"`
}

// Advanced VAD configuration
func WithCallAdvancedVadSettings(turnEndpoint, minTurn, minInterruption time.Duration, threshold float64) CallOption {
	return func(r *CallRequest) {
		r.VadSettings = &VadSettings{
			TurnEndpointDelay:           UltravoxDuration(turnEndpoint),
			MinimumTurnDuration:         UltravoxDuration(minTurn),
			MinimumInterruptionDuration: UltravoxDuration(minInterruption),
			FrameActivationThreshold:    threshold,
		}
	}
}

// Message creation helpers
func NewUserMessage(text string, medium OutputMediumType) Message {
	return Message{
		Role:   string(MessageRoleUser),
		Text:   text,
		Medium: medium,
	}
}

func NewAgentMessage(text string, medium OutputMediumType) Message {
	return Message{
		Role:   string(MessageRoleAgent),
		Text:   text,
		Medium: medium,
	}
}

func NewToolCallMessage(toolName, invocationID, arguments string) Message {
	return Message{
		Role:         string(MessageRoleToolCall),
		ToolName:     toolName,
		InvocationID: invocationID,
		Text:         arguments,
	}
}

func NewToolResultMessage(toolName, invocationID, result string) Message {
	return Message{
		Role:         string(MessageRoleToolResult),
		ToolName:     toolName,
		InvocationID: invocationID,
		Text:         result,
	}
}

// Tool creation helpers
func NewHTTPTool(name, description, baseURL, method string) *BaseToolDefinition {
	return &BaseToolDefinition{
		ModelToolName: name,
		Description:   description,
		HTTP: &BaseHTTPToolDetails{
			BaseURLPattern: baseURL,
			HTTPMethod:     method,
		},
	}
}

func NewClientTool(name, description string) *BaseToolDefinition {
	return &BaseToolDefinition{
		ModelToolName: name,
		Description:   description,
		Client:        &BaseClientToolDetails{},
	}
}

func NewDataConnectionTool(name, description string) *BaseToolDefinition {
	return &BaseToolDefinition{
		ModelToolName:  name,
		Description:    description,
		DataConnection: &BaseDataConnectionToolDetails{},
	}
}

// Parameter creation helpers
func NewDynamicParameter(name string, location ParameterLocation, schema interface{}, required bool) DynamicParameter {
	return DynamicParameter{
		Name:     name,
		Location: location,
		Schema:   schema,
		Required: required,
	}
}

func NewStaticParameter(name string, location ParameterLocation, value interface{}) StaticParameter {
	return StaticParameter{
		Name:     name,
		Location: location,
		Value:    value,
	}
}

func NewAutomaticParameter(name string, location ParameterLocation, knownValue KnownParameterValue) AutomaticParameter {
	return AutomaticParameter{
		Name:       name,
		Location:   location,
		KnownValue: knownValue,
	}
}

// Channel mode constants for data connections
const (
	ChannelModeUnspecified ChannelModeType = "CHANNEL_MODE_UNSPECIFIED"
	ChannelModeMixed       ChannelModeType = "CHANNEL_MODE_MIXED"
	ChannelModeSeparated   ChannelModeType = "CHANNEL_MODE_SEPARATED"
)

type ChannelModeType string

// Update DataConnectionAudioConfig to use the enum
func NewDataConnectionConfigWithChannelMode(websocketURL string, sampleRate int, channelMode ChannelModeType) *DataConnectionConfig {
	return &DataConnectionConfig{
		WebsocketURL: websocketURL,
		AudioConfig: &DataConnectionAudioConfig{
			SampleRate:  sampleRate,
			ChannelMode: string(channelMode),
		},
	}
}
