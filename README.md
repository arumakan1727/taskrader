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


## 使い方
GUI で操作したい場合は、`gui` サブコマンドを実行すればよいです。
```
taskrader gui
```

コマンドラインで操作するには、2つのサブコマンド 'login', 'list' を使用します。
詳細は

```
taskrader login --help
```
および
```
taskrader list --help
```

参照してください。
