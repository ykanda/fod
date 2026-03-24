# bug-33 調査レポート

## 対象
- Issue: https://github.com/ykanda/fod/issues/33
- Title: fix: Esc キーで終了時のステータスコードを 0 にする

## 症状
- `Esc` でダイアログを終了すると、プロセスの終了コードが `1` になる。

## 原因
- `pkg/tui_bubbletea.go` で `Esc` は `selector.cancel()` を呼び、`ResultCancel` として終了する。
- しかし `cmd/fod/main.go` の `action` は `ResultOK` 以外をすべて `error` として返している。
- そのため `ResultCancel` も異常終了扱いとなり、`run()` が `ExitCodeError(=1)` を返していた。

## 修正方針
- `ResultCancel` は「ユーザーが出力せずに終了する正常系」として扱い、`action` では `nil` を返す。
- `ResultNone` など想定外コードのみエラー扱いを維持する。

## テスト方針
- `ResultCode` の扱いを切り出したヘルパー関数を追加し、
  `ResultOK` / `ResultCancel` / 想定外コードの3ケースをユニットテストで検証する。
