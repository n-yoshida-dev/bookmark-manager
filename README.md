# ブックマーク管理アプリ

URLとメモを保存・一覧・検索・削除できるシンプルなフルスタックアプリ。
Claude Code でのバイブコーディング学習用（Go + React）。

## 構成
- `backend/` … Go（net/http）+ SQLite の REST API
- `frontend/` … React + TypeScript + Vite の UI

## ローカル起動

### バックエンド
```bash
cd backend
go run .
# http://localhost:8080 で起動
```

### フロントエンド
```bash
cd frontend
npm install
npm run dev
# http://localhost:5173 で起動
```

## 環境変数
- バックエンド：`PORT`（既定8080）、`DB_PATH`（既定 ./bookmarks.db）
- フロント：`VITE_API_BASE`（既定 http://localhost:8080）。公開時はデプロイ先バックエンドのURLを指定する

詳細な仕様は [SPEC.md](SPEC.md) を参照。
