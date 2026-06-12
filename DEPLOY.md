# デプロイ手順

フロント = Cloudflare Pages、バックエンド = Render の構成で公開する。
どちらも GitHub 連携で「push したら自動で再デプロイ」される。

---

## 手順1：バックエンドを Render で公開

1. https://render.com にアクセスし、**GitHub アカウントでサインアップ**（無料）
2. ダッシュボードで **New > Blueprint** を選ぶ
3. `n-yoshida-dev/bookmark-manager` リポジトリを選択
4. リポジトリ直下の `render.yaml` が自動で読み込まれる → **Apply**
5. デプロイ完了後、`https://bookmark-manager-api-xxxx.onrender.com` のような **公開URL**が発行される
6. ブラウザで `その URL + /api/health` を開き `{"status":"ok"}` が出れば成功

> メモ：無料枠は15分アクセスが無いとスリープする。次回アクセス時に数十秒かかるが正常。

---

## 手順2：フロントを Cloudflare Pages で公開

1. https://dash.cloudflare.com にアクセスし、サインアップ（無料）
2. **Workers & Pages > Create > Pages > Connect to Git** で同じリポジトリを選択
3. ビルド設定を以下にする：
   - **Framework preset**: Vite
   - **Build command**: `cd frontend && npm install && npm run build`
   - **Build output directory**: `frontend/dist`
4. **環境変数**に以下を追加（手順1で得たバックエンドURL）：
   - `VITE_API_BASE` = `https://bookmark-manager-api-xxxx.onrender.com`
5. **Save and Deploy** → `https://bookmark-manager-xxxx.pages.dev` が発行される
6. その URL を開き、ブックマークの追加・検索・削除が動けば**リリース完了**

---

## つまずいたら

- フロントから追加できない → CORS。バックエンドは全オリジン許可済みなので、まず `VITE_API_BASE` のURLが正しいか確認
- 登録したデータが消える → SQLite が無料枠で揮発するため（[KNOWLEDGE.md](KNOWLEDGE.md) 参照）。残したいなら Render の PostgreSQL へ移行
