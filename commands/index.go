package commands

import (
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"ms/movielite"
	"ms/movielite/models"
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

	conf := movielite.GetConfig(ctx.GlobalString("config"))
	db := models.ConnectDataBase(conf.DataBase)
	tx := db.Begin()
	defer tx.Close()

	rows, err := tx.Model(&models.Movie{}).Rows()
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var movie models.Movie
		var new models.Movie
		db.ScanRows(rows, &movie)
		if err := tx.Set("gorm:auto_preload", true).Where("id = ?", movie.ID).First(&new).Error; err != nil {
			log.Error(err)
		}
		log.Debugf("indexing %s", movie.Title)
		if err := new.FullTextIndex(tx); err != nil {
			log.Warn(err)
		}
	}
	tx.Commit()
	return nil
}
