# tachibanashi v0 計画（仕様書ベース / チェックリスト）

目的: 立花証券 e支店 API を「立ち話」感覚で扱える薄い Go ラッパーとして実装し、認証から通知まで順に整備する。

参照:
- spec/api_request_if_v4r7.pdf
- spec/api_request_if_master_v4r5.pdf
- spec/api_event_if_v4r7.pdf
- spec/mfds_json_api_refference.pdf

進め方:
- `tmp/plan.md` の順序を基本とする
- エンドポイントごとに「構造体/メソッド/エラー変換/テスト」を揃える
- 過度な抽象化は避け、必要性が出たら整理する

## 0) 共通/基盤
- [x] Go module / パッケージ構成
- [x] client.Config/New/DoJSON/共通エラー変換
- [x] README/ドキュメント雛形
- [ ] CI（go test / golangci-lint）
- [ ] golangci 設定

## 1) 認証 I/F（api_request_if_v4r7）
- [x] CLMAuthLoginRequest（実装）
- [ ] CLMAuthLoginRequest（テスト）
- [x] CLMAuthLoginAck（レスポンス取り込み）
- [ ] CLMAuthLoginAck（テスト）
- [x] CLMAuthLogoutRequest（実装）
- [ ] CLMAuthLogoutRequest（テスト）
- [x] CLMAuthLogoutAck（レスポンス取り込み）
- [ ] CLMAuthLogoutAck（テスト）
- [x] 仮想URL取得 / TokenStore
- [ ] 仮想URL取得 / TokenStore テスト

## 2) 注文 I/F（api_request_if_v4r7）
- [x] CLMKabuNewOrder（実装）
- [ ] CLMKabuNewOrder（テスト）
- [x] CLMKabuCorrectOrder（実装）
- [ ] CLMKabuCorrectOrder（テスト）
- [x] CLMKabuCancelOrder（実装）
- [ ] CLMKabuCancelOrder（テスト）
- [x] CLMOrderList（実装）
- [x] CLMOrderList（テスト）
- [x] CLMOrderListDetail（実装）
- [x] CLMOrderListDetail（テスト）
- [x] (spec外) CLMKabuCancelOrderAll（実装）
- [ ] (spec外) CLMKabuCancelOrderAll（テスト）

## 3) 口座/建玉/余力 I/F（api_request_if_v4r7）
- [x] CLMGenbutuKabuList（実装）
- [ ] CLMGenbutuKabuList（テスト）
- [x] CLMShinyouTategyokuList（実装）
- [ ] CLMShinyouTategyokuList（テスト）
- [x] CLMZanKaiKanougaku（実装）
- [ ] CLMZanKaiKanougaku（テスト）
- [x] CLMZanShinkiKanoIjiritu（実装）
- [ ] CLMZanShinkiKanoIjiritu（テスト）
- [x] CLMZanUriKanousuu（実装）
- [ ] CLMZanUriKanousuu（テスト）
- [x] CLMZanKaiSummary（実装）
- [ ] CLMZanKaiSummary（テスト）
- [x] CLMZanKaiKanougakuSuii（実装）
- [ ] CLMZanKaiKanougakuSuii（テスト）
- [x] CLMZanKaiGenbutuKaitukeSyousai（実装）
- [ ] CLMZanKaiGenbutuKaitukeSyousai（テスト）
- [x] CLMZanKaiSinyouSinkidateSyousai（実装）
- [ ] CLMZanKaiSinyouSinkidateSyousai（テスト）
- [x] CLMZanRealHosyoukinRitu（実装）
- [ ] CLMZanRealHosyoukinRitu（テスト）

## 4) マスタ I/F（api_request_if_v4r7 / api_request_if_master_v4r5）
### 4-1) マスタ情報ダウンロード（CLMEventDownload）
#### 4-1-0) 仕様確認 / 方針決定
- [x] 応答形式の確認（連続配信、JSON オブジェクト単位で区切り）
- [x] 接続維持・切断条件の確認（初期完了後は切断をデフォルト）
- [x] MasterStore の公開方針（参照用 API を公開 / 購読型は将来対応）
- [x] 共通項目の扱い（メッセージ単位のメタ情報として保持/警告は継続/結果コードは中断）
- [x] 更新方式の確認（更新通番優先、同値は更新日時、削除フラグは削除）
- [x] 静的/動的マスタの分類確定（静的: SystemStatus/DateZyouhou/UnyouStatus/Yobine/DaiyouKakeme/HosyoukinMst/OrderErrReason）

