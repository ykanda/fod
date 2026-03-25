# feature-30: ステータスライン固定表示

## 背景
現状の `buildView` は表示するエントリ行数に応じて描画全体の高さが変わるため、エントリが少ないとステータスラインが画面最下部に固定されません。

## 対応方針
- 画面高さを `top(1) + list + status(1) + help` に分割して計算する。
- `list` 領域は必要に応じて空行で埋める。
- `status/help` は常に下部領域に配置する。
- `showHelp=true` のときの行数計算と実描画行数を一致させる。

## 実装内容
- `pkg/tui_bubbletea.go`
  - `helpLines` を先に組み立て、表示可能な下部領域から `linePerPage` を算出。
  - リスト描画後に空行を追加して領域を埋める。
  - `status` と `help` を固定領域に描画。
- `pkg/tui_bubbletea_test.go`
  - `buildView` の総行数が `height` と一致するテストを追加。
  - help の有無で下部領域が固定されることをテスト。

## 検証
- `go test ./...` : pass
