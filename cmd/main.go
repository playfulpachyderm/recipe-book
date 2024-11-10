package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	pkg_db "recipe_book/pkg/db"
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
