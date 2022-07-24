package notion

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const posturl = "https://api.notion.com/v1/pages"

func Add(show_name string, current_season string, episode_number string, air_date string, database_id string, notion_token string) bool {
	// JSON body
	comment := "Pendente"

	timeT, _ := time.Parse("2006-01-02", air_date)
	currentTime := time.Now()
	if timeT.Before(currentTime) {
		comment = "Disponivel"
	}

	body := []byte(`{
            "parent": {
                "database_id": "` + database_id + `"
            },
            "properties": {
                "Name": {
                    "title": [
                        {
                            "text": {
                                "content": "` + show_name + ` S` + current_season + `E` + episode_number + `"
                            }
                        }
                    ]
                },
                "Estreia": {
                    "date": {
                        "start": "` + air_date + `"
                    }
                },
                "Status": {
                    "select": {
                        "name": "` + comment + `"
                    }
                },
                "Type": {
                    "select": {
                        "name": "TV Series"
                    }
                }
            }
        }`)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Notion-Version", "2022-02-22")
	r.Header.Add("Authorization", "Bearer "+notion_token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		return true
	} else {
		fmt.Println(res)
		return false
	}
}
