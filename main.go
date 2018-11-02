package main

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"log"
	"time"
	"sync"

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
	playerStateMap map[string]*PlayerState
	playerStateMu sync.Mutex
}

type Game struct {
	Id    int64 `sql:",primary" graphql:",key"`
	State string
	Data  string
	Name  string
}

// For a single game.
type PlayerState struct {
	PlayerId string
	Color PlayerColor
	X int64
	Y int64
}

type Player struct {
	Id   int64 `sql:",primary" graphql:",key"`
	Name string
}

type Message struct {
	Id   int64 `sql:",primary" graphql:",key"`
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
	object.FieldFunc("playerStates", func(ctx context.Context, g *Game) []*PlayerState {
		playerStates := []*PlayerState{}
		for _, player := range s.playerStateMap {
			playerStates = append(playerStates, player)
		}
		return playerStates
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
		var r, c = int(args.Row), int(args.Col)
		if origState := sudoku.LoadSDK(game.Data).Cell(r, c).Number(); origState != 0 {
			return errors.New("Can't change original cell")
		}

		grid := sudoku.MutableLoadSDK(game.State)
		grid.MutableCell(r, c).SetNumber(int(args.Val))
		game.State = grid.DataString()

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

func (s *Server) newPlayerState(id string) *PlayerState {
	return &PlayerState{
		PlayerId: id,
		Color: AssignablePlayerColors[rand.Intn(len(AssignablePlayerColors))],
	}
}

func (s *Server) AddPlayerToStateMap(id string) {
	s.playerStateMu.Lock()
	defer s.playerStateMu.Unlock()

	s.playerStateMap[id] = s.newPlayerState(id)
}

func (s *Server) RemovePlayerFromStateMap(id string) {
	s.playerStateMu.Lock()
	defer s.playerStateMu.Unlock()

	delete(s.playerStateMap, id)
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

func (l *subscriptionLogger) Subscribe(ctx context.Context, id string, tags map[string]string) {
	log.Println("~~ Subscribe", id, tags)
	l.server.AddPlayerToStateMap(id)
}

func (l *subscriptionLogger) Unsubscribe(ctx context.Context, id string) {
	log.Println("~~ Unsubscribe", id)
	l.server.RemovePlayerFromStateMap(id)
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

		conn := graphql.CreateConnection(r.Context(), socket, schema, graphql.WithExecutionLogger(&executionLogger{}), graphql.WithSubscriptionLogger(&subscriptionLogger{server: server}))
		conn.ServeJSONSocket()
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

	log.Println("== STARTED ==")

	server := &Server{
		db: db,
		playerStateMap: make(map[string]*PlayerState),
	}
	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	http.Handle("/graphql", handlerWithPlayerTracking(graphqlSchema, server))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	if err := http.ListenAndServe(":3030", nil); err != nil {
		panic(err)
	}
}
