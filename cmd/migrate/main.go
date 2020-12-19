package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/models"

	log "github.com/sirupsen/logrus"
)

func main() {
	count := flag.Int("count", 0, "how many movies")
	offset := flag.Int("offset", 0, "begin with offset")
	flag.Parse()
	src := flag.Arg(0)
	target := flag.Arg(1)
	//spew.Dump(flag.Args())
	log.SetLevel(log.DebugLevel)
	log.Debugf("migrating from %s to %s: %d movies with offset %d", src, target, *count, *offset)
	i := 0
	for {
		log.Debug("getting next movie to migrate")
		i++
		url := fmt.Sprintf("%s/movie?&limit=1&offset=%d&show=all", src, *offset)
		*offset++
		log.Debugf("Getting: %s", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error(err)
		}
		req.Header.Add("Accept", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		//bodyString := string(body)
		//log.Debug(bodyString)
		movies := movielite.Movielist{}
		//log.Error(string(body))
		err = json.Unmarshal(body, &movies)
		if err != nil {
			log.Errorf("next movie %s", err)
			log.Error(string(body))
		}
		//spew.Dump(movies)
		if len(movies.Data) == 0 {
			log.Debug("no more movies to migrate")
			break
		}
		err = addMovie(movies.Data[0], target)
		if err != nil {
			log.Error(err)
		}
		//spew.Dump(movies.Data[0])
		if *count > 0 && i >= *count {
			log.Debugf("all %d movies migrated", i)
			break
		}
	}
}

func addMovie(m models.Movie, target string) error {
	url := fmt.Sprintf("%s/movie", target)
	jsonValue, _ := json.Marshal(m)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	//response := executeRequest(req)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer response.Body.Close()
	log.Debug("rescan: response Status:", response.Status)
	// body, _ := ioutil.ReadAll(response.Body)
	//	log.Debug("response Body:", string(body))
	return err
}
