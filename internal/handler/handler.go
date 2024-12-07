package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"casualgames/internal/config"
	"casualgames/internal/models"
	"casualgames/internal/repo"
)

type Handler struct {
	pool *repo.Repo
	cnf  *config.Cnf
}

func NewHandler(pool *repo.Repo, cnf *config.Cnf) *Handler {
	return &Handler{
		pool: pool,
		cnf:  cnf,
	}
}

func (h *Handler) DeallocateAll(c *fiber.Ctx) error {
	err := h.pool.DeallocateAll()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return c.SendString("Deallocated all prepared statements")
}

func (h *Handler) GamesParser(c *fiber.Ctx) error {
	page := 0
	for {
		if page > 50 {
			break
		}

		response := getGamesList(page)
		games := parseGamesFromGd(response)
		for _, game := range games {
			h.pool.InsertGame(game)
		}
		log.Info().Msg("Page - " + fmt.Sprintf("%d", page))
		page++
	}

	return c.JSON("Page - " + fmt.Sprintf("%d", page))
}

func (h *Handler) GetGames(c *fiber.Ctx) error {
	page := c.Params("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON("Page is not a number")
		return nil
	}
	games := h.pool.GetGames(pageInt)
	c.JSON(games)
	return nil
}

func (h *Handler) IncrementGameRang(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	h.pool.IncrementGameRang(gameId)
	return nil
}

func getGamesList(page int) models.ResponseGd {
	url := "https://gd-website-api.gamedistribution.com/graphql"
	method := "POST"

	//convert to string page
	pageStr := fmt.Sprintf("%d", page)

	payload := `{"query":"fragment CoreGame on SearchHit {\n  objectID\n  title\n  company\n  visible\n  exclusiveGame\n  slugs {\n    name\n    __typename\n  }\n  assets {\n    name\n    __typename\n  }\n  __typename\n}\n\nquery GetGamesSearched($id: String! = \"\", $perPage: Int! = 0, $page: Int! = 0, $search: String! = \"\", $UIfilter: UIFilterInput! = {}, $filters: GameSearchFiltersFlat! = {}, $sortBy: KnownOrder, $sortByGeneric: [String!], $sortByCountryPerf: SortByCountryPerf! = {}, $sortByGenericWithDirection: [SortByGenericWithDirection!]) {\n  gamesSearched(\n    input: {collectionObjectId: $id, hitsPerPage: $perPage, page: $page, search: $search, UIfilter: $UIfilter, filters: $filters, sortBy: $sortBy, sortByCountryPerf: $sortByCountryPerf, sortByGeneric: $sortByGeneric, sortByGenericWithDirection: $sortByGenericWithDirection}\n  ) {\n    hitsPerPage\n    nbHits\n    nbPages\n    page\n    hits {\n      ...CoreGame\n      __typename\n    }\n    filters {\n      title\n      key\n      type\n      values\n      __typename\n    }\n    __typename\n  }\n}","variables":{"id":"","perPage":100,"page":` + pageStr + `,"search":"","UIfilter":{},"filters":{}}}`

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {

	}

	// Headers
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "ru,en-US;q=0.9,en;q=0.8,ru-RU;q=0.7")
	req.Header.Add("apollographql-client-name", "GDWebSite")
	req.Header.Add("apollographql-client-version", "1.0")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("dnt", "1")
	req.Header.Add("origin", "https://gamedistribution.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://gamedistribution.com/")
	req.Header.Add("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer res.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	var response models.ResponseGd
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return response
}

func parseGamesFromGd(response models.ResponseGd) []models.InsertIntoDbGame {
	var games []models.InsertIntoDbGame
	for _, hit := range response.Data.GamesSearched.Hits {
		var link string
		if len(hit.SlugsGd) > 0 {
			link = hit.SlugsGd[0].Name
		}

		game := models.InsertIntoDbGame{
			GameId:        hit.ObjectID,
			GameNameEn:    hit.Title,
			GameDeveloper: hit.Company,
			GameUrlName:   link,
		}
		games = append(games, game)
	}
	return games
}
