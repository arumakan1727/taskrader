# TaskRader
TaskRader は学情・EdStem・Teamsのの課題情報を一括して取得するアプリケーションです。
コマンドラインによる操作も、WebブラウザによるGUI操作もできます。

## 動作要件
Linux x86-64 でのみ動作を確認しました。

Teams の課題取得には Selenium WebDriver を使用します。
Ubuntu をご利用なら下記コマンドを実行すれば WebDriver をインストールできます。

```
sudo apt update && sudo apt install -y chromium-chromedriver
```

## ビルド方法
ビルドには Go コンパイラと Yarn (node モジュールのパッケージマネージャ) が必要です。

- Go のインストール: https://go.dev/doc/install
- Yarn のインストール: https://classic.yarnpkg.com/lang/en/docs/install
    - Yarn のインストールには nodejs, npm が必要です。  
        nodejs, npm のインストール方法はいくつかの方法があります:
        - Node バージョンマネージャ nvm を用いる
        - apt パッケージマネージャを用いる (`apt install nodejs`)


この README.md が配置されているディレクトリで、下記コマンドを実行すればビルドできます。
```
make
```
`bin/taskrader` に実行可能バイナリが生成されます。



## 使い方
GUI で操作したい場合は、`gui` サブコマンドを実行すればよいです。
ブラウザが起動して GUI 操作できます。
```
taskrader gui
```

コマンドラインで操作するには、2つのサブコマンド 'login', 'list' を使用します。
詳細は `taskrader login --help` および `taskrader list --help` を実行してください。

