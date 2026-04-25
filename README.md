# atcoder-cui

[![CI](https://github.com/Atnuhs/atcoder-cui/actions/workflows/ci.yml/badge.svg)](https://github.com/Atnuhs/atcoder-cui/actions/workflows/ci.yml)

Goによる競技プログラミング用ライブラリ＆作業スペース。バンドルには[go-bundler](https://github.com/Atnuhs/go-bundler)を使用。

## 構成

```text
acl/
├── main.go       # 解答を書く場所
├── watch.sh      # ファイル監視・自動ビルド・対話実行
├── *.go          # ライブラリ
└── *_test.go     # ライブラリのテスト
```

## 使い方

```bash
cd acl
./watch.sh
```

`main.go` を保存するたびに自動でバンドル＆ビルド。ビルド成功後は入力を貼り付けると即座に実行結果が表示される。
