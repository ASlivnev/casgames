package repo

import (
	"casualgames/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"casualgames/internal/config"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepository(cnf *config.Cnf) *Repo {
	pool, err := NewPgxPool(context.Background(), cnf)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	return &Repo{db: pool}
}

func NewPgxPool(ctx context.Context, cnf *config.Cnf) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cnf.Db.User, cnf.Db.Pass, cnf.Db.Host, cnf.Db.Port, cnf.Db.Name)
	log.Printf(dsn)
	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	log.Info().Msg("Database connected!")

	return pool, nil
}

func (repo *Repo) DeallocateAll() error {
	sql := `DEALLOCATE PREPARE ALL;`
	query, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] DeallocateAll : " + err.Error())
	}
	defer query.Close()

	return err
}

func (repo *Repo) InsertGame(game models.InsertIntoDbGame) {
	sqlStatement := `
		INSERT INTO casualgames.games (
			game_id, 
			game_name_en, 
			game_developer, 
			game_url_name
		) VALUES ($1, $2, $3, $4) 
		ON CONFLICT DO NOTHING`
	_, err := repo.db.Exec(
		context.Background(),
		sqlStatement,
		game.GameId,
		game.GameNameEn,
		game.GameDeveloper,
		game.GameUrlName,
	)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (repo *Repo) GetGames(page int) []models.Game {
	page = page * 100
	sqlStatement := `
		SELECT 
			id, 
			game_id, 
			game_name_en, 
			game_description_en, 
			game_developer, 
			game_url_name, 
			game_rang 
		FROM casualgames.games 
		    ORDER BY game_rang DESC
		LIMIT 100 OFFSET $1`
	rows, err := repo.db.Query(context.Background(), sqlStatement, page)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer rows.Close()

	games := make([]models.Game, 0)
	for rows.Next() {
		var game models.Game
		err = rows.Scan(
			&game.Id,
			&game.GameId,
			&game.GameNameEn,
			&game.GameDescriptionEn,
			&game.GameDeveloper,
			&game.GameUrlName,
			&game.GameRang,
		)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		games = append(games, game)
	}

	return games
}

func (repo *Repo) IncrementGameRang(gameId string) {
	sqlStatement := `
		UPDATE casualgames.games 
		SET game_rang = game_rang + 1
		WHERE game_id = $1`
	_, err := repo.db.Exec(
		context.Background(),
		sqlStatement,
		gameId,
	)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
