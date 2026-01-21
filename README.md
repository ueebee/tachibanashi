# tachibanashi

立花証券 e 支店 API の薄い Go ラッパーです。投資ボットでの利用を前提に、最小限の API 面で設計しています。

## ステータス
- v4r8 を前提に整備中
- Auth / Price（時価スナップショット）まで実装済み

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

	resp, err := cli.Price().Snapshot(
		context.Background(),
		[]string{"6501", "6502", "6503"},          // まとめて取得（最大120）
		[]string{"pDPP", "pPRP", "tDPP:T"}, // 欲しいカラム
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range resp.Prices {
		fmt.Println(entry.IssueCode, entry.Value("pDPP"), entry.Value("tDPP:T"))
	}
}
```

### 2) サンプル CLI で動作確認
`.env` に認証情報を置いて実行します（`.env.example` 参照）。

```bash
go run ./cmd/auth-check
go run ./cmd/price-snapshot
```

## 環境変数（サンプル）
`.env.example` に記載しています。

- `TACHIBANASHI_BASE_URL`（任意、未指定なら demo）
- `TACHIBANASHI_LOGIN_ID`（必須）
- `TACHIBANASHI_PASSWORD`（必須）
- `TACHIBANASHI_CODES`（price-snapshot 用、カンマ区切り）
- `TACHIBANASHI_COLUMNS`（price-snapshot 用、カンマ区切り）

## Base URL
- prod: `https://kabuka.e-shiten.jp/e_api_v4r8/`
- demo: `https://demo-kabuka.e-shiten.jp/e_api_v4r8/`

## 開発
```bash
go test ./...
```
