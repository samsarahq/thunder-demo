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

type Server struct {
	db *livesql.LiveDB
}

type Game struct {
	Id   int64 `sql:",primary" graphql:",key"`
	State	int32
	Data string
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

func (s *Server) SchemaBuilderSchema() *schemabuilder.Schema {
	schema := schemabuilder.NewSchema()

	s.registerGameQueries(schema)
	s.registerGameMutations(schema)

	return schema
}

func (s *Server) Schema() *graphql.Schema {
	return s.SchemaBuilderSchema().MustBuild()
}

func main() {
	sqlgenSchema := sqlgen.NewSchema()
	sqlgenSchema.MustRegisterType("games", sqlgen.AutoIncrement, Game{})

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
