package parser

import (
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
)

type GamePlayerStatsParser struct{}

func (*GamePlayerStatsParser) parseBasicBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string, quarter int) []model.GamePlayerBasicStat {
	stats := []model.GamePlayerBasicStat{}

	for _, rowMap := range Table(tbl) {
		gpbs := gamePlayerBasicStatsFromRow(rowMap)
		if gpbs != nil {
			gpbs.GameId = gameId
			gpbs.TeamId = teamId
			gpbs.Quarter = quarter
			stats = append(stats, *gpbs)
		}
	}

	return stats
}

func (*GamePlayerStatsParser) parseAdvancedBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string) []model.GamePlayerAdvancedStat {
	stats := []model.GamePlayerAdvancedStat{}

	for _, rowMap := range Table(tbl) {
		gpas := gamePlayerAdvancedStatsFromRow(rowMap)
		if gpas != nil {
			gpas.GameId = gameId
			gpas.TeamId = teamId
			stats = append(stats, *gpas)
		}
	}

	return stats
}

func (*GamePlayerStatsParser) parseBasicBoxScoreGameTable(tbl *colly.HTMLElement, gameId, teamId string) []model.GamePlayer {
	players := []model.GamePlayer{}

	for i, rowMap := range Table(tbl) {
		player := *gamePlayerFromRow(rowMap, i)
		player.GameId = gameId
		player.TeamId = teamId
		players = append(players, player)
	}

	return players
}

func (*GamePlayerStatsParser) parseInactivePlayersList(box *colly.HTMLElement, gameId string) []model.GamePlayer {
	gp := []model.GamePlayer{}
	var teamId string

	// we will encounter two team labels, each surrounded by span and strong, and should set the teamId when this happens
	box.ForEach("div:nth-child(1) > span, div:nth-child(1) > a", func(_ int, t *colly.HTMLElement) {
		if t.ChildText("strong") != "" {
			teamId = strings.TrimSpace(t.Text)
		}
		// all players after that label therefore belong to that team
		// note we also check for link text but ignore it, thanks to an empty link pointing to an invalid player ID in 202202110BOS
		if teamId != "" && t.Attr("href") != "" && t.Text != "" {
			gp = append(gp, model.GamePlayer{
				GameId:   gameId,
				TeamId:   teamId,
				PlayerId: ParseLastId(t.Attr("href")),
				Status:   "I",
			})
		}
	})

	return gp
}

func parseBoxScoreTableProperties(id string) (team, boxType string, quarter int) {
	parts := strings.Split(id, "-") // expected format: box-TEAM-q/ot#-basic
	team = parts[1]
	boxType = parts[3]

	if boxType == "basic" {
		if q := parts[2]; strings.HasPrefix(q, "q") {
			quarter, _ = strconv.Atoi(strings.ReplaceAll(q, "q", ""))
		} else if strings.HasPrefix(q, "ot") {
			quarter, _ = strconv.Atoi(strings.ReplaceAll(q, "ot", ""))
			quarter = quarter + 4 // adjust for OT
		} else if q == "game" {
			quarter = math.MaxInt
		}
	}

	// quarter will be 0 for half basic boxes and advanced boxes

	return
}

