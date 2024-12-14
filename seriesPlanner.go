package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/eSkiSo/seriesPlanner/db"
	"github.com/eSkiSo/seriesPlanner/moviedb"
	"github.com/eSkiSo/seriesPlanner/notion"
)

const (
	Version     = 0.2
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	italicStart = "\033[3m"
	boldStart   = "\033[1m"
	boldEnd     = "\033[0m"
)

type ShowInfo struct {
	Name              string
	Number_of_seasons int
	Poster_path       string
}

type Episode struct {
	Air_date       string
	Episode_number int
}

type SeasonInfo struct {
	Air_date string
	Episodes []Episode
}

var _config_moviedb_key = ""
var _config_notion_key = ""
var _config_notion_db = ""
var are_settings_defined = false

const image_url = "https://image.tmdb.org/t/p/w600_and_h900_bestv2"

func main() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	user_dir := filepath.FromSlash(dirname + "/")
	//check if local db is already created, if not, create it
	if _, err := os.Stat(user_dir + ".seriesPlanner.local.db"); err != nil {
		db.CreateDefaults()
	}

	settings := db.GetSettings()
	_config_moviedb_key = settings.Moviedb_key
	_config_notion_key = settings.Notion_key
	_config_notion_db = settings.Notion_database

	if _config_moviedb_key != "" && _config_notion_key != "" && _config_notion_db != "" {
		are_settings_defined = true
	}

	clearScreen()
	printTitle("")
	printMenu()
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func showMenu(msg string) {
	clearScreen()
	printTitle(msg)
	printMenu()
}

func printTitle(msg string) {
	fmt.Printf(boldStart+colorCyan+"Series Planner (Notion import) "+colorYellow+"Version %v \n\n"+boldEnd, Version)
	if msg != "" {
		fmt.Println(colorRed + " " + msg + "\n")
	}
}

func addShow() {
	clearScreen()
	printTitle("")
	var n int
	var input string
	for {
		fmt.Print("Show ID (from MovieDB) [x to cancel]: ")
		fmt.Scanln(&input)
		var err error
		if input == "x" {
			showMenu("")
		}
		n, err = strconv.Atoi(input)
		if err == nil && (n > 0) {
			break
		}
		showMenu("Invalid choice.")
	}
	show_name, current_season, _, _ := moviedb.GetInfo(input, true, _config_moviedb_key)
	if show_name == "" {
		fmt.Println(colorRed + "Show not found, invalid ID" + boldEnd)
	} else {
		result := db.AddShow(n, show_name, current_season)
		fmt.Println(colorGreen + result + boldEnd + "\n")
	}
	printGoBack()
}

func defineMovieDbKey() {
	clearScreen()
	printTitle("")
	var input string
	for {
		fmt.Print("MovieDB DB API Key [x to cancel]: ")
		fmt.Scanln(&input)
		if input == "x" {
			showMenu("")
		}
		if len(input) > 30 {
			break
		}
		showMenu("Invalid key.")
	}
	setSettings("moviedb_key", input)
	_config_moviedb_key = input
	fmt.Println(colorGreen + "New MovieDB key defined!")
	printGoBack()
}

func defineNotionKey() {
	clearScreen()
	printTitle("")
	var input string
	for {
		fmt.Print("Notion API Key [x to cancel]: ")
		fmt.Scanln(&input)
		if input == "x" {
			showMenu("")
		}
		if len(input) > 30 {
			break
		}
		showMenu("Invalid key.")
	}
	setSettings("notion_key", input)
	_config_notion_key = input
	fmt.Println(colorGreen + "New Notion key defined!")
	printGoBack()
}

func defineNotionDatabase() {
	clearScreen()
	printTitle("")
	var input string
	for {
		fmt.Print("Notion Dabase Key [x to cancel]: ")
		fmt.Scanln(&input)
		if input == "x" {
			showMenu("")
		}
		if len(input) > 30 {
			break
		}
		showMenu("Invalid key.")
	}
	setSettings("notion_database", input)
	_config_notion_db = input
	fmt.Println(colorGreen + "New Notion database defined!")
	printGoBack()
}

