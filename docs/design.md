# 基本設計（tachibanashi）

目的: 立花証券 e 支店 API の薄い Go ラッパーを提供し、投資ボットに組み込みやすい最小 API を実現する。

## 設計方針
- API 面は小さく、合成しやすい設計（pybotters 風）
- グローバルより明示的な設定を優先
- demo/prod は base URL で切替可能にする
- まずは薄く・簡潔に実装し、必要性が出た時点で抽象化する

## プロジェクト設定
- Go module: `github.com/ueebee/tachibanashi`
- Go version: 最新（現環境は `go1.24.2`）
- API バージョン: v4r8
- demo/prod の base URL は固定で採用する
  - prod: `https://kabuka.e-shiten.jp/e_api_v4r8/`
  - demo: `https://demo-kabuka.e-shiten.jp/e_api_v4r8/`

## 全体構成
中核は `Client` と共通の HTTP 送受信レイヤとし、機能はサービス単位で分割する。

```
client/
  client.go        (Client, Config, Option)
  http.go          (doJSON, 共通送受信)
auth/
request/
price/
master/
event/
model/
errors/
```

## 抽象化レイヤ（Facade）
立花証券APIの用語や構造を一般的な用語へ変換する薄い抽象化レイヤを設ける。
低レイヤは「立花証券APIそのままの薄いI/F」、高レイヤは「一般的なドメインI/F」とする。

例:
- IssueCode → Symbol/Code
- MarketPrice → Quote
- sTargetColumn → Fields

方針:
- 抽象レイヤは `model` に一般的な型（Quote/Order/Position/Balance 等）を定義
- Raw レスポンスは map 等で保持し、未知フィールドを捨てない
- 既存の薄い API は維持し、後方互換を優先

## 中核 API
```
type Client struct {
  cfg   Config
  http  *http.Client
  token TokenStore
}

type Config struct {
  BaseURL   string
  Timeout   time.Duration
  UserAgent string
  Logger    Logger
}

type Option func(*Config)

func New(cfg Config, opts ...Option) (*Client, error)

func (c *Client) Auth() *auth.Service
func (c *Client) Request() *request.Service
func (c *Client) Price() *price.Service
func (c *Client) Master() *master.Service
func (c *Client) Event() *event.Service
```

## データ型方針
- 価格は `int64`（円の整数）
- 数量は `int64`（株数）
- 型安全のため `type Price int64`, `type Quantity int64` を `model` に定義する

## 拡張ポイント
- `TokenStore` インターフェース
  - `p_no` の保存/更新/無効化を差し替え可能にする
  - デフォルトはメモリ実装（将来ファイル/DB 実装へ拡張可能）
- `HTTPClient` 注入
  - 計測、リトライ、プロキシ対応などを外部で制御可能
- `Logger` インターフェース
  - ロガー差し替えを許可し、必要なら無効化も可能

## 共通送受信レイヤ
- `client.doJSON(ctx, method, path, req, resp)` を共通経路とする
- ここで HTTP/業務エラーを `errors` に統一変換
- 必須フィールド不足は `errors.ValidationError` で返す

## 文字コード
- 参照マニュアル HTML は `Shift_JIS`（`spec/mfds_json_api_refference_src/mfds_json_api_ref_text.html` 参照）
- API レスポンス内の日本語が `Shift_JIS` の可能性があるため、共通送受信レイヤで UTF-8 へ正規化する
  - `Content-Type` の `charset` が指定されていればそれに従う
  - 未指定の場合でも UTF-8 不正なら `Shift_JIS` とみなしてデコードする

## 認証設計
- `auth.Service` が login/logout/仮想URL を提供
- `login` 成功時に `p_no` を `TokenStore` に保存
- `logout` 成功時に `p_no` を無効化

## REQUEST/PRICE/MASTER
- それぞれのサービスが対応する API を提供
- 送受信は共通レイヤに集約し、各サービスは入出力モデルのみ定義

## EVENT（WebSocket）
- `event.Client` が接続・再接続・keepalive を管理
- 1 セッション制約は内部でガード（同時接続を防止）
- 受信は raw → event 型へパース
- 未知イベントは `event.Unknown` へ落とし、ロバスト性を確保
- 基本 API は `Recv(ctx)` の pull 型とし、利便性のため `Stream(ctx)` を補助で提供する

## エラー設計
- `errors.APIError{Code, Message, Detail, Raw}`（業務エラー）
- `errors.HTTPError{Status, Body}`（HTTP レベル）
- `errors.ValidationError{Field, Reason}`
- `errors.IsRetryable(err)` を用意し、ボット側の制御を支援

## テスト方針
- `httptest` で API クライアントのスモークテスト
- `event` と `model` にパース/エンコードのユニットテストを集中
- `TokenStore` はモック差し替えで検証

## 実装順序
`tmp/plan.md` の順序を基本とし、認証 → REQUEST → PRICE → MASTER → EVENT の順で実装する。