func gamePlayerBasicStatsFromRow(rowMap map[string]*colly.HTMLElement) *model.GamePlayerBasicStat {
	if _, ok := rowMap["reason"]; !ok {
		gpbs := new(model.GamePlayerBasicStat)
		gpbs.PlayerId = ParseLastId(parseLink(rowMap["player"])) // a "reason" column indicates the player did not play
		gpbs.TimePlayed, _ = parseDuration(rowMap["mp"].Text)
		gpbs.FieldGoals, _ = strconv.Atoi(rowMap["fg"].Text)
		gpbs.FieldGoalsAttempted, _ = strconv.Atoi(rowMap["fga"].Text)
		gpbs.FieldGoalPct, _ = parseFloatStat(rowMap["fg_pct"].Text)
		gpbs.ThreePointers, _ = strconv.Atoi(rowMap["fg3"].Text)
		gpbs.ThreePointersAttempted, _ = strconv.Atoi(rowMap["fg3a"].Text)
		gpbs.ThreePointersPct, _ = parseFloatStat(rowMap["fg3_pct"].Text)
		gpbs.FreeThrows, _ = strconv.Atoi(rowMap["ft"].Text)
		gpbs.FreeThrowsAttempted, _ = strconv.Atoi(rowMap["fta"].Text)
		gpbs.FreeThrowsPct, _ = parseFloatStat(rowMap["ft_pct"].Text)
		gpbs.OffensiveRB, _ = strconv.Atoi(rowMap["orb"].Text)
		gpbs.DefensiveRB, _ = strconv.Atoi(rowMap["drb"].Text)
		gpbs.TotalRB, _ = strconv.Atoi(rowMap["trb"].Text)
		gpbs.Assists, _ = strconv.Atoi(rowMap["ast"].Text)
		gpbs.Steals, _ = strconv.Atoi(rowMap["stl"].Text)
		gpbs.Blocks, _ = strconv.Atoi(rowMap["blk"].Text)
		gpbs.Turnovers, _ = strconv.Atoi(rowMap["tov"].Text)
		gpbs.PersonalFouls, _ = strconv.Atoi(rowMap["pf"].Text)
		gpbs.Points, _ = strconv.Atoi(rowMap["pts"].Text)
		gpbs.PlusMinus, _ = strconv.Atoi(strings.Replace(rowMap["plus_minus"].Text, "+", "", 1))
		return gpbs
	}

	return nil
}

func gamePlayerAdvancedStatsFromRow(rowMap map[string]*colly.HTMLElement) *model.GamePlayerAdvancedStat {
	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		gpas := new(model.GamePlayerAdvancedStat)
		gpas.PlayerId = ParseLastId(parseLink(rowMap["player"]))
		gpas.TrueShootingPct, _ = parseFloatStat(rowMap["ts_pct"].Text)
		gpas.EffectiveFgPct, _ = parseFloatStat(rowMap["efg_pct"].Text)
		gpas.ThreePtAttemptRate, _ = parseFloatStat(rowMap["fg3a_per_fga_pct"].Text)
		gpas.FreeThrowAttemptRate, _ = parseFloatStat(rowMap["fta_per_fga_pct"].Text)
		gpas.OffensiveRbPct, _ = parseFloatStat(rowMap["orb_pct"].Text)
		gpas.DefensiveRbPct, _ = parseFloatStat(rowMap["drb_pct"].Text)
		gpas.TotalRbPct, _ = parseFloatStat(rowMap["trb_pct"].Text)
		gpas.AssistPct, _ = parseFloatStat(rowMap["ast_pct"].Text)
		gpas.StealPct, _ = parseFloatStat(rowMap["stl_pct"].Text)
		gpas.BlockPct, _ = parseFloatStat(rowMap["blk_pct"].Text)
		gpas.TurnoverPct, _ = parseFloatStat(rowMap["tov_pct"].Text)
		gpas.UsagePct, _ = parseFloatStat(rowMap["usg_pct"].Text)
		gpas.OffensiveRating, _ = strconv.Atoi(rowMap["off_rtg"].Text)
		gpas.DefensiveRating, _ = strconv.Atoi(rowMap["def_rtg"].Text)

		if _, ok := rowMap["bpm"]; ok {
			gpas.BoxPlusMinus, _ = parseFloatStat(rowMap["bpm"].Text)
		}

		return gpas
	}

	return nil
}

func gamePlayerFromRow(rowMap map[string]*colly.HTMLElement, index int) *model.GamePlayer {
	gp := new(model.GamePlayer)
	gp.PlayerId = ParseLastId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		if index < 5 { // the first 5 players are the starters
			gp.Status = "S"
		} else { // the rest are reserves
			gp.Status = "R"
		}
	} else {
		gp.Status = "D"
	}

	return gp
}