func listShows() {
	clearScreen()
	printTitle("")
	db.ListShows()
	printGoBack()
}

func editShow() {
	clearScreen()
	printTitle("")
	var n int
	var show_id int
	var input string
	for {
		fmt.Print("What show do you want to edit (Show ID) [x to cancel]: ")
		fmt.Scanln(&input)
		if input == "x" {
			showMenu("")
		}
		var err error
		n, err = strconv.Atoi(input)
		if err == nil && (n > 0) {
			break
		}
		showMenu("Invalid choice.")
	}
	show_id = n
	show := db.GetShow(show_id)
	if show.Show_name == "" {
		showMenu("Show id " + strconv.Itoa(show_id) + " not found!")
	}

	fmt.Printf("\nShow: %s \n1) Season: %s \n2) Episode: %d\n", show.Show_name, show.Season, show.Episode)

	for {
		fmt.Print("What do you want to edit [x to cancel]: ")
		fmt.Scanln(&input)
		if input == "x" {
			showMenu("")
		}
		var err error
		n, err = strconv.Atoi(input)
		if err == nil && (n > 0 && n <= 2) {
			break
		}
		showMenu("Invalid choice.")
	}
	switch n {
	case 1:
		//edit season
		for {
			fmt.Print("What season should it be on [x to cancel]: ")
			fmt.Scanln(&input)
			if input == "x" {
				showMenu("")
			}
			var err error
			n, err = strconv.Atoi(input)
			if err == nil && (n > 0) {
				break
			}
			showMenu("Invalid choice, only numbers > 0 accepted.")
		}
		response := db.UpdateShow(show_id, "Season", input)
		if response {
			showMenu(colorGreen + "Show updated with success.")
		} else {
			showMenu(colorRed + "There was an error updating show.")
		}
	case 2:
		//edit episode
		for {
			fmt.Print("What episode should it be on [x to cancel]: ")
			fmt.Scanln(&input)
			if input == "x" {
				showMenu("")
			}
			var err error
			n, err = strconv.Atoi(input)
			if err == nil && (n > -1) {
				break
			}
			showMenu("Invalid choice, only numbers > 0 accepted.")
		}
		response := db.UpdateShow(show_id, "Episode", input)
		if response {
			showMenu(colorGreen + "Show updated with success.")
		} else {
			showMenu(colorRed + "There was an error updating show.")
		}
	}

	printGoBack()
}

func removeShow() {
	clearScreen()
	printTitle("")
	var n int
	var input string
	for {
		fmt.Print("Show ID (to delete) [x to cancel]: ")
		fmt.Scanln(&input)
		var err error
		if input == "x" {
			showMenu("")
		}
		n, err = strconv.Atoi(input)
		if err == nil && (n > 0) {
			break
		}
		showMenu("Invalid choice.")
	}
	result := db.DeleteShow(n)
	if result > 0 {
		fmt.Println(colorGreen + "Show removed" + boldEnd)
	} else {
		fmt.Println(colorRed + "Show not found" + boldEnd)
	}
	printGoBack()
}

func printGoBack() {
	var input string
	for {
		fmt.Print(colorYellow + "Press any key to go back")
		fmt.Scanln(&input)
		showMenu("")
	}
}