#### 4-1-1) 通信・パース基盤
- [x] CLMEventDownload 要求モデル定義（共通項目 + sCLMID）
- [x] CLMEventDownload 送信メソッド（MASTER 仮想URL）
- [x] 受信フレームの分割・JSON デコード
- [x] 受信メッセージのディスパッチ（sCLMID でルーティング）
- [x] 初期ダウンロード受信ループ（CLMEventDownloadComplete まで）
- [x] 初期完了後の更新通知処理（UPDATE 差分反映）
- [x] 初期完了後も接続を維持して更新通知を受け続けるモード（DownloadStream）
- [x] エラー応答（p_errno/p_err）処理
- [x] 文字コード正規化（Shift_JIS → UTF-8 前提の確認）
- [x] テスト: 連続配信のパース / 完了通知 / 更新反映

#### 4-1-2) マスタストア / 更新処理
- [x] MasterStore インターフェース（Get/Upsert/Delete/All）
- [x] インメモリ実装（種別ごとに map + index）
- [x] 主キー抽出関数（マスタ種別ごと）
- [x] 更新通番/削除フラグの優先ルール（更新通番優先、同値は更新日時）
- [x] 参照用インデックス（銘柄コード/市場/口座区分など）
- [x] スナップショット取得 API（読み取り用）
- [x] テスト: Upsert/Delete とインデックス整合性
- [ ] TODO: 変更通知（購読型）インターフェース

#### 4-1-3) 運用系マスタ（spec 2-1〜2-4）
- [x] CLMSystemStatus: 項目一覧抽出 / 主キー決定
- [x] CLMSystemStatus: モデル定義 / パース / 格納 / テスト
- [x] CLMDateZyouhou: 項目一覧抽出 / 主キー決定
- [x] CLMDateZyouhou: モデル定義 / パース / 格納 / テスト
- [x] CLMUnyouStatus: 項目一覧抽出 / 主キー決定
- [x] CLMUnyouStatus: モデル定義 / パース / 格納 / テスト
- [x] CLMUnyouStatusKabu: 項目一覧抽出 / 主キー決定
- [x] CLMUnyouStatusKabu: モデル定義 / UPDATE 反映 / テスト
- [x] CLMUnyouStatusHasei: 項目一覧抽出 / 主キー決定
- [x] CLMUnyouStatusHasei: モデル定義 / UPDATE 反映 / テスト

#### 4-1-4) 銘柄系マスタ（spec 2-6〜2-11）
- [x] CLMIssueMstKabu: 項目一覧抽出 / 主キー決定
- [x] CLMIssueMstKabu: モデル定義 / UPDATE 反映 / テスト
- [x] CLMIssueSizyouMstKabu: 項目一覧抽出 / 主キー決定
- [x] CLMIssueSizyouMstKabu: モデル定義 / UPDATE 反映 / テスト
- [x] CLMIssueSizyouKiseiKabu: 項目一覧抽出 / 主キー決定
- [x] CLMIssueSizyouKiseiKabu: モデル定義 / UPDATE 反映 / テスト
- [x] CLMIssueMstSak: 項目一覧抽出 / 主キー決定
- [x] CLMIssueMstSak: モデル定義 / UPDATE 反映 / テスト
- [x] CLMIssueMstOp: 項目一覧抽出 / 主キー決定
- [x] CLMIssueMstOp: モデル定義 / UPDATE 反映 / テスト
- [x] CLMIssueSizyouKiseiHasei: 項目一覧抽出 / 主キー決定
- [x] CLMIssueSizyouKiseiHasei: モデル定義 / UPDATE 反映 / テスト

#### 4-1-5) 静的マスタ（spec 2-12〜2-15）
- [x] CLMYobine: 項目一覧抽出 / 主キー決定
- [x] CLMYobine: モデル定義 / パース / 格納 / テスト
- [x] CLMDaiyouKakeme: 項目一覧抽出 / 主キー決定
- [x] CLMDaiyouKakeme: モデル定義 / パース / 格納 / テスト
- [x] CLMHosyoukinMst: 項目一覧抽出 / 主キー決定
- [x] CLMHosyoukinMst: モデル定義 / パース / 格納 / テスト
- [x] CLMOrderErrReason: 項目一覧抽出 / 主キー決定
- [x] CLMOrderErrReason: モデル定義 / パース / 格納 / テスト
- [x] CLMEventDownloadComplete: 受信 / 状態管理 / テスト

### 4-2) マスタ問合取得
- [x] 共通: Request/Response ラッパ（MASTER 仮想URL）
- [x] 共通: 文字コード正規化（ニュース p_HDL/p_TX の Shift_JIS）
- [x] 共通: sTargetIssueCode 最大120件のバリデーション
- [x] 共通: 空文字/0 の扱い（値なし）方針整理
- [x] 共通: ユニットテスト（リクエスト/レスポンス）

