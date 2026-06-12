import { useEffect, useState } from "react";
import {
  fetchBookmarks,
  createBookmark,
  deleteBookmark,
  type Bookmark,
} from "./api";
import "./App.css";

// App はブックマーク管理アプリの画面全体。
function App() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
  const [query, setQuery] = useState("");
  const [url, setUrl] = useState("");
  const [title, setTitle] = useState("");
  const [memo, setMemo] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  // load は現在の検索語で一覧を読み込む。
  async function load() {
    setLoading(true);
    setError("");
    try {
      setBookmarks(await fetchBookmarks(query));
    } catch (e) {
      setError(e instanceof Error ? e.message : "読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  }

  // 検索語が変わるたびに一覧を読み直す。
  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  // handleAdd はフォームの内容でブックマークを追加する。
  async function handleAdd(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    try {
      await createBookmark({ url, title, memo });
      // 入力をクリアして一覧を更新する。
      setUrl("");
      setTitle("");
      setMemo("");
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "追加に失敗しました");
    }
  }

  // handleDelete は確認のうえブックマークを削除する。
  async function handleDelete(id: number) {
    if (!confirm("このブックマークを削除しますか？")) return;
    setError("");
    try {
      await deleteBookmark(id);
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "削除に失敗しました");
    }
  }

  return (
    <div className="container">
      <h1>ブックマーク管理</h1>

      {/* 追加フォーム */}
      <form className="add-form" onSubmit={handleAdd}>
        <input
          type="text"
          placeholder="タイトル"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        <input
          type="url"
          placeholder="https://..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
        <input
          type="text"
          placeholder="メモ（任意）"
          value={memo}
          onChange={(e) => setMemo(e.target.value)}
        />
        <button type="submit">追加</button>
      </form>

      {/* 検索 */}
      <input
        className="search"
        type="search"
        placeholder="タイトル・URL・メモで検索"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />

      {error && <p className="error">{error}</p>}
      {loading && <p className="muted">読み込み中...</p>}
      {!loading && bookmarks.length === 0 && (
        <p className="muted">ブックマークはまだありません。</p>
      )}

      {/* 一覧 */}
      <ul className="list">
        {bookmarks.map((b) => (
          <li key={b.id} className="item">
            <div className="item-main">
              <a href={b.url} target="_blank" rel="noreferrer">
                {b.title}
              </a>
              {b.memo && <p className="memo">{b.memo}</p>}
              <p className="url">{b.url}</p>
            </div>
            <button className="delete" onClick={() => handleDelete(b.id)}>
              削除
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
