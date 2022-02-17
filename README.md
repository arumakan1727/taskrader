# TaskRader
TaskRader は静岡大学の学務情報システム・[EdStem](https://edstem.org)・[Microsoft Teams](https://teams.microsoft.com/) の課題情報を一括して取得するアプリケーションです。
コマンドラインによる操作も、WebブラウザによるGUI操作もできます。

<p align="center"><img alt="cmd-hero" src="https://user-images.githubusercontent.com/33191176/154385178-f7174a1a-3b12-407a-be77-a9dad36dbeff.png" /></p>

<p align="center"><img alt="cmd-list2" src="https://user-images.githubusercontent.com/33191176/154385237-36a13059-51f1-42f1-83fb-4ae918fc45c6.png" /></p>

<p align="center"><img alt="page-list" src="https://user-images.githubusercontent.com/33191176/154385343-3619f071-3570-44a8-95b2-0e0460bc43f5.png" /></p>

## Why TaskRader
現在の課題を全て把握するには、複数の学情やTeams等Webページを確認する必要があります。
TaskRader を使えば、一つの画面で未提出課題の一覧を一括把握できます。

なお TaskRader は課題を統合取得して一覧表示する以外の機能はありません。

## 使用方法
CLIとGUIそれぞれの操作方法は以下のとおりです。

### CLI操作方法

#### Step 1. ログイン用の認証情報を登録する
サブコマンド `login` を用いてユーザ名やパスワード等を taskrader に登録してください。
`login` の引数にはプラットフォーム名 ( `gakujo` | `edstem` | `teams` ) を指定する必要があります。

例: EdStem へログインする場合

```
taskrader login edstem
```

<p align="center"><img alt="cmd-login" src="https://user-images.githubusercontent.com/33191176/154385258-053b228a-083b-4ad0-97e6-d80bc0d7efcd.png" /></p>

全てのプラットフォームの認証情報を登録する必要はありません (例えば EdStem のみの登録でも課題取得は可能) 。

認証情報は [PC固有の値](https://github.com/denisbrodbeck/machineid) から動的に生成される共通鍵で暗号化されてローカルに保存されます。

認証情報の登録状況や、認証情報の保存ファイルパスは `status` サブコマンドを実行すると確認できます。

```
taskrader status
```

#### Step 2. 課題一覧を取得する
サブコマンド `list` で未提出課題の一覧を縦に並んだボックス形式で閲覧できます。  
各課題の締切日は、締切までの残り時間に応じて強調表示されます。

<p align="center"><img alt="page-list" src="https://user-images.githubusercontent.com/33191176/154385343-3619f071-3570-44a8-95b2-0e0460bc43f5.png" /></p>

### GUI操作方法

#### Step 1. Webブラウザを起動する
`gui` サブコマンド (または `open`) を実行することでWebブラウザが起動し、GUIで操作できます。

本コマンドを実行後、`taskrader` プログラムは API サーバとして動作します。
ブラウザでの操作が終わったら Ctrl-C でサーバを終了してください。

```sh
# Webブラウザを起動 & APIサーバとして動作
$ taskrader gui

# Webブラウザでの操作が終わったら Ctrl-C で終了してください
^C
```

#### Step 2. ログイン用の認証情報を登録する
▼ 何も認証情報が登録されていない場合、下図のような画面が表示されるでしょう。
案内に従って「認証情報の登録画面へ」ボタンを押してください。

<p align="center"><img alt="page-empty" src="https://user-images.githubusercontent.com/33191176/154385318-e1e79c35-fb30-4315-b0c6-755d5f194e02.png" /></p>

なお、認証情報の登録画面は画面左の **歯車アイコン** からも開くことができます。

▼ 登録ボタンを押すとログインが試行され、成功すればグリーンで登録済みの表示になります。

<p align="center"><img alt="page-login2" src="https://user-images.githubusercontent.com/33191176/154385403-39797c26-cac6-4299-8016-265f8cfa9729.png" /></p>

#### Step 3. 未提出課題の一覧を表示する
画面左の **ホームアイコン** を押して課題一覧画面へ移動できます。

▼ 表右上の「課題一覧を再取得する」ボタンを押下することで未提出課題の一覧を再取得できます。

<p align="center"><img alt="page-list-loading" src="https://user-images.githubusercontent.com/33191176/154385370-a5f5e8bc-7b82-49a2-b385-eb33b5902294.png" /></p>

▼ 課題がない場合はこのような表示になります。嬉しいね。

<p align="center"><img alt="page-list-empty" src="https://user-images.githubusercontent.com/33191176/154385415-998e4a6f-91c2-4956-9f22-d9f6c7c9cb48.png" /></p>


▼ 締め切りまでの時間に応じて下記のように強調表示されます。※ テスト用データです

<p align="center"><img alt="page-status-pills" src="https://user-images.githubusercontent.com/33191176/154385427-73f6dec7-e341-4f55-80b4-108ef0b6d0d8.png" /></p>

## 動作要件
Linux x86-64 でのみ動作を確認しました。

Teams の課題取得には Selenium WebDriver を使用します。
Ubuntu をご利用なら下記コマンドを実行すれば WebDriver をインストールできます。

```
sudo apt update && sudo apt install -y chromium-chromedriver
```

## ビルド方法
ビルドには Go コンパイラ (>= v.1.17) と Yarn 必要です。

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


## 技術情報
学情の課題取得は Go 標準ライブラリの `http.Client` を用いて実装されています。
Webブラウザでのログイン時の通信を模倣しました。

EdStem の課題取得にも Go 標準ライブラリの `http.Client` を用いました。
こちらは API のエンドポイントを叩いてログインや課題取得を行っています。

Teams の課題取得は Selenium を用いて実装されています。
当初は [Microsoft Graph API](https://docs.microsoft.com/ja-jp/graph/api/educationassignment-get?view=graph-rest-1.0&tabs=http) を用いる予定でしたが、
大学のMSアカウント管理者の許可が必要だったため、打開策として Selenium を用いました。
