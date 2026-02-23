package generate

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
	A: `<a href="%s">%s</a>`,
	P: "<p>%s</p>",
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

var htmlPartStrings = [HTMLCount]string {
	HomePage: `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Nerd with a mouth - Blog</title>
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
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
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
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
    <link rel="stylesheet" href="/assets/styles/tufte.css"/>
    <link rel="stylesheet" href="/assets/styles/custom.css"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
  </head>
	<body>
	%s
	<h1 id="contact-info">Contact Information</h1>
	<table class="nwm-table">
		<tbody>
			<tr>
				<td>Email:</td>
				<td><a href="mailto:alcolka@gmail.com">adrian<em>[at]</em>nerdwithamouth<em>[dot]</em>com</a></td>
			</tr>
			<tr>
				<td>Github:</td>
				<td><a href="https://github.com/ShadowSeick">github.com/ShadowSeick</a></td>
			</tr>
			<tr>
				<td>Linkedin:</td>
				<td><a href="https://www.linkedin.com/in/adri%C3%A1n-mu%C3%B1oz-gonz%C3%A1lez-b98669136/">Adrián Muñoz González</a></td>
			</tr>
		</tbody>
	</table>
	</body>
</html>`,
	Navbar: `<header>
	<nav>
		<h1 id="logo">
			<a href="/">Nerd<span class="white">with</span>a<span class="white">mouth</span></a>
		</h1>
		<ul class="menu">
			<li><a href="/">Home</a></li>
			<li><a href="/articles/">Articles</a></li>
			<li><a href="/contact.html">Contact Info</a></li>
		</ul>
	</nav>
</header>`,
	ProjectsInfo: `<h2 id="projects-info">Projects</h2><table class="nwm-table"><tbody>%s</tbody></table>`,
	ArticlesInfo:  `<h2 id="articles-info">Articles</h2><ul class="articles">%s</ul>`,
	ArticleTemplate: `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="/assets/styles/tufte.css"/>
		<link rel="stylesheet" href="/assets/styles/custom.css"/>
		<script src="/assets/scripts/highlight/highlight.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
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

func (html HTMLPart) String() string {
	if html >= HTMLCount {
		panic("invalid html part")
	}
	return htmlPartStrings[html]
}
