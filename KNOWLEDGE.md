# KNOWLEDGE — ブックマーク管理アプリ

## 設計判断の記録
- **SQLite ドライバに modernc.org/sqlite を採用**：cgo不要のpure-Go実装。gcc等のCツールチェーンが無くても `go run` できる。公開ホスティングのビルドも楽。
- **バックエンドは標準ライブラリ net/http のみ**：最初の1本でフレームワーク依存を増やさない方針。ルーティングは Go 1.22+ の `http.ServeMux` のパターンマッチ（`GET /api/bookmarks` 等）を使用。

## ハマったこと
（随時追記）

## 2本目以降の拡張候補
- ユーザー認証（公開するなら必須）
- OGP取得でタイトル自動補完・サムネイル表示
- タグ／フォルダ分け
- PostgreSQL へ移行（複数ユーザー対応時）

## デプロイメモ
- バックエンド：Render / Fly.io の無料枠（git push連携）が候補
- フロント：Cloudflare Pages / Vercel
- SQLiteはコンテナ再起動で消える点に注意（永続ディスク or Postgres移行が必要）
