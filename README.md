# Tufte Style Static Site Generator
This is a tool to generate from markdown files static sites with [Tufte CSS](https://edwardtufte.github.io/tufte-css/). The idea is for blog pages, but it could be used for much more. What more? I don't know

## Translations
I will be following the [Markdown syntax](https://www.markdownguide.org/) style. IDs for the elements have been done following the extension style in the style page.
I have added 3 other representations that are not present in the syntax:
- Side note -> \[^side-note\]
- Margin note -> \[^margin-note\]
- Sub header -> \[^sub-header\]

## How it works
The program walks through the blog folder. In it, there are different folders containing the name of the article in a slug way and inside them there are the article in spanish (TBD) and english with the images and gifs they use. Once all articles are generated in html, we upload them into bunny CDN storage with the files they need. When we upload them they will follow the next structure:
- Images and gifs will be uploaded to */images/{article folder name}/file-name*
- HTML files will be uploaded to: */{article folder name}/{file language}*

## TODO
- [IN PROGRESS] Make the proper translations for the pages
- Create github actions to execute this when merged into main
- Finish HTTP article and publish it
- Make project public

## Improvements
- Update how I handle the creation of templates and I build pages and articles
- Naive Pratt parsing like. It is searching for strings and tokens. I tried first this implementation, with no idea I was doing it, but it seemed complicated so I opted for something that worked fast. After seeing others solutions, I realize it was not that difficult and gone with this. Of course I am not following Pratt parsing as a whole, but trying to do a tokenizer like.
- Logger instead of using raw fmt.Println

## Implemention
- Regexp. It was the simplest and the most naive approach. It is useful for my needs, but after it all works, I will research for a proper parsing implementations.
- Creation of templates is using fmt.Sprintf in it's majority with some sprinkle of strings.Builder for the difficult stuff (mainly Articles). This is naive but it works.

## Webpage
As it is a static website, I went with a CDN, exactly [Bunny CDN](bunny.net) for all the blog, images and scripts. Why? It's the cheapest and it does not play with you after a while like cloudflare and apparently it's really well regarded in the community. Besides, it was a really pleasant experience to use and port the domain and use the storage. It has all the tools I need and again it's freakingly cheap.

## Future integrations
As I like to write, I might integrate twitter and substack in the future for visibility. As I already have the parser, doing this integration should not take me that much time (as a seasoned developer, I can tell this will be more complicated and take me more time than expected).

## Goal
FINISH A PROJECT, learn Go language better and mix writting with programming, my two loves.

