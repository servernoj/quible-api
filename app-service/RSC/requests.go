package RSC

import (
	"fmt"
	"net/url"
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

func GetTeamInfo[T TeamInfoItem](query url.Values) ([]T, error) {
	client := NewClient[T]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	url := fmt.Sprintf("%s/team-info/%s?%s", client.URL, client.Sport, client.Query.Encode())
	return client.RequestRunner(url)
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
