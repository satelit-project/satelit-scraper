package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/jsonpb"

	"shitty.moe/satelit-project/satelit-scraper/proto/data"
	"shitty.moe/satelit-project/satelit-scraper/spider/anidb"
)

const userAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 " +
	"(KHTML, like Gecko) Version/13.0.2 Safari/605.1.15"

type arguments struct {
	anidbID int
	outPath string
}

func main() {
	args := parseArguments()

	log.Println("Downloading anime with id", args.anidbID)
	html, err := getAniDBTitle(args.anidbID)
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	tee := io.TeeReader(html, &buf)

	log.Println("Parsing HTML page")
	anime, err := parseAniDBTitle(args.anidbID, tee)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Saving data to", args.outPath)
	err = saveAniDBTitle(args.anidbID, args.outPath, anime, &buf)
	if err != nil {
		log.Fatalln(err)
	}
}

func parseArguments() arguments {
	var args arguments
	flag.IntVar(&args.anidbID, "anidb-id", 0, "Anime ID in AniDB")
	flag.StringVar(&args.outPath, "out-path", "", "Path to put HTML and parsed JSON")
	flag.Parse()

	if args.anidbID <= 0 {
		log.Fatalln("Invalid AniDB ID")
	}

	pathStat, err := os.Stat(args.outPath)
	if err != nil {
		log.Fatalln(err)
	}

	if !pathStat.IsDir() {
		log.Fatalln("Output path is not a directory")
	}

	return args
}

func getAniDBTitle(id int) (io.Reader, error) {
	titleURL := urlForAniDBTitle(id)
	req, _ := http.NewRequest("GET", titleURL, nil)
	req.Header.Set("User-Agent", userAgent)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	html, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(html), nil
}

func parseAniDBTitle(id int, html io.Reader) (*data.Anime, error) {
	htmlURL, err := url.Parse(urlForAniDBTitle(id))
	if err != nil {
		return nil, err
	}

	parser, err := anidb.NewParser(htmlURL, html)
	if err != nil {
		return nil, err
	}

	return parser.Anime()
}

func saveAniDBTitle(id int, path string, anime *data.Anime, html io.Reader) error {
	encoder := jsonpb.Marshaler{}
	json, err := encoder.MarshalToString(anime)
	if err != nil {
		return err
	}

	jsonName := filepath.Join(path, fmt.Sprintf("%d.json", id))
	htmlName := filepath.Join(path, fmt.Sprintf("%d.html", id))

	err = ioutil.WriteFile(jsonName, []byte(json), 0644)
	if err != nil {
		return err
	}

	fullHTML, err := ioutil.ReadAll(html)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(htmlName, fullHTML, 0644)
}

func urlForAniDBTitle(id int) string {
	return fmt.Sprintf("https://anidb.net/anime/%d", id)
}
