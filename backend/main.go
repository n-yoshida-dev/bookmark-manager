// ブックマーク管理アプリのAPIサーバー。
// Go標準ライブラリ net/http と pure-Go の SQLite ドライバ（modernc.org/sqlite）で構成する。
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Bookmark は1件のブックマークを表す。
type Bookmark struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Memo      string `json:"memo"`
	CreatedAt string `json:"created_at"`
}

// db はアプリ全体で共有するDB接続。
var db *sql.DB

// getenv は環境変数を読み、未設定なら既定値を返す。
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// main はDBを初期化し、ルーティングを登録してサーバーを起動する。
func main() {
	dbPath := getenv("DB_PATH", "./bookmarks.db")
	port := getenv("PORT", "8080")

	if err := initDB(dbPath); err != nil {
		log.Fatalf("DB初期化に失敗: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	// Go 1.22+ のパターンマッチでメソッド＋パスを登録する。
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/bookmarks", handleList)
	mux.HandleFunc("POST /api/bookmarks", handleCreate)
	mux.HandleFunc("DELETE /api/bookmarks/{id}", handleDelete)

	log.Printf("サーバー起動: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatalf("サーバー起動に失敗: %v", err)
	}
}

// initDB はDBを開き、テーブルが無ければ作成する。
func initDB(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bookmarks (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			url        TEXT NOT NULL,
			title      TEXT NOT NULL,
			memo       TEXT,
			created_at TEXT NOT NULL
		)
	`)
	return err
}

// withCORS はローカル開発用にCORSヘッダを付与するミドルウェア。
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// プリフライト（OPTIONS）にはここで応答する。
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeJSON はJSONレスポンスを書き出す共通ヘルパー。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("JSONエンコードに失敗: %v", err)
		}
	}
}

// handleHealth はヘルスチェックに応答する。
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleList はブックマーク一覧を返す。クエリ q があれば部分一致で絞り込む。
func handleList(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows *sql.Rows
	var err error
	if q == "" {
		rows, err = db.Query(`SELECT id, url, title, memo, created_at FROM bookmarks ORDER BY id DESC`)
	} else {
		like := "%" + q + "%"
		rows, err = db.Query(
			`SELECT id, url, title, memo, created_at FROM bookmarks
			 WHERE title LIKE ? OR url LIKE ? OR memo LIKE ?
			 ORDER BY id DESC`, like, like, like)
	}
	if err != nil {
		log.Printf("一覧取得に失敗: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "一覧の取得に失敗しました"})
		return
	}
	defer rows.Close()

	// nil ではなく空スライスで初期化し、0件でも [] を返す。
	list := []Bookmark{}
	for rows.Next() {
		var b Bookmark
		var memo sql.NullString
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &memo, &b.CreatedAt); err != nil {
			log.Printf("行の読み取りに失敗: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "一覧の取得に失敗しました"})
			return
		}
		b.Memo = memo.String
		list = append(list, b)
	}
	writeJSON(w, http.StatusOK, list)
}

// handleCreate はブックマークを1件追加する。
func handleCreate(w http.ResponseWriter, r *http.Request) {
	var in Bookmark
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
		return
	}
	in.URL = strings.TrimSpace(in.URL)
	in.Title = strings.TrimSpace(in.Title)
	// URLとタイトルは必須。
	if in.URL == "" || in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "URLとタイトルは必須です"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	res, err := db.Exec(
		`INSERT INTO bookmarks (url, title, memo, created_at) VALUES (?, ?, ?, ?)`,
		in.URL, in.Title, in.Memo, now)
	if err != nil {
		log.Printf("追加に失敗: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "追加に失敗しました"})
		return
	}
	id, _ := res.LastInsertId()
	in.ID = id
	in.CreatedAt = now
	writeJSON(w, http.StatusCreated, in)
}

// handleDelete はパスパラメータ id のブックマークを削除する。
func handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
		return
	}
	res, err := db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		log.Printf("削除に失敗: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "削除に失敗しました"})
		return
	}
	// 対象が無ければ404を返す。
	if n, _ := res.RowsAffected(); n == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "指定のブックマークが見つかりません"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
