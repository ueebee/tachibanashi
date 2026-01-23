package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	params, symbols := eventParamsFromEnv()

	baseURL := envOrAny(client.BaseURLDemo, "TACHIBANASHI_BASE_URL", "TACHIBANA_BASE_URL")
	cfg := client.Config{BaseURL: baseURL, EventParams: params}
	if isTrue(os.Getenv("TACHIBANASHI_EVENT_LOG")) {
		cfg.Logger = log.New(os.Stdout, "event: ", log.LstdFlags)
	}

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	_, err = cli.Auth().Login(ctx, auth.Credentials{
		LoginID:  loginID,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	if isTrue(os.Getenv("TACHIBANASHI_EVENT_LOG")) {
		logEventURL(cli, params)
	}
	defer func() {
		if err := cli.Auth().Logout(context.Background()); err != nil {
			log.Printf("logout failed: %v", err)
		}
	}()

	conn, err := cli.Event().Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Printf("event connected")

	for {
		ev, err := conn.Recv(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Fatal(err)
		}
		printEvent(ev, symbols)
	}
}

func eventParamsFromEnv() (event.Params, map[int]string) {
	var params event.Params
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
	params.IssueCodes = issueCodes
	params.MarketCodes = marketCodes

	if eno := strings.TrimSpace(os.Getenv("TACHIBANASHI_EVENT_ENO")); eno != "" {
		value, err := strconv.ParseInt(eno, 10, 64)
		if err != nil {
			log.Fatalf("invalid TACHIBANASHI_EVENT_ENO: %v", err)
		}
		params.Eno = value
	}

	cmds := parseCommandCSV(os.Getenv("TACHIBANASHI_EVENT_CMDS"))
	if len(cmds) > 0 {
		params.Cmds = cmds
	}

	return params, buildSymbolMap(params.Rows, issueCodes)
}

func parseCommandCSV(value string) []event.Command {
	items := parseCSV(value)
	if len(items) == 0 {
		return nil
	}
	out := make([]event.Command, 0, len(items))
	for _, item := range items {
		out = append(out, event.Command(strings.ToUpper(item)))
	}
	return out
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

func printEvent(ev event.Event, symbols map[int]string) {
	switch e := ev.(type) {
	case event.ST:
		fmt.Printf("st error_no=%s err=%s\n", e.ErrNo, e.Err)
	case event.KP:
		fmt.Println("kp")
	case event.EC:
		order := e.Order()
		fmt.Printf("ec order=%s status=%s symbol=%s side=%s qty=%d price=%d\n",
			order.ID, order.Status, order.Symbol, order.Side, order.Quantity, order.Price)
		if exec, ok := e.Execution(); ok {
			fmt.Printf("  exec qty=%d price=%d time=%s\n", exec.Quantity, exec.Price, exec.Time)
		}
	case event.NS:
		fmt.Printf("ns id=%s headline=%s\n", e.NewsID, truncate(e.Headline, 120))
		if e.Body != "" {
			fmt.Printf("  body=%s\n", truncate(e.Body, 200))
		}
	case event.SS:
		fmt.Printf("ss login=%s status=%s changed_at=%s\n", e.LoginKind, e.SystemStatus, e.ChangedAt)
	case event.US:
		fmt.Printf("us uc=%s uu=%s status=%s market=%s changed_at=%s\n",
			e.OperationCode, e.OperationUnit, e.OperationStatus, e.MarketCode, e.ChangedAt)
	case event.FD:
		quotes := e.Quotes(symbols)
		if len(quotes) == 0 {
			fmt.Println("fd (no rows)")
			return
		}
		for _, quote := range quotes {
			last, lastOK := quote.LastPrice()
			prev, prevOK := quote.PrevClose()
			fmt.Printf("fd %s last=%s prev=%s time=%s\n",
				quote.Symbol,
				formatMaybePrice(last, lastOK),
				formatMaybePrice(prev, prevOK),
				quote.LastTime(),
			)
		}
	case event.Unknown:
		fmt.Printf("unknown %s\n", strings.TrimSpace(string(e.Raw)))
	default:
		fmt.Printf("%s\n", ev.Kind())
	}
}

func formatMaybePrice(value model.Price, ok bool) string {
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if value == "" || max <= 0 {
		return value
	}
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
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

func logEventURL(cli *client.Client, params event.Params) {
	urls := cli.VirtualURLs()
	if urls.EventWS == "" {
		log.Printf("event ws url is empty")
		return
	}
	encoded, err := event.BuildWSURL(urls.EventWS, params)
	if err != nil {
		log.Printf("event ws url build error: %v", err)
		return
	}
	parsed, err := url.Parse(encoded)
	if err != nil {
		log.Printf("event ws url parse error: %v", err)
		return
	}
	log.Printf("event ws host=%s path=%s query=%s", parsed.Host, maskPath(parsed.Path), parsed.RawQuery)
}

func maskPath(path string) string {
	if path == "" {
		return path
	}
	if len(path) <= 12 {
		return path
	}
	return path[:6] + "..." + path[len(path)-6:]
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
