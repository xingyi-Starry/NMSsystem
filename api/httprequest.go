package api

import (
	"io"
	"net/http"
	"strings"
)

func HttpRequest(method string, url string, body string, header string) ([]byte, error) {
	client := &http.Client{}
	var data io.Reader = nil
	if method == "POST" {
		data = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, "http://localhost:1323/"+url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if header != "" {
		req.Header.Set("Authorization", "Bearer "+header)
	}
	resp, err := client.Do(req)
	if err != nil {
		// err = fmt.Errorf("err in Login: %w", err)
		err = NewNmsError("server crashed", ServerCrashErr)
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyText, nil
}
