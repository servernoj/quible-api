package BasketAPI

import (
	"context"
	"fmt"

	"github.com/quible-io/quible-api/lib/models"
)

type TeamId struct {
	ID uint `json:"id"`
}

type TeamInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	ShortName      string  `json:"shortName"`
	Abbr           string  `json:"abbr"`
	ArenaName      string  `json:"arenaName"`
	ArenaSize      int     `json:"arenaSize"`
	Color          string  `json:"color"`
	SecondaryColor string  `json:"secondaryColor"`
	Logo           *string `json:"logo"`
}

const (
	Host = "basketapi1.p.rapidapi.com"
)

func getTeamEnhancer(ctx context.Context) (func(TeamId) TeamInfo, error) {
	teamsInfo, err := models.TeamInfos().AllG(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve combined team info (for all teams): %w", err)
	}
	teamInfoById := make(map[int]*models.TeamInfo, len(teamsInfo))
	for _, teamInfo := range teamsInfo {
		teamInfoById[teamInfo.ID] = teamInfo
	}
	return func(team TeamId) TeamInfo {
		teamInfo := *teamInfoById[int(team.ID)]
		return TeamInfo{
			ID:             teamInfo.ID,
			Name:           teamInfo.Name,
			Slug:           teamInfo.Slug,
			ShortName:      teamInfo.ShortName,
			Abbr:           teamInfo.Abbr,
			ArenaName:      teamInfo.ArenaName,
			ArenaSize:      teamInfo.ArenaSize,
			Color:          teamInfo.Color,
			SecondaryColor: teamInfo.SecondaryColor,
			Logo:           teamInfo.Logo.Ptr(),
		}
	}, nil
}

func ApplyMapper[F any, T any](s []F, m func(F) T) []T {
	result := make([]T, len(s))
	for idx := range s {
		result[idx] = m(s[idx])
	}
	return result
}
