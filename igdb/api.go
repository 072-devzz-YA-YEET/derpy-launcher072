package igdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type Image struct {
	ImageID string `json:"image_id"`
}

type ApiGame struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"summary"`
	CoverID     int    `json:"cover"`
	MainCover   string
	ArtworkIDList   []int `json:"artworks"`
	ArtworkLinkList []string
	Screenshots     []Image `json:"screenshots"`
	ScreenshotLinkList []string
}

type APIManager struct {
	client *http.Client
}

func SetupHeader(request *http.Request) {
	request.Header.Set("Client-ID", os.Getenv("IGDB_CLIENT"))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("IGDB_AUTH")))
}

// Special error so we can check in main
var (
	ErrorNoCoversFound = errors.New("could not find a cover with this id")
)

func NewAPI() *APIManager {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return &APIManager{client: &http.Client{}}
}

func (a *APIManager) GetGameData(id int) (ApiGame, error) {
	header := fmt.Sprintf(`fields id, name, summary, cover, artworks, screenshots.*; where id = %d;`, id)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return ApiGame{}, err
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return ApiGame{}, err
	}
	defer response.Body.Close()

	var gameDataList []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&gameDataList)
	if jsonErr != nil {
		return ApiGame{}, err
	}

	if len(gameDataList) == 0 {
		return ApiGame{}, fmt.Errorf("no games found with id %d", id)
	}

	firstGameData := gameDataList[0]
	coverImageUrl, err := a.GetCover(firstGameData.CoverID)
	if err != nil {
		// return game without cover
		return firstGameData, err
	}
	firstGameData.MainCover = coverImageUrl
	bannerImageUrls, err := a.GetArtworkURLs(firstGameData.ArtworkIDList)
	for _, url := range bannerImageUrls {
		fmt.Println("Link " + url)
	}
	if err != nil {
		// return game without banners
		return firstGameData, err
	}
	firstGameData.ArtworkLinkList = bannerImageUrls

	for _, image := range firstGameData.Screenshots {
		imageID := image.ImageID
		imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", imageID)
		firstGameData.ScreenshotLinkList = append(firstGameData.ScreenshotLinkList, imageURL)
		fmt.Println("added " + imageURL)
	}


	return firstGameData, nil
}

func (a *APIManager) GetGames(query string) []ApiGame {
	header := fmt.Sprintf(`fields id, name, summary, cover; search "%s";`, query)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return []ApiGame{}
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return []ApiGame{}
	}
	defer response.Body.Close()

	var games []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&games)
	if jsonErr != nil {
		return []ApiGame{}
	}
	return games
}
