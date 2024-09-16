// internal/usecase_impl.go
package internal

import (
	"PortalCRG/internal/repository"
	"PortalCRG/internal/repository/entity"
	"PortalCRG/internal/util"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PortalRetroGamerImpl es una implementación de UserService.
type PortalRetroGamerImpl struct {
	UserRepository   repository.UserRepositoryMongo
	PortalRepository repository.PortalRepositoryMongo
}

// NewUserService crea una nueva instancia de UserServiceImpl.
func NewUserService(userRepository repository.UserRepositoryMongo, portalRepository repository.PortalRepositoryMongo) *PortalRetroGamerImpl {
	return &PortalRetroGamerImpl{
		UserRepository:   userRepository,
		PortalRepository: portalRepository,
	}
}

// Greet retorna un saludo simple.
func (s *PortalRetroGamerImpl) Greet() string {
	usuarios, _ := s.GetAllUsers()

	for _, user := range usuarios {
		for _, rrss := range user.RRSS {
			if rrss.Type == "youtube" && strings.Contains(rrss.URL, "youtube.com/") {
				user.AvatarYT = util.GetAvatarByURL(rrss.URL)
				s.SaveUser(*user)
				break
			}
		}

	}
	return "Hello, world!"
} // Greet retorna un saludo simple.
func (s *PortalRetroGamerImpl) UpdateUserAvatar() string {
	usuarios, _ := s.GetAllUsers()

	for _, user := range usuarios {
		for _, rrss := range user.RRSS {
			if rrss.Type == "youtube" && strings.Contains(rrss.URL, "youtube.com/") {
				user.AvatarYT = util.GetAvatarByURL(rrss.URL)
				s.SaveUser(*user)
				log.Println(user.Alias + "..[OK]")
				break
			}
		}

	}
	return "Ready"
}

func (s *PortalRetroGamerImpl) GetNewTipsFromSearch(videoID string) (*entity.PostNew, error) {

	videoURL := "https://www.youtube.com/watch?v=" + videoID
	client := &http.Client{}
	req, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		log.Println("Error al crear la solicitud HTTP: %v", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.147 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al realizar la solicitud HTTP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error al parsear el HTML: %v", err)
		return nil, err
	}
	title := doc.Find("meta[name='title']").AttrOr("content", "nada")
	// Buscar y extraer la descripción del video
	description := doc.Find("meta[name='description']").AttrOr("content", "nada")

	if title != "nada" && description != "nada" {

		tip := entity.PostNew{}

		tip.Content = description
		tip.Title = title
		tip.Type = "youtube"
		tip.URL = "https://www.youtube.com/embed/" + videoID

		return &tip, nil
	} else {
		return nil, fmt.Errorf("No encontro video")
	}
}
func (s *PortalRetroGamerImpl) loadVideosFromYoutubeChannel(channelURL, alias, search string) bool {
	realChannelURL := s.findInYoutubeURL(channelURL, search)
	fmt.Println("Busco en " + realChannelURL)
	channelURL = realChannelURL
	client := &http.Client{}
	req, err := http.NewRequest("GET", channelURL, nil)
	if err != nil {
		log.Println("Error al crear la solicitud HTTP: ", err)
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.147 Safari/537.36")

	// Enviar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al enviar la solicitud: ", err)
		return false
	}
	defer resp.Body.Close()

	// Parsear el HTML de la página
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	javascript := doc.Text()

	re := regexp.MustCompile(`\{([^}]*)\}`)
	if err != nil {
		return false
	}

	matches := re.FindAllString(javascript, 10000)
	if len(matches) > 1 {
		for _, values := range matches {

			if strings.Contains(values, "videoId") && strings.Contains(values, "videoRenderer") {
				videoID := strings.Split(values, "\"")[9]
				idVideo := "admin_" + videoID
				if videoID != "videoRenderer" {

					tipsFromDB := s.GetTipByID(idVideo)
					if tipsFromDB == nil {
						urlVideo := "https://www.youtube.com/embed/" + videoID
						tipsFromDB = s.GetTipByURL(urlVideo)
						if tipsFromDB == nil {

							tipsFromDB, err := s.GetNewTipsFromSearch(videoID)

							// ACA DEBERIAS HACER LA RUTINA
							if err != nil {
								fmt.Println("[!] URL ya existe en BASE DE DATOS " + urlVideo)
								return false
							}

							// Rutina para buscar términos del 'search' en Content y Title
							searchTerms := strings.Split(search, " ")
							found := false

							// Buscar en Content y Title
							for _, term := range searchTerms {
								if strings.Contains(strings.ToLower(tipsFromDB.Content), strings.ToLower(term)) ||
									strings.Contains(strings.ToLower(tipsFromDB.Title), strings.ToLower(term)) {
									found = true
									break
								}
							}

							// Si no encuentra coincidencias, retornamos false
							if !found {
								fmt.Println("[!] No se encontraron coincidencias en el contenido o título.")
								return false
							}

							// Actualización de valores si encuentra coincidencias
							tipsFromDB.ID = idVideo
							tipsFromDB.Author = alias
							currentTime := time.Now().UTC()
							formattedTime := currentTime.Format("2006-01-02T15:04:05.000Z")
							tipsFromDB.Date = formattedTime
							s.UserRepository.SaveTips(tipsFromDB)

							// Imprimir los resultados
							fmt.Println("[+] ID :", tipsFromDB.ID)
							fmt.Println("[+] Título del video:", tipsFromDB.Title)
							fmt.Println("[+] Descripción del video:", tipsFromDB.Content)

						}
					} else {
						fmt.Println("[!] idVideo ya existe en BASE DE DATOS " + idVideo)
						return false
					}
				}
			}
		}
	} else {
		fmt.Println("No se encontró contenido entre los corchetes.")
		return false
	}
	return true
}

func (s *PortalRetroGamerImpl) UpdateVideosTeams(search string) bool {
	usuarios, _ := s.GetAllUsers()

	var wg sync.WaitGroup

	for _, user := range usuarios {
		for _, rrss := range user.RRSS {
			if rrss.Type == "youtube" && strings.Contains(rrss.URL, "youtube.com/") {
				fmt.Println("Buscando en el canal de " + user.Alias)

				wg.Add(1)

				go func(url, alias, search string) {
					defer wg.Done()
					s.loadVideosFromYoutubeChannel(url, alias, search)
				}(rrss.URL, user.Alias, search)
			}
		}
	}

	// Espera a que todas las gorutinas finalicen
	wg.Wait()
	fmt.Println("Fin busqueda en canales para " + search)
	return true
}

func (s *PortalRetroGamerImpl) findInYoutubeURL(urlChannelTeams string, find string) string {
	videosString := "/videos"
	finalUrl := ""
	search := "/search?query=" + url.QueryEscape(find)

	if strings.Contains(urlChannelTeams, videosString) {
		finalUrl = urlChannelTeams
	} else if strings.Contains(urlChannelTeams, "/channel/") {
		finalUrl = urlChannelTeams + videosString
	} else if strings.Contains(urlChannelTeams, "/@") {
		if idx := strings.Index(urlChannelTeams, "?"); idx != -1 {
			urlChannelTeams = urlChannelTeams[:idx]
		}
		finalUrl = urlChannelTeams + videosString
	} else if strings.Contains(urlChannelTeams, "youtube.com/") {
		parts := strings.Split(urlChannelTeams, "youtube.com/")
		if len(parts) > 1 {
			finalUrl = "https://www.youtube.com/@" + parts[1] + videosString
		}
	}

	finalUrl = strings.ReplaceAll(finalUrl, videosString, search)

	return finalUrl
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (s *PortalRetroGamerImpl) AuthenticateUser(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.AuthenticateUser(alias, password)
	return usuarioOnline, err
}

func (s *PortalRetroGamerImpl) SetStatusLogin(alias, sessionToken, hash string, online bool) (bool, error) {
	usuarioOnline, err := s.UserRepository.SetUserOnline(alias, sessionToken, hash, online)
	if online {
		return usuarioOnline.Online, err
	} else {
		return false, err
	}
}

func (s *PortalRetroGamerImpl) GetStatusLogin(sessionToken, hash string) (*entity.UserOnline, error) {
	usuarioOnline, err := s.UserRepository.GetUserOnline(sessionToken, hash)
	return usuarioOnline, err
}
func (s *PortalRetroGamerImpl) GetUserByAlias(alias string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	return usuarioOnline, err
}
func (s *PortalRetroGamerImpl) GetUserByTextRefer(text string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByTextRefer(text)
	return usuarioOnline, err
}

func (s *PortalRetroGamerImpl) ChangePassword(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Password = password
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}

func (s *PortalRetroGamerImpl) SaveUser(user entity.User) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(user.Alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Name = user.Name
		usuarioOnline.RRSS = user.RRSS
		usuarioOnline.AvatarYT = user.AvatarYT
		usuarioOnline.ReferenceText = user.ReferenceText
		usuarioOnline.AboutMe = user.AboutMe
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}

func (s *PortalRetroGamerImpl) CreateUser(user *entity.User) error {
	error := s.UserRepository.SaveUser(user)
	return error
}

func (s *PortalRetroGamerImpl) GetAllUsers() ([]*entity.User, error) {

	users, err := s.PortalRepository.GetAllUsers()

	return users, err
}

func (s *PortalRetroGamerImpl) GetUserByRefer(refer string) (*entity.User, error) {

	user, err := s.PortalRepository.GetUserByAlias(refer)

	return user, err
}

func (s *PortalRetroGamerImpl) GetAllTips() ([]*entity.PostNew, error) {

	users, err := s.PortalRepository.GetAllTips()

	return users, err
}

func (s *PortalRetroGamerImpl) CreateTips(tip *entity.PostNew) error {
	error := s.UserRepository.SaveTips(tip)
	return error
}

func (s *PortalRetroGamerImpl) GetTipByID(id string) *entity.PostNew {

	tips, _ := s.PortalRepository.GetTipByID(id)
	return tips
}

func (s *PortalRetroGamerImpl) GetTipByURL(url string) *entity.PostNew {
	tips, _ := s.PortalRepository.GetTipByURL(url)
	return tips
}

func (s *PortalRetroGamerImpl) GetTipsWithPagination(skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsWithPagination(skip, limit, typeOfTips)
}

func (s *PortalRetroGamerImpl) GetTipsWithSearch(search string, skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsWithSearch(search, skip, limit, typeOfTips)
}

func (s *PortalRetroGamerImpl) DeleteTip(id, alias string) error {
	return s.PortalRepository.DeleteTipByIDandAuthor(id, alias)
}

func (s *PortalRetroGamerImpl) GetTipsByAliasWithPagination(alias string, skip, limit int64) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsByAliasWithPagination(alias, skip, limit)
}
