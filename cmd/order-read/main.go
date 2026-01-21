package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/client"
	"github.com/ueebee/tachibanashi/model"
	"github.com/ueebee/tachibanashi/request"
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

	ctx := context.Background()
	params := request.OrderParams{}

	issueCode := envFirst("TACHIBANASHI_ORDER_ISSUE_CODE", "TACHIBANASHI_ISSUE_CODE")
	if issueCode != "" {
		params["sIssueCode"] = issueCode
	}
	sikkouDay := strings.TrimSpace(os.Getenv("TACHIBANASHI_ORDER_SIKKOU_DAY"))
	if sikkouDay != "" {
		params["sSikkouDay"] = sikkouDay
	}
	status := strings.TrimSpace(os.Getenv("TACHIBANASHI_ORDER_STATUS"))
	if status != "" {
		params["sOrderSyoukaiStatus"] = status
	}

	resp, err := cli.Request().OrderList(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("orders")
	if len(resp.Entries) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, entry := range resp.Entries {
			order := entry.Order()
			side := entry.Fields.Value("sOrderBaibaiKubun")
			qty := valueOrNumber(entry.Fields, "sOrderOrderSuryou", int64(order.Quantity))
			price := valueOrNumber(entry.Fields, "sOrderOrderPrice", int64(order.Price))
			status := entry.Fields.Value("sOrderStatus")
			execDay := entry.Fields.Value("sOrderSikkouDay")
			orderTime := entry.Fields.Value("sOrderOrderDateTime")
			fmt.Printf("  %s %s side=%s qty=%s price=%s status=%s exec_day=%s order_time=%s\n",
				order.ID, order.Symbol, side, qty, price, status, execDay, orderTime)
		}
	}

	orderNumber := strings.TrimSpace(os.Getenv("TACHIBANASHI_ORDER_NUMBER"))
	eigyouDay := strings.TrimSpace(os.Getenv("TACHIBANASHI_EIGYOU_DAY"))
	if orderNumber == "" && eigyouDay == "" {
		return
	}
	if orderNumber == "" || eigyouDay == "" {
		log.Fatal("set both TACHIBANASHI_ORDER_NUMBER and TACHIBANASHI_EIGYOU_DAY to fetch detail")
	}

	detail, err := cli.Request().OrderListDetail(ctx, orderNumber, eigyouDay)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("order_detail")
	fmt.Printf("  order_number: %s\n", detail.OrderNumber)
	fmt.Printf("  eigyou_day: %s\n", detail.EigyouDay)
	fmt.Printf("  issue_code: %s\n", detail.IssueCode)
	printAttributesFiltered("  ", detail.Fields, commonAttrSkip)
}

func valueOrNumber(fields model.Attributes, key string, fallback int64) string {
	if fields != nil {
		if value := fields.Value(key); value != "" {
			return value
		}
	}
	if fallback != 0 {
		return fmt.Sprintf("%d", fallback)
	}
	return ""
}

func envFirst(names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
}

var commonAttrSkip = map[string]struct{}{
	"p_no":         {},
	"p_sd_date":    {},
	"p_rv_date":    {},
	"p_errno":      {},
	"p_err":        {},
	"sCLMID":       {},
	"sResultCode":  {},
	"sResultText":  {},
	"sWarningCode": {},
	"sWarningText": {},
}

func printAttributesFiltered(prefix string, attrs model.Attributes, skip map[string]struct{}) {
	if len(attrs) == 0 {
		fmt.Printf("%s(none)\n", prefix)
		return
	}
	keys := make([]string, 0, len(attrs))
	for key, value := range attrs {
		if value == "" {
			continue
		}
		if skip != nil {
			if _, ok := skip[key]; ok {
				continue
			}
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		fmt.Printf("%s(none)\n", prefix)
		return
	}
	for _, key := range keys {
		fmt.Printf("%s%s: %s\n", prefix, key, attrs[key])
	}
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
