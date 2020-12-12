package commands

import (
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"ms/movielight"
	"ms/movielight/models"
)

// StartCommand is used to register the start cli command
var IndexCommand = cli.Command{
	Name:  "index",
	Usage: "reindex fulltext",
	//Flags:   startFlags,
	Action: indexAction,
}

var indexFlags = []cli.Flag{
	// cli.BoolFlag{
	// 	Name:  "detach-server, d",
	// 	Usage: "detach from the console (daemon mode)",
	//	},
}

// startAction start the web server and initializes the daemon
func indexAction(ctx *cli.Context) error {

	conf := movielight.GetConfig()
	db := models.ConnectDataBase(conf.DataBase)
	tx := db.Begin()
	defer tx.Close()

	//var movies []models.Movie
	// var movie models.Movie

	// if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).Error; err != nil {
	// 	return nil
	// }
	//db.Set("gorm:auto_preload", true).Model(&models.Movie{}).Find(&movies)
	// defer rows.Close()
	// if err != nil {
	// 	log.Error(err)
	//	}

	rows, err := tx.Model(&models.Movie{}).Rows()
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		//for rows.Next() {
		var movie models.Movie
		var new models.Movie
		//ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		db.ScanRows(rows, &movie)
		if err := tx.Set("gorm:auto_preload", true).Where("id = ?", movie.ID).First(&new).Error; err != nil {
			log.Error(err)
		}
		log.Debugf("indexing %s", movie.Title)
		if err := new.FullTextIndex(tx); err != nil {
			log.Warn(err)
		}
		//	tx.Commit()
		//spew.Dump(movie)
	}
	tx.Commit()
	return nil
	// do something
}
