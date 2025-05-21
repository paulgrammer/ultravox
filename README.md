# Ultravox - Voice AI API Client for Go

A robust Go client library for the [Ultravox AI](https://ultravox.ai) voice API, enabling developers to implement voice-based AI interactions in their applications.

## Installation

```bash
go get github.com/paulgrammer/ultravox
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	"github.com/paulgrammer/ultravox"
)

func main() {
	// Initialize client
	client := ultravox.NewClient()

	// Create a new call
	call, err := client.Call(
		context.Background(),
		ultravox.WithCallSystemPrompt("You are a helpful assistant."),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Call created! Join URL: %s\n", call.JoinURL)
}
```

## Client Configuration

Configure the client with functional options:

```go
client := ultravox.NewClient(
	ultravox.WithAPIKey("your-api-key"),
	ultravox.WithSystemPrompt("You are a helpful assistant for a healthcare company..."),
	ultravox.WithModel("fixie-ai/ultravox-claude"),
	ultravox.WithVoice("Allison-English"),
	ultravox.WithMaxDuration(5 * time.Minute),
	ultravox.WithFirstSpeakerSettings(ultravox.AgentFirstSpeaker(
		false, // Not uninterruptible
		"Hello, how can I help you today?", // Text
		"", // Prompt
		0, // No delay
	)),
)
```

## Voice Providers

The client supports multiple voice synthesis providers:

```go
// Ultravox's built-in voices
client.Call(ctx, ultravox.WithCallVoice("Allison-English"))

// ElevenLabs integration
client.Call(ctx, ultravox.WithCallExternalVoice(
	ultravox.NewElevenLabsVoice("voice_id_here"),
))

// Additional provider integrations
client.Call(ctx, ultravox.WithCallExternalVoice(
	ultravox.NewCartesiaVoice("voice_id_here"),
))
client.Call(ctx, ultravox.WithCallExternalVoice(
	ultravox.NewPlayHtVoice("user_id", "voice_id"),
))
client.Call(ctx, ultravox.WithCallExternalVoice(
	ultravox.NewLmntVoice("voice_id_here"),
))
```

## Communication Protocols

Configure different communication methods:

```go
// WebSocket (default)
client.Call(ctx, ultravox.WithCallWebSocketMedium(16000, 16000))

// WebRTC
client.Call(ctx, ultravox.WithCallWebRTCMedium())

// Twilio
client.Call(ctx, ultravox.WithCallTwilioMedium())

// SIP
client.Call(ctx, ultravox.WithCallSIPOutgoing(
	"sip:user@example.com", // To
	"sip:agent@ultravox.ai", // From
	"username", // Auth username
	"password", // Auth password
))
```

## Advanced Features

### Voice Activity Detection (VAD)

Configure voice activity detection parameters:

```go
vadSettings := ultravox.NewVadSettings()
vadSettings.TurnEndpointDelay = ultravox.UltravoxDuration(500 * time.Millisecond)
vadSettings.MinimumTurnDuration = ultravox.UltravoxDuration(100 * time.Millisecond)
client.Call(ctx, ultravox.WithCallVadSettings(vadSettings))
```

### User Inactivity Handling

Define behavior for user inactivity:

```go
inactivityMessages := []ultravox.TimedMessage{
	ultravox.NewTimedMessage(10*time.Second, "Are you still there?", ultravox.EndBehaviorDefault),
	ultravox.NewTimedMessage(30*time.Second, "I'll hang up soon if I don't hear from you.", ultravox.EndBehaviorDefault),
	ultravox.NewTimedMessage(60*time.Second, "I'll end the call now. Feel free to call back later.", ultravox.EndBehaviorHangUpSoft),
}
client.Call(ctx, ultravox.WithCallInactivityMessages(inactivityMessages))
```

### Initial Message Context

Prime the conversation with initial messages:

```go
initialMessages := []ultravox.Message{
	{
		Role: "user",
		Text: "I need help with my order",
		Medium: ultravox.OutputMediumText,
	},
	{
		Role: "agent",
		Text: "I'd be happy to help with your order. Could you provide the order number?",
		Medium: ultravox.OutputMediumVoice,
	},
}
client.Call(ctx, ultravox.WithCallInitialMessages(initialMessages))
```

### Conversation Continuity

Resume previous conversations:

```go
client.Call(ctx,
	ultravox.WithCallPriorCallId("previous-call-id"),
	ultravox.WithCallEnableGreetingPrompt(true),
)
```

## Authentication

The client uses the following environment variable for authentication:

- `ULTRAVOX_API_KEY`: Your Ultravox API key
