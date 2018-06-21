package main

import (
	"context"
	"net/http"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/livesql"
	"github.com/samsarahq/thunder/sqlgen"
)

type Server struct {
	db *livesql.LiveDB
}

type Repo struct {
	Id       int64  `sql:",primary"` // from github
	FullName string // from github e.g. "facebookresearch/DensePose"
	ApiJson  []byte // binary blob of JSON from github api
}

type Event struct {
	AtMs    int64  `sql:",primary"` // timestamp of this event
	RepoId  int64  `sql:",primary"` // id of Repo in repos table
	EventId string `sql:",primary"` // id of this event (can be different for different types of events)
	ApiJson []byte // binary blob of JSON from github api
}

func (s *Server) registerRepo(schema *schemabuilder.Schema) {
	object := schema.Object("Repo", Repo{})

	object.Description = "A single repo"

}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	object := schema.Query()

	object.FieldFunc("repos", func(ctx context.Context) ([]*Repo, error) {
		var result []*Repo
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})

	// object.FieldFunc("repo", func(ctx context.Context, args struct{ RepoId int64 }) ([]*Event, error) {
	// 	//get events from db
	// 	return nil, nil
	// })
	//no
	//do we want to wrap the args into a RepoArgs struct?
	object.FieldFunc("events", func(ctx context.Context, args struct{ RepoID int64 }) ([]*Event, error) {
		var result []*Event
		if err := s.db.Query(ctx, &result, nil, nil); err != nil {
			return nil, err
		}
		return result, nil
	})
}

func (s *Server) registerMutation(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	object.FieldFunc("addRepo", func(ctx context.Context, args struct{ RepoName string }) error {
		_, err := s.db.InsertRow(ctx, &Repo{})
		return err
	})

	// object.FieldFunc("deleteEmptyStruct", func(ctx context.Context, args struct{ Id int64 }) error {
	// 	return s.db.DeleteRow(ctx, &EmptyStruct{Id: args.Id})
	// })

}

func (s *Server) importRepo(repoName string) error {
	// resp, err := http.Get("https://api.github.com/repos/" + repoName)
	// if err != nil {
	// 	return oops.Wrapf(err, "http get to github api for repo name")
	// }
	return nil
}

func (s *Server) SchemaBuilderSchema() *schemabuilder.Schema {
	schema := schemabuilder.NewSchema()

	s.registerQuery(schema)
	s.registerMutation(schema)

	return schema
}

func (s *Server) Schema() *graphql.Schema {
	return s.SchemaBuilderSchema().MustBuild()
}

