package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

const (
	basePlayerBodyElement = "body > div#wrap"
	playerInfoElement     = basePlayerBodyElement + " > div#info > div#meta > div:nth-child(2)"
)

type PlayerScraper struct {
	Config       *core.Config                               `Inject:""`
	Colly        *colly.Collector                           `Inject:""`
	PlayerParser *parser.PlayerParser                       `Inject:""`
	Repository   *repository.SimpleRepository[model.Player] `Inject:""`
	ScrapedData  []model.Player
}

func (s *PlayerScraper) Scrape(idMap map[string]struct{}) {
	core.Log.WithField("ids", len(idMap)).Info("Got Player ID(s) to scrape")
	if s.Config.Suppression.Player > 0 {
		s.suppressRecent(idMap)
	}
	s.scrape(idMap)
}

func (s *PlayerScraper) scrape(idMap map[string]struct{}) {
	core.Log.WithField("ids", len(idMap)).Info("Scraping Player ID(s)...")

	for id := range idMap {
		if strings.HasPrefix(id, "tmp-") {
			core.Log.WithField("id", id).Error("Encountered a known invalid temp Player ID, skipping...")
		} else {
			player := s.parsePlayerPage(id)
			if !player.HasErrors() {
				s.ScrapedData = append(s.ScrapedData, player)
			} else {
				player.LogErrors(fmt.Sprintf("player %s", player.ID))
			}
		}
	}

	core.Log.WithField("players", len(s.ScrapedData)).Info("Finished scraping Player page(s)!")
}

func (s *PlayerScraper) Persist() {
	if len(s.ScrapedData) > 0 {
		s.Repository.Upsert(s.ScrapedData, "Players")
	} else {
		core.Log.Info("No Players scraped to persist! Skipping...")
	}
}

func (s *PlayerScraper) suppressRecent(idMap map[string]struct{}) {
	core.Log.Infof("Checking for players updated in last %d days (for suppression)...", s.Config.Suppression.Player)
	ids, _ := s.Repository.GetRecentlyUpdated(s.Config.Suppression.Player, core.IdMapToArray(idMap), "Player")
	if ids != nil && len(ids) > 0 {
		core.SuppressIdMap(idMap, ids)
	}
}

func (s *PlayerScraper) parsePlayerPage(id string) (player model.Player) {
	c := core.CloneColly(s.Colly)

	player.ID = id

	c.OnHTML(playerInfoElement, func(div *colly.HTMLElement) {
		s.PlayerParser.PlayerInfoBox(&player, div)
	})

	c.OnError(func(r *colly.Response, err error) {
		player.CaptureError(NewScraperError(err, r.Request.URL.String()))
	})

	c.Visit(s.getUrl(id))

	return
}

func (*PlayerScraper) getUrl(id string) string {
	return fmt.Sprintf(playerIdPage, id[0:1], id)
}