#### 4-2-1) CLMMfdsGetMasterData（マスタ情報問合）
- [x] 要求モデル: sTargetCLMID / sTargetColumn
- [x] 対象CLMID一覧の整理（CLMIssueMstKabu/CLMIssueSizyouMstKabu/CLMIssueMstSak/CLMIssueMstOp/CLMIssueMstOther/CLMIssueMstIndex/CLMIssueMstFx/CLMOrderErrReason/CLMDateZyouhou）
- [x] 応答モデル: CLMID 名ごとの配列（動的キー）パース
- [x] 取得列指定のフィルタ（指定列のみ抽出）
- [x] テスト: 2種類CLMID混在レスポンスのパース

#### 4-2-2) CLMMfdsGetNewsHead（ニュースヘッダ）
- [x] 要求モデル: p_CG / p_IS / p_DT_FROM / p_DT_TO / p_REC_OFST / p_REC_LIMT
- [x] レコード取得条件（AND 条件）/ p_REC_MAX 取り扱い
- [x] 応答モデル: aCLMMfdsNewsHead の配列パース
- [x] 文字コード: p_HDL デコード
- [x] テスト: 検索条件付きの要求とレスポンス

#### 4-2-3) CLMMfdsGetNewsBody（ニュース本文）
- [x] 要求モデル: p_ID
- [x] 応答モデル: aCLMMfdsNewsBody の配列パース
- [x] 文字コード: p_HDL / p_TX デコード
- [x] テスト: 1件レスポンスのパース

#### 4-2-4) CLMMfdsGetIssueDetail（銘柄詳細）
- [x] 要求モデル: sTargetIssueCode（最大120件）
- [x] 応答モデル: aCLMMfdsIssueDetail の配列パース
- [x] テスト: 主要フィールドの型/値パース

#### 4-2-5) CLMMfdsGetSyoukinZan（証金残）
- [x] 要求モデル: sTargetIssueCode（最大120件）
- [x] 応答モデル: aCLMMfdsSyoukinZan の配列パース
- [x] テスト: 数値/日付フィールドのパース

#### 4-2-6) CLMMfdsGetShinyouZan（信用残）
- [x] 要求モデル: sTargetIssueCode（最大120件）
- [x] 応答モデル: aCLMMfdsShinyouZan の配列パース
- [x] テスト: 数値/日付フィールドのパース

#### 4-2-7) CLMMfdsGetHibuInfo（逆日歩）
- [x] 要求モデル: sTargetIssueCode（最大120件）
- [x] 応答モデル: aCLMMfdsHibuInfo の配列パース
- [x] テスト: 逆日歩フィールドのパース

## 5) 時価 I/F（api_request_if_v4r7 / mfds_json_api_refference）
- [ ] 共通: Request/Response ラッパ（PRICE 仮想URL）
- [ ] 共通: 文字コード正規化（Shift_JIS → UTF-8 前提）
- [ ] 共通: sTargetIssueCode 最大120件のバリデーション
- [ ] 共通: sTargetColumn（情報コード一覧）との整合確認
- [ ] 共通: ユニットテスト（リクエスト/レスポンス）

#### 5-1) CLMMfdsGetMarketPrice（時価スナップショット）
- [x] 要求モデル: sTargetIssueCode / sTargetColumn
- [x] 応答モデル: aCLMMfdsMarketPrice の配列パース
- [x] 抽象モデル: QuoteSnapshot への変換
- [ ] 情報コード一覧（FD）との対応表整備
- [ ] テスト: 複数銘柄・複数カラムのパース

#### 5-2) CLMMfdsGetMarketPriceHistory（蓄積情報）
- [ ] 要求モデル: sIssueCode / sSizyouC（1リクエスト1銘柄）
- [ ] 応答モデル: aCLMMfdsGetMarketPriceHistory の配列パース
- [ ] 分割関連フィールド（pSPUO/pSPUC/pSPUK）の扱い整理
- [ ] テスト: 日付昇順リストのパース

## 6) EVENT I/F（api_event_if_v4r7）
- [x] Service/Conn 骨格

#### 6-0) 仕様確認 / 方針決定
- [x] 接続方式は WebSocket（sUrlEventWebSocket）を使用
- [x] URL パラメータ仕様の整理（p_rid/p_board_no/p_gyou_no/p_issue_code/p_mkt_code/p_eno/p_evt_cmd）
- [x] p_evt_cmd 対応範囲の決定（ST/KP/FD/EC/NS/SS/US のみ）
- [x] p_eno の再送ルール整理（EC/NS/SS/US）
- [x] Base64 対象項目の確認（Shift_JIS 代替）
- [x] RR/FC 非公開通知は未対応と明記

