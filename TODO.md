# TODO

Work for later. Each item is also a GitHub issue.

- **Share the front matter parser.** `blog.go`, `reads.go` and
  `project_pages.go` each split the `---` front matter and check the
  fields they need. Move the split into one helper in `internal/site`
  that returns the raw YAML and the body, so the three loaders share it.
- **Tag and deploy 0.3.1 when ready.** The version is set in the
  `Justfile`. Tag it (`git tag 0.3.1`) and run `just deploy`. Not done
  yet on purpose.
