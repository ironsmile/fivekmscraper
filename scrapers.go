package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ironsmile/logger"
)

func getUserProfile(userID int) (uint8, string, error) {
	doc, err := goquery.NewDocument(getUserProfileURL(userID))

	if err != nil {
		return 0, "", err
	}

	ageString := strings.TrimSpace(doc.Find(profileAgeSelectorNoImage).Text())

	if ageString == "" {
		ageString = strings.TrimSpace(doc.Find(profileAgeSelectorImage).Text())
	}

	if ageString == "" {
		return 0, "", UserNotFound
	}

	ageString = strings.TrimSuffix(ageString, "г.")

	ageInt, err := strconv.Atoi(ageString)

	if err != nil {
		return 0, "", fmt.Errorf("Could not convert age (%s) to integer", ageString)
	}

	name := strings.TrimSpace(doc.Find(profileNameSelecotor).Text())
	name = strings.TrimSpace(strings.TrimLeft(name, "Персонална статистика - "))

	return uint8(ageInt), name, nil
}

func getUserStats(userID int) ([]RunData, error) {
	age, name, err := getUserProfile(userID)
	if err != nil {
		return nil, err
	}
	return getUserStatsForRuns(userID, age, name)
}

func getUserStatsForRuns(userID int, age uint8, name string) ([]RunData, error) {

	doc, err := goquery.NewDocument(getUserStatsURL(userID))

	if err != nil {
		return nil, err
	}

	tBody := doc.Find(statsRowSelector)

	if tBody.Length() < 1 {
		return []RunData{}, nil
	}

	returnData := []RunData{}

	tBody.Each(func(index int, el *goquery.Selection) {
		userData := RunData{}
		userData.ID = userID
		userData.Age = age
		userData.Name = name

		// Parse run date
		dateString := strings.TrimSpace(el.Find(statsRowDate).Text())

		// Mon Jan 2 15:04:05 MST 2006
		date, err := time.Parse(`02.01.2006`, dateString)

		if err != nil {
			logger.Errorf("Could not parse date (%d, %d): %s", userID, index, err)
			return
		}

		userData.RunDate = date

		// Parse position
		positionString := strings.TrimSpace(el.Find(statsRowPosition).Text())
		pos, err := strconv.Atoi(positionString)

		if err != nil {
			logger.Errorf("Could not parse position (%d, %d): %s", userID, index, err)
			return
		}

		userData.Position = uint32(pos)

		// Parse time
		timeString := strings.TrimSpace(el.Find(statsRowTime).Text())

		timeParsed, err := parseDuration(timeString)

		if err != nil {
			logger.Errorf("Could not parse run time (%d, %d): %s", userID, index, err)
			return
		}

		userData.Time = timeParsed

		// Parse avg speed
		avgString := strings.TrimSpace(el.Find(statsRowAvgSpeed).Text())
		avgString = strings.TrimSuffix(avgString, " км/ч")
		avgFloat, err := strconv.ParseFloat(avgString, 64)

		if err != nil {
			logger.Errorf("Could not parse avg speed (%d, %d): %s", userID, index, err)
			return
		}

		userData.AvgSpeed = float32(avgFloat)

		// Parse tempo
		tempoString := strings.TrimSpace(el.Find(statsRowTempo).Text())

		tempoParsed, err := parseDuration(tempoString)

		if err != nil {
			logger.Errorf("Could not parse run time (%d, %d): %s", userID, index, err)
			return
		}

		userData.Tempo = tempoParsed

		// Parse place
		userData.Place = strings.TrimSpace(el.Find(statsRowPlace).Text())

		returnData = append(returnData, userData)
	})

	return returnData, nil
}

func getUserProfileURL(userID int) string {
	return fmt.Sprintf("%s%s", FiveKMURL, fmt.Sprintf(UserProfilePath, userID))
}

func getUserStatsURL(userID int) string {
	return fmt.Sprintf("%s%s", FiveKMURL, fmt.Sprintf(UserStatsPath, userID))
}

func parseDuration(dur string) (time.Duration, error) {
	return time.ParseDuration(strings.Replace(dur, ":", "m", 1) + "s")
}

func findLastUserID() int {
	return binSearchUser(0, 1e5)
}

func binSearchUser(min, max int) int {
	if max-min <= 1 {
		return max
	}

	mid := (min + max) / 2

	if lessThan(mid) {
		return binSearchUser(mid, max)
	}

	return binSearchUser(min, mid)
}

func lessThan(userID int) bool {
	_, _, err := getUserProfile(userID)
	if err != UserNotFound {
		return true
	}

	for i := 1; i <= 10; i++ {
		_, _, err := getUserProfile(userID + i)
		if err != UserNotFound {
			return true
		}
	}

	return false
}
