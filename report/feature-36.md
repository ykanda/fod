# Issue #36 対応レポート

## 概要

決定動作を `Shift + Enter` に統一し、端末差異の影響を受けずに「選択して終了」を実行できるようにする。

## 実装方針

- TUI のキー入力処理で決定キーを `shift+enter` のみに統一する。
- 対象は通常モードとフィルタモードの両方とする。
- ヘルプ表示のキーバインド案内を `Shift+Enter` に更新する。
- README の説明文とキー一覧を更新する。
- `shift+enter` で終了できること、および従来キー (`ctrl+o`) で決定しないことをテストで担保する。

## 影響範囲

- `pkg/tui_bubbletea.go`
- `pkg/tui_bubbletea_test.go`
- `README.md`
