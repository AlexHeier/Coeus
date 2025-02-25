package coeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func AzureTTS(model, endpoint, apikey, apiversion string) error {

	if model == "" || endpoint == "" || apikey == "" || apiversion == "" {
		return fmt.Errorf("azuretts: all fields need to be filled")
	}

	TTSProvider = AzureTTSProvider{
		Endpoint:   endpoint,
		APIKey:     apikey,
		APIVersion: apiversion,
		Model:      model,
	}
	return nil
}

func AzureSendTTS(text, voice, format string) (base64Audio []byte, err error) {

	config := TTSProvider.(AzureTTSProvider)

	azureRes := azureTTSRequest{
		Model: config.Model,
		Audio: struct {
			Voice  string "json:\"voice\""
			Format string "json:\"format\""
		}{Voice: voice, Format: format},
		Modalities: []string{"text", "audio"},
	}

	// Appends the system message to request
	azureRes.Messages = append(azureRes.Messages,
		azureTTSMessage{
			Role: "system",
			Content: struct {
				Type       string "json:\"type\""
				Text       string "json:\"text,omitempty\""
				InputAudio struct {
					Data   string "json:\"data\""
					Format string "json:\"format\""
				} "json:\"input_audio,omitempty\""
			}{Type: "text", Text: "Convert user text to speech in the same language as the user."}},
		azureTTSMessage{
			Role: "user",
			Content: struct {
				Type       string "json:\"type\""
				Text       string "json:\"text,omitempty\""
				InputAudio struct {
					Data   string "json:\"data\""
					Format string "json:\"format\""
				} "json:\"input_audio,omitempty\""
			}{Type: "text", Text: text}})

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(azureRes)

	req, err := http.NewRequest(http.MethodPost, config.Endpoint, buf)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ttsRes azureTTSResponse

	err = json.Unmarshal(data, &ttsRes)
	if err != nil {
		return nil, err
	}

	if len(ttsRes.Choices) > 0 {
		return []byte(ttsRes.Choices[0].Message.Audio.Data), nil
	}

	return nil, fmt.Errorf("error function not implemented")
}
