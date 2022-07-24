package moviedb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//const image_url = "https://image.tmdb.org/t/p/w600_and_h900_bestv2"
const url = "https://api.themoviedb.org/3/tv/"

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

func GetInfo(id string, ignore_episodes bool, api_key string) (show_name string, current_season string, result string, Episodes []Episode) {
	api_url_info := url + id + "?api_key=" + api_key
	resp, err := http.Get(api_url_info)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var info ShowInfo
	var ListEpisodes []Episode

	json.Unmarshal([]byte(body), &info)

	show_name = info.Name
	current_season = strconv.Itoa(info.Number_of_seasons)
	//poster := image_url + info.Poster_path

	if ignore_episodes {
		return show_name, current_season, "", ListEpisodes
	}

	fmt.Printf("Show: %s Season: %s \n", show_name, current_season)

	api_url := url + id + "/season/" + current_season + "?api_key=" + api_key

	reply, err := http.Get(api_url)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	seasons_body, err := ioutil.ReadAll(reply.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var sInfo SeasonInfo

	json.Unmarshal([]byte(seasons_body), &sInfo)

	for _, v := range sInfo.Episodes {
		if v.Air_date != "" {
			ListEpisodes = append(ListEpisodes, v)
		}
	}

	//Convert the body to type string
	sb := string(body)
	return show_name, current_season, sb, ListEpisodes
}
