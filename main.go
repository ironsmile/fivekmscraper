package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ironsmile/logger"
)

var (
	UserNotFound = errors.New("User Not Found")
)

const (
	FiveKMURL       = `http://5km.5kmrun.bg`
	UserProfilePath = `/usr.php?id=%d`
	UserStatsPath   = `/stat.php?id=%d`

	OutputFilePath = "5km-stats.csv"
	MaleFemaleDict = "male-female-dict.dat"
)

func main() {
	logger.Default().Logger = log.New(os.Stdout, "", log.LstdFlags)

	outFile, err := os.Create(OutputFilePath)

	if err != nil {
		logger.Errorf("Could not create output file: %s", err)
		os.Exit(1)
	}

	defer outFile.Close()

	males, females := loadNamesDicts(MaleFemaleDict)

	logger.Logln("Searching for the last registered userID.")
	var MaxUserID = findLastUserID()
	logger.Logf("Found: %d\n", MaxUserID)

	var userIDCounter int64
	var workerWG sync.WaitGroup
	workSink := make(chan []RunData, 5)

	for i := 0; i < runtime.NumCPU(); i++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()

			for userIDCounter <= int64(MaxUserID) {
				userID := atomic.AddInt64(&userIDCounter, 1)

				userData, err := getUserStats(int(userID))

				if err == UserNotFound {
					continue
				}

				if err != nil {
					logger.Errorf("Error collecting stats: %s\n", err)
				} else {
					workSink <- userData
				}
			}

		}()
	}

	setToMale := func(names []string) {
		for _, name := range names {
			name := strings.ToLower(strings.TrimSpace(name))
			if len(name) < 3 {
				continue
			}
			males[name] = struct{}{}
		}
	}

	setToFemale := func(names []string) {
		for _, name := range names {
			name := strings.ToLower(strings.TrimSpace(name))
			if len(name) < 3 {
				continue
			}
			females[name] = struct{}{}
		}
	}

	isNameMale := func(name string) int {
		splittedNames := strings.Split(name, " ")
		var names []string

		for _, name := range splittedNames {
			if name == "" {
				continue
			}
			names = append(names, name)
		}

		for _, name := range names {
			name := strings.ToLower(strings.TrimSpace(name))

			if _, ok := males[name]; ok {
				setToMale(names)
				return 1
			}

			if _, ok := females[name]; ok {
				setToFemale(names)
				return 0
			}
		}

		if len(names) == 2 && (strings.HasSuffix(names[1], "ва") || strings.HasSuffix(names[1], "va")) {
			setToFemale(names)
			return 0
		}

		if len(names) == 2 && (strings.HasSuffix(names[1], "ов") ||
			strings.HasSuffix(names[1], "ov") ||
			strings.HasSuffix(names[1], "ев") ||
			strings.HasSuffix(names[1], "ev")) {
			setToMale(names)
			return 1
		}

		var sex string

		for sex != "m" && sex != "f" {
			fmt.Printf("The name '%s' is male [m] or female [f]? > ", name)
			if _, err := fmt.Scanf("%s", &sex); err != nil {
				if err == io.EOF {
					os.Exit(0)
				}
				logger.Errorln(err)
			}
		}

		if sex == "m" {
			setToMale(names)
			return 1
		}

		setToFemale(names)
		return 0
	}

	outputFileWritten := make(chan struct{})
	go func() {

		csvWriter := csv.NewWriter(outFile)
		csvWriter.Write([]string{
			"ID",
			"name",
			"is_male",
			"age",
			"place",
			"date",
			"time",
			"position",
			"avg_speed_kph",
			"tempo",
		})

		defer func() {
			csvWriter.Flush()
			outputFileWritten <- struct{}{}
		}()

		timeBetweenOutputs, _ := time.ParseDuration("5s")
		lastOutputTime := time.Now()
		var profilesSinceLastOutput int

		for userData := range workSink {
			for _, userRun := range userData {

				isMale := isNameMale(userRun.Name)

				csvWriter.Write([]string{
					fmt.Sprintf("%d", userRun.ID),
					userRun.Name,
					fmt.Sprintf("%d", isMale),
					fmt.Sprintf("%d", userRun.Age),
					userRun.Place,
					userRun.RunDate.Format(`2006-01-02`),
					userRun.Time.String(),
					fmt.Sprintf("%d", userRun.Position),
					fmt.Sprintf("%.2f", userRun.AvgSpeed),
					userRun.Tempo.String(),
				})
			}

			if profilesSinceLastOutput >= 100 || time.Since(lastOutputTime) >= timeBetweenOutputs {
				profilesSinceLastOutput = 0
				lastOutputTime = time.Now()

				parsedPrc := float64(float64(userIDCounter)/float64(MaxUserID)) * 100
				logger.Logf("Parsed %d/%d (%.2f%%)", userIDCounter, MaxUserID, parsedPrc)
				saveNemesDicts(MaleFemaleDict, males, females)
			}

			profilesSinceLastOutput++
		}
	}()

	workerWG.Wait()
	close(workSink)
	<-outputFileWritten
}
