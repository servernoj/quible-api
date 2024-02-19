package BasketAPI

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Options struct {
	TournamentID uint
	SeasonID     uint
}

func NewCrawler(options Options) *Crawler {
	host := "basketapi1.p.rapidapi.com"
	header := http.Header{
		"X-RapidAPI-Key":  []string{os.Getenv("ENV_RAPIDAPI_KEY")},
		"X-RapidAPI-Host": []string{host},
	}

	return &Crawler{
		Client:  *http.DefaultClient,
		Options: options,
		URL:     fmt.Sprintf("https://%s/api/basketball", host),
		Header:  header,
	}
}

type Crawler struct {
	http.Client
	Options
	URL    string
	Header http.Header
}

type Action struct {
	Name    string
	Handler func(context.Context) error
}

func (c *Crawler) Run() error {
	ctx := context.Background()
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to create an SQL transaction: %w", err)
	}
	actions := []Action{
		{"clean up old data", c.CleanUp},
		{"update team info", c.UpdateTeamInfo},
	}
	for _, action := range actions {
		log.Info().Msgf("Running %q ...\n", action.Name)
		if err := action.Handler(ctx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("unable to perform %q: %w", action.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit an SQL transaction: %w", err)
	}

	return nil
}

func (c *Crawler) CleanUp(ctx context.Context) error {
	if _, err := models.TeamInfos().DeleteAllG(ctx); err != nil {
		return fmt.Errorf("CleanUp: unable to delete records: %w", err)
	}
	return nil
}

func (c *Crawler) UpdateTeamInfo(ctx context.Context) error {
	response, err := misc.GetOne[Standings]{
		Client: c.Client,
		URL:    fmt.Sprintf("%s/tournament/%d/season/%d/standings/total", c.URL, c.TournamentID, c.SeasonID),
		UpdateRequest: func(req *http.Request) {
			req.Header = c.Header
		},
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdateTeamInfo: get list of teams: %w", err)
	}
	teamsIDs := make(map[uint]struct{})
	for _, item := range response.Standings {
		if item.Type != "total" {
			continue
		}
		for _, row := range item.Rows {
			teamsIDs[row.Team.ID] = struct{}{}
		}
	}
	teamsURLs := make([]string, len(teamsIDs))
	idx := 0
	for id := range teamsIDs {
		teamsURLs[idx] = fmt.Sprintf("%s/team/%d", c.URL, id)
		idx++
	}
	// -- retrieve list of teams
	teamsDetails, err := misc.GetList[TeamDetailsData]{
		Client:         c.Client,
		URLs:           teamsURLs,
		Concurrency:    2,
		RPS:            4,
		ExpectedStatus: 200,
		UpdateRequest: func(req *http.Request) {
			req.Header = c.Header
		},
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdateTeamInfo: get team info records: %w", err)
	}
	// logos
	images, err := models.Images(models.ImageWhere.ParentID.IsNull()).AllG(ctx)
	if err != nil {
		return fmt.Errorf("UpdateTeamInfo: unable to retrieve logos for all teams: %w", err)
	}
	logoByShortName := make(map[string]*string, len(images))
	for _, image := range images {
		imageUrl := image.ImageURL
		logoByShortName[image.DisplayName] = &imageUrl
	}
	for _, teamDetails := range teamsDetails {
		team := teamDetails.Team
		teamInfo := models.TeamInfo{
			ID:             int(team.ID),
			Name:           team.Name,
			Slug:           team.Slug,
			Abbr:           team.NameCode,
			ShortName:      team.ShortName,
			ArenaName:      team.Venue.Stadium.Name,
			ArenaSize:      team.Venue.Stadium.Capacity,
			Color:          team.TeamColors.Primary,
			SecondaryColor: team.TeamColors.Secondary,
			Logo:           null.StringFromPtr(logoByShortName[team.ShortName]),
		}
		if err := teamInfo.InsertG(ctx, boil.Infer()); err != nil {
			return fmt.Errorf("unable to insert record with id %q: %w", team.ID, err)
		}
	}

	return nil
}
