# KNOWLEDGE — ブックマーク管理アプリ

## 設計判断の記録
- **SQLite ドライバに modernc.org/sqlite を採用**：cgo不要のpure-Go実装。gcc等のCツールチェーンが無くても `go run` できる。公開ホスティングのビルドも楽。
- **バックエンドは標準ライブラリ net/http のみ**：最初の1本でフレームワーク依存を増やさない方針。ルーティングは Go 1.22+ の `http.ServeMux` のパターンマッチ（`GET /api/bookmarks` 等）を使用。

## ハマったこと
- **Cloudflareがビルド前のソースを配信して画面が真っ白**（2026-06-15）：`frontend/index.html`（`/src/main.tsx` を読む開発用）が配信されていた。原因は配信先が `frontend/dist` を指していなかったこと。リポジトリ直下に `wrangler.jsonc` を置き `assets.directory: ./frontend/dist` と `not_found_handling: single-page-application` を指定して解決。
- **VITE_API_BASE はビルド時に焼き込まれる**：環境変数を後から設定しても、再ビルドしないと反映されない。`VITE_` 系は実行時ではなくビルド時に値がJSへ埋め込まれる。設定変更後は main へpush（または再ビルド）が必要。
- **Cloudflare は Pages ではなく Workers（静的アセット）でデプロイされた**：Git連携で自動構成すると Workers になることがある。静的フロント配信としては問題なく動く。「変数とシークレット」欄は実行時用で静的アセットのみのWorkerでは使えず、ビルド変数は「設定 > ビルド > 変数とシークレット」に入れる。

## 2本目以降の拡張候補
- ユーザー認証（公開するなら必須）
- OGP取得でタイトル自動補完・サムネイル表示
- タグ／フォルダ分け
- PostgreSQL へ移行（複数ユーザー対応時）

## デプロイメモ
- バックエンド：Render / Fly.io の無料枠（git push連携）が候補
- フロント：Cloudflare Pages / Vercel
- SQLiteはコンテナ再起動で消える点に注意（永続ディスク or Postgres移行が必要）
