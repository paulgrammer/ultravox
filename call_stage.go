package ultravox

// CallStage represents a stage within a call
type CallStage struct {
	CallID               string         `json:"callId"`
	CallStageID          string         `json:"callStageId"`
	Created              string         `json:"created"`
	InactivityMessages   []TimedMessage `json:"inactivityMessages,omitempty"`
	LanguageHint         string         `json:"languageHint,omitempty"`
	Model                string         `json:"model"`
	SystemPrompt         string         `json:"systemPrompt,omitempty"`
	Temperature          float64        `json:"temperature"`
	TimeExceededMessage  string         `json:"timeExceededMessage,omitempty"`
	Voice                string         `json:"voice,omitempty"`
	ExternalVoice        *ExternalVoice `json:"externalVoice,omitempty"`
	ErrorCount           int            `json:"errorCount"`
	ExperimentalSettings interface{}    `json:"experimentalSettings,omitempty"`
	InitialState         interface{}    `json:"initialState"`
}

// CallTool represents a tool as used for a particular call
type CallTool struct {
	CallToolID string              `json:"callToolId"`
	ToolID     string              `json:"toolId,omitempty"`
	Name       string              `json:"name"`
	Definition *CallToolDefinition `json:"definition"`
}

// CallToolDefinition contains the actual tool definition for a call
type CallToolDefinition struct {
	Description         string                         `json:"description"`
	DynamicParameters   []DynamicParameter             `json:"dynamicParameters,omitempty"`
	StaticParameters    []StaticParameter              `json:"staticParameters,omitempty"`
	AutomaticParameters []AutomaticParameter           `json:"automaticParameters,omitempty"`
	Timeout             UltravoxDuration               `json:"timeout,omitempty"`
	Precomputable       bool                           `json:"precomputable,omitempty"`
	HTTP                *HTTPCallToolDetails           `json:"http,omitempty"`
	Client              *ClientCallToolDetails         `json:"client,omitempty"`
	DataConnection      *DataConnectionCallToolDetails `json:"dataConnection,omitempty"`
	DefaultReaction     AgentReactionType              `json:"defaultReaction,omitempty"`
	StaticResponse      *StaticToolResponse            `json:"staticResponse,omitempty"`
}

// HTTPCallToolDetails contains HTTP-specific details for call tools
type HTTPCallToolDetails struct {
	BaseURLPattern  string   `json:"baseUrlPattern"`
	HTTPMethod      string   `json:"httpMethod"`
	AuthHeaders     []string `json:"authHeaders,omitempty"`
	AuthQueryParams []string `json:"authQueryParams,omitempty"`
	CallTokenScopes []string `json:"callTokenScopes,omitempty"`
}

// ClientCallToolDetails contains client-specific details for call tools
type ClientCallToolDetails struct {
	// Empty for now, but included for completeness
}

// DataConnectionCallToolDetails contains data connection details for call tools
type DataConnectionCallToolDetails struct {
	// Empty for now, but included for completeness
}

// CallEvent represents an event that occurred during a call
type CallEvent struct {
	CallID        string       `json:"callId"`
	CallStageID   string       `json:"callStageId"`
	CallTimestamp string       `json:"callTimestamp"`
	Severity      SeverityType `json:"severity"`
	Type          string       `json:"type"`
	Text          string       `json:"text"`
	Extras        interface{}  `json:"extras,omitempty"`
}

// SeverityType defines the severity of an event
type SeverityType string

const (
	SeverityDebug   SeverityType = "debug"
	SeverityInfo    SeverityType = "info"
	SeverityWarning SeverityType = "warning"
	SeverityError   SeverityType = "error"
)

// InCallTimespan represents a timespan during a call
type InCallTimespan struct {
	Start UltravoxDuration `json:"start"`
	End   UltravoxDuration `json:"end"`
}

type MessageRole string
