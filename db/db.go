package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	Version     = 0.1
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	italicStart = "\033[3m"
	boldStart   = "\033[1m"
	boldEnd     = "\033[0m"
)

type Settings struct {
	gorm.Model
	Notion_database string
	Notion_key      string
	Moviedb_key     string
}

type Configurations struct {
	gorm.Model
	Type  string
	Key   string
	Value string
}

type LastEpisodes struct {
	gorm.Model
	Show_id   int
	Show_name string
	Season    string
	Episode   int
}

var user_dir = ""

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	user_dir = filepath.FromSlash(dirname + "/")
}

func AddShow(show_id int, show_name string, season string) string {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&LastEpisodes{})
	db.Create(&LastEpisodes{Show_id: show_id, Show_name: show_name, Season: season, Episode: 0})
	return "Show " + show_name + " Added"
}

func GetShow(show_id int) LastEpisodes {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&LastEpisodes{})
	var show LastEpisodes
	db.Find(&show, "show_id = ?", show_id)

	return show
}

func ListShows() bool {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&LastEpisodes{})

	rows, err := db.Model(&LastEpisodes{}).Rows()

	if err != nil {
		panic("failed to connect to database")
	}
	defer rows.Close()

	for rows.Next() {
		var ep LastEpisodes

		db.ScanRows(rows, &ep)

		fmt.Printf(colorCyan+"Show ID:"+boldEnd+" %d \n", ep.Show_id)
		fmt.Printf(colorCyan+"Show Name:"+boldEnd+" %s \n", ep.Show_name)
		fmt.Printf(colorCyan+"Current Season:"+boldEnd+" %s\n", ep.Season)
		fmt.Printf(colorCyan+"Last Imported Episode:"+boldEnd+" %d\n\n", ep.Episode)

	}
	return true
}

func GetShows() []LastEpisodes {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&LastEpisodes{})

	rows, err := db.Model(&LastEpisodes{}).Rows()

	if err != nil {
		panic("failed to connect to database")
	}
	defer rows.Close()

	var results []LastEpisodes

	for rows.Next() {
		var ep LastEpisodes

		db.ScanRows(rows, &ep)
		results = append(results, ep)
		/*fmt.Printf(colorCyan+"Show ID:"+boldEnd+" %d \n", ep.Show_id)
		fmt.Printf(colorCyan+"Show Name:"+boldEnd+" %s \n", ep.Show_name)
		fmt.Printf(colorCyan+"Current Season:"+boldEnd+" %s\n", ep.Season)
		fmt.Printf(colorCyan+"Last Imported Episode:"+boldEnd+" %d\n\n", ep.Episode)*/

	}
	return results
}

func CreateDefaults() {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Configurations{})
	db.AutoMigrate(&LastEpisodes{})

	// Create
	db.Create(&Configurations{Type: "settings", Key: "notion_database", Value: ""})
	db.Create(&Configurations{Type: "settings", Key: "notion_key", Value: ""})
	db.Create(&Configurations{Type: "settings", Key: "moviedb_key", Value: ""})

}

func DeleteShow(show_id int) int64 {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		return 0
	}
	var episodes LastEpisodes
	result := db.Where("show_id = ?", show_id).Delete(&episodes)
	return result.RowsAffected
}

func GetSettings() Settings {

	if _, err := os.Stat(user_dir + ".seriesPlanner.local.db"); err != nil {
		CreateDefaults()
	}

	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Read
	var settings Settings
	db.AutoMigrate(&Configurations{})
	rows, err := db.Model(&Configurations{}).Rows()
	for rows.Next() {
		var config Configurations

		db.ScanRows(rows, &config)

		switch config.Key {
		case "notion_database":
			settings.Notion_database = config.Value
		case "notion_key":
			settings.Notion_key = config.Value
		case "moviedb_key":
			settings.Moviedb_key = config.Value
		}
	}
	return settings
}

func UpdateSettings(key string, value string) bool {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var configs Configurations
	db.First(&configs, "type = ? AND key = ?", "settings", key)
	db.Model(&configs).Update("Value", value)
	return true
}

func UpdateShow(show_id int, key string, value string) bool {
	db, err := gorm.Open(sqlite.Open(user_dir+".seriesPlanner.local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var ep LastEpisodes
	db.First(&ep, "show_id = ?", show_id)
	db.Model(&ep).Update(key, value)
	return true
}
