# Issue #27: ci: Actions にレビューと linter の実行を追加する

## 要件
- GitHub Actions でレビューと linter の実行を追加する。
- 実行タイミングは Pull Request の作成時と、追加コミット時。

## 実装方針
- `pull_request` イベントの `opened`, `synchronize`, `reopened` で起動する workflow を新規追加する。
- レビューは `reviewdog/action-golangci-lint` を使い、PR にレビューコメントを返す。
- lint は `golangci/golangci-lint-action` を使い、CI として成否を明示する。

## 変更内容
- `.github/workflows/review-and-lint.yml` を追加。
  - `review` ジョブ:
    - `actions/checkout@v4`
    - `actions/setup-go@v5`
    - `reviewdog/action-golangci-lint@v2`
  - `lint` ジョブ:
    - `actions/checkout@v4`
    - `actions/setup-go@v5`
    - `golangci/golangci-lint-action@v6`

## 期待される効果
- PR 作成時と更新時に、自動でレビューコメントと lint 結果が得られる。
- レビュー体験と品質チェックの即時性を向上できる。
