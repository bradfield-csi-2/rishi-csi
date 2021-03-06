package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	baseURL         = "https://xkcd.com"
	jsonPath        = "info.0.json"
	dataDir         = "data/"
	defaultComicNum = 2423
)

type Comic struct {
	Num        int
	Title      string
	Transcript string
}

type resultSet map[int]bool

type searchIndex map[string]resultSet

func main() {
	downloadFlag := flag.Bool("d", false, "Download comics data")

	flag.Parse()

	if *downloadFlag {
		downloadAllComics(dataDir)
	}

	index := buildSearchIndex(dataDir)

	// Grab the first non-flag argument
	term := flag.Arg(0)
	if term == "" {
		fmt.Printf("No search term provided.\n")
		flag.PrintDefaults()
		return
	}

	results := search(term, index)
	printSearchResults(results, term, dataDir)
}

func downloadAllComics(dir string) {
	for i := getMaxComicNum(); i > 0; i-- {
		// This is an Easter Egg -- there is no Comic #404
		if i == 404 {
			continue
		}
		comicNum := strconv.Itoa(i)
		saveLoc := dir + comicNum + ".json"

		// Download a comic only if it is not already on disk
		if _, err := os.Stat(saveLoc); os.IsNotExist(err) {
			if err = downloadComic(comicUrl(comicNum), saveLoc); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func buildSearchIndex(dir string) searchIndex {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	index := make(searchIndex)
	for _, file := range files {
		comic := loadComicFromFile(dir, file.Name())
		searchableText := comic.Title + "\n" + comic.Transcript
		cleanedText := cleanText(searchableText)
		for _, word := range strings.Split(cleanedText, " ") {
			word = strings.TrimSpace(word)
			if index[word] == nil {
				index[word] = make(map[int]bool)
			}
			index[word][comic.Num] = true
		}
	}

	return index
}

func search(term string, index searchIndex) resultSet {
	cleanedTerm := cleanText(term)
	results, found := index[cleanedTerm]
	if !found {
		fmt.Printf("Search term: '%s' not found.\n", term)
	}

	return results
}

func getMaxComicNum() int {
	// The most recent comic is at: xkcd.com/info.0.json, so a blank string will get it for us
	body, err := requestComic(comicUrl(""))
	if err != nil {
		fmt.Printf("Could not get most recent comic. Defaulting to #%d\n", defaultComicNum)
		return defaultComicNum
	}

	var comic Comic
	if err = json.NewDecoder(body).Decode(&comic); err != nil {
		fmt.Printf("Error parsing JSON for most recent comic. Defaulting to #%d\n", defaultComicNum)
		return defaultComicNum
	}
	return comic.Num
}

func comicUrl(comicNum string) string {
	return strings.Join([]string{baseURL, comicNum, jsonPath}, "/")
}

func downloadComic(url string, saveLoc string) error {
	fmt.Printf("Downloading %s\n", url)
	data, err := requestComic(url)
	if err != nil {
		return err
	}

	err = saveComic(saveLoc, data)
	if err != nil {
		return fmt.Errorf("Error saving comic to %s: %s\n", saveLoc, err)
	}
	return nil
}

func requestComic(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error downloading %s: %s\n", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("Error downloading %s: %s\n", url, resp.Status)
	}

	return resp.Body, nil
}

func saveComic(location string, data io.ReadCloser) error {
	out, err := os.Create(location)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, data)
	return err
}

func loadComicFromFile(dir string, fileName string) *Comic {
	data, err := ioutil.ReadFile(dir + fileName)
	if err != nil {
		log.Fatalf("Failed to open %v", err)
	}

	var comic Comic
	if err := json.Unmarshal(data, &comic); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return &comic
}

func cleanText(text string) string {
	text = strings.ToLower(text)

	// There's always a "Title text:" description line
	text = strings.ReplaceAll(text, "title text:", "")

	// Remove all non-alphanumeric characters
	re := regexp.MustCompile(`[^a-zA-Z\d\s]`)
	text = re.ReplaceAllLiteralString(text, "")

	// Replace mutliple whitespace with single space
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllLiteralString(text, " ")

	return text
}

func printSearchResults(results resultSet, term string, dir string) {
	resultQuantifier := "result"
	numResults := len(results)
	if numResults != 1 {
		resultQuantifier += "s"
	}
	fmt.Printf("%d %s for '%s'\n", len(results), resultQuantifier, term)

	for num, _ := range results {
		comicNum := strconv.Itoa(num)
		comic := loadComicFromFile(dir, comicNum+".json")
		url := comicUrl(comicNum)
		padding := fmt.Sprintf("%*s", len(url), "=")
		fmt.Printf("\n%s\n%s\n%s\n", url, strings.ReplaceAll(padding, " ", "="), comic.Transcript)
	}
}
