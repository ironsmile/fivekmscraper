package main

const (
	profileAgeSelectorNoImage = `article.article-container:nth-child(2) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > table:nth-child(1) > tbody:nth-child(2) > tr:nth-child(1) > td:nth-child(1)`
	profileAgeSelectorImage   = `div.col-sm-9:nth-child(2) > div:nth-child(1) > div:nth-child(1) > table:nth-child(1) > tbody:nth-child(2) > tr:nth-child(1) > td:nth-child(1)`
	profileNameSelecotor      = `div.col-md-9:nth-child(1) > h2:nth-child(1)`

	statsRowSelector = `.table > tbody:nth-child(2) tr`
	statsRowPlace    = `td:nth-child(1)`
	statsRowDate     = `td:nth-child(2)`
	statsRowPosition = `td:nth-child(3)`
	statsRowTime     = `td:nth-child(4)`
	statsRowAvgSpeed = `td:nth-child(7)`
	statsRowTempo    = `td:nth-child(8)`
)
