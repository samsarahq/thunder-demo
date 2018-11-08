package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/samsarahq/thunder/sqlgen"
)

const PlayerIdCtxKey = "playerId"

type (
	PlayerColor string
	PlayerName  string
)

var (
	// Set of A100 from https://material.io/design/color/#tools-for-picking-colors
	AssignablePlayerColors = []PlayerColor{
		"#FF5252",
		"#FF4081",
		"#E040FB",
		"#7C4DFF",
		"#536DFE",
		// "#448AFF", // main player color
		"#40C4FF",
		"#18FFFF",
		"#64FFDA",
		"#69F0AE",
		"#B2FF59",
		"#EEFF41",
		"#FFFF00",
		"#FFD740",
		"#FFAB40",
		"#FF6E40",
	}
	SuperheroNames []PlayerName
)

func init() {
	superheroNamesFile, err := os.Open("./superheroes.json")
	defer superheroNamesFile.Close()
	PanicIfErr(err)

	superheroNamesByteValue, _ := ioutil.ReadAll(superheroNamesFile)

	err = json.Unmarshal([]byte(superheroNamesByteValue), &SuperheroNames)
	PanicIfErr(err)
}

// Merged player + state for primary game.
type Player struct {
	Id    int64 `sql:",primary" graphql:",key"`
	Name  PlayerName
	Color PlayerColor
	X     int64
	Y     int64
}

func PlayerId(ctx context.Context) *int64 {
	val := ctx.Value(PlayerIdCtxKey)
	if val == nil {
		return nil
	}
	return val.(*int64)
}

func WithPlayerId(ctx context.Context, playerId int64) context.Context {
	return context.WithValue(ctx, PlayerIdCtxKey, int64Ptr(playerId))
}

func NewPlayer() *Player {
	return &Player{
		Name:  SuperheroNames[rand.Intn(len(SuperheroNames))],
		Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
	}
}

func (s *Server) GetPlayerById(ctx context.Context, args struct{ Id int64 }) (*Player, error) {
	var result *Player
	err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (s *Server) GetPlayers(ctx context.Context) ([]*Player, error) {
	var result []*Player
	if err := s.db.Query(ctx, &result, nil, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) GetCurrentPlayer(ctx context.Context) (*Player, error) {
	id := PlayerId(ctx)
	log.Println("currentPlayer", id)
	if id == nil {
		return nil, nil
	}
	var result *Player
	err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": id}, nil)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (s *Server) CreatePlayer(ctx context.Context, args struct{ Name PlayerName }) error {
	_, err := s.db.InsertRow(ctx, &Player{Name: args.Name})
	return err
}

type updatePlayerSelectionArgs struct {
	PlayerId int64
	X        int64
	Y        int64
}

func (s *Server) UpdatePlayerSelection(ctx context.Context, args updatePlayerSelectionArgs) error {
	player := &Player{}
	if err := s.db.QueryRow(ctx, &player, sqlgen.Filter{"id": args.PlayerId}, nil); err != nil {
		return err
	}
	player.X = args.X
	player.Y = args.Y
	return s.db.UpdateRow(ctx, player)
}
