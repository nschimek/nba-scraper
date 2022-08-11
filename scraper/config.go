package scraper

const (
	BaseHttp       = "https://www.basketball-reference.com"
	scheduleIdPage = BaseHttp + "/leagues/NBA_%d_games-%s.html"
	playerIdPage   = BaseHttp + "/players/%s/%s.html"
	teamIdPage     = BaseHttp + "/teams/%s/%d.html"
	gameIdPage     = BaseHttp + "/boxscores/%s.html"
	injuryPage     = BaseHttp + "/friv/injuries.fcgi"
	standingPage   = BaseHttp + "/leagues/NBA_%d_standings.html"
)
