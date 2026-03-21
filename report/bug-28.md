# Issue 28 調査レポート

## 概要
フィルタに日本語（全角）を入力した際、ハイライトが最後の文字で半角分欠けて見える。

## 原因の推定
`pkg/tui_bubbletea.go` の描画処理で、
- `available := width - len([]rune(prefix))`
- `truncateRunes` / `truncateLine` が **rune数** での切り詰め

を行っている。

全角文字は表示幅が2セルなのに rune数は1のため、表示幅の計算がズレる。結果として
- 文字列が想定より長くなり
- `lipgloss.Style.Width(width)` の内部トリムで末尾の全角文字が「半分」切れたように見える

特にハイライト部分は ANSI スタイル付きで切り詰められるため、最後の一致文字だけハイライト幅が不足する症状になる。

## 修正方針
表示幅での計算に切り替える。
- `available` の計算を `runewidth.StringWidth(prefix)` ベースにする
- `truncateLine` / `truncateRunes` を `runewidth.Truncate` で表示幅基準に変更する

これにより、全角文字を含む行でも表示幅が正しく収まり、ハイライト欠けが解消される見込み。

## 影響範囲
- エントリ行の描画
- ステータス/ヘルプ行のトリム

## 追加テスト案
`truncateLine` / `truncateRunes` の全角文字を含むケースのユニットテストを追加（表示幅が指定幅を超えないことを確認）。
