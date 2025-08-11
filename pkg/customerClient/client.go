package customerClient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	baseURL string
	client  *fasthttp.Client
	timeout time.Duration
}

// New: baseURL ve timeout ile tek seferde client oluşturulur.
// Örn: New("http://localhost:8001", 5*time.Second)
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &fasthttp.Client{},
		timeout: timeout,
	}
}

func (c *Client) Get(path string, headers map[string]string, out interface{}) error {
	return c.doJSON(fasthttp.MethodGet, path, headers, nil, out)
}

func (c *Client) Post(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(fasthttp.MethodPost, path, headers, body, out)
}

func (c *Client) Put(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(fasthttp.MethodPut, path, headers, body, out)
}

func (c *Client) Patch(path string, headers map[string]string, body interface{}, out interface{}) error {
	return c.doJSON(fasthttp.MethodPatch, path, headers, body, out)
}

func (c *Client) Delete(path string, headers map[string]string) error {
	return c.doJSON(fasthttp.MethodDelete, path, headers, nil, nil)
}

// ---- helpers ----

func (c *Client) doJSON(method, path string, headers map[string]string, body interface{}, out interface{}) error {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	fullURL := c.baseURL + path
	req.SetRequestURI(fullURL)
	req.Header.SetMethod(method)

	// Varsayılan JSON içerik tipi – gerekirse headers ile override edilebilir
	req.Header.SetContentType("application/json")

	// Header’ları ekle
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Body varsa JSON’a çevir
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		req.SetBody(b)
	}

	// İstek
	if err := c.client.DoTimeout(req, res, c.timeout); err != nil {
		return fmt.Errorf("http error: %w", err)
	}

	status := res.StatusCode()
	switch method {
	case fasthttp.MethodPost:
		if status != fasthttp.StatusCreated && status != fasthttp.StatusOK {
			return fmt.Errorf("%s failed: %d", method, status)
		}
	case fasthttp.MethodPut, fasthttp.MethodPatch, fasthttp.MethodDelete, fasthttp.MethodGet:
		if status != fasthttp.StatusOK && status != fasthttp.StatusNoContent {
			return fmt.Errorf("%s failed: %d", method, status)
		}
	}

	// Body yoksa çık
	if out == nil || len(res.Body()) == 0 {
		return nil
	}

	if err := json.Unmarshal(res.Body(), out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
