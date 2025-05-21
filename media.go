package ultravox

import "time"

// Message represents a message in a conversation
type Message struct {
	Role         string          `json:"role,omitempty"`
	Text         string          `json:"text,omitempty"`
	InvocationID string          `json:"invocationId,omitempty"`
	ToolName     string          `json:"toolName,omitempty"`
	ErrorDetails string          `json:"errorDetails,omitempty"`
	Medium       OutputMediumType `json:"medium,omitempty"`
}

// TimedMessage represents a message that should be delivered after a specific duration
type TimedMessage struct {
	Duration    UltravoxDuration `json:"duration"`
	Message     string          `json:"message"`
	EndBehavior EndBehaviorType `json:"endBehavior,omitempty"`
}

// FirstSpeakerSettings defines who speaks first and related settings
type FirstSpeakerSettings struct {
	User  *UserGreeting  `json:"user,omitempty"`
	Agent *AgentGreeting `json:"agent,omitempty"`
}

// UserGreeting contains settings for when the user speaks first
type UserGreeting struct {
	Fallback *FallbackAgentGreeting `json:"fallback,omitempty"`
}

// AgentGreeting contains settings for when the agent speaks first
type AgentGreeting struct {
	Uninterruptible bool            `json:"uninterruptible,omitempty"`
	Text            string          `json:"text,omitempty"`
	Prompt          string          `json:"prompt,omitempty"`
	Delay           UltravoxDuration `json:"delay,omitempty"`
}

// FallbackAgentGreeting defines a fallback greeting if the user doesn't speak
type FallbackAgentGreeting struct {
	Delay  UltravoxDuration `json:"delay,omitempty"`
	Text   string          `json:"text,omitempty"`
	Prompt string          `json:"prompt,omitempty"`
}

// VadSettings contains voice activity detection settings
type VadSettings struct {
	TurnEndpointDelay         UltravoxDuration `json:"turnEndpointDelay,omitempty"`
	MinimumTurnDuration       UltravoxDuration `json:"minimumTurnDuration,omitempty"`
	MinimumInterruptionDuration UltravoxDuration `json:"minimumInterruptionDuration,omitempty"`
	FrameActivationThreshold  float64         `json:"frameActivationThreshold,omitempty"`
}

// CallMedium defines the medium used for the call
type CallMedium struct {
	WebRTC        *WebRTCMedium        `json:"webRtc,omitempty"`
	Twilio        *TwilioMedium        `json:"twilio,omitempty"`
	ServerWebSocket *WebSocketMedium   `json:"serverWebSocket,omitempty"`
	Telnyx        *TelnyxMedium        `json:"telnyx,omitempty"`
	Plivo         *PlivoMedium         `json:"plivo,omitempty"`
	Exotel        *ExotelMedium        `json:"exotel,omitempty"`
	SIP           *SIPMedium           `json:"sip,omitempty"`
}

// WebRTCMedium defines WebRTC-specific configuration
type WebRTCMedium struct {}

// TwilioMedium defines Twilio-specific configuration
type TwilioMedium struct {}

// WebSocketMedium defines WebSocket-specific connection parameters
type WebSocketMedium struct {
	InputSampleRate      int `json:"inputSampleRate"`
	OutputSampleRate     int `json:"outputSampleRate,omitempty"`
	ClientBufferSizeMs   int `json:"clientBufferSizeMs,omitempty"`
}

// TelnyxMedium defines Telnyx-specific configuration
type TelnyxMedium struct {}

// PlivoMedium defines Plivo-specific configuration
type PlivoMedium struct {}

// ExotelMedium defines Exotel-specific configuration
type ExotelMedium struct {}

// SIPMedium defines SIP-specific configuration
type SIPMedium struct {
	Incoming *SIPIncoming `json:"incoming,omitempty"`
	Outgoing *SIPOutgoing `json:"outgoing,omitempty"`
}

// SIPIncoming defines incoming SIP call configuration
type SIPIncoming struct {}

// SIPOutgoing defines outgoing SIP call configuration
type SIPOutgoing struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// DataConnectionConfig contains settings for data connections
type DataConnectionConfig struct {
	WebsocketURL string              `json:"websocketUrl"`
	AudioConfig  *DataConnectionAudioConfig `json:"audioConfig,omitempty"`
}

// DataConnectionAudioConfig defines audio settings for data connections
type DataConnectionAudioConfig struct {
	SampleRate   int    `json:"sampleRate,omitempty"`
	ChannelMode  string `json:"channelMode,omitempty"`
}

// AgentFirstSpeaker returns a FirstSpeakerSettings configured for agent to speak first
func AgentFirstSpeaker(uninterruptible bool, text, prompt string, delay time.Duration) *FirstSpeakerSettings {
	return &FirstSpeakerSettings{
		Agent: &AgentGreeting{
			Uninterruptible: uninterruptible,
			Text:           text,
			Prompt:         prompt,
			Delay:          UltravoxDuration(delay),
		},
	}
}

// UserFirstSpeaker returns a FirstSpeakerSettings configured for user to speak first
func UserFirstSpeaker(fallbackDelay time.Duration, fallbackText, fallbackPrompt string) *FirstSpeakerSettings {
	return &FirstSpeakerSettings{
		User: &UserGreeting{
			Fallback: &FallbackAgentGreeting{
				Delay:  UltravoxDuration(fallbackDelay),
				Text:   fallbackText,
				Prompt: fallbackPrompt,
			},
		},
	}
}

// NewVadSettings creates a new VadSettings with common defaults
func NewVadSettings() *VadSettings {
	return &VadSettings{
		TurnEndpointDelay:          UltravoxDuration(384 * time.Millisecond),
		MinimumTurnDuration:        UltravoxDuration(0),
		MinimumInterruptionDuration: UltravoxDuration(90 * time.Millisecond),
		FrameActivationThreshold:   0.1,
	}
}

// NewTimedMessage creates a new timed message
func NewTimedMessage(duration time.Duration, message string, endBehavior EndBehaviorType) TimedMessage {
	return TimedMessage{
		Duration:    UltravoxDuration(duration),
		Message:     message,
		EndBehavior: endBehavior,
	}
}

// NewDataConnectionConfig creates a new data connection configuration
func NewDataConnectionConfig(websocketURL string, sampleRate int) *DataConnectionConfig {
	return &DataConnectionConfig{
		WebsocketURL: websocketURL,
		AudioConfig: &DataConnectionAudioConfig{
			SampleRate: sampleRate,
		},
	}
}