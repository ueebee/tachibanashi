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
	"github.com/ueebee/tachibanashi/request"
)

func main() {
	_ = loadDotEnv(".env")

	loginID := mustEnvAny("TACHIBANASHI_LOGIN_ID", "TACHIBANA_USER_ID")
	password := mustEnvAny("TACHIBANASHI_PASSWORD", "TACHIBANA_PASSWORD")
	secondPassword := mustEnvAny("TACHIBANASHI_SECOND_PASSWORD", "TACHIBANASHI_ORDER_SECOND_PASSWORD")

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

	_, err = cli.Auth().Login(context.Background(), auth.Credentials{
		LoginID:  loginID,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := cli.Auth().Logout(context.Background()); err != nil {
			log.Printf("logout failed: %v", err)
		}
	}()

	issueCode := mustEnvAny("TACHIBANASHI_ORDER_ISSUE_CODE", "TACHIBANASHI_ISSUE_CODE")
	baibai, err := orderSideValue(envFirst("TACHIBANASHI_ORDER_BAIBAI_KUBUN", "TACHIBANASHI_ORDER_SIDE"))
	if err != nil {
		log.Fatal(err)
	}
	orderQty := mustEnvAny("TACHIBANASHI_ORDER_QTY")

	params := request.OrderParams{
		"sZyoutoekiKazeiC":        envOrDefault("TACHIBANASHI_ORDER_ZYOUTOEI_KAZEI", "1"),
		"sIssueCode":              issueCode,
		"sSizyouC":                envOrDefault("TACHIBANASHI_ORDER_SIZYOU_C", "00"),
		"sBaibaiKubun":            baibai,
		"sCondition":              envOrDefault("TACHIBANASHI_ORDER_CONDITION", "0"),
		"sOrderPrice":             envOrDefault("TACHIBANASHI_ORDER_PRICE", "0"),
		"sOrderSuryou":            orderQty,
		"sGenkinShinyouKubun":     envOrDefault("TACHIBANASHI_ORDER_GENKIN_SHINYOU", "0"),
		"sOrderExpireDay":         envOrDefault("TACHIBANASHI_ORDER_EXPIRE_DAY", "0"),
		"sGyakusasiOrderType":     envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_TYPE", "0"),
		"sGyakusasiZyouken":       envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_ZYOUKEN", "0"),
		"sGyakusasiPrice":         envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_PRICE", "*"),
		"sTatebiType":             envOrDefault("TACHIBANASHI_ORDER_TATEBI_TYPE", "*"),
		"sTategyokuZyoutoekiKazeiC": envOrDefault("TACHIBANASHI_ORDER_TATEGYOKU_KAZEI", "*"),
		"sSecondPassword":         secondPassword,
	}

	resp, err := cli.Request().KabuNewOrder(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("order_submit")
	fmt.Printf("  result: %s %s\n", resp.ResultCode, resp.ResultText)
	fmt.Printf("  warning: %s %s\n", resp.WarningCode, resp.WarningText)
	fmt.Printf("  order_number: %s\n", resp.OrderNumber)
	fmt.Printf("  eigyou_day: %s\n", resp.EigyouDay)
}

func orderSideValue(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("missing TACHIBANASHI_ORDER_BAIBAI_KUBUN or TACHIBANASHI_ORDER_SIDE")
	}
	switch strings.ToLower(value) {
	case "buy", "b", "kai":
		return "3", nil
	case "sell", "s", "uri":
		return "1", nil
	default:
		return value, nil
	}
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func envFirst(names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
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
