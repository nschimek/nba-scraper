package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type GamePlayerStatsParser struct{}

func (*GamePlayerStatsParser) parseBasicBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string, quarter int) []model.GamePlayerBasicStat {
	stats := []model.GamePlayerBasicStat{}

	for _, rowMap := range Table(tbl) {
		gpbs, err := gamePlayerBasicStatsFromRow(rowMap)
		if (gpbs != model.GamePlayerBasicStat{} && err == nil) {
			gpbs.GameId = gameId
			gpbs.TeamId = teamId
			gpbs.Quarter = quarter
			stats = append(stats, gpbs)
		} else if err != nil {
			core.Log.Warnf("error encountered parsing player basic stats for team %s (quarter %d): %s", teamId, quarter, err.Error())
		}
	}

	if len(stats) == 0 {
		core.Log.Warnf("did parse any player basic stats for team %s (quarter %d)", teamId, quarter)
	}

	return stats
}

func (*GamePlayerStatsParser) parseAdvancedBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string) []model.GamePlayerAdvancedStat {
	stats := []model.GamePlayerAdvancedStat{}

	for _, rowMap := range Table(tbl) {
		gpas, err := gamePlayerAdvancedStatsFromRow(rowMap)
		if (gpas != model.GamePlayerAdvancedStat{} && err == nil) {
			gpas.GameId = gameId
			gpas.TeamId = teamId
			stats = append(stats, gpas)
		} else if err != nil {
			core.Log.Warnf("error encountered parsing player advanced stats for team %s: %s", teamId, err.Error())
		}
	}

	if len(stats) == 0 {
		core.Log.Warnf("did parse any player advanced stats for team %s", teamId)
	}

	return stats
}

func (*GamePlayerStatsParser) parseBasicBoxScoreGameTable(tbl *colly.HTMLElement, gameId, teamId string) []model.GamePlayer {
	players := []model.GamePlayer{}

	for i, rowMap := range Table(tbl) {
		player, err := gamePlayerFromRow(rowMap, i)
		if err != nil {
			player.GameId = gameId
			player.TeamId = teamId
			players = append(players, player)
		} else {
			core.Log.Warnf("error encountered parsing player list from basic box score for team %s: %s", teamId, err.Error())
		}
	}

	if len(players) == 0 {
		core.Log.Warnf("did parse any players who played from basic box score for team %s", teamId)
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
			playerId, err := ParseLastId(t.Attr("href"))
			if err != nil {
				gp = append(gp, model.GamePlayer{
					GameId:   gameId,
					TeamId:   teamId,
					PlayerId: playerId,
					Status:   "I",
				})
			} else {
				core.Log.Warnf("could not get player ID of inactive player for team %s, ignoring...", teamId)
			}
		}
	})

	return gp
}

