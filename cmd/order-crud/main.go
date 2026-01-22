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

	ctx := context.Background()
	issueCode := mustEnvAny("TACHIBANASHI_ORDER_ISSUE_CODE", "TACHIBANASHI_ISSUE_CODE")
	baibai, err := orderSideValue(envFirst("TACHIBANASHI_ORDER_BAIBAI_KUBUN", "TACHIBANASHI_ORDER_SIDE"))
	if err != nil {
		log.Fatal(err)
	}
	orderQty := mustEnvAny("TACHIBANASHI_ORDER_QTY")

	newParams := request.OrderParams{
		"sZyoutoekiKazeiC":          envOrDefault("TACHIBANASHI_ORDER_ZYOUTOEI_KAZEI", "1"),
		"sIssueCode":                issueCode,
		"sSizyouC":                  envOrDefault("TACHIBANASHI_ORDER_SIZYOU_C", "00"),
		"sBaibaiKubun":              baibai,
		"sCondition":                envOrDefault("TACHIBANASHI_ORDER_CONDITION", "0"),
		"sOrderPrice":               envOrDefault("TACHIBANASHI_ORDER_PRICE", "0"),
		"sOrderSuryou":              orderQty,
		"sGenkinShinyouKubun":       envOrDefault("TACHIBANASHI_ORDER_GENKIN_SHINYOU", "0"),
		"sOrderExpireDay":           envOrDefault("TACHIBANASHI_ORDER_EXPIRE_DAY", "0"),
		"sGyakusasiOrderType":       envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_TYPE", "0"),
		"sGyakusasiZyouken":         envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_ZYOUKEN", "0"),
		"sGyakusasiPrice":           envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_PRICE", "*"),
		"sTatebiType":               envOrDefault("TACHIBANASHI_ORDER_TATEBI_TYPE", "*"),
		"sTategyokuZyoutoekiKazeiC": envOrDefault("TACHIBANASHI_ORDER_TATEGYOKU_KAZEI", "*"),
		"sSecondPassword":           secondPassword,
	}

	newResp, err := cli.Request().KabuNewOrder(ctx, newParams)
	if err != nil {
		log.Fatal(err)
	}
	printOrderResponse("order_submit", newResp)

	orderNumber := strings.TrimSpace(newResp.OrderNumber)
	eigyouDay := strings.TrimSpace(newResp.EigyouDay)
	if orderNumber == "" {
		orderNumber = strings.TrimSpace(os.Getenv("TACHIBANASHI_ORDER_NUMBER"))
	}
	if eigyouDay == "" {
		eigyouDay = strings.TrimSpace(os.Getenv("TACHIBANASHI_EIGYOU_DAY"))
	}
	if orderNumber == "" || eigyouDay == "" {
		log.Fatal("order number/eigyou day not found; set TACHIBANASHI_ORDER_NUMBER and TACHIBANASHI_EIGYOU_DAY to continue")
	}

	detail, err := cli.Request().OrderListDetail(ctx, orderNumber, eigyouDay)
	if err != nil {
		log.Fatal(err)
	}
	printOrderDetail(detail)

	correctParams := request.OrderParams{
		"sOrderNumber":      orderNumber,
		"sEigyouDay":        eigyouDay,
		"sCondition":        envOrDefault("TACHIBANASHI_ORDER_CONDITION", "*"),
		"sOrderPrice":       envOrDefault("TACHIBANASHI_ORDER_PRICE", "*"),
		"sOrderSuryou":      envOrDefault("TACHIBANASHI_ORDER_QTY", "*"),
		"sOrderExpireDay":   envOrDefault("TACHIBANASHI_ORDER_EXPIRE_DAY", "*"),
		"sGyakusasiZyouken": envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_ZYOUKEN", "*"),
		"sGyakusasiPrice":   envOrDefault("TACHIBANASHI_ORDER_GYAKUSASI_PRICE", "*"),
		"sSecondPassword":   secondPassword,
	}

	correctResp, err := cli.Request().KabuCorrectOrder(ctx, correctParams)
	if err != nil {
		log.Printf("order_correct failed: %v", err)
	} else {
		printOrderResponse("order_correct", correctResp)
	}

	detail, err = cli.Request().OrderListDetail(ctx, orderNumber, eigyouDay)
	if err != nil {
		log.Printf("order_detail (after correct) failed: %v", err)
	} else {
		printOrderDetail(detail)
	}

	cancelParams := request.OrderParams{
		"sOrderNumber":    orderNumber,
		"sEigyouDay":      eigyouDay,
		"sSecondPassword": secondPassword,
	}

	cancelResp, err := cli.Request().KabuCancelOrder(ctx, cancelParams)
	if err != nil {
		log.Printf("order_cancel failed: %v", err)
	} else {
		printOrderResponse("order_cancel", cancelResp)
	}
}

func printOrderResponse(label string, resp *request.OrderResponse) {
	fmt.Println(label)
	fmt.Printf("  result: %s %s\n", resp.ResultCode, resp.ResultText)
	fmt.Printf("  warning: %s %s\n", resp.WarningCode, resp.WarningText)
	fmt.Printf("  order_number: %s\n", resp.OrderNumber)
	fmt.Printf("  eigyou_day: %s\n", resp.EigyouDay)
}

func printOrderDetail(detail *request.OrderListDetailResponse) {
	fmt.Println("order_detail")
	fmt.Printf("  order_number: %s\n", detail.OrderNumber)
	fmt.Printf("  eigyou_day: %s\n", detail.EigyouDay)
	fmt.Printf("  issue_code: %s\n", detail.IssueCode)
	printAttributesFiltered("  ", detail.Fields, commonAttrSkip)
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
