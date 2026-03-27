# feature-12: 一括全選択機能

## Issue summary
- `--multi` 起動時に、ノーマルモードで `Ctrl + A` を押すと選択可能なアイテムを一括選択できるようにする。
- ファイルモードではファイル、ディレクトリモードではディレクトリを対象にする。
- ヘルプ表示と `README.md` に `Ctrl + A` の説明を追加する。

## Design
- `Selector` インターフェースに `selectAll()` を追加し、`SelectorCommon` で共通実装する。
- `selectAll()` は `Multi == true` のときだけ動作し、現在表示中（フィルタ適用後）のエントリを全選択する。
- キーバインド処理は `pkg/tui_bubbletea.go` のノーマルモードに `ctrl+a` を追加する。
- ヘルプ行に `Ctrl+A select all items (--multi)` を追加する。

## Test plan
- `pkg/tui_bubbletea_test.go` で `Ctrl+A` の入力をノーマルモード/フィルタモードで検証する。
- `pkg/selector_common_test.go` で `selectAll()` の multi 有効/無効時の挙動を検証する。
