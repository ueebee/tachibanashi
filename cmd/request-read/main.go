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
	"github.com/ueebee/tachibanashi/model"
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

	issueCode := strings.TrimSpace(os.Getenv("TACHIBANASHI_ISSUE_CODE"))
	ctx := context.Background()

	buyingPower, err := cli.Request().BuyingPower(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("buying_power")
	fmt.Printf("  update: %s\n", buyingPower.Raw.SummaryUpdate)
	fmt.Printf("  genkabu_kaituke: %s\n", buyingPower.Raw.SummaryGenkabuKaituke)
	fmt.Printf("  nseityou_tousi: %s\n", buyingPower.Raw.SummaryNseityouTousiKanougaku)
	fmt.Printf("  husokukin_flag: %s\n", buyingPower.Raw.HusokukinHasseiFlg)

	summary, err := cli.Request().ZanKaiSummary(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("summary")
	printSummaryValue(summary.Fields, "sUpdateDate", "update")
	printSummaryValue(summary.Fields, "sOisyouHasseiFlg", "oisyou_flag")
	printSummaryValue(summary.Fields, "sGenbutuKabuKaituke", "genbutu_kaituke")
	printSummaryValue(summary.Fields, "sSinyouSinkidate", "sinyou_sinkidate")
	printSummaryValue(summary.Fields, "sSyukkin", "syukkin")
	printSummaryValue(summary.Fields, "sNseityouTousiKanougaku", "nseityou_tousi")
	printSummaryValue(summary.Fields, "sHosyouKinritu", "hosyou_kinritu")
	printSummaryValue(summary.Fields, "sFusokugaku", "fusokugaku")

	cash, err := cli.Request().CashPositions(ctx, issueCode)
	if err != nil {
		log.Fatal(err)
	}
	printPositions("cash_positions", cash.Positions, "sUriOrderZanKabuSuryou", "sUriOrderGaisanBokaTanka", "sUriOrderGaisanHyoukaSoneki")

	margin, err := cli.Request().MarginPositions(ctx, issueCode)
	if err != nil {
		log.Fatal(err)
	}
	printPositions("margin_positions", margin.Positions, "sOrderTategyokuSuryou", "sOrderTategyokuTanka", "sOrderGaisanHyoukaSoneki")
}

func printSummaryValue(fields model.Attributes, key, label string) {
	value := fields.Value(key)
	if value == "" {
		return
	}
	fmt.Printf("  %s: %s\n", label, value)
}

func printPositions(title string, positions []model.Position, qtyKey, avgKey, pnlKey string) {
	fmt.Println(title)
	if len(positions) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, pos := range positions {
		qty := valueOrNumber(pos.Raw, qtyKey, int64(pos.Quantity))
		avg := valueOrNumber(pos.Raw, avgKey, int64(pos.AvgPrice))
		pnl := valueOrNumber(pos.Raw, pnlKey, pos.UnrealPnL)
		fmt.Printf("  %s qty=%s avg=%s pnl=%s\n", pos.Symbol, qty, avg, pnl)
	}
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
