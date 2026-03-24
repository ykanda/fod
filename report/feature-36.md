# Issue #36 対応レポート

## 概要

`Ctrl + O` でのみ実行できていた「選択して終了」の操作を、`Ctrl + Enter` でも実行できるようにする。

## 実装方針

- TUI のキー入力処理で `ctrl+enter` を `ctrl+o` と同じ分岐に追加する。
- 対象は通常モードとフィルタモードの両方とする。
- ヘルプ表示のキーバインド案内を `Ctrl+O, Ctrl+Enter` に更新する。
- README の説明文とキー一覧を更新する。
- `ctrl+enter` で終了できることとヘルプ文言反映をテストで担保する。

## 影響範囲

- `pkg/tui_bubbletea.go`
- `pkg/tui_bubbletea_test.go`
- `README.md`
