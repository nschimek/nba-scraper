package scraper

const (
	BaseHttp        = "https://www.basketball-reference.com"
	baseLeaguesPath = "leagues"
	playerIdPage    = BaseHttp + "/players/%s/%s.html"
	teamIdPage      = BaseHttp + "/teams/%s/%d.html"
	injuryPage      = BaseHttp + "/friv/injuries.fcgi"
)
