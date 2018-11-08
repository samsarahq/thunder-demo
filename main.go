package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/jkomoros/sudoku"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/livesql"
	"github.com/samsarahq/thunder/sqlgen"
)

const (
	DbName = "sudoku"
)

type Server struct {
	db *livesql.LiveDB
}

type Game struct {
	Id     int64 `sql:",primary" graphql:",key"`
	State  string
	Data   string
	Name   string
	Solved bool
}

type Message struct {
	Id   int64 `sql:",primary" graphql:",key"`
	SentBy PlayerName
	Color PlayerColor
	Text string
}

func checkPuzzle(puzzle string) bool {
	grid := sudoku.LoadSDK(puzzle)
	return grid.Solved()
}

func (s *Server) registerGameQueries(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("game", func(ctx context.Context, args struct{ Id int64 }) (*Game, error) {
		var result *Game
		if err := s.db.QueryRow(ctx, &result, sqlgen.Filter{"id": args.Id}, nil); err != nil {
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
	// TODO: track mutations to player focus in backend
	object.FieldFunc("players", func(ctx context.Context, g *Game) ([]*Player, error) {
		var result []*Player
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})
}

func (s *Server) registerPlayerQueries(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("currentPlayer", s.GetCurrentPlayer)
	object.FieldFunc("player", s.GetPlayerById)
	object.FieldFunc("players", s.GetPlayers)
}

func (s *Server) registerGameMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("createGame", func(ctx context.Context, args struct{ Name string }) (*Game, error) {
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
	})

	type updateGameArgs struct {
		Id  int64
		Row int16
		Col int16
		Val int16
	}
	object.FieldFunc("updateGame", func(ctx context.Context, args updateGameArgs) error {
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
	})
}

func (s *Server) registerPlayerMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("createPlayer", s.CreatePlayer)
	object.FieldFunc("updatePlayerSelection", s.UpdatePlayerSelection)
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

	object.FieldFunc("addMessage", func(ctx context.Context, args struct{ Text string; SentBy PlayerName; Color PlayerColor }) error {
		_, err := s.db.InsertRow(ctx, &Message{Text: args.Text, SentBy: args.SentBy, Color: args.Color})
		return err
	})

	object.FieldFunc("deleteMessage", func(ctx context.Context, args struct{ Id int64 }) error {
		return s.db.DeleteRow(ctx, &Message{Id: args.Id})
	})
}


func int64Ptr(i int64) *int64 { return &i }

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


type executionLogger struct {}
func (s *executionLogger) StartExecution(ctx context.Context, tags map[string]string, initial bool) {}
func (s *executionLogger) FinishExecution(ctx context.Context, tags map[string]string, delay time.Duration) {}
func (s *executionLogger) Error(ctx context.Context, err error, tags map[string]string) {
	log.Printf("error:%v\n%s", tags, err)
}

type subscriptionLogger struct{
	server *Server
}

func PanicIfErr(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}

func (l *subscriptionLogger) Subscribe(ctx context.Context, id string, tags map[string]string) {
	intId, err := strconv.ParseInt(id, 10, 64)
	PanicIfErr(err)

	log.Println("~~ Subscribe", intId)
}

func (l *subscriptionLogger) Unsubscribe(ctx context.Context, id string) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic("error parsing subscription id")
	}
	log.Println("~~ Unsubscribe", intId)
}

func handlerWithPlayerTracking(schema *graphql.Schema, server *Server) http.Handler {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrader.Upgrade: %v", err)
			return
		}
		defer socket.Close()

		ctx := r.Context()
		response, err := server.db.InsertRow(ctx, NewPlayer())
		PanicIfErr(err)

		playerId, err := response.LastInsertId()
		PanicIfErr(err)
		ctx = WithPlayerId(ctx, playerId)
		log.Println("~~ created playerId", playerId)

		conn := graphql.CreateConnection(ctx, socket, schema, graphql.WithExecutionLogger(&executionLogger{}), graphql.WithSubscriptionLogger(&subscriptionLogger{server: server}))
		conn.ServeJSONSocket()

		err = server.db.DeleteRow(ctx, &Player{ Id: playerId })
		PanicIfErr(err)
		log.Println("~~ deleted playerId", playerId)
	})
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

	server := &Server{
		db: db,
	}
	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	http.Handle("/graphql", handlerWithPlayerTracking(graphqlSchema, server))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	log.Println("== STARTED ==")
	if err := http.ListenAndServe(":3030", nil); err != nil {
		panic(err)
	}
}
