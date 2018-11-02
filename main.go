package main

import (
	"context"
	"database/sql"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/livesql"
	"github.com/samsarahq/thunder/sqlgen"
	"github.com/jkomoros/sudoku"
)

const (
	DbName = "sudoku"
)

type (
	PlayerColor string
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
)

type Server struct {
	db *livesql.LiveDB
}

type Game struct {
	Id   int64 `sql:",primary" graphql:",key"`
	State	int32
	Data string
	Name string
}

// For a single game.
type PlayerState struct {
	PlayerId int64
	Color PlayerColor
	X int64
	Y int64
}

type Player struct {
	Id int64 `sql:",primary" graphql:",key"`
	Name string
}

type Message struct {
	Id   int64 `sql:",primary" graphql:",key"`
	Text string
}

func checkPuzzle(puzzle string) (bool) {
	grid := sudoku.LoadSDK(puzzle)
	return grid.Solved()
}

func (s *Server) registerGameQueries(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("game", func(ctx context.Context, args struct{ Id int64 }) (*Game, error) {
		var result *Game
		if err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil); err !=nil {
			return nil, err
		}
		return result, nil
	})

	object.FieldFunc("games", func(ctx context.Context) ([]*Game, error) {
		var result []*Game
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})

	// Game Field Funcs
	object = schema.Object("Game", Game{})
	object.FieldFunc("playerStates", func(ctx context.Context, g *Game) ([]*PlayerState, error) {
		return []*PlayerState{
			&PlayerState{
				PlayerId: 1,
				Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
				X: 3,
				Y: 2,
			},
			&PlayerState{
				PlayerId: 1,
				Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
				X: 6,
				Y: 7,
			},
			&PlayerState{
				PlayerId: 1,
				Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
				X: 9,
				Y: 8,
			},
			&PlayerState{
				PlayerId: 1,
				Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
				X: 1,
				Y: 8,
			},
		}, nil
	})
}

func (s *Server) registerPlayerQueries(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("player", func(ctx context.Context, args struct{ Id int64 }) (*Player, error) {
		var result *Player
		err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	object.FieldFunc("players", func(ctx context.Context) ([]*Player, error) {
		var result []*Player
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})
}

func (s *Server) registerGameMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("createGame", func(ctx context.Context, args struct{ Name string }) (*Game, error) {
		grid := sudoku.GenerateGrid(sudoku.DefaultGenerationOptions())
		game := Game{Name: args.Name, Data: grid.DataString()}
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
	})


	type updateGameArgs struct {
		Id int64
		Row int16
		Col int16
		Val int16
	}
	object.FieldFunc("updateGame", func(ctx context.Context, args updateGameArgs) error {
		var game *Game
		if err := s.db.QueryRow(ctx, &game, sqlgen.Filter{"id": args.Id}, nil); err != nil {
			return err
		}

		grid := sudoku.MutableLoadSDK(game.Data)
		grid.MutableCell(int(args.Row),int(args.Col)).SetNumber(int(args.Val))
		game.Data = grid.DataString()

		err := s.db.UpdateRow(ctx, game)
		return err
	})
}

func (s *Server) registerPlayerMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("createPlayer", func(ctx context.Context, args struct{ Name string }) error {
		_, err := s.db.InsertRow(ctx, &Player{Name: args.Name})
		return err
	})
}

func (s *Server) registerMessageQuery(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("messages", func(ctx context.Context) ([]*Message, error) {
		var result []*Message
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})
}

func (s *Server) registerMessageMutation(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("addMessage", func(ctx context.Context, args struct{ Text string }) error {
		_, err := s.db.InsertRow(ctx, &Message{Text: args.Text})
		return err
	})

	object.FieldFunc("deleteMessage", func(ctx context.Context, args struct{ Id int64 }) error {
		return s.db.DeleteRow(ctx, &Message{Id: args.Id})
	})
}

func (s *Server) SchemaBuilderSchema() *schemabuilder.Schema {
	schema := schemabuilder.NewSchema()

	s.registerGameQueries(schema)
	s.registerGameMutations(schema)
	s.registerPlayerQueries(schema)
	s.registerPlayerMutations(schema)
	s.registerMessageQuery(schema)
	s.registerMessageMutation(schema)

	return schema
}

func (s *Server) Schema() *graphql.Schema {
	return s.SchemaBuilderSchema().MustBuild()
}

func main() {
	sqlgenSchema := sqlgen.NewSchema()
	sqlgenSchema.MustRegisterType("games", sqlgen.AutoIncrement, Game{})
	sqlgenSchema.MustRegisterType("players", sqlgen.AutoIncrement, Player{})
	sqlgenSchema.MustRegisterType("messages", sqlgen.AutoIncrement, Message{})

	db, err := livesql.Open("localhost", 3307, "root", "", DbName, sqlgenSchema)
	if err != nil {
		panic(err)
	}

	server := &Server{db: db}
	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	http.Handle("/graphql", graphql.Handler(graphqlSchema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	if err := http.ListenAndServe(":3030", nil); err != nil {
		panic(err)
	}
}
