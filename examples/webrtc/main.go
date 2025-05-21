// ultravox-webrtc demonstrates how to bridge Ultravox AI and WebRTC
package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/paulgrammer/ultravox"
	"github.com/paulgrammer/ultravox/examples/webrtc/web"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"github.com/zaf/g711"
	"gopkg.in/hraban/opus.v2"
)

const (
	// Server configuration
	ServerPort = "8080"

	// Audio sample rates
	OutputSampleRate = 8000 // 8kHz sampling rate
	InputSampleRate  = 8000 // 8kHz sampling rate

	// RTP parameters
	RTPPacketSize = 1500

	// WebRTC parameters
	ICETimeout = 30 * time.Second
)

// UltravoxConnection manages the connection to Ultravox API
type UltravoxConnection struct {
	wsConn     *websocket.Conn
	wsLock     sync.Mutex
	joinURL    string
	ctx        context.Context
	cancel     context.CancelFunc
	audioTrack *webrtc.TrackLocalStaticRTP

	// Client websocket connection (for sending events back to client)
	clientWs *websocket.Conn
}

// WebRTCConnection manages the WebRTC connection
type WebRTCConnection struct {
	peerConnection *webrtc.PeerConnection
	audioTrack     *webrtc.TrackLocalStaticRTP
	done           chan struct{}
}

// SDP message structure for exchanging offers and answers
type SDPMessage struct {
	Type webrtc.SDPType            `json:"type"`
	SDP  webrtc.SessionDescription `json:"sdp"`
}

// UltravoxEvent types
type (
	TranscriptEvent struct {
		Type  string `json:"type"`
		Role  string `json:"role"`
		Final bool   `json:"final"`
		Text  string `json:"text"`
		Delta string `json:"delta"`
	}

	ErrorEvent struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}

	StateEvent struct {
		Type  string `json:"type"`
		State string `json:"state"`
	}
)

// WebSocket upgrader for client connections
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// Global variable to track the active Ultravox connection
var activeUltravoxConnection *UltravoxConnection
var activeUltravoxLock sync.Mutex

func main() {
	// Create context with cancellation for handling shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Set up HTTP router
	router := mux.NewRouter()

	// Static file server for web assets
	webContent, err := fs.Sub(web.StaticFiles, ".")
	if err != nil {
		log.Fatalf("Failed to set up static file server: %v", err)
	}

	// Set up API routes
	router.HandleFunc("/api/sdp/offer", handleSDPOffer).Methods("POST")
	router.HandleFunc("/ws", handleWebSocketConnection)

	// Set up static file server
	staticFS, err := fs.Sub(webContent, "static")
	if err != nil {
		log.Fatalf("Failed to set up static file sub-filesystem: %v", err)
	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Serve index.html
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(webContent, "index.html")
		if err != nil {
			http.Error(w, "Failed to read index.html", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + ServerPort,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server started on http://localhost:%s", ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server shutdown complete")
}

// handleSDPOffer handles SDP offers from clients
func handleSDPOffer(w http.ResponseWriter, r *http.Request) {
	// Parse incoming SDP message
	var offerMsg SDPMessage
	if err := json.NewDecoder(r.Body).Decode(&offerMsg); err != nil {
		http.Error(w, "Failed to parse SDP offer", http.StatusBadRequest)
		return
	}

	// Setup WebRTC
	webrtcConn, err := setupWebRTC()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to setup WebRTC: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the remote SessionDescription
	if err = webrtcConn.peerConnection.SetRemoteDescription(offerMsg.SDP); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set remote description: %v", err), http.StatusInternalServerError)
		return
	}

	// Create answer
	answer, err := webrtcConn.peerConnection.CreateAnswer(nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create answer: %v", err), http.StatusInternalServerError)
		return
	}

	// Set local description
	if err = webrtcConn.peerConnection.SetLocalDescription(answer); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set local description: %v", err), http.StatusInternalServerError)
		return
	}

	// Wait for ICE gathering to complete
	gatherComplete := webrtc.GatheringCompletePromise(webrtcConn.peerConnection)
	<-gatherComplete

	// Create response
	responseMsg := SDPMessage{
		Type: webrtc.SDPTypeAnswer,
		SDP:  *webrtcConn.peerConnection.LocalDescription(),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseMsg); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleWebSocketConnection handles WebSocket connections from clients
func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	// Store the client WebSocket connection in the active Ultravox connection
	activeUltravoxLock.Lock()
	if activeUltravoxConnection != nil {
		activeUltravoxConnection.clientWs = conn
	}
	activeUltravoxLock.Unlock()

	// Simple ping-pong to keep connection alive
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from client WebSocket: %v", err)
			break
		}

		// Handle client messages (could be used for DTMF or other control messages)
		if messageType == websocket.TextMessage {
			log.Printf("Received client message: %s", string(message))
		}
	}

	// Remove client connection when it's closed
	activeUltravoxLock.Lock()
	if activeUltravoxConnection != nil {
		activeUltravoxConnection.clientWs = nil
	}
	activeUltravoxLock.Unlock()
}

