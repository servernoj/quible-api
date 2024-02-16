package BasketAPI

import (
	"context"
	"fmt"

	"github.com/quible-io/quible-api/lib/models"
)

const (
	Host = "basketapi1.p.rapidapi.com"
)

func GetTeamEnhancer(ctx context.Context) (func(TeamId) TeamInfo, error) {
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
