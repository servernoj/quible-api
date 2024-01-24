package RSC

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

// 1. Data coming from RSC API doesn't contain "team color" which has to be "injected" into response
// of the RSC API for every team. Also, the data for `arena` returned by RSC API is outdated and
// has to be re-defined from local DB (table `teams`)
// 2. Logo info is found in `images` table (result of crawling ESPN dataset)
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
	// -- data from `images` table indexed by DisplayName
	images, err := models.Images(models.ImageWhere.ParentID.IsNull()).AllG(ctx)
	if err != nil {
		log.Printf("unable to retrieve extra teams info from `images`: %v", err)
		return []TeamInfoItemExtended{}, nil
	}
	imagesByDisplayName := make(map[string]*models.Image, len(images))
	for _, image := range images {
		imagesByDisplayName[image.DisplayName] = image
	}
	// -- data from `teams` table indexed by RSC ID
	teams, err := models.Teams().AllG(ctx)
	if err != nil {
		log.Printf("unable to retrieve extra teams info from `teams`: %v", err)
		return []TeamInfoItemExtended{}, nil
	}
	teamsByRSCID := make(map[int]*models.Team, len(teams))
	for _, team := range teams {
		teamsByRSCID[team.RSCID] = team
	}
	// -- augmenting RSC API response with data from `teams` and `images`
	for idx := range teamInfoItems {
		teamInfoItem := &teamInfoItems[idx]
		if team, ok := teamsByRSCID[teamInfoItem.TeamID]; ok {
			// table `teams` in local DB contains up to date values for `arena` and `color`
			// for every team record. Those are injected on-fly, every time when team info
			// is requested from RSC API
			teamInfoItem.Arena = team.Arena
			teamInfoItem.Color = team.Color
		}
		if image, ok := imagesByDisplayName[teamInfoItem.Mascot]; ok {
			teamInfoItem.Logo = &image.ImageURL
		} else {
			teamInfoItem.Logo = nil
			log.Printf("unable find logo for %q\n", teamInfoItem.Mascot)
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

func GetPlayerInfo(query url.Values) ([]PlayerInfoItemExtended, error) {
	client := NewClient[PlayerInfoItemExtended]()
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	url := fmt.Sprintf("%s/player-info/%s?%s", client.URL, client.Sport, client.Query.Encode())
	playerInfoItems, err := client.RequestRunner(url)
	if err != nil {
		return []PlayerInfoItemExtended{}, err
	}
	ctx := context.Background()
	// -- data from `images` table indexed by DisplayName
	where := models.ImageWhere.ParentID.IsNotNull()
	if query.Has("team_id") {
		rsc_id, err := strconv.Atoi(query.Get("team_id"))
		if err == nil {
			team, err := models.Teams(models.TeamWhere.RSCID.EQ(rsc_id)).OneG(ctx)
			if team != nil && err == nil {
				teamImage, err := models.Images(
					qm.Expr(
						models.ImageWhere.ParentID.IsNull(),
						models.ImageWhere.DisplayName.EQ(team.DisplayName.String),
					),
				).OneG(ctx)
				if teamImage != nil && err == nil {
					where = models.ImageWhere.ParentID.EQ(null.StringFrom(teamImage.ID))
				}
			}
		}
	}
	images, err := models.Images(where).AllG(ctx)
	if err != nil {
		log.Printf("unable to retrieve extra players info from `images`: %v", err)
		return []PlayerInfoItemExtended{}, nil
	}
	imagesByDisplayName := make(map[string]*models.Image, len(images))
	for _, image := range images {
		imagesByDisplayName[image.DisplayName] = image
	}
	// -- augmenting RSC API response with data from `images`
	for idx := range playerInfoItems {
		playerInfoItem := &playerInfoItems[idx]
		if image, ok := imagesByDisplayName[playerInfoItem.Player]; ok {
			playerInfoItem.Headshot = &image.ImageURL
		} else {
			playerInfoItem.Headshot = nil
		}
	}
	return playerInfoItems, nil
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
