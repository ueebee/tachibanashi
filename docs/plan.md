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
- [ ] CLMEventDownload 要求モデル定義（共通項目 + sCLMID）
- [ ] CLMEventDownload 送信メソッド（MASTER 仮想URL）
- [ ] 受信フレームの分割・JSON デコード
- [ ] 受信メッセージのディスパッチ（sCLMID でルーティング）
- [ ] 初期ダウンロード受信ループ（CLMEventDownloadComplete まで）
- [ ] 初期完了後の更新通知処理（UPDATE 差分反映）
- [ ] TODO: 初期完了後も接続を維持して更新通知を受け続けるモード（将来対応）
- [ ] エラー応答（p_errno/p_err）処理
- [ ] 文字コード正規化（Shift_JIS → UTF-8 前提の確認）
- [ ] テスト: 連続配信のパース / 完了通知 / 更新反映

#### 4-1-2) マスタストア / 更新処理
- [ ] MasterStore インターフェース（Get/Upsert/Delete/All）
- [ ] インメモリ実装（種別ごとに map + index）
- [ ] 主キー抽出関数（マスタ種別ごと）
- [ ] 更新通番/削除フラグの優先ルール（更新通番優先、同値は更新日時）
- [ ] 参照用インデックス（銘柄コード/市場/口座区分など）
- [ ] スナップショット取得 API（読み取り用）
- [ ] テスト: Upsert/Delete とインデックス整合性
- [ ] TODO: 変更通知（購読型）インターフェース

#### 4-1-3) 運用系マスタ（spec 2-1〜2-4）
- [ ] CLMSystemStatus: 項目一覧抽出 / 主キー決定
- [ ] CLMSystemStatus: モデル定義 / パース / 格納 / テスト
- [ ] CLMDateZyouhou: 項目一覧抽出 / 主キー決定
- [ ] CLMDateZyouhou: モデル定義 / パース / 格納 / テスト
- [ ] CLMUnyouStatus: 項目一覧抽出 / 主キー決定
- [ ] CLMUnyouStatus: モデル定義 / パース / 格納 / テスト
- [ ] CLMUnyouStatusKabu: 項目一覧抽出 / 主キー決定
- [ ] CLMUnyouStatusKabu: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMUnyouStatusHasei: 項目一覧抽出 / 主キー決定
- [ ] CLMUnyouStatusHasei: モデル定義 / UPDATE 反映 / テスト

#### 4-1-4) 銘柄系マスタ（spec 2-6〜2-11）
- [ ] CLMIssueMstKabu: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueMstKabu: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMIssueSizyouMstKabu: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueSizyouMstKabu: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMIssueSizyouKiseiKabu: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueSizyouKiseiKabu: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMIssueMstSak: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueMstSak: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMIssueMstOp: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueMstOp: モデル定義 / UPDATE 反映 / テスト
- [ ] CLMIssueSizyouKiseiHasei: 項目一覧抽出 / 主キー決定
- [ ] CLMIssueSizyouKiseiHasei: モデル定義 / UPDATE 反映 / テスト

#### 4-1-5) 静的マスタ（spec 2-12〜2-15）
- [ ] CLMYobine: 項目一覧抽出 / 主キー決定
- [ ] CLMYobine: モデル定義 / パース / 格納 / テスト
- [ ] CLMDaiyouKakeme: 項目一覧抽出 / 主キー決定
- [ ] CLMDaiyouKakeme: モデル定義 / パース / 格納 / テスト
- [ ] CLMHosyoukinMst: 項目一覧抽出 / 主キー決定
- [ ] CLMHosyoukinMst: モデル定義 / パース / 格納 / テスト
- [ ] CLMOrderErrReason: 項目一覧抽出 / 主キー決定
- [ ] CLMOrderErrReason: モデル定義 / パース / 格納 / テスト
- [ ] CLMEventDownloadComplete: 受信 / 状態管理 / テスト

### 4-2) マスタ問合取得
- [ ] 共通: Request/Response ラッパ（MASTER 仮想URL）
- [ ] 共通: 文字コード正規化（ニュース p_HDL/p_TX の Shift_JIS）
- [ ] 共通: sTargetIssueCode 最大120件のバリデーション
- [ ] 共通: 空文字/0 の扱い（値なし）方針整理
- [ ] 共通: ユニットテスト（リクエスト/レスポンス）

#### 4-2-1) CLMMfdsGetMasterData（マスタ情報問合）
- [ ] 要求モデル: sTargetCLMID / sTargetColumn
- [ ] 対象CLMID一覧の整理（CLMIssueMstKabu/CLMIssueSizyouMstKabu/CLMIssueMstSak/CLMIssueMstOp/CLMIssueMstOther/CLMIssueMstIndex/CLMIssueMstFx/CLMOrderErrReason/CLMDateZyouhou）
- [ ] 応答モデル: CLMID 名ごとの配列（動的キー）パース
- [ ] 取得列指定のフィルタ（指定列のみ抽出）
- [ ] テスト: 2種類CLMID混在レスポンスのパース

