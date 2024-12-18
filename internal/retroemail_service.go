package internal

import (
	"PortalCRG/internal/repository/entity"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

const HTML_RECOMENDACIONES_YOUTUBE string = `
    <div class="video-container">
        <a href="https://retromasters.up.railway.app/#/notice?id={{TIPS_ID}}" target="_blank">
            <strong>{{URL_TITLE}}</strong><br>
			{{DESCRIPCION}}
			<br>
            <img src="https://img.youtube.com/vi/{{ID_YOUTUBE}}/0.jpg"  class="video-thumbnail">
        </a>
    </div>
`

const HTML_RECOMENDACIONES_RETRO string = `
    <div class="video-container">
        <a href="{{URL_TIPS}}" target="_blank">
            <strong>{{URL_TITLE}}</strong><br>{{DESCRIPCION}}
        </a>
    </div>
`

type RetroEmailService struct {
	password string
	from     string
}

func (r *RetroEmailService) Init() {
	r.password = os.Getenv("PASSWORD_EMAIL")
}

func (r *RetroEmailService) MakeRecomendacion(recomendaciones []entity.PostNew) string {
	recoHTML := ""
	for _, recomendacion := range recomendaciones {
		format := ""
		if recomendacion.Type == "youtube" {
			format = HTML_RECOMENDACIONES_YOUTUBE
			format = strings.ReplaceAll(format, "{{URL_TIPS}}", recomendacion.URL)
			format = strings.ReplaceAll(format, "{{TIPS_ID}}", recomendacion.ID)
			format = strings.ReplaceAll(format, "{{URL_TITLE}}", recomendacion.Title)
			format = strings.ReplaceAll(format, "{{DESCRIPCION}}", recomendacion.Content)
			parts := strings.Split(recomendacion.URL, "/")
			idYoutube := parts[len(parts)-1]
			format = strings.ReplaceAll(format, "{{ID_YOUTUBE}}", idYoutube)
			recoHTML = recoHTML + format
		} else {
			format = HTML_RECOMENDACIONES_RETRO
			format = strings.ReplaceAll(format, "{{URL_TIPS}}", recomendacion.URL)
			format = strings.ReplaceAll(format, "{{DESCRIPCION}}", recomendacion.Content)
			format = strings.ReplaceAll(format, "{{URL_TITLE}}", recomendacion.Title)
			recoHTML = recoHTML + format
		}
	}
	return recoHTML
}

func (r *RetroEmailService) EnviarNotificacionComentarios(alias, email string, comentario entity.CommentRetro, tips *entity.PostNew, recomendaciones []entity.PostNew) {
	// Configuración del remitente
	fromName := "Retro Master" // Nombre del remitente
	from := "retromasterstreamers@gmail.com"
	password := os.Getenv("API_GMAIL")

	// Destinatario
	to := []string{email}

	// Configuración del servidor SMTP
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Leer el archivo HTML desde la ruta "./correo.html"
	body, err := ioutil.ReadFile("./correo.html")
	if err != nil {
		log.Fatalf("Error al leer el archivo HTML: %v", err)
	}

	bodyStr := string(body)

	bodyStr = strings.ReplaceAll(bodyStr, "{{AMIGO}}", alias)
	bodyStr = strings.ReplaceAll(bodyStr, "{{TIPS_TITULO}}", tips.Title)
	//bodyStr = strings.ReplaceAll(bodyStr, "{{TIPS_TEXTO}}", tips.Content)
	bodyStr = strings.ReplaceAll(bodyStr, "{{COMENTARIO_TEXTO}}", comentario.Comment)
	bodyStr = strings.ReplaceAll(bodyStr, "{{COMENTARIO_AUTHOR}}", comentario.Author)
	bodyStr = strings.ReplaceAll(bodyStr, "{{TIPS_ID}}", tips.ID)
	bodyStr = strings.ReplaceAll(bodyStr, "{{HTML_RECOMENDACIONES}}", r.MakeRecomendacion(recomendaciones))

	// Construcción del mensaje completo con encabezados personalizados
	subject := "Novedades en RetroMaster, " + tips.Title
	headers := map[string]string{
		"MIME-Version":     "1.0",
		"Content-Type":     "text/html; charset=\"UTF-8\"",
		"Subject":          subject,
		"From":             fmt.Sprintf("\"%s\" <%s>", fromName, from),
		"To":               strings.Join(to, ", "),
		"List-Unsubscribe": "<mailto:retromasterstreamers+nomas@gmail.com>",
		"X-Priority":       "1",        // Alta prioridad
		"Importance":       "High",     // Importancia alta
		"Sensitivity":      "Personal", // Sensibilidad personal
		"ARC-Seal":         "i=1; a=rsa-sha256; t=1734526822; cv=none; d=google.com; s=arc-20240605;",
		"ARC-Authentication-Results": `i=1; mx.google.com;
	dkim=pass header.i=@gmail.com header.s=20230601 header.b=JPcYLo+X;
	spf=pass smtp.mailfrom=retromasterstreamers@gmail.com;
	dmarc=pass header.from=gmail.com;`,
		"Received-SPF": "pass (google.com: domain of retromasterstreamers@gmail.com designates 209.85.220.41 as permitted sender) client-ip=209.85.220.41;",
		"DKIM-Signature": `v=1; a=rsa-sha256; c=relaxed/relaxed;
	d=gmail.com; s=20230601; t=1734526821; bh=fOj0A/B+hIsGKNJxcxdC7uRhI+nH0mlpcwwvnciGYRQ=;
	b=JPcYLo+X80q/2L5NxlCEUuJEHIx3exB+owgUoE/0raK3K2euaLZbP+11n08yP5cZpR;`,
		"Return-Path": "<retromasterstreamers@gmail.com>",
	}
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + bodyStr // Agregar el cuerpo del mensaje

	// Autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Enviar el email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		log.Fatalf("Error al enviar el email: %v", err)
	}

	log.Println("¡Email enviado con éxito! %v", email)
	r.send()
}

func (r *RetroEmailService) send() {

}

func NewRetroEmailService() *RetroEmailService {
	return &RetroEmailService{}
}
