# Issue #36 対応レポート

## 概要

決定動作は `Shift + Enter` を基本としつつ、端末差異で `Shift + Enter` が判別できない環境向けに `Ctrl + O` をフォールバックとして併用する。

## 実装方針

- TUI のキー入力処理で決定キーを `shift+enter` と `ctrl+o` の両方に対応する。
- 対象は通常モードとフィルタモードの両方とする。
- ヘルプ表示のキーバインド案内を `Shift+Enter, Ctrl+O` に更新する。
- README の説明文とキー一覧を更新する。
- `shift+enter` / `ctrl+o` の両方で終了できることをテストで担保する。

## 影響範囲

- `pkg/tui_bubbletea.go`
- `pkg/tui_bubbletea_test.go`
- `README.md`