#### 4-2-2) CLMMfdsGetNewsHead（ニュースヘッダ）
- [ ] 要求モデル: p_CG / p_IS / p_DT_FROM / p_DT_TO / p_REC_OFST / p_REC_LIMT
- [ ] レコード取得条件（AND 条件）/ p_REC_MAX 取り扱い
- [ ] 応答モデル: aCLMMfdsNewsHead の配列パース
- [ ] 文字コード: p_HDL デコード
- [ ] テスト: 検索条件付きの要求とレスポンス

#### 4-2-3) CLMMfdsGetNewsBody（ニュース本文）
- [ ] 要求モデル: p_ID
- [ ] 応答モデル: aCLMMfdsNewsBody の配列パース
- [ ] 文字コード: p_HDL / p_TX デコード
- [ ] テスト: 1件レスポンスのパース

#### 4-2-4) CLMMfdsGetIssueDetail（銘柄詳細）
- [ ] 要求モデル: sTargetIssueCode（最大120件）
- [ ] 応答モデル: aCLMMfdsIssueDetail の配列パース
- [ ] テスト: 主要フィールドの型/値パース

#### 4-2-5) CLMMfdsGetSyoukinZan（証金残）
- [ ] 要求モデル: sTargetIssueCode（最大120件）
- [ ] 応答モデル: aCLMMfdsSyoukinZan の配列パース
- [ ] テスト: 数値/日付フィールドのパース

#### 4-2-6) CLMMfdsGetShinyouZan（信用残）
- [ ] 要求モデル: sTargetIssueCode（最大120件）
- [ ] 応答モデル: aCLMMfdsShinyouZan の配列パース
- [ ] テスト: 数値/日付フィールドのパース

#### 4-2-7) CLMMfdsGetHibuInfo（逆日歩）
- [ ] 要求モデル: sTargetIssueCode（最大120件）
- [ ] 応答モデル: aCLMMfdsHibuInfo の配列パース
- [ ] テスト: 逆日歩フィールドのパース

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
- [ ] URL パラメータ仕様の整理（p_rid/p_board_no/p_gyou_no/p_issue_code/p_mkt_code/p_eno/p_evt_cmd）
- [ ] p_evt_cmd 対応範囲の決定（ST/KP/FD/EC/NS/SS/US のみ）
- [ ] p_eno の再送ルール整理（EC/NS/SS/US）
- [ ] Base64 対象項目の確認（Shift_JIS 代替）
- [ ] RR/FC 非公開通知は未対応と明記

#### 6-1) URL / パラメータ検証
- [ ] WS URL ビルダ（パラメータ省略とデフォルト）
- [ ] パラメータ検証（p_evt_cmd/最大120銘柄/ボード組合せ）
- [ ] board/row 指定の組合せサポート（p_board_no/p_gyou_no/p_issue_code/p_mkt_code）

#### 6-2) 接続・再接続・セッション
- [ ] 接続/切断/再接続（コンテキスト停止・バックオフ）
- [ ] keepalive（KP 受信/送信扱い）
- [ ] 1セッション制約（重複接続ガード）
- [ ] p_eno レジューム（再送開始番号の管理）

#### 6-3) フレームデコード / 共通処理
- [ ] 受信フレームのデコード（^A/^B/^C 区切り）
- [ ] 共通項目パース（p_no/p_date/p_cmd）
- [ ] Base64 デコード（WebSocket は Shift_JIS を扱えないため）
- [ ] イベントディスパッチ（p_cmd でルーティング）
- [ ] 共通テスト（区切りパース/Base64/必須項目）

#### 6-4) ST / KP
- [ ] ST: 通知モデル定義（必須項目）
- [ ] ST: パース/テスト
- [ ] KP: 通知モデル定義
- [ ] KP: パース/テスト

#### 6-5) FD（時価情報）
- [ ] FD: 情報コード一覧（p_*/t_*）のパース
- [ ] FD: 初回スナップショット/差分更新の扱い
- [ ] FD: Quote/QuoteSnapshot 変換
- [ ] FD: テスト（初回/差分/複数銘柄）

#### 6-6) EC（注文約定通知）
- [ ] EC: 通知モデル定義（親注文番号/注文種別含む）
- [ ] EC: パース/テスト
- [ ] EC: Order/Execution へのマッピング

#### 6-7) NS（ニュース通知）
- [ ] NS: 通知モデル定義
- [ ] NS: Base64 → 文字列復号（見出し/本文）
- [ ] NS: テスト（複数カテゴリ/銘柄）

#### 6-8) SS / US（システム/運用ステータス通知）
- [ ] SS: 通知モデル定義/パース/テスト
- [ ] US: 通知モデル定義/パース/テスト
- [ ] SS/US: 4-1 マスタ更新との関係整理

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
- [ ] Event WS サンプル
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
