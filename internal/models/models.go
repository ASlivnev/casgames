package models

type Game struct {
	Id                *int    `json:"id"`
	GameId            *string `json:"game_id"`
	GameNameEn        *string `json:"game_name_en"`
	GameDescriptionEn *string `json:"game_description_en"`
	GameDeveloper     *string `json:"game_developer"`
	GameUrlName       *string `json:"game_url_name"`
	GameRang          *int    `json:"game_rang"`
}

type InsertIntoDbGame struct {
	GameId        string `json:"game_id"`
	GameNameEn    string `json:"game_name_en"`
	GameDeveloper string `json:"game_developer"`
	GameUrlName   string `json:"game_url_name"`
}

type ResponseGd struct {
	Data DataGd `json:"data"`
}

type DataGd struct {
	GamesSearched GamesSearchedGd `json:"gamesSearched"`
}

type GamesSearchedGd struct {
	HitsPerPage int     `json:"hitsPerPage"`
	NbHits      int     `json:"nbHits"`
	NbPages     int     `json:"nbPages"`
	Page        int     `json:"page"`
	Hits        []HitGd `json:"hits"`
}

type HitGd struct {
	ObjectID      string   `json:"objectID"`
	Title         string   `json:"title"`
	Company       string   `json:"company"`
	Visible       bool     `json:"visible"`
	ExclusiveGame int      `json:"exclusiveGame"`
	SlugsGd       []SlugGd `json:"slugs"`
}

type SlugGd struct {
	Name string `json:"name"`
}
