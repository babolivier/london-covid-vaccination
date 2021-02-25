package miner

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/babolivier/london-covid-vaccination/common"
	"github.com/babolivier/london-covid-vaccination/storage"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

const (
	rootURL        = "https://www.england.nhs.uk/statistics/statistical-work-areas/covid-19-vaccinations/"
	linksTextStart = "COVID-19 daily announced vaccinations"
	dateLayout     = "2 January 2006"

	// Interval between two calls to Miner.mine.
	interval = 1 * time.Hour
)

var (
	// Files from before Jan 18th are using a different format that doesn't detail the
	// statistics per region.
	filesToIgnore = []string{
		"COVID-19-Daily-announced-vaccinations-17-January-2021-1.xlsx",
		"COVID-19-daily-announced-vaccinations-16-January-2021.xlsx",
		"COVID-19-daily-announced-vaccinations-15-January-2021.xlsx",
		"COVID-19-daily-announced-vaccinations-14-January-2021.xlsx",
		"COVID-19-daily-announced-vaccinations-13-January-2021.xlsx",
		"COVID-19-daily-announced-vaccinations-12-January-2021.xlsx",
		"COVID-19-daily-announced-vaccinations-11-January-2021.xlsx",
	}
)

// Miner periodically crawls the page on the NHS website listing daily reports, then
// downloads and processes the ones it doesn't know about.
type Miner struct {
	db *storage.Database
}

// NewMiner returns a new instance of Miner.
func NewMiner(db *storage.Database) *Miner {
	return &Miner{db: db}
}

// Start starts the Miner in a separate goroutine.
func (m *Miner) Start() {
	go m.run()
}

// run runs the Miner in a loop, with a delay of a given duration between each iteration.
func (m *Miner) run() {
	var err error
	for {
		logrus.Info("Starting mining")
		if err = m.mine(); err != nil {
			logrus.WithError(err).Error("Mining failed")
		} else {
			logrus.Info("Mining succeeded")
		}
		time.Sleep(interval)
	}
}

// mine crawls the page on the NHS website listing daily reports, then
// downloads and processes the ones it doesn't know about.
func (m *Miner) mine() error {
	// Retrieve known files from the database.
	knownFileNames, err := m.db.GetKnownFileNames()
	if err != nil {
		return err
	}

	// Add the names of the files we need to ignore to the list of known files - that way
	// they'll get ignored as well.
	for _, fileName := range filesToIgnore {
		knownFileNames[fileName] = true
	}

	// Crawl through the page and retrieve links we're interested in.
	urls := make([]string, 0)
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Text, linksTextStart) {
			urls = append(urls, e.Attr("href"))
		}
	})
	if err = c.Visit(rootURL); err != nil {
		return err
	}

	// For each link, check if we've already processed the file (in which case skip to
	// the next file), otherwise download it, parse it and extract the data about London
	// from it.
	var dailyStats *common.DailyStats
	for _, url := range urls {
		dailyStats = new(common.DailyStats)

		_, dailyStats.FileName = filepath.Split(url)

		if _, exists := knownFileNames[dailyStats.FileName]; exists {
			continue
		}

		logrus.Infof("Processing %s\n", dailyStats.FileName)

		var data [][][]string
		data, err = parseXlsxFile(url)
		if err != nil {
			return err
		}

		// The date is in cell C7.
		dailyStats.PubDate, err = parseDate(data[0][6][2])
		if err != nil {
			return err
		}

		// The number of first doses in London is in cell D13.
		dailyStats.FirstDose, err = strconv.Atoi(data[0][15][3])
		if err != nil {
			return err
		}
		// The number of second doses in London is in cell D14.
		dailyStats.SecondDose, err = strconv.Atoi(data[0][15][4])
		if err != nil {
			return err
		}

		// Store the stats we've just extracted.
		if err = m.db.StoreStats(dailyStats); err != nil {
			return err
		}
	}

	return nil
}

// parseXlsxFile downloads and parses the XLSX file.
func parseXlsxFile(url string) ([][][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	xlsxFile, err := xlsx.OpenBinary(body)
	if err != nil {
		return nil, err
	}

	slice, err := xlsxFile.ToSlice()
	if err != nil {
		return nil, err
	}

	return slice, err
}

// parseDate parses the string from the cell the publication date is located in.
func parseDate(rawDate string) (time.Time, error) {
	// We need to do some extra processing because most files use the format "25th" for
	// the day of the month and time.Parse doesn't understand that.
	parts := strings.Split(rawDate, " ")
	day := parts[0]
	lastDayChar := day[len(day)-1]

	// Some files (e.g. Jan 18th) don't use the suffix for the day of the month, so be a
	// bit flexible about whether to correct it.
	if lastDayChar < '0' || lastDayChar > '9' {
		day = day[:len(day)-2]
		// Correct the day of the month.
		rawDate = strings.Replace(rawDate, parts[0], day, 1)
	}

	return time.Parse(dateLayout, rawDate)
}
