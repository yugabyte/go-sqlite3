# yugabyte/go-sqlite3

A fork of [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) that can also compile
[sqlite-vec](https://github.com/yugabyte/sqlite-vec) into SQLite. Everything else is upstream's.

## What it adds

Build with the `sqlite_vec` tag, and every connection has sqlite-vec's `vec0` virtual table and
`vec_*` functions:

```sh
go build -tags "sqlite_fts5 sqlite_vec" ./...
```

| File | What |
|---|---|
| `sqlite3_opt_vec.go` | Registers sqlite-vec with `sqlite3_auto_extension` at init. |
| `sqlite3_opt_vec.c` | Compiles `sqlite-vec/sqlite-vec.c` with `SQLITE_CORE`, against this package's SQLite. |
| `sqlite3_opt_vec_test.go` | `vec_version()` matches `sqlite-vec/VERSION`, plus KNN within a partition and a DELETE. |
| `sqlite-vec/` | sqlite-vec's source and licenses as `upgrade/sqlite-vec.sh` writes them, `VERSION` (the version and the commit it came from), and `sqlite3.h`, a shim that points sqlite-vec at `sqlite3-binding.h`. |
| `upgrade/sqlite-vec.sh` | Syncs `sqlite-vec/` from `yugabyte/sqlite-vec` at a tag or commit. |

Without the tag, nothing changes. The tag doesn't work with `libsqlite3`, because sqlite-vec has to
compile against the SQLite this package compiles.

Go modules leave out git submodules, so sqlite-vec's files are committed here and not linked in as
a submodule.

## Branches and tags

- `master` mirrors upstream's `master`.
- `yb` is upstream's release tag that consumers use, plus the commits above. Tags on it are
  `v<upstream version>-yb.<n>`, e.g. `v1.14.52-yb.1`. A tag is never moved: cut `-yb.<n+1>` instead.

Consumers keep importing `github.com/mattn/go-sqlite3` and pin the fork in `go.mod`:

```
replace github.com/mattn/go-sqlite3 => github.com/yugabyte/go-sqlite3 v1.14.52-yb.1
```

## Upgrading

- **sqlite-vec:** push the release's tag to `yugabyte/sqlite-vec` if it isn't there yet. Run
  `upgrade/sqlite-vec.sh <tag>` and `go test -tags "sqlite_fts5 sqlite_vec" ./...`. Commit to
  `yb`, and tag.
- **Upstream (SQLite):** rebase `yb`'s own commits onto upstream's new release tag. They only add
  files, so they don't conflict. Run the tests with and without the tags, push `yb` with
  `--force-with-lease`, and tag `v<new upstream version>-yb.1`. Older tags stay reachable because
  tags hold their commits.
