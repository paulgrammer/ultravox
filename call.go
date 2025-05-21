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
