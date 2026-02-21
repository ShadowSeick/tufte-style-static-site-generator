package domain

type Project struct {
	Name string
	Link string
	Description string
}

// Ideally this would be something that I get from github or something. Let's build something easy for now
var PublicProjects = []Project{
	{
		Name: "Tufte style static syte generator",
		Link: "https://github.com/ShadowSeick/tufte-style-static-site-generator",
		Description: "Static site generator in golang. This tool tries to build my blog page from text.<br>Why? Because text will outlive us, I don't like to depend on third party libraries and I want to proof myself I can do a tool from scratch start to finish.<br>Influenced by Tufte CSS and GingerBill.",
	},
}
