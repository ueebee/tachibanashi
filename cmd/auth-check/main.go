package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/client"
)

func main() {
	_ = loadDotEnv(".env")

	loginID := mustEnvAny("TACHIBANASHI_LOGIN_ID", "TACHIBANA_USER_ID")
	password := mustEnvAny("TACHIBANASHI_PASSWORD", "TACHIBANA_PASSWORD")

	baseURL := envOrAny(client.BaseURLDemo, "TACHIBANASHI_BASE_URL", "TACHIBANA_BASE_URL")
	cfg := client.Config{BaseURL: baseURL}

	if timeout := os.Getenv("TACHIBANASHI_TIMEOUT"); timeout != "" {
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			log.Fatalf("invalid TACHIBANASHI_TIMEOUT: %v", err)
		}
		cfg.Timeout = dur
	}

	if ua := os.Getenv("TACHIBANASHI_USER_AGENT"); ua != "" {
		cfg.UserAgent = ua
	}

	if isTrue(os.Getenv("TACHIBANASHI_INSECURE_TLS")) {
		cfg.HTTPClient = &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.Auth().Login(context.Background(), auth.Credentials{
		LoginID:  loginID,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("login ok: result=%s msg=%s p_no=%s\n", resp.ResultCode, resp.ResultText, resp.PNo)
	fmt.Printf("urls: request=%s master=%s price=%s event=%s\n", resp.Request, resp.Master, resp.Price, resp.Event)

	if err := cli.Auth().Logout(context.Background()); err != nil {
		log.Printf("logout failed: %v", err)
		return
	}

	fmt.Println("logout ok")
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = trimQuotes(value)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

func trimQuotes(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func mustEnvAny(names ...string) string {
	for _, name := range names {
		value := strings.TrimSpace(os.Getenv(name))
		if value != "" {
			return value
		}
	}
	log.Fatalf("missing %s", strings.Join(names, " or "))
	return ""
}

func envOrAny(fallback string, names ...string) string {
	for _, name := range names {
		value := strings.TrimSpace(os.Getenv(name))
		if value != "" {
			return value
		}
	}
	return fallback
}

func isTrue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}
