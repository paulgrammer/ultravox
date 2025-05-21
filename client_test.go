package ultravox_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/paulgrammer/ultravox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHTTPClient implements the HTTPClient interface for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []ultravox.Option
		wantErr bool
	}{
		{
			name: "Default configuration",
			opts: []ultravox.Option{
				ultravox.WithAPIKey("test-api-key"),
			},
			wantErr: false,
		},
		{
			name: "Custom configuration",
			opts: []ultravox.Option{
				ultravox.WithAPIKey("test-api-key"),
				ultravox.WithSystemPrompt("Custom prompt"),
				ultravox.WithModel("custom-model"),
				ultravox.WithVoice("custom-voice"),
				ultravox.WithFirstSpeaker(ultravox.FirstSpeakerUser),
				ultravox.WithHTTPTimeout(30 * time.Second),
				ultravox.WithAPIBaseURL("https://custom-api.example.com"),
				ultravox.WithTemperature(0.7),
				ultravox.WithInitialOutputMedium(ultravox.OutputMediumText),
				ultravox.WithRecordingEnabled(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ultravox.NewClient(tt.opts...)
			assert.NotNil(t, client, "Client should not be nil")
		})
	}
}

func TestClient_WithHTTPClient(t *testing.T) {
	client := ultravox.NewClient(ultravox.WithAPIKey("test-api-key"))

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{}")),
			}, nil
		},
	}

	client.WithHTTPClient(mockClient)

	// Since we can't directly test the property, we'll verify it worked in the next test
	assert.NotNil(t, client, "Client should not be nil after setting HTTP client")
}

