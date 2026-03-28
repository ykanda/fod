# bug #45 調査レポート

## 概要
Issue #45 では、モード表示とヘルプガイドをステータスラインに統合し、`[F]/[D]` 表示を `[File]/[Dir]` に変更し、`--mode` ヘルプ説明を正確化することが求められている。

## 原因
- `pkg/tui_bubbletea.go` の `buildView` は、`showHelp=false` 時に `"[<mode>] ? help"` をステータス行とは別の最下行に描画していた。
- モード表示文字列は `pkg/selector_file.go` と `pkg/selector_directory.go` の `getMode()` でそれぞれ `"F"` / `"D"` が返されていた。
- `cmd/fod/main.go` の `--mode` フラグに `Usage` がなく、アプリケーションヘルプ上の説明が不足していた。

## 対応方針
- ステータス行の右側に `"[<mode>] ? help"` を右寄せで合成する。
- モード表示を `File` / `Dir` に変更する。
- `--mode` の Usage を `d|dir|directory` と `f|file` の説明が明確に分かる文言へ更新する。
- 表示高さと既存のヘルプ展開挙動を壊さないよう、テストを更新して回帰を防ぐ。
