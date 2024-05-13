package util

import (
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func GetAvatarByURL(url string) string {
	avatarURL := ""

	client := &http.Client{}

	// Crear una solicitud HTTP con el User-Agent personalizado
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error al crear la solicitud HTTP:", err)
	}

	// Simular el User-Agent de Chrome en Windows
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36")

	// Realizar la solicitud HTTP
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Error al hacer la solicitud HTTP:", err)
	}
	defer response.Body.Close()

	// Parsear el HTML
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error al parsear el HTML:", err)
	}

	body := doc.Find("body")
	html, err := body.Html()
	if err != nil {
		log.Fatal("Error al obtener el contenido HTML del elemento <body>:", err)
	}

	re := regexp.MustCompile(`"url":"(https://yt3.googleusercontent.com[^"]+)"`)

	matches := re.FindAllStringSubmatch(html, -1)
	i := 0

	for _, match := range matches {
		// La URL se encuentra en el segundo grupo capturado
		avatarURL = match[1]

		if i == 5 {
			break
		}
	}

	return avatarURL
}