func TestClient_Call(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		callOpts       []ultravox.CallOption
		wantErr        bool
	}{
		{
			name: "Successful call creation",
			mockResponse: `{
				"callId": "call-123",
				"joinUrl": "wss://example.com/join/call-123",
				"created": "2023-05-20T12:34:56Z",
				"maxDuration": "3600s",
				"joinTimeout": "300s",
				"initialOutputMedium": "MESSAGE_MEDIUM_VOICE",
				"recordingEnabled": false,
				"errorCount": 0
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "API error",
			mockResponse:   `{"error": "Something went wrong"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "Invalid JSON response",
			mockResponse:   `{invalid json}`,
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name:           "Missing join URL",
			mockResponse:   `{"callId": "call-123"}`,
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "With call options",
			mockResponse: `{
				"callId": "call-123",
				"joinUrl": "wss://example.com/join/call-123",
				"created": "2023-05-20T12:34:56Z",
				"maxDuration": "3600s",
				"joinTimeout": "300s",
				"initialOutputMedium": "MESSAGE_MEDIUM_VOICE",
				"recordingEnabled": true,
				"errorCount": 0
			}`,
			mockStatusCode: http.StatusOK,
			callOpts: []ultravox.CallOption{
				ultravox.WithCallSystemPrompt("Override prompt"),
				ultravox.WithCallModel("override-model"),
				ultravox.WithCallVoice("override-voice"),
				ultravox.WithCallFirstSpeaker(ultravox.FirstSpeakerUser),
				ultravox.WithCallWebSocketMedium(24000, 24000),
				ultravox.WithCallTemperature(0.7),
				ultravox.WithCallRecordingEnabled(true),
			},
			wantErr: false,
		},
		{
			name: "With first speaker settings",
			mockResponse: `{
				"callId": "call-123",
				"joinUrl": "wss://example.com/join/call-123",
				"created": "2023-05-20T12:34:56Z",
				"maxDuration": "3600s",
				"joinTimeout": "300s",
				"initialOutputMedium": "MESSAGE_MEDIUM_VOICE",
				"recordingEnabled": false,
				"errorCount": 0
			}`,
			mockStatusCode: http.StatusOK,
			callOpts: []ultravox.CallOption{
				ultravox.WithCallFirstSpeakerSettings(ultravox.AgentFirstSpeaker(
					true, "Hello there!", "", 0,
				)),
			},
			wantErr: false,
		},
		{
			name: "With WebRTC medium",
			mockResponse: `{
				"callId": "call-123",
				"joinUrl": "wss://example.com/join/call-123",
				"created": "2023-05-20T12:34:56Z",
				"maxDuration": "3600s",
				"joinTimeout": "300s",
				"initialOutputMedium": "MESSAGE_MEDIUM_VOICE",
				"recordingEnabled": false,
				"errorCount": 0
			}`,
			mockStatusCode: http.StatusOK,
			callOpts: []ultravox.CallOption{
				ultravox.WithCallWebRTCMedium(),
			},
			wantErr: false,
		},
		{
			name: "With external voice",
			mockResponse: `{
				"callId": "call-123",
				"joinUrl": "wss://example.com/join/call-123",
				"created": "2023-05-20T12:34:56Z",
				"maxDuration": "3600s",
				"joinTimeout": "300s",
				"initialOutputMedium": "MESSAGE_MEDIUM_VOICE",
				"recordingEnabled": false,
				"errorCount": 0
			}`,
			mockStatusCode: http.StatusOK,
			callOpts: []ultravox.CallOption{
				ultravox.WithCallExternalVoice(ultravox.NewElevenLabsVoice("voice-id-123")),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Verify request
					assert.Equal(t, "POST", req.Method)
					assert.Contains(t, req.URL.String(), "/calls")
					assert.Equal(t, "test-api-key", req.Header.Get("X-API-Key"))
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

					// For tests with call options, verify request body
					if len(tt.callOpts) > 0 {
						body, err := io.ReadAll(req.Body)
						require.NoError(t, err)

						var requestBody map[string]interface{}
						err = json.Unmarshal(body, &requestBody)
						require.NoError(t, err)

						// Check specific options based on test case
						if tt.name == "With call options" {
							assert.Equal(t, "Override prompt", requestBody["systemPrompt"])
							assert.Equal(t, "override-model", requestBody["model"])
							assert.Equal(t, "override-voice", requestBody["voice"])
							assert.Equal(t, "FIRST_SPEAKER_USER", requestBody["firstSpeaker"])
							assert.Equal(t, 0.7, requestBody["temperature"])
							assert.Equal(t, true, requestBody["recordingEnabled"])

							medium := requestBody["medium"].(map[string]interface{})
							serverWebSocket := medium["serverWebSocket"].(map[string]interface{})
							assert.Equal(t, float64(24000), serverWebSocket["inputSampleRate"])
							assert.Equal(t, float64(24000), serverWebSocket["outputSampleRate"])
						} else if tt.name == "With first speaker settings" {
							settings := requestBody["firstSpeakerSettings"].(map[string]interface{})
							agent := settings["agent"].(map[string]interface{})
							assert.Equal(t, true, agent["uninterruptible"])
							assert.Equal(t, "Hello there!", agent["text"])
						} else if tt.name == "With WebRTC medium" {
							medium := requestBody["medium"].(map[string]interface{})
							_, hasWebRTC := medium["webRtc"]
							assert.True(t, hasWebRTC)
						} else if tt.name == "With external voice" {
							externalVoice := requestBody["externalVoice"].(map[string]interface{})
							elevenLabs := externalVoice["elevenLabs"].(map[string]interface{})
							assert.Equal(t, "voice-id-123", elevenLabs["voiceId"])
						}
					}

					return &http.Response{
						StatusCode: tt.mockStatusCode,
						Body:       io.NopCloser(bytes.NewBufferString(tt.mockResponse)),
					}, nil
				},
			}

			client := ultravox.NewClient(ultravox.WithAPIKey("test-api-key"))
			client.WithHTTPClient(mockClient)

			ctx := context.Background()
			call, err := client.Call(ctx, tt.callOpts...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, call)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, call)
				assert.Equal(t, "call-123", call.CallID)
				assert.Equal(t, "wss://example.com/join/call-123", call.JoinURL)
			}
		})
	}
}

func TestCallWithPriorCallIdAndGreetingPrompt(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify URL includes query parameters
			url := req.URL.String()
			assert.Contains(t, url, "priorCallId=prior-call-123")
			assert.Contains(t, url, "enableGreetingPrompt=true")

			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"callId": "call-123",
					"joinUrl": "wss://example.com/join/call-123",
					"created": "2023-05-20T12:34:56Z",
					"maxDuration": "3600s",
					"joinTimeout": "300s"
				}`)),
			}, nil
		},
	}

	client := ultravox.NewClient(ultravox.WithAPIKey("test-api-key"))
	client.WithHTTPClient(mockClient)

	ctx := context.Background()
	call, err := client.Call(ctx,
		ultravox.WithCallPriorCallId("prior-call-123"),
		ultravox.WithCallEnableGreetingPrompt(true),
	)

	assert.NoError(t, err)
	assert.NotNil(t, call)
}

func TestCall_WithVadSettings(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			var requestBody map[string]interface{}
			err = json.Unmarshal(body, &requestBody)
			require.NoError(t, err)

			// Check VAD settings
			vadSettings := requestBody["vadSettings"].(map[string]interface{})
			assert.Equal(t, "0.5s", vadSettings["turnEndpointDelay"])
			assert.Equal(t, "0.1s", vadSettings["minimumTurnDuration"])

			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"callId": "call-123",
					"joinUrl": "wss://example.com/join/call-123",
					"created": "2023-05-20T12:34:56Z",
					"maxDuration": "3600s",
					"joinTimeout": "300s"
				}`)),
			}, nil
		},
	}

	client := ultravox.NewClient(ultravox.WithAPIKey("test-api-key"))
	client.WithHTTPClient(mockClient)

	vadSettings := ultravox.NewVadSettings()
	vadSettings.TurnEndpointDelay = ultravox.UltravoxDuration(500 * time.Millisecond)
	vadSettings.MinimumTurnDuration = ultravox.UltravoxDuration(100 * time.Millisecond)

	ctx := context.Background()
	call, err := client.Call(ctx, ultravox.WithCallVadSettings(vadSettings))

	assert.NoError(t, err)
	assert.NotNil(t, call)
}

