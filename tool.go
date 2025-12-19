package ultravox

// SelectedTool represents a tool selected for a particular call
type SelectedTool struct {
	ToolID              string                 `json:"toolId,omitempty" yaml:"toolId,omitempty"`
	ToolName            string                 `json:"toolName,omitempty" yaml:"toolName,omitempty"`
	TemporaryTool       *BaseToolDefinition    `json:"temporaryTool,omitempty" yaml:"temporaryTool,omitempty"`
	NameOverride        string                 `json:"nameOverride,omitempty" yaml:"nameOverride,omitempty"`
	DescriptionOverride string                 `json:"descriptionOverride,omitempty" yaml:"descriptionOverride,omitempty"`
	AuthTokens          map[string]string      `json:"authTokens,omitempty" yaml:"authTokens,omitempty"`
	ParameterOverrides  map[string]interface{} `json:"parameterOverrides,omitempty" yaml:"parameterOverrides,omitempty"`
	TransitionID        string                 `json:"transitionId,omitempty" yaml:"transitionId,omitempty"`
}

// BaseToolDefinition defines a tool that can be used during a call
type BaseToolDefinition struct {
	ModelToolName       string                         `json:"modelToolName" yaml:"modelToolName"`
	Description         string                         `json:"description" yaml:"description"`
	DynamicParameters   []DynamicParameter             `json:"dynamicParameters,omitempty" yaml:"dynamicParameters,omitempty"`
	StaticParameters    []StaticParameter              `json:"staticParameters,omitempty" yaml:"staticParameters,omitempty"`
	AutomaticParameters []AutomaticParameter           `json:"automaticParameters,omitempty" yaml:"automaticParameters,omitempty"`
	Requirements        *ToolRequirements              `json:"requirements,omitempty" yaml:"requirements,omitempty"`
	Timeout             UltravoxDuration               `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Precomputable       bool                           `json:"precomputable,omitempty" yaml:"precomputable,omitempty"`
	HTTP                *BaseHTTPToolDetails           `json:"http,omitempty" yaml:"http,omitempty"`
	Client              *BaseClientToolDetails         `json:"client,omitempty" yaml:"client,omitempty"`
	DataConnection      *BaseDataConnectionToolDetails `json:"dataConnection,omitempty" yaml:"dataConnection,omitempty"`
	DefaultReaction     AgentReactionType              `json:"defaultReaction,omitempty" yaml:"defaultReaction,omitempty"`
	StaticResponse      *StaticToolResponse            `json:"staticResponse,omitempty" yaml:"staticResponse,omitempty"`
}

// DynamicParameter represents a parameter that can be set by the model
type DynamicParameter struct {
	Name     string            `json:"name" yaml:"name"`
	Location ParameterLocation `json:"location" yaml:"location"`
	Schema   interface{}       `json:"schema" yaml:"schema"`
	Required bool              `json:"required,omitempty" yaml:"required,omitempty"`
}

// StaticParameter represents a parameter that is unconditionally added
type StaticParameter struct {
	Name     string            `json:"name" yaml:"name"`
	Location ParameterLocation `json:"location" yaml:"location"`
	Value    interface{}       `json:"value" yaml:"value"`
}

// AutomaticParameter represents a parameter automatically set by the system
type AutomaticParameter struct {
	Name       string              `json:"name" yaml:"name"`
	Location   ParameterLocation   `json:"location" yaml:"location"`
	KnownValue KnownParameterValue `json:"knownValue" yaml:"knownValue"`
}

// ToolRequirements defines requirements for using a tool
type ToolRequirements struct {
	HTTPSecurityOptions        *SecurityOptions `json:"httpSecurityOptions,omitempty" yaml:"httpSecurityOptions,omitempty"`
	RequiredParameterOverrides []string         `json:"requiredParameterOverrides,omitempty" yaml:"requiredParameterOverrides,omitempty"`
}

// SecurityOptions defines different security requirement options
type SecurityOptions struct {
	Options []SecurityRequirements `json:"options" yaml:"options"`
}

// SecurityRequirements defines a set of security requirements
type SecurityRequirements struct {
	Requirements                 map[string]SecurityRequirement `json:"requirements,omitempty" yaml:"requirements,omitempty"`
	UltravoxCallTokenRequirement *UltravoxCallTokenRequirement  `json:"ultravoxCallTokenRequirement,omitempty" yaml:"ultravoxCallTokenRequirement,omitempty"`
}

// SecurityRequirement defines a single security requirement
type SecurityRequirement struct {
	QueryAPIKey  *QueryAPIKeyRequirement  `json:"queryApiKey,omitempty" yaml:"queryApiKey,omitempty"`
	HeaderAPIKey *HeaderAPIKeyRequirement `json:"headerApiKey,omitempty" yaml:"headerApiKey,omitempty"`
	HTTPAuth     *HTTPAuthRequirement     `json:"httpAuth,omitempty" yaml:"httpAuth,omitempty"`
}

// QueryAPIKeyRequirement adds an API key to query string
type QueryAPIKeyRequirement struct {
	Name string `json:"name" yaml:"name"`
}

// HeaderAPIKeyRequirement adds an API key to header
type HeaderAPIKeyRequirement struct {
	Name string `json:"name" yaml:"name"`
}

// HTTPAuthRequirement adds HTTP authentication header
type HTTPAuthRequirement struct {
	Scheme string `json:"scheme" yaml:"scheme"`
}

// UltravoxCallTokenRequirement defines call token requirements
type UltravoxCallTokenRequirement struct {
	Scopes []string `json:"scopes" yaml:"scopes"`
}

// BaseHTTPToolDetails defines details for HTTP tools
type BaseHTTPToolDetails struct {
	BaseURLPattern string `json:"baseUrlPattern" yaml:"baseUrlPattern"`
	HTTPMethod     string `json:"httpMethod" yaml:"httpMethod"`
}

// BaseClientToolDetails defines details for client-implemented tools
type BaseClientToolDetails struct {
	// Empty for now, but included for completeness
}

// BaseDataConnectionToolDetails defines details for data connection tools
type BaseDataConnectionToolDetails struct {
	// Empty for now, but included for completeness
}

// StaticToolResponse defines a predefined static response
type StaticToolResponse struct {
	ResponseText string `json:"responseText" yaml:"responseText"`
}

// Enums for tool-related types
type ParameterLocation string

const (
	ParameterLocationUnspecified ParameterLocation = "PARAMETER_LOCATION_UNSPECIFIED"
	ParameterLocationQuery       ParameterLocation = "PARAMETER_LOCATION_QUERY"
	ParameterLocationPath        ParameterLocation = "PARAMETER_LOCATION_PATH"
	ParameterLocationHeader      ParameterLocation = "PARAMETER_LOCATION_HEADER"
	ParameterLocationBody        ParameterLocation = "PARAMETER_LOCATION_BODY"
)

type KnownParameterValue string

const (
	KnownParamUnspecified         KnownParameterValue = "KNOWN_PARAM_UNSPECIFIED"
	KnownParamCallID              KnownParameterValue = "KNOWN_PARAM_CALL_ID"
	KnownParamConversationHistory KnownParameterValue = "KNOWN_PARAM_CONVERSATION_HISTORY"
	KnownParamOutputSampleRate    KnownParameterValue = "KNOWN_PARAM_OUTPUT_SAMPLE_RATE"
	KnownParamCallState           KnownParameterValue = "KNOWN_PARAM_CALL_STATE"
)

type AgentReactionType string

const (
	AgentReactionUnspecified AgentReactionType = "AGENT_REACTION_UNSPECIFIED"
	AgentReactionSpeaks      AgentReactionType = "AGENT_REACTION_SPEAKS"
	AgentReactionListens     AgentReactionType = "AGENT_REACTION_LISTENS"
	AgentReactionSpeaksOnce  AgentReactionType = "AGENT_REACTION_SPEAKS_ONCE"
)
