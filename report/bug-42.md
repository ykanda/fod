# issue #42 調査レポート

## 問題概要
- `Ctrl + H` でドットファイル表示の切り替えは実装済みだが、ヘルプ表示にキー説明が存在しない。
- 利用者が機能を発見しづらく、操作性にギャップがある。

## 原因
- `pkg/tui_bubbletea.go` の `buildView(..., showHelp=true)` で組み立てるヘルプ行に `Ctrl+H` の記載が抜けていた。

## 対応方針
- ヘルプ行に `Ctrl+H toggle hidden file filter` を追加する。
- ヘルプ行が 1 行増えるため、表示高さを前提としたテストの期待値を更新する。

## 変更内容
- `pkg/tui_bubbletea.go`
  - `showHelp` 時のヘルプ一覧に `Ctrl+H` 行を追加。
- `pkg/tui_bubbletea_test.go`
  - ヘルプ表示時の固定行位置に関する期待値を `len(lines)-6` から `len(lines)-7` に更新。
