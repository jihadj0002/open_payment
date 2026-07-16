package bkash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TokenManager struct {
	baseURL    string
	appKey     string
	appSecret  string
	username   string
	password   string
	httpClient *http.Client

	mu          sync.RWMutex
	idToken     string
	expiresAt   time.Time
	stopRefresh chan struct{}
}

type GrantTokenResponse struct {
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenManager(baseURL, appKey, appSecret, username, password string) *TokenManager {
	return &TokenManager{
		baseURL:    strings.TrimRight(baseURL, "/"),
		appKey:     appKey,
		appSecret:  appSecret,
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		stopRefresh: make(chan struct{}),
	}
}

func (tm *TokenManager) Start() error {
	if err := tm.grantToken(); err != nil {
		return fmt.Errorf("initial token grant: %w", err)
	}
	go tm.refreshLoop()
	return nil
}

func (tm *TokenManager) Stop() {
	close(tm.stopRefresh)
}

func (tm *TokenManager) GetToken() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.idToken
}

func (tm *TokenManager) grantToken() error {
	body := fmt.Sprintf(`{"app_key":"%s","app_secret":"%s"}`, tm.appKey, tm.appSecret)
	req, err := http.NewRequest(http.MethodPost, tm.baseURL+"/tokenized/checkout/token/grant", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("create grant request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("username", tm.username)
	req.Header.Set("password", tm.password)

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do grant request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read grant response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("grant token failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var tokenResp GrantTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return fmt.Errorf("parse grant response: %w", err)
	}

	tm.mu.Lock()
	tm.idToken = tokenResp.IDToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	tm.mu.Unlock()

	return nil
}

func (tm *TokenManager) refreshLoop() {
	ticker := time.NewTicker(50 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := tm.grantToken(); err != nil {
				continue
			}
		case <-tm.stopRefresh:
			return
		}
	}
}
