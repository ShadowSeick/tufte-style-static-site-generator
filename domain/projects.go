package domain

type Project struct {
	Name        string
	Link        string
	Description string
}

// PublicProjects Ideally this would be something that I get from github or something. Let's build something easy for now
var PublicProjects = map[Language][]Project{
	English: {
		{
			Name:        "Tufte style static syte generator",
			Link:        "https://github.com/ShadowSeick/tufte-style-static-site-generator",
			Description: "Static site generator in golang. This tool builds my blog page from text.<br>Why? Because text do not depend on any tool, is easy to read and to use no matter of the platform nor the program used to open it.<br>I don't like to depend on third party libraries and I want to proof myself I can do a tool from scratch.<br>Influenced by Tufte CSS and GingerBill.",
		},
	},
	Spanish: {
		{
			Name:        "Tufte style static syte generator",
			Link:        "https://github.com/ShadowSeick/tufte-style-static-site-generator",
			Description: "Generador de sitios estáticos en Go. Esta herramienta construye mi blog a partir de texto.<br>¿Por qué? Porque el texto no depende de ninguna herramienta, es facil de leer y facil de utilizar sin importar la plataforma donde corra ni el programa que se utilice para abrirlo.<br>No me gusta depender de librerías de terceros y quiero demostrarme a mí mismo que puedo crear una herramienta desde cero.<br>Influenciado por Tufte CSS y GingerBill.",
		},
	},
}
