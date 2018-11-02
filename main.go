package main

import (
	"context"
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

var ErrNoRows = "sql: no rows in result set"

type Server struct {
	db *livesql.LiveDB
}

type Game struct {
	Id   int64 `sql:",primary" graphql:",key"`
	State	int32
	Data string
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
}

func (s *Server) registerPlayerQueries(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("player", func(ctx context.Context, args struct{ Id int64 }) (*Player, error) {
		var result *Player
		err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil)
		if err.Error() == ErrNoRows {
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

	object.FieldFunc("createGame", func(ctx context.Context, args struct{ Data string }) error {
		grid := sudoku.GenerateGrid(sudoku.DefaultGenerationOptions())
		_, err := s.db.InsertRow(ctx, &Game{Data: grid.DataString()})
		return err
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

	db, err := livesql.Open("localhost", 3307, "root", "", "sudoku", sqlgenSchema)
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