func main() {
	sqlgenSchema := sqlgen.NewSchema()
	sqlgenSchema.MustRegisterType("repos", sqlgen.UniqueId, Repo{})
	sqlgenSchema.MustRegisterType("events", sqlgen.UniqueId, Event{})

	db, err := livesql.Open("localhost", 3307, "root", "", "github", sqlgenSchema)
	if err != nil {
		panic(err)
	}

	// importDummyData(db)

	server := &Server{db: db}
	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	http.Handle("/graphql", graphql.Handler(graphqlSchema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	if err := http.ListenAndServe(":3030", nil); err != nil {
		panic(err)
	}
}

func importDummyData(db *livesql.LiveDB) {
	dummyRepos := []Repo{
		{
			Id:       int64(1),
			FullName: "samsarahq/thunder",
		},
		{
			Id:       int64(2),
			FullName: "facebookresearch/DeepPose",
		},
	}
	for _, repo := range dummyRepos {
		_, err := db.InsertRow(context.TODO(), &repo)
		if err != nil {
			panic(err)
		}
	}

	dummyEvents := []Event{
		{
			AtMs:    int64(1529612944138),
			RepoId:  1,
			EventId: "39c44e49e5b90018e206203bb857ef9a9c2ce8f6",
			ApiJson: []byte(commitJson1),
		},
		{
			AtMs:    int64(1529612884138),
			RepoId:  2,
			EventId: "7300a44587b7c2394535a3f2bfeb594627b6ac91",
			ApiJson: []byte(commitJson2),
		},
	}
	for _, event := range dummyEvents {
		_, err := db.InsertRow(context.TODO(), &event)
		if err != nil {
			panic(err)
		}
	}
}

var commitJson1 = `
{
	"sha": "39c44e49e5b90018e206203bb857ef9a9c2ce8f6",
	"node_id": "MDY6Q29tbWl0MTM2MDY5MzExOjM5YzQ0ZTQ5ZTViOTAwMThlMjA2MjAzYmI4NTdlZjlhOWMyY2U4ZjY=",
	"commit": {
		"author": {
			"name": "Ilija Radosavovic",
			"email": "ilija.radosavovic@gmail.com",
			"date": "2018-06-21T15:02:53Z"
		},
		"committer": {
			"name": "GitHub",
			"email": "noreply@github.com",
			"date": "2018-06-21T15:02:53Z"
		},
		"message": "Merge pull request #23 from ImmortalTurtle/fix-links-in-install-md\n\nFix links annotations in INSTALL.md",
		"tree": {
			"sha": "a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68",
			"url": "https://api.github.com/repos/facebookresearch/DensePose/git/trees/a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68"
		},
		"url": "https://api.github.com/repos/facebookresearch/DensePose/git/commits/39c44e49e5b90018e206203bb857ef9a9c2ce8f6",
		"comment_count": 0,
		"verification": {
			"verified": true,
			"reason": "valid",
			"signature": "-----BEGIN PGP SIGNATURE-----\n\nwsBcBAABCAAQBQJbK74dCRBK7hj4Ov3rIwAAdHIIAKWFbcfuy1GUE+UG5UWgckHA\nGsJq68AZhZGC7vucwkT8CDuZblXw73vMZ2MtKmsO0J7Ku/xJi+3snotqEBlrjBh7\nJ7PAIjZv2xorl7QM0lx5DLZjQu6n6NEC7wMQZBEXhWHgzs0wbWJgF5HHWOjX+0rk\nP7DJjuL3Nua5e6cmKSfnrCiPmZjCb5NsShPUBIehSUQ6udCjQYOc2wPC1NiJo07s\nmOTJNbRhEEmqUGem1Hi7IrmedC3lmEZesSuJ42R9QhEe1in0T8S3YxnGQ6Q/0NWp\njHuk/DWet1xuU1a3OGxPk4QU4S+gKH5cBkObL6Ys5s6ueLv4yfuF/KKteBLtsCc=\n=UKnG\n-----END PGP SIGNATURE-----\n",
			"payload": "tree a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68\nparent 50661502e5fb14182d3e94d564bfde70b0fe061b\nparent 7300a44587b7c2394535a3f2bfeb594627b6ac91\nauthor Ilija Radosavovic <ilija.radosavovic@gmail.com> 1529593373 -0600\ncommitter GitHub <noreply@github.com> 1529593373 -0600\n\nMerge pull request #23 from ImmortalTurtle/fix-links-in-install-md\n\nFix links annotations in INSTALL.md"
		}
	},
	"url": "https://api.github.com/repos/facebookresearch/DensePose/commits/39c44e49e5b90018e206203bb857ef9a9c2ce8f6",
	"html_url": "https://github.com/facebookresearch/DensePose/commit/39c44e49e5b90018e206203bb857ef9a9c2ce8f6",
	"comments_url": "https://api.github.com/repos/facebookresearch/DensePose/commits/39c44e49e5b90018e206203bb857ef9a9c2ce8f6/comments",
	"author": {
		"login": "ir413",
		"id": 6415258,
		"node_id": "MDQ6VXNlcjY0MTUyNTg=",
		"avatar_url": "https://avatars2.githubusercontent.com/u/6415258?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/ir413",
		"html_url": "https://github.com/ir413",
		"followers_url": "https://api.github.com/users/ir413/followers",
		"following_url": "https://api.github.com/users/ir413/following{/other_user}",
		"gists_url": "https://api.github.com/users/ir413/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/ir413/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/ir413/subscriptions",
		"organizations_url": "https://api.github.com/users/ir413/orgs",
		"repos_url": "https://api.github.com/users/ir413/repos",
		"events_url": "https://api.github.com/users/ir413/events{/privacy}",
		"received_events_url": "https://api.github.com/users/ir413/received_events",
		"type": "User",
		"site_admin": false
	},
	"committer": {
		"login": "web-flow",
		"id": 19864447,
		"node_id": "MDQ6VXNlcjE5ODY0NDQ3",
		"avatar_url": "https://avatars3.githubusercontent.com/u/19864447?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/web-flow",
		"html_url": "https://github.com/web-flow",
		"followers_url": "https://api.github.com/users/web-flow/followers",
		"following_url": "https://api.github.com/users/web-flow/following{/other_user}",
		"gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
		"organizations_url": "https://api.github.com/users/web-flow/orgs",
		"repos_url": "https://api.github.com/users/web-flow/repos",
		"events_url": "https://api.github.com/users/web-flow/events{/privacy}",
		"received_events_url": "https://api.github.com/users/web-flow/received_events",
		"type": "User",
		"site_admin": false
	},
	"parents": [
		{
			"sha": "50661502e5fb14182d3e94d564bfde70b0fe061b",
			"url": "https://api.github.com/repos/facebookresearch/DensePose/commits/50661502e5fb14182d3e94d564bfde70b0fe061b",
			"html_url": "https://github.com/facebookresearch/DensePose/commit/50661502e5fb14182d3e94d564bfde70b0fe061b"
		},
		{
			"sha": "7300a44587b7c2394535a3f2bfeb594627b6ac91",
			"url": "https://api.github.com/repos/facebookresearch/DensePose/commits/7300a44587b7c2394535a3f2bfeb594627b6ac91",
			"html_url": "https://github.com/facebookresearch/DensePose/commit/7300a44587b7c2394535a3f2bfeb594627b6ac91"
		}
	]
}`

var commitJson2 = `
{
	"sha": "7300a44587b7c2394535a3f2bfeb594627b6ac91",
	"node_id": "MDY6Q29tbWl0MTM2MDY5MzExOjczMDBhNDQ1ODdiN2MyMzk0NTM1YTNmMmJmZWI1OTQ2MjdiNmFjOTE=",
	"commit": {
		"author": {
			"name": "Egor Nemchinov",
			"email": "immortalturtle@users.noreply.github.com",
			"date": "2018-06-21T14:57:20Z"
		},
		"committer": {
			"name": "GitHub",
			"email": "noreply@github.com",
			"date": "2018-06-21T14:57:20Z"
		},
		"message": "Fix typo in INSTALL.md\n\n'densepose_coco_2014_valminusminus.json' -> 'densepose_coco_2014_valminusminival.json'",
		"tree": {
			"sha": "a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68",
			"url": "https://api.github.com/repos/facebookresearch/DensePose/git/trees/a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68"
		},
		"url": "https://api.github.com/repos/facebookresearch/DensePose/git/commits/7300a44587b7c2394535a3f2bfeb594627b6ac91",
		"comment_count": 0,
		"verification": {
			"verified": true,
			"reason": "valid",
			"signature": "-----BEGIN PGP SIGNATURE-----\n\nwsBcBAABCAAQBQJbK7zQCRBK7hj4Ov3rIwAAdHIIAAuDad/FKbc78iDP9H+mkpmy\nJWCJO0hFyyH95fxm8lg3wITOnGyl7dsdqZGJsuuoSyaUUsdkbHj+H2okuoXrcYgQ\nk2g9ZR9cDNBLVXMPkEVbYDIE1xjoG9RO+cMYCkNQqvvz3yopU6VkYNeNt25DVZvo\nfXmZiSbNLS4BD+P4qLd2RZZ/ECiY67oK+9FHMHDBtZRfwOwl5+MsI6Tyf/5nR4xr\nkhbeG/NIC7GQiLbCefH5o4nv/XcHfiUKPRdVDSXlhJYZD9SS/VWmoG0Ty9wb4eOn\ns+wZERS7Idc14VdB4tgctQu9azZG2M47dLd/71g4193/iKRi4kr4SCXadGZ1EdQ=\n=y+wV\n-----END PGP SIGNATURE-----\n",
			"payload": "tree a25ad178e7e9a50fb87ebc13ca5d134c4d5fae68\nparent 280617f8a4df7b067b050fc4d27b1b1b6c3f1f1a\nauthor Egor Nemchinov <ImmortalTurtle@users.noreply.github.com> 1529593040 +0300\ncommitter GitHub <noreply@github.com> 1529593040 +0300\n\nFix typo in INSTALL.md\n\n'densepose_coco_2014_valminusminus.json' -> 'densepose_coco_2014_valminusminival.json'"
		}
	},
	"url": "https://api.github.com/repos/facebookresearch/DensePose/commits/7300a44587b7c2394535a3f2bfeb594627b6ac91",
	"html_url": "https://github.com/facebookresearch/DensePose/commit/7300a44587b7c2394535a3f2bfeb594627b6ac91",
	"comments_url": "https://api.github.com/repos/facebookresearch/DensePose/commits/7300a44587b7c2394535a3f2bfeb594627b6ac91/comments",
	"author": {
		"login": "ImmortalTurtle",
		"id": 22173703,
		"node_id": "MDQ6VXNlcjIyMTczNzAz",
		"avatar_url": "https://avatars0.githubusercontent.com/u/22173703?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/ImmortalTurtle",
		"html_url": "https://github.com/ImmortalTurtle",
		"followers_url": "https://api.github.com/users/ImmortalTurtle/followers",
		"following_url": "https://api.github.com/users/ImmortalTurtle/following{/other_user}",
		"gists_url": "https://api.github.com/users/ImmortalTurtle/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/ImmortalTurtle/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/ImmortalTurtle/subscriptions",
		"organizations_url": "https://api.github.com/users/ImmortalTurtle/orgs",
		"repos_url": "https://api.github.com/users/ImmortalTurtle/repos",
		"events_url": "https://api.github.com/users/ImmortalTurtle/events{/privacy}",
		"received_events_url": "https://api.github.com/users/ImmortalTurtle/received_events",
		"type": "User",
		"site_admin": false
	},
	"committer": {
		"login": "web-flow",
		"id": 19864447,
		"node_id": "MDQ6VXNlcjE5ODY0NDQ3",
		"avatar_url": "https://avatars3.githubusercontent.com/u/19864447?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/web-flow",
		"html_url": "https://github.com/web-flow",
		"followers_url": "https://api.github.com/users/web-flow/followers",
		"following_url": "https://api.github.com/users/web-flow/following{/other_user}",
		"gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
		"organizations_url": "https://api.github.com/users/web-flow/orgs",
		"repos_url": "https://api.github.com/users/web-flow/repos",
		"events_url": "https://api.github.com/users/web-flow/events{/privacy}",
		"received_events_url": "https://api.github.com/users/web-flow/received_events",
		"type": "User",
		"site_admin": false
	},
	"parents": [
		{
			"sha": "280617f8a4df7b067b050fc4d27b1b1b6c3f1f1a",
			"url": "https://api.github.com/repos/facebookresearch/DensePose/commits/280617f8a4df7b067b050fc4d27b1b1b6c3f1f1a",
			"html_url": "https://github.com/facebookresearch/DensePose/commit/280617f8a4df7b067b050fc4d27b1b1b6c3f1f1a"
		}
	]
}`
