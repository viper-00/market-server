package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	URL     string
	Headers map[string]string
	client  http.Client
	Timeout time.Duration
}

func (c *Client) HTTPGetUnique(dest interface{}) (err error) {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create get request: %v", err.Error())
	}

	if len(c.Headers) > 0 {
		for key, value := range c.Headers {
			req.Header.Set(key, value)
		}
	}

	// setup timeout
	c.client.Timeout = c.Timeout

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect the %s: %s", c.URL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to connect the %s: %d", c.URL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to readall the body of response: %s", err.Error())
	}

	switch v := dest.(type) {
	case *string:
		*v = string(data)
	case *int:
		intValue, err := strconv.Atoi(string(data))
		if err != nil {
			return fmt.Errorf("failed to convert response body to int: %s", err.Error())
		}
		*v = intValue
	case *int64:
		int64Value, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert response body to int: %s", err.Error())
		}
		*v = int64Value
	default:
		return fmt.Errorf("unsupported type: %T", dest)
	}

	return nil
}

func (c *Client) HTTPGet(dest interface{}) (err error) {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create get request: %v", err.Error())
	}

	if len(c.Headers) > 0 {
		for key, value := range c.Headers {
			req.Header.Set(key, value)
		}
	}

	// default header
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// setup timeout
	c.client.Timeout = c.Timeout

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect the %s: %s", c.URL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to connect the %s: %d", c.URL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to readall the body of response: %s", err.Error())
	}

	err = json.Unmarshal([]byte(data), &dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the body of response: %s", err.Error())
	}

	return nil
}

func (c *Client) HTTPPost(source interface{}, dest interface{}) (err error) {
	body, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to encode source: %s", err.Error())
	}

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create post request: %s", err.Error())
	}

	if len(c.Headers) > 0 {
		for key, value := range c.Headers {
			req.Header.Set(key, value)
		}
	}

	// default header
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// setup timeout
	c.client.Timeout = c.Timeout

	response, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect the %s: %s", c.URL, err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to connect the %s: %d", c.URL, response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&dest)
	if err != nil {
		return fmt.Errorf("failed to decode the body of response: %s", err.Error())
	}

	return nil
}
