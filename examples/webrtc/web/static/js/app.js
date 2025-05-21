// Ultravox WebRTC Demo - Main Application Logic

document.addEventListener("DOMContentLoaded", () => {
  // DOM Elements
  const callButton = document.getElementById("call-button");
  const hangupButton = document.getElementById("hangup-button");
  const connectionStatus = document.getElementById("connection-status");
  const connectionStatusIndicator = document.getElementById(
    "connection-status-indicator"
  );
  const remoteAudio = document.getElementById("remote-audio");
  const transcriptsContainer = document.getElementById("transcripts");
  const logsContainer = document.getElementById("logs");
  const microphoneSelect = document.getElementById("microphone-select");
  const speakerSelect = document.getElementById("speaker-select");
  const volumeControl = document.getElementById("volume-control");
  const muteToggle = document.getElementById("mute-toggle");

  // WebRTC Configuration
  const configuration = {
    iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  };

  // State variables
  let peerConnection = null;
  let localStream = null;
  let websocket = null;
  let isCallActive = false;

  // Track active transcripts by role
  const activeTranscripts = {
    user: null,
    assistant: null,
  };

  // Initialize audio device selection
  initAudioDevices();

  // Set up event listeners
  callButton.addEventListener("click", startCall);
  hangupButton.addEventListener("click", endCall);
  volumeControl.addEventListener("input", updateVolume);
  muteToggle.addEventListener("change", toggleMute);
  microphoneSelect.addEventListener("change", changeMicrophone);
  speakerSelect.addEventListener("change", changeSpeaker);

  // Initialize audio devices
  async function initAudioDevices() {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();

      // Filter for audio input devices (microphones)
      const microphones = devices.filter(
        (device) => device.kind === "audioinput"
      );
      microphones.forEach((mic) => {
        const option = document.createElement("option");
        option.value = mic.deviceId;
        option.text = mic.label || `Microphone ${microphoneSelect.length + 1}`;
        microphoneSelect.appendChild(option);
      });

      // Filter for audio output devices (speakers)
      const speakers = devices.filter(
        (device) => device.kind === "audiooutput"
      );
      speakers.forEach((speaker) => {
        const option = document.createElement("option");
        option.value = speaker.deviceId;
        option.text = speaker.label || `Speaker ${speakerSelect.length + 1}`;
        speakerSelect.appendChild(option);
      });

      // If no devices found, show some defaults
      if (microphones.length === 0) {
        const option = document.createElement("option");
        option.text = "Default Microphone";
        microphoneSelect.appendChild(option);
      }

      if (speakers.length === 0) {
        const option = document.createElement("option");
        option.text = "Default Speaker";
        speakerSelect.appendChild(option);
      }
    } catch (error) {
      logMessage(`Error getting audio devices: ${error.message}`, "error");
    }
  }

  // Start a call
  async function startCall() {
    try {
      // Update UI
      updateConnectionStatus("connecting");
      callButton.disabled = true;

      // Get user media (microphone)
      const constraints = {
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
        video: false,
      };

      // If a specific microphone is selected, use it
      if (microphoneSelect.value) {
        constraints.audio.deviceId = { exact: microphoneSelect.value };
      }

      localStream = await navigator.mediaDevices.getUserMedia(constraints);

      // Create WebRTC peer connection
      peerConnection = new RTCPeerConnection(configuration);

      // Add local stream tracks to the peer connection
      localStream.getTracks().forEach((track) => {
        peerConnection.addTrack(track, localStream);
      });

      // Set up ICE candidate handling
      peerConnection.onicecandidate = (event) => {
        if (event.candidate === null) {
          logMessage("ICE gathering complete", "info");
        }
      };

      // Handle ICE connection state changes
      peerConnection.oniceconnectionstatechange = () => {
        logMessage(
          `ICE connection state: ${peerConnection.iceConnectionState}`,
          "info"
        );

        if (
          peerConnection.iceConnectionState === "connected" ||
          peerConnection.iceConnectionState === "completed"
        ) {
          updateConnectionStatus("connected");
          callButton.classList.remove("pulse-animation");
          callButton.style.display = "none";
          hangupButton.style.display = "flex";
          isCallActive = true;
        } else if (
          peerConnection.iceConnectionState === "failed" ||
          peerConnection.iceConnectionState === "disconnected" ||
          peerConnection.iceConnectionState === "closed"
        ) {
          updateConnectionStatus("error");
          if (isCallActive) {
            endCall();
          } else {
            resetUI();
          }
        }
      };

      // Handle incoming audio stream
      peerConnection.ontrack = (event) => {
        logMessage("Remote track received", "success");
        remoteAudio.srcObject = event.streams[0];

        // Apply the current volume setting
        remoteAudio.volume = parseFloat(volumeControl.value);

        // If a specific speaker is selected and the browser supports it, use it
        if (
          speakerSelect.value &&
          typeof remoteAudio.setSinkId === "function"
        ) {
          try {
            remoteAudio.setSinkId(speakerSelect.value);
          } catch (error) {
            logMessage(
              `Error setting audio output device: ${error.message}`,
              "error"
            );
          }
        }
      };

      // Create offer
      const offer = await peerConnection.createOffer();
      await peerConnection.setLocalDescription(offer);

      // Wait for ICE gathering to complete
      await waitForIceGatheringComplete(peerConnection);

      // Send offer to server
      const response = await fetch("/api/sdp/offer", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          type: "offer",
          sdp: peerConnection.localDescription,
        }),
      });

      if (!response.ok) {
        throw new Error(
          `Server responded with ${response.status}: ${response.statusText}`
        );
      }

      // Get answer from server
      const answerData = await response.json();
      await peerConnection.setRemoteDescription(
        new RTCSessionDescription(answerData.sdp)
      );

      // Connect to WebSocket for events
      connectWebSocket();

      logMessage("Call setup complete", "success");
    } catch (error) {
      logMessage(`Error starting call: ${error.message}`, "error");
      updateConnectionStatus("error");
      resetUI();

      // Clean up if there was an error
      if (localStream) {
        localStream.getTracks().forEach((track) => track.stop());
        localStream = null;
      }

      if (peerConnection) {
        peerConnection.close();
        peerConnection = null;
      }
    }
  }

  // End the call
  function endCall() {
    logMessage("Ending call", "info");

    // Close WebSocket connection
    if (websocket) {
      websocket.close();
      websocket = null;
    }

    // Stop all tracks from local stream
    if (localStream) {
      localStream.getTracks().forEach((track) => track.stop());
      localStream = null;
    }

    // Close peer connection
    if (peerConnection) {
      peerConnection.close();
      peerConnection = null;
    }

    // Reset UI
    resetUI();
    isCallActive = false;

    // Clear active transcripts
    activeTranscripts.user = null;
    activeTranscripts.assistant = null;
  }

  // Reset UI to initial state
  function resetUI() {
    callButton.disabled = false;
    callButton.style.display = "flex";
    hangupButton.style.display = "none";
    updateConnectionStatus("disconnected");
  }

  // Update connection status indicator
  function updateConnectionStatus(status) {
    connectionStatusIndicator.className = "w-3 h-3 rounded-full mr-2";

    switch (status) {
      case "disconnected":
        connectionStatusIndicator.classList.add("status-disconnected");
        connectionStatus.textContent = "Disconnected";
        break;
      case "connecting":
        connectionStatusIndicator.classList.add("status-connecting");
        connectionStatus.textContent = "Connecting...";
        break;
      case "connected":
        connectionStatusIndicator.classList.add("status-connected");
        connectionStatus.textContent = "Connected";
        break;
      case "error":
        connectionStatusIndicator.classList.add("status-error");
        connectionStatus.textContent = "Error";
        break;
    }
  }

  // Connect to the WebSocket server for events
  function connectWebSocket() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    websocket = new WebSocket(wsUrl);

    websocket.onopen = () => {
      logMessage("WebSocket connected", "success");
    };

    websocket.onclose = () => {
      logMessage("WebSocket disconnected", "info");
    };

    websocket.onerror = (error) => {
      logMessage(`WebSocket error: ${error}`, "error");
    };

    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        handleUltravoxEvent(data);
      } catch (error) {
        logMessage(
          `Error parsing WebSocket message: ${error.message}`,
          "error"
        );
      }
    };
  }

  // Handle events from Ultravox
  function handleUltravoxEvent(event) {
    switch (event.type) {
      case "transcript":
        addTranscript(event);
        break;
      case "state":
        logMessage(`Call state: ${event.state}`, "info");
        break;
      case "error":
        logMessage(`Call error: ${event.error}`, "error");
        break;
      default:
        logMessage(`Received event: ${event.type}`, "info");
    }
  }

  // Add a transcript message to the UI
  function addTranscript(event) {
    // Only add final transcripts to avoid cluttering the UI
    if (!event.final) {
      return;
    }

    const transcriptDiv = document.createElement("div");
    transcriptDiv.className = "transcript-bubble fade-in";

    if (event.role === "user") {
      transcriptDiv.classList.add("transcript-user");
    } else {
      transcriptDiv.classList.add("transcript-assistant");
    }

    transcriptDiv.textContent = event.text;
    transcriptsContainer.appendChild(transcriptDiv);

    // Scroll to bottom
    transcriptsContainer.scrollTop = transcriptsContainer.scrollHeight;
  }

  // Log a message to the UI
  function logMessage(message, level = "info") {
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement("div");
    logEntry.className = `log-entry log-${level}`;
    logEntry.textContent = `[${timestamp}] ${message}`;

    logsContainer.appendChild(logEntry);
    logsContainer.scrollTop = logsContainer.scrollHeight;
  }

  // Update the volume of the remote audio
  function updateVolume() {
    if (remoteAudio) {
      remoteAudio.volume = parseFloat(volumeControl.value);
    }
  }

  // Toggle mute state of local audio tracks
  function toggleMute() {
    if (localStream) {
      const audioTracks = localStream.getAudioTracks();
      audioTracks.forEach((track) => {
        track.enabled = !muteToggle.checked;
      });

      logMessage(
        muteToggle.checked ? "Microphone muted" : "Microphone unmuted",
        "info"
      );
    }
  }

  // Change microphone device
  async function changeMicrophone() {
    if (!localStream || !peerConnection) {
      return; // Not in a call
    }

    try {
      // Stop current tracks
      localStream.getAudioTracks().forEach((track) => {
        track.stop();
        localStream.removeTrack(track);
        peerConnection.getSenders().forEach((sender) => {
          if (sender.track === track) {
            peerConnection.removeTrack(sender);
          }
        });
      });

      // Get new track with selected device
      const constraints = {
        audio: {
          deviceId: { exact: microphoneSelect.value },
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
        video: false,
      };

      const newStream = await navigator.mediaDevices.getUserMedia(constraints);
      const newTrack = newStream.getAudioTracks()[0];

      // Add new track to stream and peer connection
      localStream.addTrack(newTrack);
      peerConnection.addTrack(newTrack, localStream);

      // Apply mute state
      newTrack.enabled = !muteToggle.checked;

      logMessage(`Microphone changed to: ${newTrack.label}`, "success");
    } catch (error) {
      logMessage(`Error changing microphone: ${error.message}`, "error");
    }
  }

  // Change speaker device
  function changeSpeaker() {
    if (remoteAudio && typeof remoteAudio.setSinkId === "function") {
      try {
        remoteAudio
          .setSinkId(speakerSelect.value)
          .then(() => {
            logMessage(
              `Speaker changed to: ${speakerSelect.selectedOptions[0].text}`,
              "success"
            );
          })
          .catch((error) => {
            logMessage(`Error changing speaker: ${error.message}`, "error");
          });
      } catch (error) {
        logMessage(`Error changing speaker: ${error.message}`, "error");
      }
    } else {
      logMessage(
        "Your browser does not support output device selection",
        "warning"
      );
    }
  }

  // Helper function to wait for ICE gathering to complete
  function waitForIceGatheringComplete(pc) {
    return new Promise((resolve) => {
      if (pc.iceGatheringState === "complete") {
        resolve();
        return;
      }

      const checkState = () => {
        if (pc.iceGatheringState === "complete") {
          pc.removeEventListener("icegatheringstatechange", checkState);
          resolve();
        }
      };

      pc.addEventListener("icegatheringstatechange", checkState);

      // Timeout after 5 seconds just in case
      setTimeout(resolve, 5000);
    });
  }
});