func TestCall_WithInactivityMessages(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			var requestBody map[string]interface{}
			err = json.Unmarshal(body, &requestBody)
			require.NoError(t, err)

			// Check inactivity messages
			inactivityMessages := requestBody["inactivityMessages"].([]interface{})
			assert.Len(t, inactivityMessages, 2)

			message1 := inactivityMessages[0].(map[string]interface{})
			assert.Equal(t, "10s", message1["duration"])
			assert.Equal(t, "Are you still there?", message1["message"])

			message2 := inactivityMessages[1].(map[string]interface{})
			assert.Equal(t, "30s", message2["duration"])
			assert.Equal(t, "I'll end the call soon if I don't hear from you.", message2["message"])
			assert.Equal(t, "END_BEHAVIOR_HANG_UP_SOFT", message2["endBehavior"])

			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"callId": "call-123",
					"joinUrl": "wss://example.com/join/call-123",
					"created": "2023-05-20T12:34:56Z",
					"maxDuration": "3600s",
					"joinTimeout": "300s"
				}`)),
			}, nil
		},
	}

	client := ultravox.NewClient(ultravox.WithAPIKey("test-api-key"))
	client.WithHTTPClient(mockClient)

	inactivityMessages := []ultravox.TimedMessage{
		ultravox.NewTimedMessage(10*time.Second, "Are you still there?", ultravox.EndBehaviorDefault),
		ultravox.NewTimedMessage(30*time.Second, "I'll end the call soon if I don't hear from you.", ultravox.EndBehaviorHangUpSoft),
	}

	ctx := context.Background()
	call, err := client.Call(ctx, ultravox.WithCallInactivityMessages(inactivityMessages))

	assert.NoError(t, err)
	assert.NotNil(t, call)
}

func TestCallOptions(t *testing.T) {
	// Create a call request to test modifications
	request := &ultravox.CallRequest{
		SystemPrompt: "default",
		Model:        "default-model",
		Voice:        "default-voice",
	}

	// Test each option individually to ensure it modifies the correct field
	t.Run("WithCallSystemPrompt", func(t *testing.T) {
		opt := ultravox.WithCallSystemPrompt("new prompt")
		opt(request)
		assert.Equal(t, "new prompt", request.SystemPrompt)
	})

	t.Run("WithCallTemperature", func(t *testing.T) {
		opt := ultravox.WithCallTemperature(0.8)
		opt(request)
		assert.Equal(t, 0.8, request.Temperature)
	})

	t.Run("WithCallModel", func(t *testing.T) {
		opt := ultravox.WithCallModel("new-model")
		opt(request)
		assert.Equal(t, "new-model", request.Model)
	})

	t.Run("WithCallVoice", func(t *testing.T) {
		opt := ultravox.WithCallVoice("new-voice")
		opt(request)
		assert.Equal(t, "new-voice", request.Voice)
	})

	t.Run("WithCallExternalVoice", func(t *testing.T) {
		externalVoice := ultravox.NewElevenLabsVoice("voice-id-123")
		opt := ultravox.WithCallExternalVoice(externalVoice)
		opt(request)
		assert.Equal(t, externalVoice, request.ExternalVoice)
	})

	t.Run("WithCallFirstSpeakerSettings", func(t *testing.T) {
		settings := ultravox.AgentFirstSpeaker(true, "Hello", "", 0)
		opt := ultravox.WithCallFirstSpeakerSettings(settings)
		opt(request)
		assert.Equal(t, settings, request.FirstSpeakerSettings)
	})

	t.Run("WithCallWebRTCMedium", func(t *testing.T) {
		opt := ultravox.WithCallWebRTCMedium()
		opt(request)
		assert.NotNil(t, request.Medium)
		assert.NotNil(t, request.Medium.WebRTC)
	})

	t.Run("WithCallTimeExceededMessage", func(t *testing.T) {
		opt := ultravox.WithCallTimeExceededMessage("Time's up, goodbye!")
		opt(request)
		assert.Equal(t, "Time's up, goodbye!", request.TimeExceededMessage)
	})

	t.Run("WithCallMetadata", func(t *testing.T) {
		metadata := map[string]string{"customer_id": "123", "session_id": "abc"}
		opt := ultravox.WithCallMetadata(metadata)
		opt(request)
		assert.Equal(t, metadata, request.Metadata)
	})
}

func TestClientOptions(t *testing.T) {
	// Create a base config to modify
	config := &ultravox.Config{
		APIKey:      "default-key",
		APIBaseURL:  "default-url",
		HTTPTimeout: 10 * time.Second,
		CallRequest: ultravox.CallRequest{
			SystemPrompt: "default-prompt",
			Model:        "default-model",
			Voice:        "default-voice",
			FirstSpeaker: ultravox.FirstSpeakerAgent,
		},
	}

	// Test each option individually
	t.Run("WithAPIKey", func(t *testing.T) {
		opt := ultravox.WithAPIKey("new-key")
		opt(config)
		assert.Equal(t, "new-key", config.APIKey)
	})

	t.Run("WithAPIBaseURL", func(t *testing.T) {
		opt := ultravox.WithAPIBaseURL("new-url")
		opt(config)
		assert.Equal(t, "new-url", config.APIBaseURL)
	})

	t.Run("WithSystemPrompt", func(t *testing.T) {
		opt := ultravox.WithSystemPrompt("new-prompt")
		opt(config)
		assert.Equal(t, "new-prompt", config.SystemPrompt)
	})

	t.Run("WithTemperature", func(t *testing.T) {
		opt := ultravox.WithTemperature(0.5)
		opt(config)
		assert.Equal(t, 0.5, config.Temperature)
	})

	t.Run("WithModel", func(t *testing.T) {
		opt := ultravox.WithModel("new-model")
		opt(config)
		assert.Equal(t, "new-model", config.Model)
	})

	t.Run("WithVoice", func(t *testing.T) {
		opt := ultravox.WithVoice("new-voice")
		opt(config)
		assert.Equal(t, "new-voice", config.Voice)
	})

	t.Run("WithFirstSpeaker", func(t *testing.T) {
		opt := ultravox.WithFirstSpeaker(ultravox.FirstSpeakerUser)
		opt(config)
		assert.Equal(t, ultravox.FirstSpeakerUser, config.FirstSpeaker)
	})

	t.Run("WithHTTPTimeout", func(t *testing.T) {
		opt := ultravox.WithHTTPTimeout(20 * time.Second)
		opt(config)
		assert.Equal(t, 20*time.Second, config.HTTPTimeout)
	})

	t.Run("WithFirstSpeakerSettings", func(t *testing.T) {
		settings := ultravox.AgentFirstSpeaker(true, "Hello", "", 0)
		opt := ultravox.WithFirstSpeakerSettings(settings)
		opt(config)
		assert.Equal(t, settings, config.FirstSpeakerSettings)
	})

	t.Run("WithExternalVoice", func(t *testing.T) {
		voice := ultravox.NewElevenLabsVoice("voice-id-123")
		opt := ultravox.WithExternalVoice(voice)
		opt(config)
		assert.Equal(t, voice, config.ExternalVoice)
	})

	t.Run("WithVadSettings", func(t *testing.T) {
		settings := ultravox.NewVadSettings()
		opt := ultravox.WithVadSettings(settings)
		opt(config)
		assert.Equal(t, settings, config.VadSettings)
	})

	t.Run("WithInitialOutputMedium", func(t *testing.T) {
		opt := ultravox.WithInitialOutputMedium(ultravox.OutputMediumText)
		opt(config)
		assert.Equal(t, ultravox.OutputMediumText, config.InitialOutputMedium)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("AgentFirstSpeaker", func(t *testing.T) {
		settings := ultravox.AgentFirstSpeaker(true, "Hello", "Greet the user warmly", 500*time.Millisecond)
		assert.NotNil(t, settings)
		assert.NotNil(t, settings.Agent)
		assert.Nil(t, settings.User)
		assert.True(t, settings.Agent.Uninterruptible)
		assert.Equal(t, "Hello", settings.Agent.Text)
		assert.Equal(t, "Greet the user warmly", settings.Agent.Prompt)
		assert.Equal(t, ultravox.UltravoxDuration(500*time.Millisecond), settings.Agent.Delay)
	})

	t.Run("UserFirstSpeaker", func(t *testing.T) {
		settings := ultravox.UserFirstSpeaker(5*time.Second, "Are you there?", "")
		assert.NotNil(t, settings)
		assert.NotNil(t, settings.User)
		assert.Nil(t, settings.Agent)
		assert.NotNil(t, settings.User.Fallback)
		assert.Equal(t, ultravox.UltravoxDuration(5*time.Second), settings.User.Fallback.Delay)
		assert.Equal(t, "Are you there?", settings.User.Fallback.Text)
	})

	t.Run("NewVadSettings", func(t *testing.T) {
		settings := ultravox.NewVadSettings()
		assert.NotNil(t, settings)
		assert.Equal(t, ultravox.UltravoxDuration(384*time.Millisecond), settings.TurnEndpointDelay)
		assert.Equal(t, ultravox.UltravoxDuration(0), settings.MinimumTurnDuration)
		assert.Equal(t, ultravox.UltravoxDuration(90*time.Millisecond), settings.MinimumInterruptionDuration)
		assert.Equal(t, 0.1, settings.FrameActivationThreshold)
	})

	t.Run("NewTimedMessage", func(t *testing.T) {
		message := ultravox.NewTimedMessage(30*time.Second, "Test message", ultravox.EndBehaviorHangUpSoft)
		assert.Equal(t, ultravox.UltravoxDuration(30*time.Second), message.Duration)
		assert.Equal(t, "Test message", message.Message)
		assert.Equal(t, ultravox.EndBehaviorHangUpSoft, message.EndBehavior)
	})

	t.Run("External Voice Creators", func(t *testing.T) {
		// Test ElevenLabs voice
		elevenLabsVoice := ultravox.NewElevenLabsVoice("voice-id-123")
		assert.NotNil(t, elevenLabsVoice.ElevenLabs)
		assert.Equal(t, "voice-id-123", elevenLabsVoice.ElevenLabs.VoiceID)

		// Test Cartesia voice
		cartesiaVoice := ultravox.NewCartesiaVoice("voice-id-456")
		assert.NotNil(t, cartesiaVoice.Cartesia)
		assert.Equal(t, "voice-id-456", cartesiaVoice.Cartesia.VoiceID)

		// Test PlayHt voice
		playHtVoice := ultravox.NewPlayHtVoice("user-id", "voice-id-789")
		assert.NotNil(t, playHtVoice.PlayHt)
		assert.Equal(t, "user-id", playHtVoice.PlayHt.UserID)
		assert.Equal(t, "voice-id-789", playHtVoice.PlayHt.VoiceID)

		// Test LMNT voice
		lmntVoice := ultravox.NewLmntVoice("voice-id-abc")
		assert.NotNil(t, lmntVoice.Lmnt)
		assert.Equal(t, "voice-id-abc", lmntVoice.Lmnt.VoiceID)

		// Test Generic voice
		body := map[string]string{"param": "value"}
		genericVoice := ultravox.NewGenericVoice("https://example.com/tts", body)
		assert.NotNil(t, genericVoice.Generic)
		assert.Equal(t, "https://example.com/tts", genericVoice.Generic.URL)
		assert.Equal(t, body, genericVoice.Generic.Body)
	})

	t.Run("NewDataConnectionConfig", func(t *testing.T) {
		config := ultravox.NewDataConnectionConfig("wss://example.com/data", 16000)
		assert.NotNil(t, config)
		assert.Equal(t, "wss://example.com/data", config.WebsocketURL)
		assert.NotNil(t, config.AudioConfig)
		assert.Equal(t, 16000, config.AudioConfig.SampleRate)
	})
}

// TestIntegration_Call is an integration test that should be skipped by default
// It can be run with: go test -tags=integration
func TestIntegration_Call(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a valid API key in the environment
	client := ultravox.NewClient()

	ctx := context.Background()
	call, err := client.Call(ctx, ultravox.WithCallSystemPrompt("Say hello and introduce yourself briefly."))

	require.NoError(t, err)
	assert.NotEmpty(t, call.CallID)
	assert.NotEmpty(t, call.JoinURL)
	t.Logf("Created call with ID: %s and Join URL: %s", call.CallID, call.JoinURL)
}
