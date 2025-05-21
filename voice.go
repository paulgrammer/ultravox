package ultravox

// ExternalVoice contains configurations for external voice providers
type ExternalVoice struct {
	ElevenLabs *ElevenLabsVoice `json:"elevenLabs,omitempty"`
	Cartesia   *CartesiaVoice   `json:"cartesia,omitempty"`
	PlayHt     *PlayHtVoice     `json:"playHt,omitempty"`
	Lmnt       *LmntVoice       `json:"lmnt,omitempty"`
	Generic    *GenericVoice    `json:"generic,omitempty"`
}

// ElevenLabsVoice defines configuration for ElevenLabs voice service
type ElevenLabsVoice struct {
	VoiceID                  string `json:"voiceId"`
	Model                    string `json:"model,omitempty"`
	Speed                    float64 `json:"speed,omitempty"`
	UseSpeakerBoost          bool   `json:"useSpeakerBoost,omitempty"`
	Style                    float64 `json:"style,omitempty"`
	SimilarityBoost          float64 `json:"similarityBoost,omitempty"`
	Stability                float64 `json:"stability,omitempty"`
	PronunciationDictionaries []PronunciationDictionary `json:"pronunciationDictionaries,omitempty"`
	OptimizeStreamingLatency int    `json:"optimizeStreamingLatency,omitempty"`
	MaxSampleRate            int    `json:"maxSampleRate,omitempty"`
}

// PronunciationDictionary references a pronunciation dictionary in ElevenLabs
type PronunciationDictionary struct {
	DictionaryID string `json:"dictionaryId"`
	VersionID    string `json:"versionId,omitempty"`
}

// CartesiaVoice defines configuration for Cartesia voice service
type CartesiaVoice struct {
	VoiceID  string   `json:"voiceId"`
	Model    string   `json:"model,omitempty"`
	Speed    float64  `json:"speed,omitempty"`
	Emotion  string   `json:"emotion,omitempty"`
	Emotions []string `json:"emotions,omitempty"`
}

// PlayHtVoice defines configuration for PlayHT voice service
type PlayHtVoice struct {
	UserID                 string  `json:"userId"`
	VoiceID                string  `json:"voiceId"`
	Model                  string  `json:"model,omitempty"`
	Speed                  float64 `json:"speed,omitempty"`
	Quality                string  `json:"quality,omitempty"`
	Temperature            float64 `json:"temperature,omitempty"`
	Emotion                float64 `json:"emotion,omitempty"`
	VoiceGuidance          float64 `json:"voiceGuidance,omitempty"`
	StyleGuidance          float64 `json:"styleGuidance,omitempty"`
	TextGuidance           float64 `json:"textGuidance,omitempty"`
	VoiceConditioningSeconds float64 `json:"voiceConditioningSeconds,omitempty"`
}

// LmntVoice defines configuration for LMNT voice service
type LmntVoice struct {
	VoiceID       string  `json:"voiceId"`
	Model         string  `json:"model,omitempty"`
	Speed         float64 `json:"speed,omitempty"`
	Conversational bool   `json:"conversational,omitempty"`
}

// GenericVoice defines configuration for a generic voice service
type GenericVoice struct {
	URL                  string            `json:"url"`
	Headers              map[string]string `json:"headers,omitempty"`
	Body                 interface{}       `json:"body,omitempty"`
	ResponseSampleRate   int               `json:"responseSampleRate,omitempty"`
	ResponseWordsPerMinute int             `json:"responseWordsPerMinute,omitempty"`
	ResponseMimeType     string            `json:"responseMimeType,omitempty"`
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