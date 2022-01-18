# Contribution Guide

## ディレクトリ構成

```
./
├─ assignment/    ... 課題の統合取得処理
│  ├── concurrency_test.go ... 並行処理のテスト
│  ├── fetch.go            ... 課題取得処理; clients/{gakujo,edstem,teams} モジュールを呼び出す形になるはず
│  ├── fetch_test.go       ... 課題取得処理の動作確認用
│  └── model.go            ... 課題の共通モデル
├─ cmd/            ... コマンドの定義
│  ├── list.go
│  ├── login.go
│  └── root.go
├─ cred/        ... 認証情報; credentialの略
│  └── credential.go
├─ clients/     ... 通信処理を行うクライアント
│  ├─ edstem/
│  │  ├── login.go
│  │  └── model.go
│  ├─ gakujo/
│  │  ├── login.go
│  │  └── model.go
│  └─ teams/
│     ├── login.go
│     └── model.go
├── env.sample
├── .env
├── go.mod
├── go.sum
└── main.go
```


## ブランチの命名
他のブランチ名と重複しなければOKです

- 対応する issue があれば `#{issue番号}-{適当な名前}` 
- なければ `{適当な名前}`
- さらに、その先頭に以下の接頭辞をつける
    - 新機能 or 機能向上 の場合: `feat/`
    - バグ修正の場合: `bugfix/`
    - 小さな修正の場合: `fix/`
    - リファクタリング: `refactor/`

例:
- `feat/#1-gakujo-fetch`
- `fix/spell-miss`


## その他 Tips

### `go test` でテスト成功・失敗に関わらず print の内容を表示したい
テストが成功した場合は print の結果が表示されない。
`-v` をつければ、成功した時も表示される。


### `go test` のキャッシュを無視して再テストしたい
`go test` は前回のテスト結果を保存して、使いまわす事がある。
これを阻止したい場合、2とおりの解決方法がある:

- `go clean -testcache` を実行してから `go test` をする
- `go test -count=1` を実行する
