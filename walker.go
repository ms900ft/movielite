package movielight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"ms/movielight/models"
)

type Walker struct {
	Config *Config
}

type payload struct {
	Path string `json:"FullPath"`
}

func (w *Walker) Run() error {
	log.Debug("Starting walker")
	var err error
	directories := viper.GetStringSlice("Directories")
	for _, dir := range directories {
		log.Debugf("Scanning directory: %s\n", dir)
	}
	watchdir := viper.GetString("WatchDir")
	log.Debug("watching " + watchdir)
	//flag.Parse()
	//root := flag.Arg(0)
	if len(directories) > 0 {
		for _, dir := range directories {
			dirname := dir + string(filepath.Separator)
			if _, err = os.Stat(dirname); err != nil {
				log.Fatalf("directory %s not found", dirname)
			}
			fileList := []string{}
			err = filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
				path = toUtf8Nfc(path)
				fileList = append(fileList, path)
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			//spew.Dump(fileList)
			for _, file := range fileList {
				err = sendfile(file)
				if err != nil {
					log.Errorf("sending file %s", err)
				}
			}
		}
	}
	//begin to watch
	if watchdir != "" {
		err = watchdirectory(watchdir)
		if err != nil {
			log.Errorf("watch dir %s", err)
		}
	}
	return err
}

func sendfile(file string) error {
	//log.Debugf("got files: %s", file)
	extens := map[string]bool{".mov": true, ".avi": true, ".mp4": true, "mkv": true, "m4v": true}
	extName := path.Ext(file)
	fName := path.Base(file)
	if strings.HasPrefix(fName, ".") {
		log.Debugf("not sending dotname %s", fName)
		return nil
	}
	if extens[extName] {
		log.Debugf("file matches %s", file)
		//change permission
		//mediathekview set the x bit
		err := os.Chmod(file, 0644)
		if err != nil {
			log.Errorf("sending file %s", err)
		}

		log.Debugf("search for %s", fName)

		//fmt.Println(extName)
		files := searchFile(fName)
		if len(files) > 0 {
			log.Debugf("found: %s\n", fName)
			fileExits := map[string]bool{}
			var toUpdate models.File
			for _, rfile := range files {
				if _, err := os.Stat(rfile.FullPath); err == nil {
					fileExits[rfile.FullPath] = true
				} else {
					fileExits[rfile.FullPath] = false
					toUpdate = rfile
				}
			}
			if !fileExits[file] {
				log.Debugf("update path to: %s", file)
				if toUpdate.ID > 0 {
					toUpdate.FullPath = file
					//spew.Dump(toUpdate)
					err = update(toUpdate)
					if err != nil {
						log.Error(err)
						return err
					}
				}
			} else {
				log.Debugf("path not changed: %s\n", file)
			}
		} else {
			//spew.Dump(files)
			err := send(file)
			if err != nil {
				log.Errorf("sending file %s", err)
			}
		}
	}
	return nil
}

func searchFile(name string) []models.File {
	surl := viper.GetString("MovieServerUrl")
	name, err := URLEncoded(name)
	if err != nil {
		log.Errorf("search file %s", err)
	}
	url := fmt.Sprintf("%s/files?f=%s", surl, name)
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
	bodyString := string(body)
	log.Debug(bodyString)
	files := []models.File{}
	err = json.Unmarshal(body, &files)
	if err != nil {
		log.Error(err)
	}
	return files
}

func send(path string) error {
	surl := viper.GetString("MovieServerUrl")
	pl := payload{Path: path}
	url := surl + "/api/file"
	jsonValue, _ := json.Marshal(pl)
	//pl := []byte(`{"FullPath": "xxxxxxx"}`)
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
	log.Debug("response Status:", response.Status)
	body, _ := ioutil.ReadAll(response.Body)
	log.Debug("response Body:", string(body))
	return err
}

func update(file models.File) error {
	surl := viper.GetString("MovieServerUrl")
	url := fmt.Sprintf("%s/file/%d", surl, file.ID)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(file)
	if err != nil {
		log.Error(err)
	}
	req, err := http.NewRequest("PUT", url, b)
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
	log.Debug("response Status:", response.Status)
	body, _ := ioutil.ReadAll(response.Body)
	log.Debug("response Body:", string(body))
	return err
}

func watchdirectory(dir string) error {
	log.Debugf("watching directory %s", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					path := toUtf8Nfc(event.Name)
					log.Debugf("sending file %s", path)
					err = sendfile(path)
					if err != nil {
						log.Errorf("sending file %s", err)
					}
				}
				// if event.Op&fsnotify.Write == fsnotify.Write {
				// 	//log.Println("modified file:", event.Name)

				// }
			case err = <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Error(err)
	}
	<-done
	return err
}