// setupWebRTC initializes the WebRTC connection
func setupWebRTC() (*WebRTCConnection, error) {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Create a done channel to block until done
	done := make(chan struct{})

	// Create the WebRTC connection object
	webrtcConn := &WebRTCConnection{
		peerConnection: peerConnection,
		done:           done,
	}

	// Create a PCM audio track - using PCMU codec for G.711 µ-law
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMU, ClockRate: OutputSampleRate}, "audio", "ultravox-webrtc")
	if err != nil {
		return nil, fmt.Errorf("failed to create audio track: %w", err)
	}

	if _, err = peerConnection.AddTrack(audioTrack); err != nil {
		return nil, fmt.Errorf("failed to add audio track: %w", err)
	}
	webrtcConn.audioTrack = audioTrack

	// Setup peer connection handlers
	setupPeerConnectionHandlers(peerConnection, audioTrack, done)

	return webrtcConn, nil
}

// setupPeerConnectionHandlers sets up handlers for the WebRTC peer connection
func setupPeerConnectionHandlers(peerConnection *webrtc.PeerConnection, audioTrack *webrtc.TrackLocalStaticRTP, done chan struct{}) {
	// Handle ICE connection state changes
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			// Start Ultravox connection when WebRTC connects
			uvConn := &UltravoxConnection{
				audioTrack: audioTrack,
				wsLock:     sync.Mutex{},
			}
			setActiveUltravoxConnection(uvConn)
			go startUltravoxConnection(uvConn)
		} else if connectionState == webrtc.ICEConnectionStateDisconnected ||
			connectionState == webrtc.ICEConnectionStateFailed ||
			connectionState == webrtc.ICEConnectionStateClosed {
			close(done)
		}
	})

	// Handle incoming tracks (audio from browser)
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track has started, of type %d: %s", remoteTrack.PayloadType(), remoteTrack.Codec().MimeType)

		// Only handle audio tracks
		if strings.HasPrefix(remoteTrack.Codec().MimeType, "audio/") {
			log.Println("Got audio track from remote peer")
			go handleIncomingAudio(remoteTrack)
		}
	})
}

// handleIncomingAudio processes incoming audio from WebRTC
func handleIncomingAudio(track *webrtc.TrackRemote) {
	// Create appropriate audio decoder based on codec
	decoder, err := createAudioDecoder(track.Codec().MimeType)
	if err != nil {
		log.Printf("Failed to create audio decoder: %v", err)
		return
	}

	buffer := make([]byte, RTPPacketSize)
	rtpPacket := &rtp.Packet{}

	for {
		n, _, readErr := track.Read(buffer)
		if readErr != nil {
			log.Printf("Error reading from track: %v", readErr)
			return
		}

		if err := rtpPacket.Unmarshal(buffer[:n]); err != nil {
			log.Printf("Error unmarshalling RTP packet: %v", err)
			continue
		}

		// Process the packet based on codec
		pcmData, err := processAudioPacket(rtpPacket.Payload, track.Codec().MimeType, decoder)
		if err != nil {
			log.Printf("Error processing audio packet: %v", err)
			continue
		}

		// Find the active Ultravox connection
		activeUVConn := findActiveUltravoxConnection()
		if activeUVConn != nil && activeUVConn.wsConn != nil {
			activeUVConn.wsLock.Lock()
			if err := activeUVConn.wsConn.WriteMessage(websocket.BinaryMessage, pcmData); err != nil {
				log.Printf("Error sending audio to Ultravox: %v", err)
			}
			activeUVConn.wsLock.Unlock()
		}
	}
}

// setActiveUltravoxConnection sets the active Ultravox connection
func setActiveUltravoxConnection(conn *UltravoxConnection) {
	activeUltravoxLock.Lock()
	defer activeUltravoxLock.Unlock()
	activeUltravoxConnection = conn
}

// findActiveUltravoxConnection returns the active Ultravox connection
func findActiveUltravoxConnection() *UltravoxConnection {
	activeUltravoxLock.Lock()
	defer activeUltravoxLock.Unlock()
	return activeUltravoxConnection
}

