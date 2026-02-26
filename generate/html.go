package generate

import (
	"fmt"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
)

type HTMLElement uint8

const (
	Li HTMLElement = iota
	Tr
	Td
	Br
	A
	P
	H3
	HTMLElementCount
)

// This needs to be changed to the start element and the end part as well
var htmlElementString = [HTMLElementCount]string{
	Li: "<li>%s</li>",
	Tr: "<tr>%s</tr>",
	Td: "<td>%s</td>",
	Br: "<br>",
	A:  `<a href="%s">%s</a>`,
	P:  "<p>%s</p>",
	H3: "<h3>%s</h3>",
}

func (el HTMLElement) String() string {
	if el >= HTMLElementCount {
		panic("invalid HTML element")
	}
	return htmlElementString[el]
}

type HTMLPart uint8

const (
	HomePage HTMLPart = iota
	ArticlesPage
	ContactPage
	Navbar
	ProjectsInfo
	ArticlesInfo
	ArticleTemplate
	HTMLCount
)

var htmlPartStrings = [HTMLCount]string{
	HomePage: `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Nerd with a mouth - Blog</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png"/>
		<link rel="icon" type="image/png" sizes="32x32" href="/assets/favicon-32x32.png"/>
		<link rel="icon" type="image/png" sizes="16x16" href="/assets/favicon-16x16.png"/>
		<link rel="manifest" href="/assets/site.webmanifest"/>
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
  </head>
	<body>
	%s
	</body>
</html>`,
	ArticlesPage: `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Nerd with a mouth - Articles</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png"/>
		<link rel="icon" type="image/png" sizes="32x32" href="/assets/favicon-32x32.png"/>
		<link rel="icon" type="image/png" sizes="16x16" href="/assets/favicon-16x16.png"/>
		<link rel="manifest" href="/assets/site.webmanifest"/>
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
  </head>
	<body>
	%s
	<h1 id="articles">Articles</h1>
	<ul class="articles">
		%s
	</ul>
	</body>
</html>`,
	ContactPage: `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Nerd with a mouth - Contact Information</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png"/>
		<link rel="icon" type="image/png" sizes="32x32" href="/assets/favicon-32x32.png"/>
		<link rel="icon" type="image/png" sizes="16x16" href="/assets/favicon-16x16.png"/>
		<link rel="manifest" href="/assets/site.webmanifest"/>
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
  </head>
	<body>
	%s
	<h2 id="contact-info">Contact Information</h2>
	<table class="nwm-table">
		<tbody>
			<tr>
				<td>Email:</td>
				<td><a href="mailto:alcolka@gmail.com">adrian<em>[at]</em>nerdwithamouth<em>[dot]</em>com</a></td>
			</tr>
			<tr>
				<td>Github:</td>
				<td><a target="_blank" href="https://github.com/ShadowSeick">github.com/ShadowSeick</a></td>
			</tr>
			<tr>
				<td>Linkedin:</td>
				<td><a target="_blank" href="https://linkedin.com/in/adrián-muñoz-gonzález-b98669136/">Adrián Muñoz González</a></td>
			</tr>
		</tbody>
	</table>
	</body>
</html>`,
	Navbar: `<header>
	<nav>
		<h1 id="logo">
			<a href="/%s/">Nerd<span class="white">With</span>A<span class="white">Mouth</span></a>
		</h1>
		<ul class="menu">
			<li><a href="/%s/">Home</a></li>
			<li><a href="/%s/articles/">Articles</a></li>
			<li><a href="/%s/contact.html">Contact</a></li>
		</ul>
	</nav>
</header>`,
	ProjectsInfo: `<h2 id="projects-info">Projects</h2><table class="nwm-table"><tbody>%s</tbody></table>`,
	ArticlesInfo: `<h2 id="articles-info">Articles</h2><ul class="articles">%s</ul>`,
	ArticleTemplate: `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png"/>
		<link rel="icon" type="image/png" sizes="32x32" href="/assets/favicon-32x32.png"/>
		<link rel="icon" type="image/png" sizes="16x16" href="/assets/favicon-16x16.png"/>
		<link rel="manifest" href="/assets/site.webmanifest"/>
		<link rel="stylesheet" href="/assets/styles/tufte.css"/>
		<link rel="stylesheet" href="/assets/styles/custom.css"/>
		<script src="/assets/scripts/highlight/highlight.js"></script>
		<script>hljs.highlightAll();</script>
	</head>
	<body>
		%s
		<article>
		%s
		</article>
	</body>
</html>`,
}

func (html HTMLPart) String(language domain.Language) string {
	if html >= HTMLCount {
		panic("invalid html part")
	}
	switch html {
	case Navbar:
		lang := language.String()
		return fmt.Sprintf(htmlPartStrings[html], lang, lang, lang, lang)
	}
	return htmlPartStrings[html]
}
