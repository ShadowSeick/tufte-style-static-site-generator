const path = window.location.pathname; // e.g., /en/index.html or /es/index.html
const parts = path.split("/").filter((p) => p);

const currentLang = parts[0] || "en";
const otherLang = currentLang === "en" ? "es" : "en";

// Get the <a> elements
const enLink = document.getElementById("lang-en");
const esLink = document.getElementById("lang-es");

const switchLangURL = `/${otherLang}/${parts.slice(1).join("/") || "index.html"}`;
switch (currentLang) {
  case "en":
    enLink.classList.add("active");
    esLink.href = switchLangURL;
    break;
  case "es":
    esLink.classList.add("active");
    enLink.href = switchLangURL;
    break;
  default:
    console.error("not handled language");
}

if (
  parts.length === 3 &&
  parts[1] === "articles" &&
  parts[2] !== "index.html"
) {
  document.getElementsByClassName("languages")[0].classList.add("hidden");
}
