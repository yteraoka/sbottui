# sbottui

SwitchBot デバイス・シーンをターミナルから操作できる TUI アプリケーションです。

## 機能

- SwitchBot デバイス一覧の表示と操作
- シーンの一覧表示と実行
- デバイス種別ごとの操作オーバーレイ
- API レスポンスのローカルキャッシュ（`~/.cache/sbottui/`）

## 対応デバイス

### 物理デバイス

| デバイス | 操作 |
|---------|------|
| カラー電球（Color Bulb, Strip Light） | 電源 on/off、明るさ（1〜100%）、色温度（2700〜6500K） |
| プラグミニ（Plug Mini US/JP） | 電源 on/off |

### IR デバイス

| デバイス | 操作 |
|---------|------|
| エアコン | 電源、温度（16〜30℃）、モード（auto/cool/dry/fan/heat）、風量（auto/low/medium/high） |
| テレビ | 電源、音量 up/down、チャンネル up/down、ミュート |
| ライト | 電源 on/off |
| 扇風機 | 電源 on/off |

### シーン

選択すると即時実行します（オーバーレイなし）。

## インストール

### ビルド済みバイナリ

[Releases](https://github.com/yteraoka/sbottui/releases) からお使いのプラットフォーム向けのバイナリをダウンロードしてください。

### ソースからビルド

Go 1.23 以上が必要です。

```bash
git clone https://github.com/yteraoka/sbottui.git
cd sbottui
go build -o sbottui .
```

## 使い方

SwitchBot アプリから API トークンとシークレットを取得し、環境変数にセットして実行します。

```bash
export SWITCHBOT_TOKEN=<your-token>
export SWITCHBOT_CLIENT_SECRET=<your-client-secret>
./sbottui
```

### バージョン確認

```bash
./sbottui --version
# または
./sbottui -v
```

出力例:
```
version: v0.1.0
commit:  abc1234
date:    2026-02-25T12:34:56Z
```

## キー操作

### 一覧画面

| キー | 動作 |
|------|------|
| `↑` / `k` | 上に移動 |
| `↓` / `j` | 下に移動 |
| `←` | 名前順でソート |
| `→` | 種別順でソート |
| `Enter` | デバイス操作 / シーン実行 |
| `r` | キャッシュをクリアして再取得 |
| `q` / `Ctrl+C` | 終了 |

### 操作オーバーレイ

| キー | 動作 |
|------|------|
| `↑` / `k` | 前の項目へ |
| `↓` / `j` | 次の項目へ |
| `←` | 値を減少 / 前の選択肢へ |
| `→` | 値を増加 / 次の選択肢へ |
| `Enter` | コマンド送信 |
| `Esc` | 一覧に戻る |

## キャッシュ

デバイス・シーン一覧は `~/.cache/sbottui/` にキャッシュされます。TTL はなく、`r` キーで手動クリア＆再取得できます。

## ライセンス

[MIT](LICENSE)
