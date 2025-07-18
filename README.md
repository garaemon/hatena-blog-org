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

## Orgファイルの書式

### タイトルとカテゴリの指定

タイトルとカテゴリはOrgファイルの冒頭で指定できます：

```org
#+title: ブログ記事のタイトル
#+filetags: :プログラミング:Go:技術ブログ:

記事の内容をここに書きます。

* 見出し1

これは記事の内容です。

** 見出し2

- リスト項目1
- リスト項目2
```

### タイトルの指定方法

- `#+title:` ディレクティブでタイトルを指定（大文字小文字は区別しません）
- 指定しない場合は「Untitled」になります

```org
#+title: 記事のタイトル
# または
#+TITLE: 記事のタイトル
```

### カテゴリの指定方法

- `#+filetags:` ディレクティブでカテゴリを指定（大文字小文字は区別しません）
- 複数の区切り文字に対応：スペース、コロン、混合形式

```org
# スペース区切り
#+filetags: カテゴリ1 カテゴリ2 カテゴリ3

# コロン区切り（Orgモード標準）
#+filetags: :カテゴリ1:カテゴリ2:カテゴリ3:

# 混合形式
#+filetags: カテゴリ1:カテゴリ2 カテゴリ3:カテゴリ4
```

- `-category`オプションで指定したカテゴリも追加されます

## サンプルorgファイル

```org
#+title: Goでのブログ投稿ツール作成
#+filetags: :Go:CLI:ブログ:技術:

Goを使ってはてなブログに投稿するCLIツールを作成しました。

* 概要

このツールの特徴は以下の通りです。

** 機能

- orgファイルの自動変換
- カテゴリの自動設定
- 下書き投稿対応

** コード例

#+begin_src go
package main

import "fmt"

func main() {
    fmt.Println("Hello, はてなブログ!")
}
#+end_src

* まとめ

便利なツールができました。
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