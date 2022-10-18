package webpage

import (
	"autify/v1/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

var metaDataFilename = "meta.json"

// DownloadPages downloads web pages from the internet, save the HTML file to disk and save the web page metadata to a json file.
// If the page already exist then it updates the last fetch time.
func DownloadPages(urls []string) error {
	for _, url := range urls {
		filename := path.Base(url) + ".html"
		// If page does not exist, download
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			err := download(url, filename)
			if err != nil {
				return fmt.Errorf("download error : %w", err)
			}
			// Parse HTML file to generate the page metadata
			pMeta, err := parseFile(filename)
			if err != nil {
				return fmt.Errorf("page parse error: %w", err)
			}
			pMeta.Site = url
			pMeta.Last_Fetch = getCurTime()
			// Save webpage  metadata to disk
			err = saveMetaToDisk(*pMeta)
			if err != nil {
				return fmt.Errorf("unable to save metadata to disk: %w", err)
			}
		}
		// Update page last fect time
		err := updateLastFetchTime(url)
		if err != nil {
			return fmt.Errorf("update failed: %w", err)
		}
	}
	return nil
}

// GetPageMeta fetchs page metadata from the  meta.json file and returns it if it exist otherwise return not found error
func GetPageMeta(url string) (models.WebPage, error) {
	byteValue, err := ioutil.ReadFile(metaDataFilename)
	if err != nil {
		return models.WebPage{}, err
	}

	var pages map[string]models.WebPage
	err = json.Unmarshal(byteValue, &pages)
	if err != nil {
		return models.WebPage{}, err
	}
	if _, ok := pages[url]; !ok {
		return models.WebPage{}, fmt.Errorf("page not found, download page first then try again")
	}
	return pages[url], nil
}

// Setup creates meta.json file if it doesn't already exist, this file will store the web pages metadata
func Setup() error {
	if _, err := os.Stat(metaDataFilename); os.IsNotExist(err) {
		_, err := os.Create(metaDataFilename)
		if err != nil {
			return err
		}
		byteValue, err := json.Marshal(map[string]*models.WebPage{})
		if err != nil {
			return err
		}
		// Write  empty map in form of byte to json file
		err = ioutil.WriteFile(metaDataFilename, byteValue, 0644)
		return err
	}
	return nil
}

func download(url, filename string) error {
	color.Blue("Downloading: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	// save page to disk
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// HTML file parser
func parseFile(filename string) (*models.WebPage, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file error : %w", err)
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		return nil, err
	}
	imgCount := 0
	linkCount := 0
	var parse_html func(n *html.Node)
	parse_html = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, element := range n.Attr {
				// ensure link tag have a hypertext reference before incrementing linkCount
				if element.Key == "href" {
					linkCount++
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, element := range n.Attr {
				// ensure image tag have  a source before incrementing imgCount
				if element.Key == "src" {
					imgCount++
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse_html(c)
		}
	}
	parse_html(doc)
	pMeta := &models.WebPage{
		Num_Links: linkCount,
		Images:    imgCount,
	}
	return pMeta, nil
}

func getCurTime() string {
	currentTime := time.Now()
	return currentTime.Format("Mon Jan 01 2006 3:4:5 UTC")
}

func updateLastFetchTime(url string) error {
	byteValue, err := ioutil.ReadFile(metaDataFilename)
	if err != nil {
		return err
	}

	var pages map[string]models.WebPage
	err = json.Unmarshal(byteValue, &pages)
	if err != nil {
		return err
	}
	if pMeta, ok := pages[url]; ok {
		// update last fetch time
		pMeta.Last_Fetch = getCurTime()
		// reset entry
		pages[url] = pMeta
	}
	byteValue, err = json.Marshal(pages)
	if err != nil {
		return err
	}
	// Write back to file
	err = ioutil.WriteFile(metaDataFilename, byteValue, 0644)
	return err
}

// save web page metadata to disk as JSON
func saveMetaToDisk(meta models.WebPage) error {

	byteValue, err := ioutil.ReadFile(metaDataFilename)
	if err != nil {
		return err
	}
	
	var pages map[string]models.WebPage
	err = json.Unmarshal(byteValue, &pages)
	if err != nil {
		return err
	}
	pages[meta.Site] = meta
	byteValue, err = json.Marshal(pages)
	if err != nil {
		return err
	}
	// Write back to file
	err = ioutil.WriteFile(metaDataFilename, byteValue, 0644)
	return err
}