func parseBoxScoreTableProperties(id string) (team, boxType string, quarter int, err error) {
	parts := strings.Split(id, "-") // expected format: box-TEAM-q/ot#-basic

	if len(parts) != 4 {
		return "", "", 0, errors.New("box score ID not in expected format")
	}

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

func gamePlayerBasicStatsFromRow(rowMap map[string]*colly.HTMLElement) (model.GamePlayerBasicStat, error) {
	var err error

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates a player did not play
		gpbs := new(model.GamePlayerBasicStat)
		gpbs.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))
		gpbs.TimePlayed, _ = parseDuration(getColumnText(rowMap, "mp"))
		gpbs.FieldGoals, _ = strconv.Atoi(getColumnText(rowMap, "fg"))
		gpbs.FieldGoalsAttempted, _ = strconv.Atoi(getColumnText(rowMap, "fga"))
		gpbs.FieldGoalPct, _ = parseFloatStat(getColumnText(rowMap, "fg_pct"))
		gpbs.ThreePointers, _ = strconv.Atoi(getColumnText(rowMap, "fg3"))
		gpbs.ThreePointersAttempted, _ = strconv.Atoi(getColumnText(rowMap, "fg3a"))
		gpbs.ThreePointersPct, _ = parseFloatStat(getColumnText(rowMap, "fg3_pct"))
		gpbs.FreeThrows, _ = strconv.Atoi(getColumnText(rowMap, "ft"))
		gpbs.FreeThrowsAttempted, _ = strconv.Atoi(getColumnText(rowMap, "fta"))
		gpbs.FreeThrowsPct, _ = parseFloatStat(getColumnText(rowMap, "ft_pct"))
		gpbs.OffensiveRB, _ = strconv.Atoi(getColumnText(rowMap, "orb"))
		gpbs.DefensiveRB, _ = strconv.Atoi(getColumnText(rowMap, "drb"))
		gpbs.TotalRB, _ = strconv.Atoi(getColumnText(rowMap, "trb"))
		gpbs.Assists, _ = strconv.Atoi(getColumnText(rowMap, "ast"))
		gpbs.Steals, _ = strconv.Atoi(getColumnText(rowMap, "stl"))
		gpbs.Blocks, _ = strconv.Atoi(getColumnText(rowMap, "blk"))
		gpbs.Turnovers, _ = strconv.Atoi(getColumnText(rowMap, "tov"))
		gpbs.PersonalFouls, _ = strconv.Atoi(getColumnText(rowMap, "pf"))
		gpbs.Points, _ = strconv.Atoi(getColumnText(rowMap, "pts"))
		gpbs.PlusMinus, _ = strconv.Atoi(strings.Replace(getColumnText(rowMap, "plus_minus"), "+", "", 1))
		return *gpbs, err
	}

	return model.GamePlayerBasicStat{}, nil
}

func gamePlayerAdvancedStatsFromRow(rowMap map[string]*colly.HTMLElement) (model.GamePlayerAdvancedStat, error) {
	var err error

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		gpas := new(model.GamePlayerAdvancedStat)
		gpas.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))
		gpas.TrueShootingPct, _ = parseFloatStat(getColumnText(rowMap, "ts_pct"))
		gpas.EffectiveFgPct, _ = parseFloatStat(getColumnText(rowMap, "efg_pct"))
		gpas.ThreePtAttemptRate, _ = parseFloatStat(getColumnText(rowMap, "fg3a_per_fga_pct"))
		gpas.FreeThrowAttemptRate, _ = parseFloatStat(getColumnText(rowMap, "fta_per_fga_pct"))
		gpas.OffensiveRbPct, _ = parseFloatStat(getColumnText(rowMap, "orb_pct"))
		gpas.DefensiveRbPct, _ = parseFloatStat(getColumnText(rowMap, "drb_pct"))
		gpas.TotalRbPct, _ = parseFloatStat(getColumnText(rowMap, "trb_pct"))
		gpas.AssistPct, _ = parseFloatStat(getColumnText(rowMap, "ast_pct"))
		gpas.StealPct, _ = parseFloatStat(getColumnText(rowMap, "stl_pct"))
		gpas.BlockPct, _ = parseFloatStat(getColumnText(rowMap, "blk_pct"))
		gpas.TurnoverPct, _ = parseFloatStat(getColumnText(rowMap, "tov_pct"))
		gpas.UsagePct, _ = parseFloatStat(getColumnText(rowMap, "usg_pct"))
		gpas.OffensiveRating, _ = strconv.Atoi(getColumnText(rowMap, "off_rtg"))
		gpas.DefensiveRating, _ = strconv.Atoi(getColumnText(rowMap, "def_rtg"))

		if bpm, ok := rowMap["bpm"]; ok {
			gpas.BoxPlusMinus, _ = parseFloatStat(bpm.Text)
		}

		return *gpas, err
	}

	return model.GamePlayerAdvancedStat{}, nil
}

func gamePlayerFromRow(rowMap map[string]*colly.HTMLElement, index int) (model.GamePlayer, error) {
	var err error

	gp := new(model.GamePlayer)
	gp.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		if index < 5 { // the first 5 players are the starters
			gp.Status = "S"
		} else { // the rest are reserves
			gp.Status = "R"
		}
	} else {
		gp.Status = "D"
	}

	return *gp, err
}
