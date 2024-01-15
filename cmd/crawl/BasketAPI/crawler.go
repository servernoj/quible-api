package BasketAPI

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quible-io/quible-api/cmd/crawl/common"
	"github.com/quible-io/quible-api/lib/models"
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
		log.Printf("Running %q ...\n", action.Name)
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
	response, err := common.GetOne[Standings]{
		Client: c.Client,
		URL:    fmt.Sprintf("%s/tournament/%d/season/%d/standings/total", c.URL, c.TournamentID, c.SeasonID),
		UpdateRequest: func(req *http.Request) {
			req.Header = c.Header
		},
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdateTeams: get refs: %w", err)
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
	log.Printf("%+v\nlength: %d\n", teamsIDs, len(teamsIDs))
	// // -- retrieve list of teams
	// teams, err := common.GetList[TeamDetailsData]{
	// 	Client:      c.Client,
	// 	URLs:        teamsURLs,
	// 	Concurrency: 10,
	// }.Do()
	// if err != nil {
	// 	return fmt.Errorf("UpdateTeams: get items: %w", err)
	// }
	return nil
}
