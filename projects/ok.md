---
name: "Ok - The Version Control Helper"
years: "2006-Current"
summary: "A small tool that finds changed git repos in my home folder, and the project I rewrite every few months to learn new languages."
---

For the longest time I have been using a small utility that i created sometime in the early 2000s. Truthfully, there isn't much to it. This application will look at all _git_ projects in my home folder, if anything is changed it will list it. This alone solves a huge problem for me. I am forever jumping from project to project, so it's easy for me to forget to check something in. The other part of this is that with a quick command i can _add_, _commit_ and _push_ any matching repos. This is helpful for projects where I use git for more of a "backup" then a true version control system. For example, i use this for dotfiles[^1].

The _meta_ project here is that I rewrite this app every few months. If I'm learning a new programming language I usually try it out here first. Off the top of my head this has been writting in: Go, Rust, Haskell, and Zig. Since I use this on Mac, Linux and Windows it's a good test for cross-platform compatibility as well.

Around 2025 I started using [jj](https://github.com/jj-vcs/jj) full-time. Support has been added to my `ok` tool

[^1]: I use [gnu stow](https://www.gnu.org/software/stow/) to manage my files.
