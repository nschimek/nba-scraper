package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type GamePlayer struct {
	TeamId, PlayerId, Status string
}

type GamePlayerBasicStats struct {
	TeamId, PlayerId                                                                                        string
	Quarter                                                                                                 int
	TimePlayed                                                                                              time.Duration
	FieldGoals, FieldGoalsAttempted, ThreePointers, ThreePointersAttempted, FreeThrows, FreeThrowsAttempted int
	FieldGoalPct, ThreePointersPct, FreeThrowsPct                                                           float64
	OffensiveRB, DefensiveRB, TotalRB, Assists, Steals, Blocks, Turnovers, PersonalFouls, Points, PlusMinus int
}

type GamePlayerAdvancedStats struct {
	TeamId, PlayerId                                                                                                      string
	TrueShootingPct, EffectiveFgPct, ThreePtAttemptRate, FreeThrowAttemptRate, OffensiveRbPct, DefensiveRbPct, TotalRbPct float64
	AssistPct, StealPct, BlockPct, TurnoverPct, UsagePct, BoxPlusMinus                                                    float64
	OffensiveRating, DefensiveRating                                                                                      int
}

func parseBasicBoxScoreTable(tbl *colly.HTMLElement, teamId string, quarter int) []GamePlayerBasicStats {
	stats := []GamePlayerBasicStats{}

	for _, rowMap := range Table(tbl) {
		gpbs := gamePlayerBasicStatsFromRow(rowMap, teamId, quarter)
		if gpbs != nil {
			stats = append(stats, *gpbs)
		}
	}

	return stats
}

func parseAdvancedBoxScoreTable(tbl *colly.HTMLElement, teamId string) []GamePlayerAdvancedStats {
	stats := []GamePlayerAdvancedStats{}

	for _, rowMap := range Table(tbl) {
		gpas := gamePlayerAdvancedStatsFromRow(rowMap, teamId)
		if gpas != nil {
			stats = append(stats, *gpas)
		}
	}

	return stats
}

func parseBasicBoxScoreGameTable(tbl *colly.HTMLElement, teamId string) []GamePlayer {
	players := []GamePlayer{}

	for i, rowMap := range Table(tbl) {
		players = append(players, *gamePlayerFromRow(rowMap, teamId, i))
	}

	return players
}

func parseInactivePlayersList(box *colly.HTMLElement) []GamePlayer {
	gp := []GamePlayer{}
	var teamId string

	// we will encounter two team labels, each surrounded by span and strong, and should set the teamId when this happens
	box.ForEach("div:nth-child(1) > span, div:nth-child(1) > a", func(_ int, t *colly.HTMLElement) {
		if t.ChildText("strong") != "" {
			teamId = t.Text
		}
		// all players after that label therefore belong to that team
		if t.Attr("href") != "" {
			gp = append(gp, GamePlayer{
				TeamId:   teamId,
				PlayerId: parsePlayerId(t.Attr("href")),
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

func gamePlayerBasicStatsFromRow(rowMap map[string]*colly.HTMLElement, teamId string, quarter int) *GamePlayerBasicStats {
	if _, ok := rowMap["reason"]; !ok {
		gpbs := new(GamePlayerBasicStats)
		gpbs.TeamId = teamId
		gpbs.Quarter = quarter
		gpbs.PlayerId = parsePlayerId(parseLink(rowMap["player"])) // a "reason" column indicates the player did not play
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

func gamePlayerAdvancedStatsFromRow(rowMap map[string]*colly.HTMLElement, teamId string) *GamePlayerAdvancedStats {
	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		gpas := new(GamePlayerAdvancedStats)
		gpas.TeamId = teamId
		gpas.PlayerId = parsePlayerId(parseLink(rowMap["player"]))
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
		gpas.BoxPlusMinus, _ = parseFloatStat(rowMap["bpm"].Text)
		return gpas
	}

	return nil
}

func gamePlayerFromRow(rowMap map[string]*colly.HTMLElement, teamId string, index int) *GamePlayer {
	gp := new(GamePlayer)
	gp.TeamId = teamId
	gp.PlayerId = parsePlayerId(parseLink(rowMap["player"]))

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
