// File: pkg/customerClient/client.go
package customerClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

// New: baseURL ve timeout ile net/http client oluşturur.
// Örn: New("http://localhost:8001", 5*time.Second)
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *Client) Get(path string, headers map[string]string, out interface{}) error {
	return c.doJSON(http.MethodGet, path, headers, nil, out)
}

func (c *Client) Post(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(http.MethodPost, path, headers, body, out)
}

func (c *Client) Put(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(http.MethodPut, path, headers, body, out)
}

func (c *Client) Patch(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(http.MethodPatch, path, headers, body, out)
}

func (c *Client) Delete(path string, headers map[string]string) error {
	return c.doJSON(http.MethodDelete, path, headers, nil, nil)
}

// ---- helpers ----

func (c *Client) doJSON(method, path string, headers map[string]string, body interface{}, out interface{}) error {
	fullURL := c.baseURL + path

	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		r = bytes.NewReader(b)
	}

	// İstersen burada context’li versiyon kullanabilirsin: http.NewRequestWithContext(ctx, ...)
	req, err := http.NewRequest(method, fullURL, r)
	if err != nil {
		return err
	}

	// Varsayılan JSON içerik tipi – header’lar override edebilir
	contentTypeSet := false
	for k, v := range headers {
		req.Header.Set(k, v)
		if http.CanonicalHeaderKey(k) == "Content-Type" {
			contentTypeSet = true
		}
	}
	if !contentTypeSet {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	// Status kontrolü (POST: 200/201; diğerleri: 200/204)
	switch method {
	case http.MethodPost:
		if status != http.StatusOK && status != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("%s failed: %d - %s", method, status, string(bodyBytes))
		}
	default:
		if status != http.StatusOK && status != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("%s failed: %d - %s", method, status, string(bodyBytes))
		}
	}

	// 204 No Content veya out == nil ise decode etme
	if status == http.StatusNoContent || out == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
