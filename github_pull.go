package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/samsarahq/thunder/livesql"
	"github.com/samsarahq/thunder/sqlgen"
	"encoding/hex"
	//"database/sql"
)

type githubResponse struct {
	Commit map[string]map[string]interface{} `json:"commit"`
}

type githubResponseHash struct {
	Hash string `json:"sha"`
}

type githubResponseId struct {
	Id int64 `json:"id"`
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func lol(t time.Time) {
	url := fmt.Sprintf("https://api.github.com/repos/facebookresearch/DensePose/commits?access_token=%v", "07daf881d73d37b6c72f024d4bb2f749b603f0d1")
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	bodyString := string(body)
	bodyString = strings.Replace(bodyString, ",\"comment_count\":0", "", -1)
	body = []byte(bodyString)

	response := make([]githubResponse, 0)

	responseHash := make([]githubResponseHash, 0)

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err)
	}

	if err := json.Unmarshal(body, &responseHash); err != nil {
		log.Println(err)
	}

	url = fmt.Sprintf("https://api.github.com/repos/facebookresearch/DensePose")

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client = &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	body, _ = ioutil.ReadAll(resp.Body)

	var nice githubResponseId

	if err := json.Unmarshal(body, &nice); err != nil {
		log.Println(err)
	}

	sqlgenSchema := sqlgen.NewSchema()

	db, err := livesql.Open("localhost", 3307, "root", "", "github", sqlgenSchema)
	if err != nil {
		panic(err)
	}

	repoSql := fmt.Sprintf("INSERT INTO repos (id) VALUES (%v)", nice.Id)
	fmt.Println(repoSql)
	if _, err := db.QueryExecer(context.TODO()).Exec(repoSql); err != nil {
		log.Println("found duplicate repo sql")
	}
  

	count := 0
	for _, resp := range response {
		layout := "2006-01-02T15:04:05.000Z"
		date := resp.Commit["author"]["date"]
		if str, ok := date.(string); ok {
			datee := strings.Replace(str, "Z", ".000Z", -1)
			t, _ := time.Parse(layout, datee)
			millis := t.UnixNano() / 1000000
			jsonapi, _ := json.Marshal(resp)
			jsonapistr := hex.EncodeToString(jsonapi) 
			eventSql := fmt.Sprintf("INSERT INTO events (at_ms, repo_id, event_id, api_json) VALUES (%v, %v, '%s', 0x%s)", millis, nice.Id, responseHash[count].Hash, jsonapistr)
			fmt.Println(eventSql)
			if _, err := db.QueryExecer(context.TODO()).Exec(eventSql); err != nil {
				log.Println("found duplicate event sql")
				continue;
			}
		}
		count++
	}
}

func main() {
	
		doEvery(10000*time.Millisecond, lol)

}
