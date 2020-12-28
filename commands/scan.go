package commands

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/models"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
)

var directory string

// StartCommand is used to register the start cli command
var ScanCommand = cli.Command{
	Name:   "scan",
	Usage:  "scan directory for movies",
	Flags:  scanFlags,
	Action: scanAction,
}

var scanFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "directory, d",
		Usage:       "scan directory for movies",
		Destination: &directory,
		Required:    true,
	},
}

// startAction start the web server and initializes the daemon
func scanAction(ctx *cli.Context) error {

	conf := movielite.GetConfig(ctx.GlobalString("config"))
	expiresAt := time.Now().Add(time.Minute * 1000000).Unix()
	tk := &models.Token{
		Name: "admin",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	if conf.Secret == "" {
		log.Fatal("no secret found")
	}
	tokenString, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		log.Fatal(err)
	}

	w := movielite.Walker{Config: conf}
	w.Token = tokenString
	err = w.Run(directory)
	if err != nil {
		log.Fatalf("can't scan for movies: %s", err)
	}
	return nil
}
