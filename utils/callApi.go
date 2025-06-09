package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"math/rand"
)

var (
	BaseURL1 = "https://cdn-api.swiftfederation.com"
	BaseURL2 = "https://base-api.swiftfederation.com"
)

func CallApi(method, url, uri string, requestBody interface{}, profileName string) ([]byte, error) {
	// Get credentials - use default profile if profileName is empty
	var accessKey, accessKeySecret string
	var err error
	
	if profileName == "" {
		// Try to get default profile
		accessKey, accessKeySecret, err = GetDefaultProfile()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve default profile: %w", err)
		}
	} else {
		accessKey, accessKeySecret, err = DisplayProfiles(profileName, true)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve profile %s: %w", profileName, err)
		}
	}

	if accessKey == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("no valid credentials found. Please add a profile using 'cdnctl config add'")
	}

	// Handle request body
	var requestBodyJson string
	switch v := requestBody.(type) {
	case string:
		requestBodyJson = v
	case []byte:
		requestBodyJson = string(v)
	case nil:
		requestBodyJson = "{}"
	default:
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBodyJson = string(requestBodyBytes)
	}

	// Generate authorization header
	authHeader, date, nonce, err := generateAuthorizeHeader(method, uri, requestBodyJson, accessKey, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authorization header: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(requestBodyJson)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("X-SFD-Date", date)
	req.Header.Set("X-SFD-Nonce", nonce)

	// Debug output
	fmt.Printf("Making %s request to: %s\n", method, url)
	fmt.Printf("Request body: %s\n", requestBodyJson)

	// Send the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	return body, nil
}

func generateAuthorizeHeader(method, uri, requestBodyJson, accessKey, accessKeySecret string) (string, string, string, error) {
	date := time.Now().UTC().Format("20060102T150405Z")
	nonce := fmt.Sprintf("%d", rand.Intn(90000)+10000)

	// Create signing string
	signingString := strings.Join([]string{
		method,
		uri,
		date,
		nonce,
		accessKey,
		requestBodyJson,
	}, "\n")

	// Generate HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(accessKeySecret))
	h.Write([]byte(signingString))
	signature := hex.EncodeToString(h.Sum(nil))

	authHeader := fmt.Sprintf("HMAC-SHA256 %s:%s", accessKey, signature)
	return authHeader, date, nonce, nil
}