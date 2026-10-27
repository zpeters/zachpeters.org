---
name: "Scratch Notes"
years: "2023-Current"
date: "2026-10-25T19:24:58-0500"
summary: "A minimal note-taking app written in Go and SQLite, with Litestream backups to remote storage."
---

[Scratch Notes](https://scratch.zachpeters.org) is a minimal note-taking application. This was written in Go using [Fiber](https://docs.gofiber.io/) and [SQLite](https://sqlite.org/index.html). I used this project as an excuse to learn more about litestream so there is an interesting bit about how this works. When the application boots up[^1] it will check for an existing backup on remote storage, if that exists and there isn't a local database it will restore it. Then it will continue with regular replications. For some reason with this I really struggled to get litestream working with a local [_minio_](https://min.io) instance. So I have documented my config for future reference[^2]

[^1]: Example [_Dockerfile_](https://gist.github.com/zpeters/08d325dc3c435f31ab8f0b5ea63260e8) and [_start.sh_](https://gist.github.com/zpeters/ecf742fe4a4538c434fd866c49fcaa8e) script
[^2]: [_litestream.yml_](https://gist.github.com/zpeters/a2b967a2647297161e799e1bb7c005d7)