// createAudioDecoder creates an appropriate decoder based on codec type
func createAudioDecoder(mimeType string) (interface{}, error) {
	switch mimeType {
	case webrtc.MimeTypeOpus:
		return opus.NewDecoder(InputSampleRate, 1) // 8kHz sample rate, mono
	case webrtc.MimeTypePCMA, webrtc.MimeTypePCMU:
		// g711 library handles both µ-law and A-law
		return nil, nil // No stateful decoder needed for G.711
	default:
		return nil, fmt.Errorf("unsupported codec: %s", mimeType)
	}
}

// processAudioPacket converts audio data based on codec type
func processAudioPacket(payload []byte, mimeType string, decoder interface{}) ([]byte, error) {
	switch mimeType {
	case webrtc.MimeTypePCMA:
		// Convert A-law to PCM
		pcmData := make([]byte, len(payload)*2)
		for i, sample := range payload {
			pcmSample := g711.DecodeAlawFrame(sample)
			binary.LittleEndian.PutUint16(pcmData[i*2:], uint16(pcmSample))
		}
		return pcmData, nil

	case webrtc.MimeTypePCMU:
		// Convert µ-law to PCM
		pcmData := make([]byte, len(payload)*2)
		for i, sample := range payload {
			pcmSample := g711.DecodeUlawFrame(sample)
			binary.LittleEndian.PutUint16(pcmData[i*2:], uint16(pcmSample))
		}
		return pcmData, nil

	case webrtc.MimeTypeOpus:
		// Decode Opus to PCM
		opusDecoder, ok := decoder.(*opus.Decoder)
		if !ok {
			return nil, fmt.Errorf("invalid opus decoder")
		}

		// Max frame size for Opus
		pcm := make([]int16, 5760)
		n, err := opusDecoder.Decode(payload, pcm)
		if err != nil {
			return nil, err
		}

		// Convert to byte array
		pcmData := make([]byte, n*2)
		for i := 0; i < n; i++ {
			binary.LittleEndian.PutUint16(pcmData[i*2:], uint16(pcm[i]))
		}
		return pcmData, nil

	default:
		return nil, fmt.Errorf("unsupported codec: %s", mimeType)
	}
}

// startUltravoxConnection initializes and manages the Ultravox connection
func startUltravoxConnection(uvConn *UltravoxConnection) {
	// Create a new Ultravox client
	uv := ultravox.NewClient()

	// Configure Ultravox call options
	call, err := configureAndStartUltravoxCall(uv)
	if err != nil {
		log.Fatalf("Failed to start Ultravox call: %v", err)
		return
	}

	// Log call information
	logCallInfo(call)

	// Create context with cancellation for the WebSocket connection
	uvConn.ctx, uvConn.cancel = context.WithCancel(context.Background())
	defer uvConn.cancel()

	// Connect to Ultravox WebSocket
	uvConn.joinURL = call.JoinURL
	handleUltravoxWebSocket(uvConn)
}

// configureAndStartUltravoxCall configures and starts a call with Ultravox
func configureAndStartUltravoxCall(uv *ultravox.Client) (*ultravox.Call, error) {
	// Configure first speaker settings
	firstSpeakerSettings := ultravox.AgentFirstSpeaker(
		false,                                // Not uninterruptible
		"Hello! How can I assist you today?", // Text
		"",                                   // No prompt (using text directly)
		0,                                    // No delay
	)

	// Configure VAD settings
	vadSettings := ultravox.NewVadSettings()
	vadSettings.TurnEndpointDelay = ultravox.UltravoxDuration(400 * time.Millisecond)

	// Set up inactivity messages
	inactivityMessages := []ultravox.TimedMessage{
		ultravox.NewTimedMessage(5*time.Second, "Are you still there? I'm here to help if you need anything.", ultravox.EndBehaviorDefault),
		ultravox.NewTimedMessage(15*time.Second, "I'll wait a bit longer in case you want to continue our conversation.", ultravox.EndBehaviorDefault),
		ultravox.NewTimedMessage(20*time.Second, "Since I haven't heard from you, I'll be ending our call now. Feel free to call back anytime if you need assistance!", ultravox.EndBehaviorHangUpSoft),
	}

	// Start new call with options
	return uv.Call(
		context.Background(),
		ultravox.WithCallSystemPrompt("You are a helpful assistant. Provide concise, helpful information to user queries. Be warm and friendly but brief in your responses."),
		ultravox.WithCallMaxDuration(5*time.Minute),
		ultravox.WithCallFirstSpeakerSettings(firstSpeakerSettings),
		ultravox.WithCallVadSettings(vadSettings),
		ultravox.WithCallInactivityMessages(inactivityMessages),
		ultravox.WithCallRecordingEnabled(false),
		ultravox.WithCallWebSocketMedium(InputSampleRate, OutputSampleRate),
	)
}

