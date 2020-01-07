package help

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// DoRequestObj Create a new request and auto marshal/demarshal the body and response
func DoRequestObj(ctx context.Context, method, URL string, body interface{}, response interface{}) (*http.Response, error) {
	bodyReader := &bytes.Buffer{}

	// Create body
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader.Write(bodyBytes)
	}

	// Make request
	req, err := http.NewRequest(method, URL, bodyReader)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Decode response
	if response != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(response)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
