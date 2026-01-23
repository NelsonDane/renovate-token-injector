// Nelson Dane
// Modified From: https://github.com/orgs/community/discussions/48186

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type githubTokenResponse struct {
	Token string `json:"token"`
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Fprintf(os.Stderr, "Missing %s\n", key)
		os.Exit(1)
	}
	return val
}

func loadPrivateKey(pemStr string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		fmt.Fprintln(os.Stderr, "Failed to parse PEM private key")
		os.Exit(1)
	}

	// Check if PKCS1
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key
	}

	// Try PKCS8
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse private key: %v\n", err)
		os.Exit(1)
	}

	// Assert it's RSA
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		fmt.Fprintln(os.Stderr, "Private key is not RSA")
		os.Exit(1)
	}

	return rsaKey
}

func main() {
	// Ensure needed env vars
	privateKeyPEM := mustEnv("PRIVATE_KEY")
	clientID := mustEnv("CLIENT_ID")
	installID := mustEnv("INSTALL_ID")
	outPath := os.Getenv("OUT_PATH")
	if outPath == "" {
		outPath = "/out/config.env"
	}

	// Create JWT from private key
	privateKey := loadPrivateKey(privateKeyPEM)
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": clientID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedJWT, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sign JWT: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created JWT")

	// Send Request for GitHub Token
	url := fmt.Sprintf(
		"https://api.github.com/app/installations/%s/access_tokens",
		installID,
	)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+signedJWT)
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
		os.Exit(1)
	}

	// Parse Response
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "Bad response (%d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}
	var tokenResp githubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	// Write to output file
	if err := os.WriteFile(
		outPath,
		[]byte("GH_TOKEN='"+tokenResp.Token+"'\n"),
		0644,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved GitHub token to %s\n", outPath)
}