func printMenu() {
	fmt.Println(colorYellow + "1) " + colorCyan + "Run seriesPlanner")
	fmt.Println("--")
	fmt.Println(colorYellow + "2) " + colorCyan + "List Shows")
	fmt.Println(colorYellow + "3) " + colorCyan + "Add Show")
	fmt.Println(colorYellow + "4) " + colorCyan + "Edit Show")
	fmt.Println(colorYellow + "5) " + colorCyan + "Remove Show")
	fmt.Println("--")
	fmt.Println(colorYellow + "6) " + colorCyan + "Check Settings")
	fmt.Println(colorYellow + "7) " + colorCyan + "Add Notion Database Key")
	fmt.Println(colorYellow + "8) " + colorCyan + "Add Notion API Key")
	fmt.Println(colorYellow + "9) " + colorCyan + "Add Movie DB API Key")
	fmt.Println("--")
	fmt.Println(colorYellow + "q) " + colorRed + "Quit \n" + boldEnd)
	var n int
	for {
		fmt.Print("Your choice: ")
		var input string
		fmt.Scanln(&input)
		if input == "q" {
			os.Exit(0)
		}
		var err error
		n, err = strconv.Atoi(input)
		if err == nil && (0 <= n && n <= 9) {
			break
		}
		showMenu("Invalid choice.")

	}
	switch n {
	case 1:
		if are_settings_defined {
			run()
		} else {
			fmt.Println(colorRed + "Settings are not yet defined, please set them first.")
		}

	case 2:
		listShows()
	case 3:
		if are_settings_defined {
			addShow()
		} else {
			fmt.Println(colorRed + "Settings are not yet defined, please set them first.")
		}

	case 4:
		if are_settings_defined {
			editShow()
		} else {
			fmt.Println(colorRed + "Settings are not yet defined, please set them first.")
		}

	case 5:
		removeShow()
	case 6:
		checkSettings()
	case 7:
		defineNotionDatabase()
	case 8:
		defineNotionKey()
	case 9:
		defineMovieDbKey()
	case 0:
		os.Exit(0)
	}
}

func run() {
	//Get episodes saved
	shows := db.GetShows()
	//Loop: Get Show and check episode dates
	for _, show := range shows {
		fmt.Println(colorBlue + "Checking " + show.Show_name + "..." + boldEnd)
		_, _, _, result := moviedb.GetInfo(strconv.Itoa(show.Show_id), false, _config_moviedb_key)

		update_show := 0
		current_s := 0
		for _, ep := range result {
			cuSE, _ := strconv.Atoi(fmt.Sprintf("%d%d", ep.Season_number, ep.Episode_number))
			savedSE, _ := strconv.Atoi(fmt.Sprintf("%s%d", show.Season, show.Episode))

			if cuSE > savedSE {
				fmt.Println("Add show on date " + ep.Air_date)
				//Add to notion
				res := notion.Add(show.Show_name, strconv.Itoa(ep.Season_number), strconv.Itoa(ep.Episode_number), ep.Air_date, _config_notion_db, _config_notion_key)
				if res {
					fmt.Println(colorGreen + "Show added" + boldEnd)
					update_show = ep.Episode_number
					current_s = ep.Season_number

				} else {
					fmt.Println(colorRed + "Show not added" + boldEnd)
				}
			}
		}
		if update_show > 0 {
			// Update local db
			update_response := db.UpdateShowSE(show.Show_id, strconv.Itoa(current_s), strconv.Itoa(update_show))
			if update_response {
				fmt.Println(colorGreen + "Show updated with success in SeriesPlanner." + boldEnd)
			} else {
				fmt.Println(colorRed + "There was an error updating show in SeriesPlanner." + boldEnd)
			}
		}
	}

	printGoBack()
}

func checkSettings() {
	clearScreen()
	printTitle("")
	fmt.Println(colorYellow + "Checking configs..." + boldEnd)

	mk := _config_moviedb_key
	nk := _config_notion_key
	nd := _config_notion_db

	if mk != "" {
		mk = mk[len(mk)-5:] + ": " + colorGreen + "OK" + boldEnd
	} else {
		mk = "     : " + colorRed + "NOK" + boldEnd
	}
	if nk != "" {
		nk = nk[len(nk)-5:] + ": " + colorGreen + "OK" + boldEnd
	} else {
		nk = "     : " + colorRed + "NOK" + boldEnd
	}
	if nd != "" {
		nd = nd[len(nd)-5:] + ": " + colorGreen + "OK" + boldEnd
	} else {
		nd = "     : " + colorRed + "NOK" + boldEnd
	}
	fmt.Printf("MovieDB Key: ...%12s \nNotion Key : ...%12s \nNotion DB  : ...%12s \n\n", mk, nk, nd)
	printGoBack()
}

func setSettings(key string, value string) bool {
	return db.UpdateSettings(key, value)
}
