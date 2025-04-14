package games

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	localGameList = GameList{}
	CacheGames    = make(map[string][]byte)
	CacheMux      sync.Mutex
)

func LoadGameList(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&localGameList); err != nil {
		return err
	}

	return nil
}

func CacheURLGamesContent(url, key string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	CacheMux.Lock()
	CacheGames[key] = body
	CacheMux.Unlock()

	return nil
}

func InsertCoin(emulator, name string) ([]byte, error) {

	CacheMux.Lock()
	data, ok := CacheGames[name]
	CacheMux.Unlock()

	if !ok {
		// Search for the game in the game list
		var gameURL string
		for _, game := range localGameList.Games {
			if game.Name == name && game.Emulator == emulator {
				gameURL = game.URL
				break
			}
		}

		if gameURL == "" {
			return nil, fmt.Errorf("No encontro Juego")
		}

		// Cache the content
		err := CacheURLGamesContent(gameURL, name)
		if err != nil {
			return nil, fmt.Errorf("No encontro Juego")
		}

		CacheMux.Lock()
		data = CacheGames[name]
		CacheMux.Unlock()
	}

	return data, nil
}
