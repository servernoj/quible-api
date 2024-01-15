package espn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"

	"github.com/quible-io/quible-api/cmd/crawl/common"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	BASE_URL = "https://sports.core.api.espn.com/v2/sports/basketball/leagues/nba"
)

func NewCrawler() *Crawler {
	return &Crawler{
		Client: *http.DefaultClient,
		URL:    BASE_URL,
	}
}

type Crawler struct {
	http.Client
	URL string
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
		{"clean up old images", c.CleanUp},
		{"update teams", c.UpdateTeams},
		{"update players", c.UpdatePlayers},
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
	if _, err := models.Images().DeleteAllG(ctx); err != nil {
		return fmt.Errorf("CleanUp: unable to delete records: %w", err)
	}
	return nil
}

func (c *Crawler) UpdateTeams(ctx context.Context) error {
	query := url.Values{
		"limit": {"100"},
	}
	teamsRefs, err := common.GetOne[ResponseWithItems]{
		Client: c.Client,
		URL:    fmt.Sprintf("%s/teams?%s", c.URL, query.Encode()),
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdateTeams: get refs: %w", err)
	}
	teamsURLs := []string{}
	for _, item := range teamsRefs.Items {
		teamsURLs = append(teamsURLs, item.Ref)
	}
	// -- retrieve list of teams
	teams, err := common.GetList[TeamItem]{
		Client:      c.Client,
		URLs:        teamsURLs,
		Concurrency: 10,
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdateTeams: get items: %w", err)
	}
	// -- extract and save logo from each team
	for _, team := range teams {
		if len(team.Logos) == 0 {
			continue
		}
		goodLogo := team.Logos[0]
		for _, logo := range team.Logos {
			if len(logo.Rel) > 0 && slices.Contains(logo.Rel, "default") {
				goodLogo = logo
				break
			}
		}
		teamRecord := &models.Image{
			ID:          team.Ref,
			DisplayName: team.ShortDisplayName,
			Slug:        team.Slug,
			ImageURL:    goodLogo.Href,
		}
		if err := teamRecord.InsertG(ctx, boil.Infer()); err != nil {
			return fmt.Errorf("UpdateTeams: unable to insert team logo for %q: %w", team.ShortDisplayName, err)
		}
	}

	return nil
}

func (c *Crawler) UpdatePlayers(ctx context.Context) error {
	query := url.Values{
		"limit": {"1000"},
	}
	playersRefs, err := common.GetOne[ResponseWithItems]{
		Client: c.Client,
		URL:    fmt.Sprintf("%s/athletes?%s", c.URL, query.Encode()),
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdatePlayers: get refs: %w", err)
	}
	playersURLs := []string{}
	for _, item := range playersRefs.Items {
		playersURLs = append(playersURLs, item.Ref)
	}
	// -- retrieve list of players
	players, err := common.GetList[PlayerItem]{
		Client:      c.Client,
		URLs:        playersURLs,
		Concurrency: 50,
	}.Do()
	if err != nil {
		return fmt.Errorf("UpdatePlayers: get items: %w", err)
	}
	// -- extract and save logo from each player
	for _, player := range players {
		if player.Headshot == nil {
			continue
		}
		playerRecord := &models.Image{
			ID:          player.Ref,
			ParentID:    null.StringFrom(player.Team.Ref),
			DisplayName: player.DisplayName,
			Slug:        player.Slug,
			ImageURL:    player.Headshot.Href,
		}
		if err := playerRecord.InsertG(ctx, boil.Infer()); err != nil {
			return fmt.Errorf("UpdateTeams: unable to insert team logo for %q: %w", player.DisplayName, err)
		}
	}

	return nil
}
