// バックエンドAPIとの通信をまとめたモジュール。
// APIのベースURLは環境変数 VITE_API_BASE で差し替え可能（未設定ならローカルの8080）。

// Bookmark はAPIが返す1件のブックマーク。
export type Bookmark = {
  id: number;
  url: string;
  title: string;
  memo: string;
  created_at: string;
};

// NewBookmark は追加時に送る入力データ。
export type NewBookmark = {
  url: string;
  title: string;
  memo: string;
};

const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

// fetchBookmarks は一覧を取得する。検索語 q があれば絞り込む。
export async function fetchBookmarks(q: string): Promise<Bookmark[]> {
  const url = q
    ? `${API_BASE}/api/bookmarks?q=${encodeURIComponent(q)}`
    : `${API_BASE}/api/bookmarks`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error("一覧の取得に失敗しました");
  }
  return res.json();
}

// createBookmark はブックマークを1件追加する。
export async function createBookmark(input: NewBookmark): Promise<Bookmark> {
  const res = await fetch(`${API_BASE}/api/bookmarks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    // サーバーが返すエラーメッセージがあれば拾う。
    const data = await res.json().catch(() => null);
    throw new Error(data?.error ?? "追加に失敗しました");
  }
  return res.json();
}

// deleteBookmark は指定IDのブックマークを削除する。
export async function deleteBookmark(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/api/bookmarks/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error("削除に失敗しました");
  }
}
