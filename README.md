# hatena-blog-org

はてなブログにorgファイルを投稿するためのCLIツールです。pandocを使用してorgファイルをmarkdownに変換し、はてなブログAtomPub APIを使用して投稿します。

## 必要な環境

- Go 1.18以上
- pandoc

## インストール

```bash
go build -o hatena-blog-org
```

## 使用方法

### 基本的な使用方法

```bash
./hatena-blog-org -file article.org -id your-hatena-id -key your-api-key -domain your-blog-domain
```

### オプション

- `-file`: 投稿するorgファイルのパス（必須）
- `-id`: はてなID（必須）
- `-key`: APIキー（必須）
- `-domain`: ブログドメイン（必須）
- `-category`: カテゴリー（任意）
- `-draft`: 下書きとして投稿（任意）
- `-config`: 設定ファイルのパス（任意）
- `-interactive`: 対話モード（任意）

### 設定ファイルの使用

設定ファイル（JSON形式）を使用して認証情報を保存できます：

```json
{
  "hatena_id": "your-hatena-id",
  "api_key": "your-api-key",
  "blog_domain": "your-blog-domain"
}
```

デフォルトの設定ファイルパス：
- `~/.config/hatena-blog-org/config.json`

### 対話モード

```bash
./hatena-blog-org -interactive
```

対話モードでは、必要な情報を順次入力できます。

## サンプルorgファイル

```org
* ブログ記事のタイトル

これは記事の内容です。

** 見出し2

- リスト項目1
- リスト項目2

** コードブロック

#+begin_src go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
#+end_src
```

## はてなブログAPIの設定

1. はてなブログの設定画面で「詳細設定」を開く
2. 「AtomPub」セクションでAPIを有効化
3. 「ルートエンドポイント」に表示されるURLを確認
4. アカウント設定でAPIキーを確認

## テスト

```bash
go test -v
```

## 注意事項

- APIキーは機密情報です。適切に管理してください
- 設定ファイルは適切な権限（600）で保存されます
- pandocがインストールされていない場合、変換テストはスキップされます