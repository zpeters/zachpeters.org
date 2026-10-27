---
title: My 2025 "Stack"
date: "2025-01-01T00:00:00-0500"
---

<span class="newthought">A few of the tools</span> I use every day, and why I like them.

## Bruno

[Bruno is an open source IDE for exploring and testing APIs](https://usebruno.com/). This could be compared to software like <em>Postman</em> and <em>Insomnia</em>. Over the past years my feeling is that both of the prior have become overly slow and bloated. While you can <em>export</em> data, I never felt like I had full control of the configurations and data I was creating.

With Bruno you get a pretty minimalist (yet powerful) interface. All basic http operations are supported. There is a full scripting and a basic testing interface. This is the perfect mix of what I'm looking for. The added advantage over similar software is that the configurations are <em>simple plain text</em>. Again, the [power of plain text](https://wiki.c2.com/?PowerOfPlainText).

## Helix

[Helix is a <em>post-modern</em> text editor](https://helix-editor.com/) that is written with speed and comfort in mind. For me I see this as the best possible version of Vim. The learning curve is more gentle than "normal" Vim thanks to the built-in menu system. It took a while for me to get used to the "reverse" selection[^helix], but it has been slowly growing on me. At this point `vim` is now an alias to `hx` for me!

## Logseq

It's hard for me to classify [Logseq](https://logseq.com/). In many ways it's a very simple application - you edit markdown files. With that simplicity there is a lot of power. You can very easily transform a file into a calendar entry, a checklist, a project plan, etc. The text based format for me is the clincher. Notion does have a much flashier interface, but when I used it I never felt like I <em>owned</em> my data. With Logseq, everything is a text file that I am free to copy elsewhere or take to another program if I want to switch.

What I really like about Logseq over other editors is the <em>outline first</em>[^outline] paradigm. This takes a while to wrap your head around, but once you adapt to that notion, it is very easy to make simple and efficient outlines and then expand on notes, links to specific blocks, etc. [Logseq - A Powerful Tool for Thought](https://tfthacker.medium.com/logseq-a-powerful-tool-for-thought-9058dec80dbe) has a really good summary of the strengths and power of logseq.

## VHS

As you can see from the demo below, [VHS](https://github.com/charmbracelet/vhs?tab=readme-ov-file#vhs-command-reference) is a <em>screen recorder for the terminal</em>. For terminal based apps a lot of the time, a text description is great, but a screencast is really great for showing off how commands can work together and for showing off tui apps. There is another app called [asciinema](https://asciinema.org/) but I enjoy this a little more.

<figure>
  <img src="/assets/img/vhs.gif" alt="animated gif of working in helix to edit a webpage">
</figure>

There are some more advanced commands, but this is what I normally use: `vhs record > my_recording.tape` start hacking and exit when you are done, then `vhs my_recording.tape` renders the actions to a gif[^gifsicle].

The other thing I find really cool is that the [<em>tape</em>](/assets/vhs-demo.tape) file is just a [text file with some markup and commands](https://github.com/charmbracelet/vhs?tab=readme-ov-file#set-theme). That means that you can manually edit it and re-render it at any point, there are even [themes](https://github.com/charmbracelet/vhs?tab=readme-ov-file#set-theme)!

[^helix]: Helix is `selection → action` instead of `action → selection` like vim. See [Migrating from Vim](https://docs.helix-editor.com/from-vim.html).
[^outline]: Hierarchical organization and structure are prioritized over "raw" markdown. Instead of a free-form Word document, imagine a permanent bullet-point list for everything on the page - each being able to be linked to, referenced, and embedded from anywhere else
[^gifsicle]: Bonus tip to reduce the size of the gif: `gifsicle --resize-width 300 --colors 256 out.gif > resized.gif`
