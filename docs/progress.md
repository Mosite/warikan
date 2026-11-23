# Warikan 学習進捗

このファイルは、セッションをまたいで Claude Code に現在地を伝えるための記録である。
**Claude Code だけが更新し、学習者は触れない**(`CLAUDE.md` §8)。区切りごとに Claude Code が更新し、コミットとPRは学習者が行う。

最終更新:2026-11-23

## 現在地

- 開始日:2026-11-17
- 現在:第1月 W2 開始(W1 は 2026-11-23 に完了)
- 次にやること:
  1. Dockerマルチステージビルドで `backend/` のイメージを作り、手元で起動して `/health` を確かめる
  2. `/health` だけの状態で Fly.io へデプロイする(`fly.toml` の `kill_timeout` を `shutdownTimeout` 10秒より長くする)
  3. 以降は `CLAUDE.md` §7 W2 の順(CD、slog、Docker Compose の PostgreSQL、golang-migrate)

## プロジェクトの状態

- リポジトリ:`Mosite/warikan`(公開)
- 構成(モノレポ):
  - `backend/`:Goモジュール `github.com/Mosite/warikan/backend`(Go 1.27.4、chi v5)
    - `cmd/api/main.go`:起動処理(環境変数 `PORT`、既定 8080、`ReadHeaderTimeout` 設定済み)。グレースフルシャットダウン対応済み(`signal.NotifyContext` で SIGINT/SIGTERM を受け、`srv.Shutdown` を `shutdownTimeout` 10秒で実行。シャットダウン用の context は `context.Background()` から作る)。自動テストはなく、シグナルでの停止(終了コード0)とポート使用中での失敗(終了コード1)を手動で確認済み
    - `internal/httpapi/health.go` と `health_test.go`:`GET /health` → `200 {"status":"ok"}`(テストのリクエストは `httptest.NewRequestWithContext(t.Context(), ...)`)
    - `.golangci.yml`:golangci-lint v2 の設定。`standard` に `bodyclose`・`errorlint`・`gosec`・`noctx` を追加。errcheck は `check-blank` と `check-type-assertions` を有効化
  - `.github/workflows/backend-ci.yml`:PRと `main` へのpushで動く。ジョブは2つで並列実行
    - `backend-test`:gofmt チェック、`go mod tidy -diff`、`go vet`、`go test -race`
    - `backend-lint`:`golangci/golangci-lint-action@v9`、golangci-lint `v2.15.0` 固定、`working-directory: backend`
- マージ済みPR:#1〜#9(#3 `/health`、#5 CI、#7 golangci-lint、#9 グレースフルシャットダウン)
- 本番URL:なし(W2で作る)
- 手元の環境:Go は Homebrew の 1.27.4 に一本化済み(旧 `/usr/local/go` は削除)。golangci-lint は Homebrew の v2.15.0(built with go1.27.2。`go.mod` の 1.27.4 とは言語バージョンが同じなので問題なく動くことを確認済み)。`gh` は認証済み。`gh` のエディタは vim。プロンプトは Starship
- lint の実行:`backend/` で `golangci-lint run ./...`。設定の検証は `golangci-lint config verify`
- 使い捨ての実験は `~/go-sandbox/` で行う(リポジトリの外)。ビルドした実行ファイルも `go build -o ~/go-sandbox/api ./cmd/api` のように外へ出す

## リポジトリの設定(確認済み)

- マージ方式:squash のみ。コミットの1行目はPRタイトル、本文はPRの本文(`PR_TITLE` / `PR_BODY`)。マージ済みブランチは自動削除
- Secret scanning と Push protection:有効
- ルールセット `main-protection`(ID 27447374、Active、bypass list は空):削除禁止、force push 禁止、linear history 必須、PR必須(承認数0、squashのみ)、status check `backend-test` と `backend-lint` の成功必須(どちらも GitHub Actions、`integration_id` 15368 に紐付け。up to date 必須)。ジョブの `name:` を変えるときはルールセットも同時に更新する
- 動作確認済み:`main` への直接 push は拒否される。CIが落ちたPRはマージできない(`--admin` も不可)

## 完了したマイルストーン

