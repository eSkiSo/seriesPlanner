# seriesPlanner
_GO project_

![Main Menu](/screenshot.png)

## Adds series from MovieDB to Notion

### TODO

- [x] Run sP (get list of shows, update Notion with new episodes)

- [x] List Shows

- [x] Add Shows

- [x] Edit Show (last season and episode tracked)

- [x] Remove Show

- [x] Manage Settings (list of api keys and notion database id)

- [ ] Make it cross-season (if you set season 1 and its on season 3 it should get all episodes from S1 to S3)


### Requirements

* TheMovieDB Account ( https://www.themoviedb.org )
* TheMovieDB API key ( https://www.themoviedb.org/settings/api )
* Notion Account ( https://www.notion.so )
* Notion API Key ( https://www.notion.so/my-integrations )


### Compile code

```
GOOS=windows GOARCH=amd64 go build -o bin/windows/seriesPlanner.exe seriesPlanner.go
GOOS=linux GOARCH=amd64 go build -o bin/linux/seriesPlanner seriesPlanner.go
GOOS=darwin GOARCH=amd64 go build -o bin/macos/seriesPlanner seriesPlanner.go
```