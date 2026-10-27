# zachpeters.org

The source for [zachpeters.org](https://zachpeters.org): a small Go (Gin)
site styled with [Tufte CSS](https://edwardtufte.github.io/tufte-css/).

This guide tells you how to add blog posts and projects so that they show
up correctly, including on the timelines that open the Blog and Projects
pages.

**Design rule:** follow Tufte's style. Use text and margin notes, not
boxes and color. Use italic, not bold, for emphasis.

## Run the site on your computer

The tools come from the Nix flake. With direnv, they load when you enter
the folder:

```sh
direnv allow     # once
just dev         # runs the site at http://localhost:8080 and restarts it when a file changes
```

Other recipes:

| Command     | What it does                                     |
| ----------- | ------------------------------------------------ |
| `just run`  | Build and run once                               |
| `just lint` | gofmt, go vet and golangci-lint                  |
| `just vuln` | govulncheck                                      |
| `go test ./...` | Unit tests                                   |

## Blog posts

### 1. Make the file

Put one markdown file per post in `blog/`. The filename gives the date,
and the front matter gives the title:

```
blog/2026-10-18T21:30:00-0500_Recommended-Apps.md
     └──── date and time ────┘ └─ any text ─┘
```

```markdown
---
title: Recommended Apps
---

The text of the post starts here.
```

- **Date:** `YYYY-MM-DDTHH:MM:SS` plus the UTC offset (`-0500` for
  US Central daylight time, `-0600` for standard time).
- **Underscore:** one `_` between the date and the rest of the name. The
  rest of the name is not used for the title, so it can hold anything.
- **Title:** the `title:` line in the front matter. It can contain any
  characters, including `-` and `_`.
- **URL:** the title in lowercase, with punctuation removed and words
  joined by `-`, for example `/blog/recommended-apps`.

### 2. Write the first sentence with care

The first sentence (up to the first `.`, `!` or `?` that comes before a
space) is the excerpt under the title on the Blog page. Make it a complete
sentence that tells the reader what the post is about.

### 3. Open with small caps (optional)

Tufte opens a section with a few words in small caps:

```markdown
<span class="newthought">A few of the tools</span> I use every day, and why I like them.
```

### 4. Sidenotes

Use markdown footnotes. The site shows each one as a numbered note in the
margin, next to its reference:

```markdown
It took a while to get used to the "reverse" selection[^helix], but it has grown on me.

[^helix]: Helix is `selection → action` instead of `action → selection` like vim.
```

- The note name (`helix`) is only a label. The site numbers the notes.
- Put the note text anywhere in the file. The end of the file is easiest.
- Keep a note to one paragraph. Links, `code` and *italic* work in it.

For a note without a number, use a margin note. There is no markdown
syntax for it, so write the HTML in the paragraph:

```markdown
I grew the culture for 48 hours<label for="mn-incubate" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-incubate" class="margin-toggle"/><span class="marginnote">At 37 °C, the standard for human pathogens.</span> before the next test.
```

- Each `id` (here `mn-incubate`) must be different on the page. The `for`
  value must be the same as the `id`.
- Keep the label, input and span on one line, with no spaces between them.
- On a narrow screen, the note hides behind a &#8853; that the reader taps.

### 5. Headings

The post title is the page's top heading, so start at `##`:

```markdown
## Case Studies
### Ok
```

Use `##` for sections and `###` for subsections. Do not go deeper.

### 6. Images

Put images in `assets/img/<post-or-project>/` as `.webp`. In the post, use
HTML so the caption goes in the margin:

```html
<figure>
  <label for="mn-demo" class="margin-toggle">&#8853;</label>
  <input type="checkbox" id="mn-demo" class="margin-toggle" />
  <span class="marginnote">Short caption for the margin</span>
  <img src="/assets/img/vhs.gif" alt="What the image shows">
</figure>
```

Each `id` (here `mn-demo`) must be different on the page. For an image as
wide as the page, use `<figure class="fullwidth">` with a
`<figcaption>` instead of the margin note.

### 7. Publish

The site reads posts when it starts, so a new post appears after a restart
or a deploy. It then shows:

- on the Blog page timeline and list,
- in "Recent Writing" on the home page (the 5 newest posts),
- in the RSS feed at `/rss`.

## Projects

Every project is one file, `projects/<id>.md`. Software
and hardware projects are the same kind of thing, so there is no split.
The Projects page (`/projects`) lists them with a summary. Each project has
its own page at `/projects/<id>`. The old `/software` and `/hardware`
addresses redirect to the Projects page.

1. **Make the file** `projects/<id>.md`. Start it with front matter, then
   write the page in markdown:

   ```markdown
   ---
   name: "Dropsearch"
   years: "2024-Current"
   summary: "One or two sentences the Projects list shows under the title."
   ---

   What it is and why you built it.
   ```

   - **`name` and `years`:** required. The page heading is the name and the
     years in brackets, for example "Dropsearch (2024-Current)". The table
     below says what `years` can be. A file with other years still has a
     page, but it is not on the timeline.
   - **`summary`:** optional. Without it, the project is listed with no
     summary.
   - **The file name** is the `<id>`. It must be different from all other
     project files. The timeline, the list and the page address all use it.

   | years          | Meaning                          | On the timeline                         |
   | -------------- | -------------------------------- | --------------------------------------- |
   | `2022`         | One year                         | A dot in that year                      |
   | `2017-2019`    | Finished, over several years     | A line with a dot at each end           |
   | `2023-Current` | Still active, started in 2023    | A dot, then a line to "now"             |
   | `Ongoing`      | Still active, start not known    | A dashed start, then a line to "now"    |

   `Present` and `Now` work in place of `Current`.

2. **Sidenotes** work as in blog posts: `text[^1]` in the text, and
   `[^1]: The note.` at the end of the file.

3. **Figures, tables and iframes** are raw HTML. Put a blank line before and
   after each block, and do not put a blank line inside it. A blank line
   ends the HTML block early, and the rest of it is shown as broken
   markdown. Add images in `assets/img/<id>/` as `.webp`, with margin
   captions (see "Images" above).

4. **Nothing else.** The Projects pages read the files each time they are
   opened, so a new file shows up after a save and refresh, with no
   restart. Projects appear newest first on the list, and by start year on
   the timeline. If a file has an error, the log names it, and the Projects
   pages show no projects until it is fixed.

The list page is `pages/projects.html`, and the project page
wrapper is `templates/include/project-page.html`. The site reads every `*.md`
file in `projects/` as a project, so keep other files out of that folder.

## Reads

The Reads page (`/reads`) lists articles by other people that you
recommend. Put one markdown file per article in `reads/`:

```markdown
---
title: Choose Boring Technology
author: Dan McKinley
year: 2015
url: https://mcfunley.com/choose-boring-technology
archive: http://archive.today/2024.01.04-012806/https://mcfunley.com/choose-boring-technology
added: 2026-10-19
---

Why you like it, in a few sentences. Sidenotes[^1] work as in blog posts.

[^1]: Like this one.
```

- **Filename:** short, lowercase, words joined by `-`, for example
  `choose-boring-technology.md`. It is the anchor on the page
  (`/reads#choose-boring-technology`).
- **`title`, `url`, `added`:** necessary. `added` is the date you added it
  (`YYYY-MM-DD`). The page shows the newest first.
- **`author`, `year`, `archive`:** optional. They show in the margin.
  `year` is when the article was written. It also puts the read on the
  chart at the top of the page; a read without a `year` is not on it.

The site reads these files when it starts, so a new read appears after a
restart or a deploy. If a file has an error, the site does not start, and
the log names the file.

## Other parts of the page

- **Sparkline** (top right): your GitHub contributions per week for the
  past year. The server fetches them from GitHub when it starts and then
  every 6 hours. If GitHub cannot be reached, the sparkline is not shown.
- **Margin dates:** the Blog list and the home page show each post's date
  in the margin. You do not need to write them.

## Deploy

This Mac is arm64 and the cluster is x86_64, so the image builds on the
remote x86 podman through the `cluster` connection.

1. Set a new `version` at the top of the `Justfile`. Do not reuse a tag.
2. Run `just deploy`. It builds and pushes
   `registry.home.zach/zachpeters.org:<version>`.
3. Set the same tag in
   `~/Projects/infrastructure/kubernetes/zachpeters.org.yaml`, then:

   ```sh
   kubectl apply -f ~/Projects/infrastructure/kubernetes/zachpeters.org.yaml
   kubectl rollout status deploy/zachpeters
   ```

   Keel does not move to a new tag by itself (`keel.sh/match-tag`), so
   this step is necessary.