- 第1月W1:開発環境構築、Gitブランチ運用の取り決めとブランチ保護、Goプロジェクト雛形、`/health` とハンドラーテスト、GitHub Actions のCI(gofmt / vet / test)稼働、golangci-lint の導入(手元とCI、必須チェック化。2026-11-23)、Go言語基礎の確認(エラー処理・interface・context。題材としてグレースフルシャットダウンを実装。2026-11-23)。W1 完了(2026-11-23)

## 今月のチェックリスト(第1月)

口頭試問で説明できた項目だけにチェックを付け、通った日付を添える。月が変わったら `CLAUDE.md` §7 から次の月の項目を写す。

- [x] CIが落ちるとマージできない理由を説明できる(2026-11-23 合格。初回は「GitHub Actions がマージを止めている」と答え、ワークフロー(実行と報告)とルールセット(マージの拒否)を混同していた。`integration_id` は報告者の GitHub App の ID で、報告者は GitHub が認証情報から決めるので偽装できない、まで確認した)
- [ ] セッション方式とJWTの違いと選定理由を説明できる
- [ ] CSRFが成立する条件と、自分の実装がそれをどう防いでいるかを説明できる
- [ ] マイグレーションを手書きSQLで往復(up/down)できる
- [ ] Dockerマルチステージビルドの各ステージの役割を説明できる

## 未解消の [must] 指摘

なし。

## CLAUDE.md の未対応事項

2026-11-17 の /init で出た改善提案のうち、残っているもの。対応したら行を消す。

- [should] 端数処理ルールを例示ではなく仕様として1つに決める(§3 ドメインルール)。第2月W2の債務計算より前に必要
- [should] 通貨の範囲(円のみか、グループごとに1通貨か)を決める(§3 ドメインルール)。第1月W3の支出テーブル設計より前に必要
- [should] ディレクトリ構成・ブランチ運用・各種コマンド(起動、lint、テスト、単一テスト、マイグレーション)を `CLAUDE.md` に追記する。ディレクトリ構成・ブランチ運用・起動・テスト・lint は決まった。マイグレーションが決まったら、学習者の依頼を受けてまとめて追記する

## 決定事項

ADR にするほどではないが、後のセッションが知っておくべき取り決めを書く。

- 2026-11-17:`docs/progress.md` の初版は Claude Code が作成した
- 2026-11-17:`CLAUDE.md` を学習者の依頼で Claude Code が修正した(ルール番号の振り直し、操作の線引きの明記、標準ヘッダー、進捗ファイルの読み込み、ADR第1号を第1月W2に追加)
- 2026-11-17:学習内容の見直しを `CLAUDE.md` に反映した(Claude Code が学習者の依頼で編集)。主な変更は次のとおり
  - モバイルは iOS を主(第4月〜第5月W2)、Android を従(第5月W3〜W4)とした。これは Claude Code による仮決めで、第4月の開始前なら入れ替えてよい
  - 必須と「余裕があれば」の範囲、および削る順序を §3 に定めた
  - 第1月のデプロイを W4 から W2 に前倒しし、認証を W3、支出登録を W4 に移した
  - Terraform は「余裕があれば」に下げた(Fly.io 公式プロバイダーが2024年3月にアーカイブされたため)
  - アーキテクチャ導入の前に、痛みを3行で書くことを条件にした(§5)
  - 予備週を第3月の後と第5月の後に1週ずつ置いた
- 2026-11-17:リポジトリは公開にした(GitHub Pro は使わない)。理由は、Actions が無料・無制限で、第4月の macOS ランナーのコストを気にせずに済むこと、シークレット管理を最初から厳格にする動機になること
- 2026-11-17:コミットの作者は `Mosite <Mosite@users.noreply.github.com>`。個人のメールアドレスを公開しないため
- 2026-11-17:ブランチ運用は GitHub Flow。ブランチ名は `feat/`・`fix/`・`docs/`・`chore/`・`ci/`・`test/`・`refactor/` の接頭辞を付ける。マージは squash merge のみ。PRタイトルは Conventional Commits の接頭辞を付ける
- 2026-11-17:個人用のメモ(`self-diary.md`)は `.gitignore` ではなく `.git/info/exclude` で除外する
- 2026-11-17:監修の進め方を変更した(`CLAUDE.md` §2)
  - ペースはゆっくり。1回に渡す作業は1〜3個で、作業ごとに理由を添える(ルール5)
  - 調査のためのコマンドも学習者が実行し、結果を貼り付ける。Claude Code はファイルを読むことだけ行う(ルール2)
  - 学習者に「なぜ」を書かせる練習はしない。PRの本文やコミットメッセージの文章は問わない。例外は ADR のみ(ルール5、§6)
