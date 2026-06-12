# SPEC — ブックマーク管理アプリ

## 技術スタック
- フロントエンド：React + TypeScript + Vite
- バックエンド：Go（標準ライブラリ net/http）
- DB：SQLite（pure-Go ドライバ modernc.org/sqlite、cgo不要）

## データ構造（bookmarks テーブル）
| カラム | 型 | 説明 |
|--------|----|----|
| id | INTEGER PRIMARY KEY AUTOINCREMENT | ID |
| url | TEXT NOT NULL | ブックマークURL |
| title | TEXT NOT NULL | タイトル |
| memo | TEXT | メモ（任意） |
| created_at | TEXT NOT NULL | 作成日時（RFC3339） |

## API 仕様（ベースパス /api）
| メソッド | パス | 説明 | リクエスト | レスポンス |
|---------|------|------|-----------|-----------|
| GET | /api/bookmarks?q={検索語} | 一覧（q指定時はtitle/url/memo部分一致） | - | `Bookmark[]` |
| POST | /api/bookmarks | 追加 | `{url, title, memo}` | `Bookmark`（201） |
| DELETE | /api/bookmarks/{id} | 削除 | - | 204 |
| GET | /api/health | ヘルスチェック | - | `{"status":"ok"}` |

### Bookmark JSON
```json
{ "id": 1, "url": "https://example.com", "title": "例", "memo": "メモ", "created_at": "2026-06-12T23:00:00Z" }
```

## 設定（環境変数）
| 変数 | 既定値 | 説明 |
|------|--------|------|
| PORT | 8080 | サーバーのポート |
| DB_PATH | ./bookmarks.db | SQLiteファイルのパス |

## CORS
- ローカル開発のためフロント（http://localhost:5173）からのアクセスを許可
