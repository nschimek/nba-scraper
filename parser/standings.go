package parser

type Standings struct {
	TeamId                                                       string
	Season                                                       int
	Overall, Home, Road, East, West                              WinLoss
	Atlantic, Central, Southeast, Northwest, Pacific, Southwest  WinLoss
	PreAllStar, PostAllStar, MarginLess3, MarginGreater10        WinLoss
	October, November, December, January, February, March, April WinLoss
}

type WinLoss struct {
	wins, losses int
}
