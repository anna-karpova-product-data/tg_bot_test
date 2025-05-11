package speech

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type DeepgramClient struct {
	apiKey string
	client *http.Client
}

type TranscriptionResponse struct {
	Results struct {
		Channels []struct {
			Alternatives []struct {
				Transcript string `json:"transcript"`
			} `json:"alternatives"`
		} `json:"channels"`
	} `json:"results"`
}

func NewDeepgramClient(apiKey string) *DeepgramClient {
	return &DeepgramClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (d *DeepgramClient) TranscribeAudio(audioPath string) (string, error) {
	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open the audio file
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}
	defer file.Close()

	// Create form file
	part, err := writer.CreateFormFile("audio", "audio.ogg")
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}

	// Copy the file content to the form
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file to form: %v", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.deepgram.com/v1/listen?language=ru&model=nova-2&smart_format=true&punctuate=true&detect_language=false&encoding=linear16&sample_rate=16000", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Token "+d.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Log the response for debugging
	log.Printf("Deepgram API Response: %s", string(respBody))

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deepgram API returned non-200 status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var transcription TranscriptionResponse
	if err := json.Unmarshal(respBody, &transcription); err != nil {
		return "", fmt.Errorf("error parsing response: %v, response body: %s", err, string(respBody))
	}

	// Extract the transcript
	if len(transcription.Results.Channels) > 0 && len(transcription.Results.Channels[0].Alternatives) > 0 {
		return transcription.Results.Channels[0].Alternatives[0].Transcript, nil
	}

	return "", fmt.Errorf("no transcript found in response: %s", string(respBody))
} 