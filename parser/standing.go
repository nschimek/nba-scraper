package parser

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type StandingParser struct {
	Config *core.Config `Inject:""`
}

type WinLoss struct {
	Wins, Losses int
}

func (p *StandingParser) StandingsTable(tbl *colly.HTMLElement) []model.TeamStanding {
	standings := []model.TeamStanding{}

	for _, rowMap := range Table(tbl) {
		standing, err := standingFromRow(rowMap)
		standing.CaptureError(err)
		standing.Season = p.Config.Season
		standings = append(standings, standing)
	}

	return standings
}

func standingFromRow(rowMap map[string]*colly.HTMLElement) (model.TeamStanding, error) {
	var err error
	standing := new(model.TeamStanding)

	standing.Rank, _ = strconv.Atoi(getColumnText(rowMap, "ranker"))
	standing.TeamId, err = ParseTeamId(parseLink(rowMap["team_name"]))
	standing.Overall = parseWinLoss(getColumnText(rowMap, "Overall"), "overall")
	standing.Home = parseWinLoss(getColumnText(rowMap, "Home"), "home")
	standing.Road = parseWinLoss(getColumnText(rowMap, "Road"), "road")
	standing.East = parseWinLoss(getColumnText(rowMap, "E"), "east")
	standing.West = parseWinLoss(getColumnText(rowMap, "W"), "west")
	standing.Atlantic = parseWinLoss(getColumnText(rowMap, "A"), "atlantic")
	standing.Central = parseWinLoss(getColumnText(rowMap, "C"), "central")
	standing.Southeast = parseWinLoss(getColumnText(rowMap, "SE"), "southeast")
	standing.Northwest = parseWinLoss(getColumnText(rowMap, "NW"), "northwest")
	standing.Pacific = parseWinLoss(getColumnText(rowMap, "P"), "pacific")
	standing.Southwest = parseWinLoss(getColumnText(rowMap, "SW"), "southwest")
	standing.PreAllStar = parseWinLoss(getColumnText(rowMap, "Pre"), "pre-allstar")
	standing.PostAllStar = parseWinLoss(getColumnText(rowMap, "Post"), "post-allstar")
	standing.MarginLess3 = parseWinLoss(getColumnText(rowMap, "3"), "marginLess3")
	standing.MarginGreater10 = parseWinLoss(getColumnText(rowMap, "10"), "marginGreater10")
	standing.October = parseWinLoss(getColumnText(rowMap, "Oct"), "october")
	standing.November = parseWinLoss(getColumnText(rowMap, "Nov"), "november")
	standing.December = parseWinLoss(getColumnText(rowMap, "Dec"), "december")
	standing.January = parseWinLoss(getColumnText(rowMap, "Jan"), "january")
	standing.February = parseWinLoss(getColumnText(rowMap, "Feb"), "february")
	standing.March = parseWinLoss(getColumnText(rowMap, "Mar"), "march")
	standing.April = parseWinLoss(getColumnText(rowMap, "Apr"), "april")

	return *standing, err
}

func parseWinLoss(s string, label string) model.WinLoss {
	wl := new(model.WinLoss)

	if p := strings.Split(s, "-"); s != "" && len(p) == 2 {
		wl.Wins, _ = strconv.Atoi(p[0])
		wl.Losses, _ = strconv.Atoi(p[1])
	} else {
		core.Log.WithField("column", label).Warn("could not parse win/loss record (probably blank), using 0-0")
		wl.Wins, wl.Losses = 0, 0
	}

	return *wl
}
