# Warikan 学習進捗

このファイルは、セッションをまたいで Claude Code に現在地を伝えるための記録である。
**Claude Code だけが更新し、学習者は触れない**(`CLAUDE.md` §8)。区切りごとに Claude Code が更新し、コミットとPRは学習者が行う。古い振り返りは削除せず、末尾の「振り返りの履歴」へ移す。

最終更新:2026-11-22

## 現在地

- 開始日:2026-11-17
- 現在:第1月 W1 後半
- 次にやること:
  1. `golangci-lint` を手元に入れて動かす(`brew install golangci-lint`。`golangci-lint --version` で `built with go1.27` 以上かを確認し、`gh release view --repo golangci/golangci-lint` で最新版を確認する)
  2. 設定ファイルを作り、CI(`.github/workflows/backend-ci.yml`)に追加する。CIでのバージョンは固定する
  3. Go言語基礎(エラー処理、interface、context)を、書いたコードを題材に確認する
  4. W1の振り返り(チェックリスト「CIが落ちるとマージできない理由」の口頭試問を含む)

## プロジェクトの状態

- リポジトリ:`Mosite/warikan`(公開)
- 構成(モノレポ):
  - `backend/`:Goモジュール `github.com/Mosite/warikan/backend`(Go 1.27.4、chi v5)
    - `cmd/api/main.go`:起動処理(環境変数 `PORT`、既定 8080、`ReadHeaderTimeout` 設定済み)
    - `internal/httpapi/health.go` と `health_test.go`:`GET /health` → `200 {"status":"ok"}`
  - `.github/workflows/backend-ci.yml`:PRと `main` へのpushで gofmt チェック、`go mod tidy -diff`、`go vet`、`go test -race` を実行。ジョブ名 `backend-test`
- マージ済みPR:#1〜#5(#3 `/health`、#5 CI)
- 本番URL:なし(W2で作る)
- 手元の環境:Go は Homebrew の 1.27.4 に一本化済み(旧 `/usr/local/go` は削除)。`gh` は認証済み。`gh` のエディタは vim

## リポジトリの設定(確認済み)

- マージ方式:squash のみ。コミットの1行目はPRタイトル、本文はPRの本文(`PR_TITLE` / `PR_BODY`)。マージ済みブランチは自動削除
- Secret scanning と Push protection:有効
- ルールセット `main-protection`(ID 27447374、Active、bypass list は空):削除禁止、force push 禁止、linear history 必須、PR必須(承認数0、squashのみ)、status check `backend-test` の成功必須(up to date 必須)
- 動作確認済み:`main` への直接 push は拒否される。CIが落ちたPRはマージできない(`--admin` も不可)

## 完了したマイルストーン

- 第1月W1:開発環境構築、Gitブランチ運用の取り決めとブランチ保護、Goプロジェクト雛形、`/health` とハンドラーテスト、GitHub Actions のCI(gofmt / vet / test)稼働。残りは `golangci-lint` と Go言語基礎の確認

## 今月のチェックリスト(第1月)

口頭試問で説明できた項目だけにチェックを付け、通った日付を添える。月が変わったら `CLAUDE.md` §7 から次の月の項目を写す。

- [ ] CIが落ちるとマージできない理由を説明できる(実演は済み。口頭試問はW1の振り返りで行う)
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
- [should] ディレクトリ構成・ブランチ運用・各種コマンド(起動、lint、テスト、単一テスト、マイグレーション)を `CLAUDE.md` に追記する。ディレクトリ構成・ブランチ運用・起動・テストは決まった。lint とマイグレーションが決まったら、学習者の依頼を受けてまとめて追記する

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

## ADR 一覧

まだない。書いたら「番号・題名・日付」を1行ずつ追加する。

## 直近の振り返り

まだない。週末に学習者がチャットで次の3問に答え、Claude Code が記録する。

1. 今週「なぜ」を説明できるようになったことは何か
2. 指摘された初心者っぽさのうち最も刺さったものは何か
3. 来週に持ち越す疑問は何か

## 振り返りの履歴

まだない。
