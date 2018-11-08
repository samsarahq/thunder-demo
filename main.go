package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
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

func (s *Server) registerGameQueries(schema *schemabuilder.Schema) {
	object := schema.Query()
	object.FieldFunc("game", s.GetGameById)
	object.FieldFunc("games", s.GetAllGames)

	// Game Field Funcs
	object = schema.Object("Game", Game{})
	object.FieldFunc("players", s.GetPlayersForGame)
}

func (s *Server) registerPlayerQueries(schema *schemabuilder.Schema) {
	object := schema.Query()
	object.FieldFunc("currentPlayer", s.GetCurrentPlayer)
	object.FieldFunc("player", s.GetPlayerById)
	object.FieldFunc("players", s.GetPlayers)
}

func (s *Server) registerGameMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()
	object.FieldFunc("createGame", s.CreateGame)
	object.FieldFunc("updateGame", s.UpdateGame)
}

func (s *Server) registerPlayerMutations(schema *schemabuilder.Schema) {
	object := schema.Mutation()
	object.FieldFunc("createPlayer", s.CreatePlayer)
	object.FieldFunc("updatePlayerSelection", s.UpdatePlayerSelection)
}

func (s *Server) registerMessageQuery(schema *schemabuilder.Schema) {
	object := schema.Query()
	object.FieldFunc("messages", s.GetMessages)
}

func (s *Server) registerMessageMutation(schema *schemabuilder.Schema) {
	object := schema.Mutation()
	object.FieldFunc("addMessage", s.AddMessage)
	object.FieldFunc("deleteMessage", s.DeleteMessage)
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

type executionLogger struct{}

func (s *executionLogger) StartExecution(ctx context.Context, tags map[string]string, initial bool) {}
func (s *executionLogger) FinishExecution(ctx context.Context, tags map[string]string, delay time.Duration) {
}
func (s *executionLogger) Error(ctx context.Context, err error, tags map[string]string) {
	log.Printf("error:%v\n%s", tags, err)
}

type subscriptionLogger struct {
	server *Server
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

		err = server.db.DeleteRow(ctx, &Player{Id: playerId})
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
	PanicIfErr(err)

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