URL パラメータまとめ（EVENT WebSocket / HTTP 共通）
- p_rid: アプリ機能識別。API 利用は 0（時価配信なし）/ 22（時価配信あり）。他は株価ボード画面向け。
- p_board_no: ボード番号。API は 1000、画面系は 1-10/120 など。
- p_gyou_no: 行番号（1-120）。必要時のみ指定、カンマ区切りで複数可。
- p_issue_code: 銘柄コード（最大 120）。必要時のみ、p_gyou_no と同数でカンマ区切り。
- p_mkt_code: 市場コード（最大 120）。必要時のみ、p_issue_code と同数でカンマ区切り。
- p_eno: 再送開始番号。指定番号の次から送信、0 は全件。
- p_evt_cmd: 通知種別のカンマ区切り（ST/KP/FD/EC/NS/SS/US）。

p_eno 再送ルール
- 対象通知: EC/NS/SS/US のみ（ST/KP/FD は無関係）
- 0 を指定すると当日未削除通知を全再送（通知削除機能は非公開）
- 0 以外は指定番号の「次」から送信（p_ENO はユニークだが連番ではない）
- 再接続時は直近の p_ENO を引き継ぐ想定（重複回避のベストエフォート）

Base64 対象項目（WebSocket 版）
- EC: p_IN（銘柄名称）
- NS: p_HDL（ニュースタイトル）, p_TX（ニュース本文）
- FD: Shift_JIS を含む値は x_ プレフィックスの16進表現で送信（Base64 ではない）

#### 6-1) URL / パラメータ検証
- [x] WS URL ビルダ（パラメータ省略とデフォルト）
- [x] パラメータ検証（p_evt_cmd/最大120銘柄/ボード組合せ）
- [x] board/row 指定の組合せサポート（p_board_no/p_gyou_no/p_issue_code/p_mkt_code）

#### 6-2) 接続・再接続・セッション
- [x] 接続/切断/再接続（コンテキスト停止・バックオフ）
- [x] keepalive（KP 受信/送信扱い）
- [x] 1セッション制約（重複接続ガード）
- [x] p_eno レジューム（再送開始番号の管理）

#### 6-3) フレームデコード / 共通処理
- [x] 受信フレームのデコード（^A/^B/^C 区切り）
- [x] 共通項目パース（p_no/p_date/p_cmd）
- [x] Base64 デコード（WebSocket は Shift_JIS を扱えないため）
- [x] イベントディスパッチ（p_cmd でルーティング）
- [x] 共通テスト（区切りパース/Base64/必須項目）

#### 6-4) ST / KP
- [x] ST: 通知モデル定義（必須項目）
- [x] ST: パース/テスト
- [x] KP: 通知モデル定義
- [x] KP: パース/テスト

#### 6-5) FD（時価情報）
- [x] FD: 情報コード一覧（p_*/t_*）のパース
- [x] FD: 初回スナップショット/差分更新の扱い
- [x] FD: Quote/QuoteSnapshot 変換
- [x] FD: テスト（初回/差分/複数銘柄）

#### 6-6) EC（注文約定通知）
- [x] EC: 通知モデル定義（親注文番号/注文種別含む）
- [x] EC: パース/テスト
- [x] EC: Order/Execution へのマッピング

#### 6-7) NS（ニュース通知）
- [x] NS: 通知モデル定義
- [x] NS: Base64 → 文字列復号（見出し/本文）
- [x] NS: テスト（複数カテゴリ/銘柄）

#### 6-8) SS / US（システム/運用ステータス通知）
- [x] SS: 通知モデル定義/パース/テスト
- [x] US: 通知モデル定義/パース/テスト
- [x] SS/US: 4-1 マスタ更新との関係整理（US: p_UC=11 / p_UU=1100 を検知し「マスタ再DL推奨」として露出、初期は自動DLしない）

## 7) API 抽象化レイヤ
- [x] Quote 抽象化（model.Quote / QuoteSnapshot）
- [x] Attributes / Order / Position / Balance 型
- [ ] 注文/建玉/残高の mapper 追加
- [ ] 口座/余力の抽象スナップショット拡充
- [ ] Facade API（薄いAPI + 抽象APIの二層）

## 8) 例 + テスト
- [x] CLI: auth-check
- [x] CLI: price-snapshot
- [x] CLI: request-read
- [x] CLI: order-read
- [x] CLI: order-submit
- [x] CLI: order-correct / order-cancel
- [x] CLI: order-crud
- [x] Event WS サンプル
- [x] CLI: board-stream
- [ ] httptest を使ったスモークテスト

## 決定事項
- Go module: github.com/ueebee/tachibanashi
- Go version: 最新（現環境は go1.24.2）
- API バージョン: v4r8
- demo/prod の base URL は固定で採用
  - prod: https://kabuka.e-shiten.jp/e_api_v4r8/
  - demo: https://demo-kabuka.e-shiten.jp/e_api_v4r8/
- Event API: Recv(ctx) を基本に、Stream(ctx) を補助で提供
- TokenStore: メモリ実装のみで開始
- 価格/数量の型: int64（Price/Quantity を model に定義）
