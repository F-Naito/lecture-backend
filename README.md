# lecture-backend

レクチャー用に作成した、**Go + SQLite のシンプルな REST API** です。  
本（Book）を管理するための CRUD API を提供します。

- 学習・教材目的
- ローカル実行前提
- 本番利用は想定していません

---

## 概要

この API は、以下の操作を提供します。

- 本の一覧取得
- 本の詳細取得
- 本の作成
- 本の更新
- 本の削除

REST API / CRUD の基本を学ぶための最小構成です。

---

## データモデル

### Book

```ts
Book {
  id: number
  title: string
  author: string
  createdAt: string
  updatedAt: string
}