- 2026-11-19:コードは Claude Code が完成形を示し、学習者が写経する方式に変更した(ルール1・3・6)
- 2026-11-19:モノレポ構成とし、Go のコードは `backend/` に置く。実行ファイルの入口は `cmd/<名前>/`、非公開のパッケージは `internal/` の下
- 2026-11-22:CI のランナーは `ubuntu-latest`。`staging` ブランチは作らない(環境を分ける必要が出たら、ブランチではなく GitHub Environments で分ける)
- 2026-11-22:`docs/progress.md` は Claude Code だけが更新し、学習者は触れない(`CLAUDE.md` §2 ルール2、§8)
- 2026-11-23:golangci-lint は CI で `v2.15.0` に固定する(`latest` にすると、コードを変えていないのにCIが落ちることがあるため)。上げるときはPRで意図して上げ、手元も揃える
- 2026-11-23:DB向けの linter(`sqlclosecheck`・`rowserrcheck` など)は、DBのコードが入るW2以降に追加する
- 2026-11-23:lint の指摘は除外設定で逃がさず、コードを直す。テストのリクエストには `t.Context()` を渡す
- 2026-11-23:CIの checkout と setup-go は `test` と `lint` の2ジョブで重複しているが、許容する。共通化は3つ目のジョブが増えて痛みが出たときに検討する
- 2026-11-23:`shutdownTimeout` は10秒。W2 で `fly.toml` を書くとき、Fly.io が強制終了するまでの猶予(`kill_timeout`)をこれより長くする
- 2026-11-23:`cmd/api` には自動テストを書かない(起動の組み立てだけで、テストの仕組みが大がかりになるため)。`run()` に分岐が増えて手動確認がつらくなったら再検討する
- 2026-11-23:週末の振り返り3問(`CLAUDE.md` §8)は行わない。学習者の判断で、時間がかかるため完了を優先する。チェックリストの口頭試問は続ける。代わりに、Claude Code がセッション中に見えたつまずきと理解を、このファイルの「つまずいた点と理解した点」に記録する(学習者の作業は増やさない)。同日、学習者の依頼で `CLAUDE.md` を修正した(§2 ルール9、§7 第1月W4・第3月W4、§8 のセッション開始時・進捗の記録・週末の振り返り、§10 の依頼例)

### Go言語基礎の確認で、つまずいた点と理解した点(2026-11-23)

後の口頭試問や復習で掘り下げる材料として残す。

- エラー処理:`%v` でラップすると見た目は変わらず鎖だけが失われる点を、最初は「メッセージが変わる」と誤解していた。`Error()` が内側の `Error()` を連結して組み立てること、`%T` で `*fmt.wrapError` と `*errors.errorString` に分かれることを実験で確認して理解した
- interface:実験で `hello` を `hello()` と写し間違えた(関数値と関数呼び出しの区別)。`http.HandlerFunc` が型変換によるアダプターであること、chi の `Get` が `HandlerFunc` を受け取るので包まずに渡せることを理解した。interface は「要求する側」、具体的な型は「実装する側」という用語の区別を確認した
- context:キャンセルは協力的で、受け取る側が `ctx.Done()` を見なければ止まらないことを理解した
- その他に質問が出た点:`:=` と `=` の違い、チャネル演算子 `<-`、`go build` の引数、`git push -u` と `gh pr create` の推測の違い

## ADR 一覧

まだない。書いたら「番号・題名・日付」を1行ずつ追加する。

## 振り返り

週末の振り返り3問は 2026-11-23 にやめた(「決定事項」を参照)。過去の回答はない。
