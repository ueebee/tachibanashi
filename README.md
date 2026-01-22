# tachibanashi

立花証券 e 支店 API の薄い Go ラッパーです。投資ボットでの利用を前提に、最小限の API 面で設計しています。

## ステータス
- v4r8 を前提に整備中
- Auth / Request（口座・建玉・余力の読み取り）/ Price（時価スナップショット）まで実装済み

## 使い方（例）

### 1) ログイン → 時価スナップショット → ログアウト
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/client"
)

func main() {
	cli, err := client.New(client.Config{
		BaseURL: client.BaseURLDemo, // or client.BaseURLProd
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.Auth().Login(context.Background(), auth.Credentials{
		LoginID:  "your_login_id",
		Password: "your_password",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Auth().Logout(context.Background())

	resp, err := cli.Price().QuoteSnapshot(
		context.Background(),
		[]string{"6501", "6502", "6503"},          // まとめて取得（最大120）
		[]string{"pDPP", "pPRP", "tDPP:T"}, // 欲しいカラム
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, quote := range resp.Quotes {
		fmt.Println(quote.Symbol, quote.Value("pDPP"), quote.Value("tDPP:T"))
	}
}
```

### 2) 注文一覧と新規注文（REQUEST I/F）
```go
package main

import (
	"context"
	"log"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/client"
	"github.com/ueebee/tachibanashi/request"
)

func main() {
	cli, err := client.New(client.Config{BaseURL: client.BaseURLDemo})
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.Auth().Login(context.Background(), auth.Credentials{
		LoginID:  "your_login_id",
		Password: "your_password",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Auth().Logout(context.Background())

	orders, err := cli.Request().Orders(context.Background(), request.OrderParams{
		"sIssueCode": "6501",
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = orders

	_, err = cli.Request().KabuNewOrder(context.Background(), request.OrderParams{
		"sZyoutoekiKazeiC":        "1",
		"sIssueCode":              "6501",
		"sSizyouC":                "00",
		"sBaibaiKubun":            "3",
		"sCondition":              "0",
		"sOrderPrice":             "0",
		"sOrderSuryou":            "100",
		"sGenkinShinyouKubun":     "0",
		"sOrderExpireDay":         "0",
		"sGyakusasiOrderType":     "0",
		"sGyakusasiZyouken":       "0",
		"sGyakusasiPrice":         "*",
		"sTatebiType":             "*",
		"sTategyokuZyoutoekiKazeiC": "*",
		"sSecondPassword":         "your_second_password",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

### 3) サンプル CLI で動作確認
`.env` に認証情報を置いて実行します（`.env.example` 参照）。

```bash
go run ./cmd/auth-check
go run ./cmd/price-snapshot
go run ./cmd/request-read
go run ./cmd/order-read
go run ./cmd/order-submit
go run ./cmd/order-correct
go run ./cmd/order-cancel
go run ./cmd/order-crud
```

## 環境変数（サンプル）
`.env.example` に記載しています。

- `TACHIBANASHI_BASE_URL`（任意、未指定なら demo）
- `TACHIBANASHI_LOGIN_ID`（必須）
- `TACHIBANASHI_PASSWORD`（必須）
- `TACHIBANASHI_CODES`（price-snapshot 用、カンマ区切り）
- `TACHIBANASHI_COLUMNS`（price-snapshot 用、カンマ区切り）
- `TACHIBANASHI_ISSUE_CODE`（request-read 用、任意）
- `TACHIBANASHI_GENBUTU_HITUKE_INDEX`（request-read 用、任意、3〜5）
- `TACHIBANASHI_SINYOU_HITUKE_INDEX`（request-read 用、任意、0〜5）
- `TACHIBANASHI_ORDER_ISSUE_CODE`（order-read 用、任意）
- `TACHIBANASHI_ORDER_SIKKOU_DAY`（order-read 用、任意）
- `TACHIBANASHI_ORDER_STATUS`（order-read 用、任意）
- `TACHIBANASHI_ORDER_NUMBER`（order-read 詳細用、任意）
- `TACHIBANASHI_EIGYOU_DAY`（order-read 詳細用、任意）
- `TACHIBANASHI_ORDER_BAIBAI_KUBUN`（order-submit 用、必須）
- `TACHIBANASHI_ORDER_QTY`（order-submit 用、必須）
- `TACHIBANASHI_ORDER_PRICE`（order-submit 用、任意、未指定は成行）
- `TACHIBANASHI_SECOND_PASSWORD`（order-submit/order-correct/order-cancel 用、必須）
- `TACHIBANASHI_ORDER_NUMBER`（order-correct/order-cancel 用、必須）
- `TACHIBANASHI_EIGYOU_DAY`（order-correct/order-cancel 用、必須）
- `TACHIBANASHI_ORDER_CONDITION`（order-correct 用、任意、未指定は変更なし）
- `TACHIBANASHI_ORDER_EXPIRE_DAY`（order-correct 用、任意、未指定は変更なし）
- `TACHIBANASHI_ORDER_GYAKUSASI_ZYOUKEN`（order-correct 用、任意、未指定は変更なし）
- `TACHIBANASHI_ORDER_GYAKUSASI_PRICE`（order-correct 用、任意、未指定は変更なし）
- `TACHIBANASHI_TIMEOUT`（任意）
- `TACHIBANASHI_USER_AGENT`（任意）
- `TACHIBANASHI_INSECURE_TLS`（任意、true/1 で検証無効）

## 動作確認（demo）
- `.env.example` を ` .env` にコピーして認証情報を設定
- `go run ./cmd/auth-check` でログイン疎通
- `go run ./cmd/price-snapshot` で時価取得
- `go run ./cmd/request-read` で口座・建玉・余力の読み取り
- `go run ./cmd/order-read` で注文一覧の読み取り
- `go run ./cmd/order-submit` で注文を送信
- `go run ./cmd/order-correct` で注文訂正
- `go run ./cmd/order-cancel` で注文取消
- `go run ./cmd/order-crud` で注文の CRUD フローを一括確認
- 余力詳細を出す場合は `TACHIBANASHI_GENBUTU_HITUKE_INDEX` / `TACHIBANASHI_SINYOU_HITUKE_INDEX` を設定
- 注文詳細を出す場合は `TACHIBANASHI_ORDER_NUMBER` / `TACHIBANASHI_EIGYOU_DAY` を設定

## Base URL
- prod: `https://kabuka.e-shiten.jp/e_api_v4r8/`
- demo: `https://demo-kabuka.e-shiten.jp/e_api_v4r8/`

## 注意事項
- API の日本語メッセージは Shift_JIS の可能性があるため、ライブラリ側で UTF-8 へ正規化して扱います

## 開発
```bash
go test ./...
```
