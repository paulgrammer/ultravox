package ultravox

// ExternalVoice contains configurations for external voice providers
type ExternalVoice struct {
	ElevenLabs *ElevenLabsVoice `json:"elevenLabs,omitempty" yaml:"elevenLabs,omitempty"`
	Cartesia   *CartesiaVoice   `json:"cartesia,omitempty" yaml:"cartesia,omitempty"`
	PlayHt     *PlayHtVoice     `json:"playHt,omitempty" yaml:"playHt,omitempty"`
	Lmnt       *LmntVoice       `json:"lmnt,omitempty" yaml:"lmnt,omitempty"`
	Generic    *GenericVoice    `json:"generic,omitempty" yaml:"generic,omitempty"`
}

// ElevenLabsVoice defines configuration for ElevenLabs voice service
type ElevenLabsVoice struct {
	VoiceID                   string                    `json:"voiceId" yaml:"voiceId"`
	Model                     string                    `json:"model,omitempty" yaml:"model,omitempty"`
	Speed                     float64                   `json:"speed,omitempty" yaml:"speed,omitempty"`
	UseSpeakerBoost           bool                      `json:"useSpeakerBoost,omitempty" yaml:"useSpeakerBoost,omitempty"`
	Style                     float64                   `json:"style,omitempty" yaml:"style,omitempty"`
	SimilarityBoost           float64                   `json:"similarityBoost,omitempty" yaml:"similarityBoost,omitempty"`
	Stability                 float64                   `json:"stability,omitempty" yaml:"stability,omitempty"`
	PronunciationDictionaries []PronunciationDictionary `json:"pronunciationDictionaries,omitempty" yaml:"pronunciationDictionaries,omitempty"`
	OptimizeStreamingLatency  int                       `json:"optimizeStreamingLatency,omitempty" yaml:"optimizeStreamingLatency,omitempty"`
	MaxSampleRate             int                       `json:"maxSampleRate,omitempty" yaml:"maxSampleRate,omitempty"`
}

// PronunciationDictionary references a pronunciation dictionary in ElevenLabs
type PronunciationDictionary struct {
	DictionaryID string `json:"dictionaryId" yaml:"dictionaryId"`
	VersionID    string `json:"versionId,omitempty" yaml:"versionId,omitempty"`
}

// CartesiaVoice defines configuration for Cartesia voice service
type CartesiaVoice struct {
	VoiceID  string   `json:"voiceId" yaml:"voiceId"`
	Model    string   `json:"model,omitempty" yaml:"model,omitempty"`
	Speed    float64  `json:"speed,omitempty" yaml:"speed,omitempty"`
	Emotion  string   `json:"emotion,omitempty" yaml:"emotion,omitempty"`
	Emotions []string `json:"emotions,omitempty" yaml:"emotions,omitempty"`
}

// PlayHtVoice defines configuration for PlayHT voice service
type PlayHtVoice struct {
	UserID                   string  `json:"userId" yaml:"userId"`
	VoiceID                  string  `json:"voiceId" yaml:"voiceId"`
	Model                    string  `json:"model,omitempty" yaml:"model,omitempty"`
	Speed                    float64 `json:"speed,omitempty" yaml:"speed,omitempty"`
	Quality                  string  `json:"quality,omitempty" yaml:"quality,omitempty"`
	Temperature              float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	Emotion                  float64 `json:"emotion,omitempty" yaml:"emotion,omitempty"`
	VoiceGuidance            float64 `json:"voiceGuidance,omitempty" yaml:"voiceGuidance,omitempty"`
	StyleGuidance            float64 `json:"styleGuidance,omitempty" yaml:"styleGuidance,omitempty"`
	TextGuidance             float64 `json:"textGuidance,omitempty" yaml:"textGuidance,omitempty"`
	VoiceConditioningSeconds float64 `json:"voiceConditioningSeconds,omitempty" yaml:"voiceConditioningSeconds,omitempty"`
}

// LmntVoice defines configuration for LMNT voice service
type LmntVoice struct {
	VoiceID        string  `json:"voiceId" yaml:"voiceId"`
	Model          string  `json:"model,omitempty" yaml:"model,omitempty"`
	Speed          float64 `json:"speed,omitempty" yaml:"speed,omitempty"`
	Conversational bool    `json:"conversational,omitempty" yaml:"conversational,omitempty"`
}

// GenericVoice defines configuration for a generic voice service
type GenericVoice struct {
	URL                    string            `json:"url" yaml:"url"`
	Headers                map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Body                   interface{}       `json:"body,omitempty" yaml:"body,omitempty"`
	ResponseSampleRate     int               `json:"responseSampleRate,omitempty" yaml:"responseSampleRate,omitempty"`
	ResponseWordsPerMinute int               `json:"responseWordsPerMinute,omitempty" yaml:"responseWordsPerMinute,omitempty"`
	ResponseMimeType       string            `json:"responseMimeType,omitempty" yaml:"responseMimeType,omitempty"`
}

// NewElevenLabsVoice creates a new ElevenLabs voice configuration
func NewElevenLabsVoice(voiceID string) *ExternalVoice {
	return &ExternalVoice{
		ElevenLabs: &ElevenLabsVoice{
			VoiceID: voiceID,
		},
	}
}

// NewCartesiaVoice creates a new Cartesia voice configuration
func NewCartesiaVoice(voiceID string) *ExternalVoice {
	return &ExternalVoice{
		Cartesia: &CartesiaVoice{
			VoiceID: voiceID,
		},
	}
}

// NewPlayHtVoice creates a new PlayHT voice configuration
func NewPlayHtVoice(userID, voiceID string) *ExternalVoice {
	return &ExternalVoice{
		PlayHt: &PlayHtVoice{
			UserID:  userID,
			VoiceID: voiceID,
		},
	}
}

// NewLmntVoice creates a new LMNT voice configuration
func NewLmntVoice(voiceID string) *ExternalVoice {
	return &ExternalVoice{
		Lmnt: &LmntVoice{
			VoiceID: voiceID,
		},
	}
}

// NewGenericVoice creates a new generic voice configuration
func NewGenericVoice(url string, body interface{}) *ExternalVoice {
	return &ExternalVoice{
		Generic: &GenericVoice{
			URL:  url,
			Body: body,
		},
	}
}
