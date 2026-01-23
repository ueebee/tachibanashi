package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/client"
	"github.com/ueebee/tachibanashi/event"
	"github.com/ueebee/tachibanashi/model"
)

func main() {
	_ = loadDotEnv(".env")

	loginID := mustEnvAny("TACHIBANASHI_LOGIN_ID", "TACHIBANA_USER_ID")
	password := mustEnvAny("TACHIBANASHI_PASSWORD", "TACHIBANA_PASSWORD")

	params, symbols := boardParamsFromEnv()

	baseURL := envOrAny(client.BaseURLDemo, "TACHIBANASHI_BASE_URL", "TACHIBANA_BASE_URL")
	cfg := client.Config{BaseURL: baseURL, EventParams: params}

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

	levels := parseIntEnvWithDefault("TACHIBANASHI_BOARD_LEVELS", 5)
	if levels < 1 {
		levels = 1
	}
	if levels > 10 {
		levels = 10
	}

	refresh := parseDurationEnv("TACHIBANASHI_BOARD_REFRESH", time.Second)
	if refresh <= 0 {
		refresh = time.Second
	}

	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	_, err = cli.Auth().Login(ctx, auth.Credentials{
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

	events, errs := cli.Event().Stream(ctx)
	quoteBook := event.NewQuoteBook()
	dirty := false

	ticker := time.NewTicker(refresh)
	defer ticker.Stop()

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			if fd, ok := ev.(event.FD); ok {
				quoteBook.Apply(fd)
				dirty = true
			}
		case err, ok := <-errs:
			if ok && err != nil && !errors.Is(err, context.Canceled) {
				log.Fatal(err)
			}
			return
		case <-ticker.C:
			if dirty {
				printBoard(quoteBook.SnapshotRows(), symbols, levels)
				dirty = false
			}
		}
	}
}

func boardParamsFromEnv() (event.Params, map[int]string) {
	params := event.Params{
		Cmds: []event.Command{event.CommandFD},
	}
	params.RID = parseIntEnv("TACHIBANASHI_EVENT_RID")
	params.BoardNo = parseIntEnv("TACHIBANASHI_EVENT_BOARD_NO")
	params.Rows = parseIntCSV(os.Getenv("TACHIBANASHI_EVENT_ROWS"))

	issueCodes := parseCSV(os.Getenv("TACHIBANASHI_EVENT_CODES"))
	marketCodes := parseCSV(os.Getenv("TACHIBANASHI_EVENT_MARKETS"))
	if len(issueCodes) > 0 && len(marketCodes) == 0 {
		marketCodes = make([]string, len(issueCodes))
		for i := range marketCodes {
			marketCodes[i] = "00"
		}
	}
	if len(params.Rows) == 0 && len(issueCodes) > 0 {
		params.Rows = make([]int, len(issueCodes))
		for i := range params.Rows {
			params.Rows[i] = i + 1
		}
	}
	params.IssueCodes = issueCodes
	params.MarketCodes = marketCodes

	if eno := strings.TrimSpace(os.Getenv("TACHIBANASHI_EVENT_ENO")); eno != "" {
		value, err := strconv.ParseInt(eno, 10, 64)
		if err != nil {
			log.Fatalf("invalid TACHIBANASHI_EVENT_ENO: %v", err)
		}
		params.Eno = value
	}

	return params, buildSymbolMap(params.Rows, issueCodes)
}

func printBoard(rows []event.QuoteRow, symbols map[int]string, levels int) {
	if len(rows) == 0 {
		fmt.Println("board: (no data)")
		return
	}
	for _, row := range rows {
		label := resolveLabel(row.Row, symbols)
		last, lastOK := row.Quote.LastPrice()
		fmt.Printf("%s last=%s time=%s\n", label, formatMaybePrice(last, lastOK), row.Quote.LastTime())
		fmt.Printf("  ask %s\n", formatLevels(row.Quote, "pGAP", "pGAV", levels))
		fmt.Printf("  bid %s\n", formatLevels(row.Quote, "pGBP", "pGBV", levels))
	}
}

func formatLevels(quote model.Quote, pricePrefix, sizePrefix string, levels int) string {
	if levels <= 0 {
		return ""
	}
	parts := make([]string, 0, levels)
	for i := 1; i <= levels; i++ {
		price := quote.Value(fmt.Sprintf("%s%d", pricePrefix, i))
		size := quote.Value(fmt.Sprintf("%s%d", sizePrefix, i))
		if price == "" && size == "" {
			parts = append(parts, "-")
			continue
		}
		if size == "" {
			size = "-"
		}
		if price == "" {
			price = "-"
		}
		parts = append(parts, fmt.Sprintf("%s(%s)", price, size))
	}
	return strings.Join(parts, " ")
}

func resolveLabel(row int, symbols map[int]string) string {
	if row == 0 {
		return "row=0"
	}
	if symbols == nil {
		return fmt.Sprintf("row=%d", row)
	}
	if symbol := strings.TrimSpace(symbols[row]); symbol != "" {
		return symbol
	}
	return fmt.Sprintf("row=%d", row)
}

func formatMaybePrice(value model.Price, ok bool) string {
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func buildSymbolMap(rows []int, codes []string) map[int]string {
	if len(codes) == 0 {
		return nil
	}
	out := make(map[int]string, len(codes))
	if len(rows) == len(codes) && len(rows) > 0 {
		for i, row := range rows {
			out[row] = codes[i]
		}
		return out
	}
	for i, code := range codes {
		out[i+1] = code
	}
	return out
}

func parseIntEnv(name string) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", name, err)
	}
	return parsed
}

func parseIntEnvWithDefault(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", name, err)
	}
	return parsed
}

func parseDurationEnv(name string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", name, err)
	}
	return parsed
}

func parseIntCSV(value string) []int {
	parts := parseCSV(value)
	if len(parts) == 0 {
		return nil
	}
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		parsed, err := strconv.Atoi(part)
		if err != nil {
			log.Fatalf("invalid integer %q: %v", part, err)
		}
		out = append(out, parsed)
	}
	return out
}

func parseCSV(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	return out
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
