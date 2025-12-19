package ultravox

// CallStage represents a stage within a call
type CallStage struct {
	CallID               string         `json:"callId" yaml:"callId"`
	CallStageID          string         `json:"callStageId" yaml:"callStageId"`
	Created              string         `json:"created" yaml:"created"`
	InactivityMessages   []TimedMessage `json:"inactivityMessages,omitempty" yaml:"inactivityMessages,omitempty"`
	LanguageHint         string         `json:"languageHint,omitempty" yaml:"languageHint,omitempty"`
	Model                string         `json:"model" yaml:"model"`
	SystemPrompt         string         `json:"systemPrompt,omitempty" yaml:"systemPrompt,omitempty"`
	Temperature          float64        `json:"temperature" yaml:"temperature"`
	TimeExceededMessage  string         `json:"timeExceededMessage,omitempty" yaml:"timeExceededMessage,omitempty"`
	Voice                string         `json:"voice,omitempty" yaml:"voice,omitempty"`
	ExternalVoice        *ExternalVoice `json:"externalVoice,omitempty" yaml:"externalVoice,omitempty"`
	ErrorCount           int            `json:"errorCount" yaml:"errorCount"`
	ExperimentalSettings interface{}    `json:"experimentalSettings,omitempty" yaml:"experimentalSettings,omitempty"`
	InitialState         interface{}    `json:"initialState" yaml:"initialState"`
}

// CallTool represents a tool as used for a particular call
type CallTool struct {
	CallToolID string              `json:"callToolId" yaml:"callToolId"`
	ToolID     string              `json:"toolId,omitempty" yaml:"toolId,omitempty"`
	Name       string              `json:"name" yaml:"name"`
	Definition *CallToolDefinition `json:"definition" yaml:"definition"`
}

// CallToolDefinition contains the actual tool definition for a call
type CallToolDefinition struct {
	Description         string                         `json:"description" yaml:"description"`
	DynamicParameters   []DynamicParameter             `json:"dynamicParameters,omitempty" yaml:"dynamicParameters,omitempty"`
	StaticParameters    []StaticParameter              `json:"staticParameters,omitempty" yaml:"staticParameters,omitempty"`
	AutomaticParameters []AutomaticParameter           `json:"automaticParameters,omitempty" yaml:"automaticParameters,omitempty"`
	Timeout             UltravoxDuration               `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Precomputable       bool                           `json:"precomputable,omitempty" yaml:"precomputable,omitempty"`
	HTTP                *HTTPCallToolDetails           `json:"http,omitempty" yaml:"http,omitempty"`
	Client              *ClientCallToolDetails         `json:"client,omitempty" yaml:"client,omitempty"`
	DataConnection      *DataConnectionCallToolDetails `json:"dataConnection,omitempty" yaml:"dataConnection,omitempty"`
	DefaultReaction     AgentReactionType              `json:"defaultReaction,omitempty" yaml:"defaultReaction,omitempty"`
	StaticResponse      *StaticToolResponse            `json:"staticResponse,omitempty" yaml:"staticResponse,omitempty"`
}

// HTTPCallToolDetails contains HTTP-specific details for call tools
type HTTPCallToolDetails struct {
	BaseURLPattern  string   `json:"baseUrlPattern" yaml:"baseUrlPattern"`
	HTTPMethod      string   `json:"httpMethod" yaml:"httpMethod"`
	AuthHeaders     []string `json:"authHeaders,omitempty" yaml:"authHeaders,omitempty"`
	AuthQueryParams []string `json:"authQueryParams,omitempty" yaml:"authQueryParams,omitempty"`
	CallTokenScopes []string `json:"callTokenScopes,omitempty" yaml:"callTokenScopes,omitempty"`
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
	CallID        string       `json:"callId" yaml:"callId"`
	CallStageID   string       `json:"callStageId" yaml:"callStageId"`
	CallTimestamp string       `json:"callTimestamp" yaml:"callTimestamp"`
	Severity      SeverityType `json:"severity" yaml:"severity"`
	Type          string       `json:"type" yaml:"type"`
	Text          string       `json:"text" yaml:"text"`
	Extras        interface{}  `json:"extras,omitempty" yaml:"extras,omitempty"`
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
	Start UltravoxDuration `json:"start" yaml:"start"`
	End   UltravoxDuration `json:"end" yaml:"end"`
}

type MessageRole string
