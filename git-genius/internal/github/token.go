package github

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"git-genius/internal/system"
)

const (
	geniusDir = ".git/.genius"
	tokenFile = geniusDir + "/token"
	apiURL    = "https://api.github.com/user"
)

type userResponse struct {
	Login string `json:"login"`
}

// -------------------- TOKEN FILE --------------------

// Get reads the stored GitHub token
func Get() string {
	data, _ := os.ReadFile(tokenFile)
	return string(data)
}

// Save stores the GitHub token securely
func Save(token string) error {
	if token == "" {
		return errors.New("empty token")
	}
	os.MkdirAll(geniusDir, 0700)
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// Delete removes the stored token (used when invalid)
func Delete() {
	_ = os.Remove(tokenFile)
}

// -------------------- VALIDATION --------------------

// Validate checks token validity using GitHub API
// Returns GitHub username if valid
func Validate() (string, error) {
	token := Get()
	if token == "" {
		return "", errors.New("no GitHub token found")
	}

	// Offline mode → skip validation
	if !system.Online {
		return "offline-mode", nil
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("User-Agent", "git-genius")

	resp, err := client.Do(req)
	if err != nil {
		system.LogError("github api request failed", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid or expired GitHub token")
	}

	var user userResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	if user.Login == "" {
		return "", errors.New("unable to read github username")
	}

	return user.Login, nil
}