// logCallInfo logs information about the Ultravox call
func logCallInfo(call *ultravox.Call) {
	log.Printf("Call created successfully!")
	log.Printf("Call ID: %s", call.CallID)
	log.Printf("Join URL: %s", call.JoinURL)
	log.Printf("Max Duration: %s", call.MaxDuration.String())
	log.Printf("Join Timeout: %s", call.JoinTimeout.String())
}

// handleUltravoxWebSocket manages the WebSocket connection to Ultravox
func handleUltravoxWebSocket(uvConn *UltravoxConnection) {
	var err error
	uvConn.wsConn, _, err = websocket.DefaultDialer.Dial(uvConn.joinURL, nil)
	if err != nil {
		log.Fatalf("WebSocket connection error: %v", err)
	}
	defer uvConn.wsConn.Close()

	// Set up audio parameters
	sequenceNumber := uint16(0)
	timestamp := uint32(0)
	ssrc := uint32(12345) // Consistent SSRC identifier

	for {
		select {
		case <-uvConn.ctx.Done():
			return
		default:
			messageType, message, err := uvConn.wsConn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			switch messageType {
			case websocket.TextMessage:
				// Handle JSON messages from Ultravox
				handleUltravoxJsonMessage(uvConn, message)
			case websocket.BinaryMessage:
				// Process audio data from Ultravox and send to WebRTC
				processUltravoxAudio(uvConn, message, &sequenceNumber, &timestamp, ssrc)
			default:
				log.Printf("Received unexpected WebSocket message type: %d", messageType)
			}
		}
	}
}

// processUltravoxAudio processes audio data from Ultravox and sends it to WebRTC
func processUltravoxAudio(uvConn *UltravoxConnection, pcmData []byte, sequenceNumber *uint16, timestamp *uint32, ssrc uint32) {
	// Convert from PCM 16-bit to PCMU (G.711 µ-law) using g711 library
	muLawData := make([]byte, len(pcmData)/2)
	for i := 0; i < len(pcmData)/2; i++ {
		// Read 16-bit PCM sample (little-endian)
		sample := int16(binary.LittleEndian.Uint16(pcmData[i*2:]))
		// Convert to µ-law
		muLawData[i] = g711.EncodeUlawFrame(sample)
	}

	// Calculate timestamp increment (for 8kHz audio)
	tsIncrement := uint32(len(muLawData))

	// Create RTP packet
	packet := &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0, // 0 = PCMU (G.711 µ-law)
			SequenceNumber: *sequenceNumber,
			Timestamp:      *timestamp,
			SSRC:           ssrc,
		},
		Payload: muLawData,
	}

	// Update sequence number and timestamp
	*sequenceNumber++
	*timestamp += tsIncrement

	// Marshal and send the packet
	raw, err := packet.Marshal()
	if err != nil {
		log.Printf("Failed to marshal RTP packet: %v", err)
		return
	}

	if _, err := uvConn.audioTrack.Write(raw); err != nil {
		log.Printf("Failed to write to track: %v", err)
	}
}

// handleUltravoxJsonMessage processes JSON messages from Ultravox and forwards them to the client
func handleUltravoxJsonMessage(uvConn *UltravoxConnection, message []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		log.Println(string(message))
		return
	}

	eventType, ok := event["type"].(string)
	if !ok {
		log.Printf("Unknown JSON message: %s", string(message))
		return
	}

	// Forward the event to the client if the WebSocket connection is established
	if uvConn.clientWs != nil {
		if err := uvConn.clientWs.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error forwarding event to client: %v", err)
		}
	}

	// Process the event locally
	switch eventType {
	case "transcript":
		var transcriptEvent TranscriptEvent
		if err := json.Unmarshal(message, &transcriptEvent); err != nil {
			log.Printf("Error parsing transcript event: %v", err)
			return
		}

		if transcriptEvent.Final {
			log.Printf("Transcript [%s]: %s", transcriptEvent.Role, transcriptEvent.Text)
		}

	case "error":
		var errorEvent ErrorEvent
		if err := json.Unmarshal(message, &errorEvent); err != nil {
			log.Printf("Error parsing error event: %v", err)
			return
		}
		log.Printf("Ultravox Error: %s", errorEvent.Error)

	case "state":
		var stateEvent StateEvent
		if err := json.Unmarshal(message, &stateEvent); err != nil {
			log.Printf("Error parsing state event: %v", err)
			return
		}
		log.Printf("Ultravox State: %s", stateEvent.State)

	default:
		log.Printf("Received unknown event type: %s", eventType)
		log.Println(string(message))
	}
}
