package main

import "context"

type Message struct {
	Id     int64 `sql:",primary" graphql:",key"`
	SentBy PlayerName
	Color  PlayerColor
	Text   string
}

func (s *Server) GetMessages(ctx context.Context) ([]*Message, error) {
	var result []*Message
	if err := s.db.Query(ctx, &result, nil, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) AddMessage(ctx context.Context, args struct {
	Text   string
	SentBy PlayerName
	Color  PlayerColor
}) error {
	_, err := s.db.InsertRow(ctx, &Message{Text: args.Text, SentBy: args.SentBy, Color: args.Color})
	return err
}

func (s *Server) DeleteMessage(ctx context.Context, args struct{ Id int64 }) error {
	return s.db.DeleteRow(ctx, &Message{Id: args.Id})
}
