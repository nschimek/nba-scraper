package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type GamePlayer struct {
	GameId, TeamId, PlayerId, Status string
}

type GamePlayerBasicStats struct {
	GameId, TeamId, PlayerId                                                                                string
	Quarter                                                                                                 int
	TimePlayed                                                                                              time.Duration
	FieldGoals, FieldGoalsAttempted, ThreePointers, ThreePointersAttempted, FreeThrows, FreeThrowsAttempted int
	FieldGoalPct, ThreePointersPct, FreeThrowsPct                                                           float64
	OffensiveRB, DefensiveRB, TotalRB, Assists, Steals, Blocks, Turnovers, PersonalFouls, Points, PlusMinus int
}

type GamePlayerAdvancedStats struct {
	GameId, TeamId, PlayerId                                                                                              string
	TrueShootingPct, EffectiveFgPct, ThreePtAttemptRate, FreeThrowAttemptRate, OffensiveRbPct, DefensiveRbPct, TotalRbPct float64
	AssistPct, StealPct, BlockPct, TurnoverPct, UsagePct, BoxPlusMinus                                                    float64
	OffensiveRating, DefensiveRating                                                                                      int
}

func ParseBasicBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string, quarter int) []GamePlayerBasicStats {
	stats := []GamePlayerBasicStats{}

	for _, rowMap := range Table(tbl) {
		stats = append(stats, gamePlayerBasicStatsFromRow(rowMap, gameId, teamId, quarter))
	}

	return stats
}

func ParseAdvancedBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string) []GamePlayerAdvancedStats {
	stats := []GamePlayerAdvancedStats{}

	for _, rowMap := range Table(tbl) {
		stats = append(stats, gamePlayerAdvancedStatsFromRow(rowMap, gameId, teamId))
	}

	return stats
}

func ParseBasicBoxScoreGameTable(tbl *colly.HTMLElement, gameId, teamId string) []GamePlayer {
	players := []GamePlayer{}

	for i, rowMap := range Table(tbl) {
		players = append(players, gamePlayerFromRow(rowMap, gameId, teamId, i))
	}

	return players
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

func gamePlayerBasicStatsFromRow(rowMap map[string]*colly.HTMLElement, gameId, teamId string, quarter int) (gpbs GamePlayerBasicStats) {
	gpbs.GameId = gameId
	gpbs.TeamId = teamId
	gpbs.Quarter = quarter
	gpbs.PlayerId = parsePlayerId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
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
	}

	return
}

func gamePlayerAdvancedStatsFromRow(rowMap map[string]*colly.HTMLElement, gameId, teamId string) (gpas GamePlayerAdvancedStats) {
	gpas.GameId = gameId
	gpas.TeamId = teamId
	gpas.PlayerId = parsePlayerId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
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
	}

	return
}

func gamePlayerFromRow(rowMap map[string]*colly.HTMLElement, gameId, teamId string, index int) (gp GamePlayer) {
	gp.GameId = gameId
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

	return
}
