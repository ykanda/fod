# bug-29 調査レポート

## 対象 Issue
- #29 bug: 隠しファイル（dotfile）フィルターが動作していない

## 原因
- `pkg/filter_dotfile.go` で dotfile 判定を `strings.HasPrefix(entry.Path, ".")` で行っていました。
- `Entry.Path` は `pkg/entry.go` で絶対パスとして保持されるため、`/tmp/.hidden` のような値になり、先頭文字は `/` です。
- その結果、`Ctrl+H` でフィルターを有効にしても絶対パスの隠しファイルが除外されませんでした。

## 対応方針
- 判定対象をパス全体ではなくファイル名に変更します。
- `filepath.Base(entry.Path)` を取得し、`strings.HasPrefix(name, ".")` で判定します。
- 絶対パスの隠しファイルを使ったテストを追加して再発防止します。
