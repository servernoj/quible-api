package RSC

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/quible-io/quible-api/lib/models"
)

func GetScheduleSeason[T ScheduleItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	season := client.GetSeason()
	url := fmt.Sprintf("%s/schedule-season/%s/%s?%s", client.URL, season, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

func GetDailySchedule[T ScheduleItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	date := client.GetDate()
	url := fmt.Sprintf("%s/schedule/%s/%s?%s", client.URL, date, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

type TeamInfoItemExtended struct {
	TeamInfoItem
	Color string `json:"color"`
}

func GetTeamInfo(query url.Values) ([]TeamInfoItemExtended, error) {
	client := NewClient[TeamInfoItemExtended]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	url := fmt.Sprintf("%s/team-info/%s?%s", client.URL, client.Sport, client.Query.Encode())
	teamInfoItems, err := client.RequestRunner(url)
	if err != nil {
		return []TeamInfoItemExtended{}, err
	}
	ctx := context.Background()
	teams, err := models.Teams().AllG(ctx)
	if err != nil {
		log.Printf("unable to retrieve extra teams infor from DB: %v", err)
		return []TeamInfoItemExtended{}, nil
	}
	teamsByRSCID := make(map[int]*models.Team, len(teams))
	for _, team := range teams {
		teamsByRSCID[team.RSCID] = team
	}
	for idx := range teamInfoItems {
		teamInfoItem := &teamInfoItems[idx]
		if team, ok := teamsByRSCID[teamInfoItem.TeamID]; ok {
			teamInfoItem.Arena = team.Arena
			teamInfoItem.Color = team.Color
		}
	}

	return teamInfoItems, nil
}

func GetTeamStats[T TeamSeasonStatItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	season := client.GetSeason()
	url := fmt.Sprintf("%s/team-stats/%s/%s?%s", client.URL, season, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

func GetPlayerInfo[T PlayerInfoItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	url := fmt.Sprintf("%s/player-info/%s?%s", client.URL, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

func GetPlayerStats[T PlayerSeasonStatItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	season := client.GetSeason()
	url := fmt.Sprintf("%s/player-stats/%s/%s?%s", client.URL, season, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

func GetInjuries[T InjuryItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	url := fmt.Sprintf("%s/injuries/%s?%s", client.URL, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}

func GetLiveFeed[T LiveFeedItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	date := client.GetDate()
	url := fmt.Sprintf("%s/live/%s/%s?%s", client.URL, date, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
}
