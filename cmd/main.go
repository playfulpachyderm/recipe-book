package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	pkg_db "recipe_book/pkg/db"
	"recipe_book/pkg/web"
)

const DB_FILENAME = "food.db"

var db_path string = ""

func main() {
	flag.StringVar(&db_path, "db", "sample_data/data", "database path")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("subcommand needed\n")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		init_db()
	case "webserver":
		fs := flag.NewFlagSet("", flag.ExitOnError)
		should_auto_open := fs.Bool("auto-open", false, "")
		addr := fs.String("addr", "localhost:3080", "port to listen on") // Random port that's probably not in use

		if err := fs.Parse(args[1:]); err != nil {
			panic(err)
		}
		start_webserver(*addr, *should_auto_open)
	default:
		fmt.Printf(COLOR_RED+"invalid subcommand: %q\n"+COLOR_RESET, args[0])
		os.Exit(1)
	}
}

func init_db() {
	db_filename := filepath.Join(db_path, DB_FILENAME)
	_, err := pkg_db.DBCreate(db_filename)
	if err != nil {
		fmt.Println(COLOR_RED + err.Error() + COLOR_RESET)
		os.Exit(1)
	}
	fmt.Println(COLOR_GREEN + "Successfully created the db" + COLOR_RESET)
}

const (
	COLOR_RESET  = "\033[0m"
	COLOR_BLACK  = "\033[30m"
	COLOR_RED    = "\033[31m"
	COLOR_GREEN  = "\033[32m"
	COLOR_YELLOW = "\033[33m"
	COLOR_BLUE   = "\033[34m"
	COLOR_PURPLE = "\033[35m"
	COLOR_CYAN   = "\033[36m"
	COLOR_GRAY   = "\033[37m"
	COLOR_WHITE  = "\033[97m"
)

func start_webserver(addr string, should_auto_open bool) {
	db, err := pkg_db.DBConnect(filepath.Join(db_path, DB_FILENAME))
	if err != nil {
		fmt.Println(COLOR_RED + "opening database: " + err.Error() + COLOR_RESET)
		os.Exit(1)
	}

	app := web.NewApp(db)
	app.Run(addr, should_auto_open)
}
