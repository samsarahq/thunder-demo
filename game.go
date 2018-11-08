package main

import (
	"context"
	"errors"

	"github.com/jkomoros/sudoku"
	"github.com/samsarahq/thunder/sqlgen"
)

type Game struct {
	Id     int64 `sql:",primary" graphql:",key"`
	State  string
	Data   string
	Name   string
	Solved bool
}

func checkPuzzle(puzzle string) bool {
	grid := sudoku.LoadSDK(puzzle)
	return grid.Solved()
}

func (s *Server) GetGameById(ctx context.Context, args struct{ Id int64 }) (*Game, error) {
	var result *Game
	if err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) GetAllGames(ctx context.Context) ([]*Game, error) {
	var result []*Game
	if err := s.db.Query(ctx, &result, nil, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) GetPlayersForGame(ctx context.Context, g *Game) ([]*Player, error) {
	var result []*Player
	if err := s.db.Query(ctx, &result, nil, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) CreateGame(ctx context.Context, args struct{ Name string }) (*Game, error) {
	grid := sudoku.GenerateGrid(sudoku.DefaultGenerationOptions())
	gridString := grid.DataString()
	game := Game{Name: args.Name, Data: gridString, State: gridString}
	res, err := s.db.InsertRow(ctx, &game)
	if err != nil {
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	game.Id = lastInsertId
	return &game, nil
}

type updateGameArgs struct {
	Id  int64
	Row int16
	Col int16
	Val int16
}

func (s *Server) UpdateGame(ctx context.Context, args updateGameArgs) error {
	var game *Game
	if err := s.db.QueryRow(ctx, &game, sqlgen.Filter{"id": args.Id}, nil); err != nil {
		return err
	}
	if game.Solved {
		return nil
	}
	var r, c = int(args.Row), int(args.Col)
	if origState := sudoku.LoadSDK(game.Data).Cell(r, c).Number(); origState != 0 {
		return errors.New("Can't change original cell")
	}

	grid := sudoku.MutableLoadSDK(game.State)
	grid.MutableCell(r, c).SetNumber(int(args.Val))
	if grid.Solved() {
		game.Solved = true
	}
	game.State = grid.DataString()

	err := s.db.UpdateRow(ctx, game)
	return err
}
