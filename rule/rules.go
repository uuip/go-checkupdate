package rule

import (
	"bytes"
	"regexp"

	. "checkupdate/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/coreos/go-semver/semver"
	"github.com/go-resty/resty/v2"
)

var numRe, _ = regexp.Compile(`[.\d]*\d`)

func parseByCss(app *VerModel, resp *resty.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	return numVersion(doc.Find(cssRules[app.Name]).First().Text()), nil
}

func parseDevManView(resp *resty.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	element := doc.Find("h4").FilterFunction(func(i int, selection *goquery.Selection) bool {
		return selection.Text() == "Versions History"
	}).Next().Find("li:nth-child(1)").First().Text()
	return numVersion(element), nil
}

func parseVmware(resp *resty.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	var vs []*semver.Version
	doc.Find("metadata>version").Each(func(i int, selection *goquery.Selection) {
		vs = append(vs, semver.New(selection.Text()))
	})
	semver.Sort(vs)
	return numVersion(vs[len(vs)-1].String()), nil
}

func parseWinrar(resp *resty.Response) (string, error) {
	re, _ := regexp.Compile("^WinRAR.*elease")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	element := doc.Find("b").FilterFunction(func(i int, selection *goquery.Selection) bool {
		return re.MatchString(selection.Text())
	}).First().Text()
	return numVersion(element), nil
}

func parseFaststone(resp *resty.Response) (string, error) {
	re, _ := regexp.Compile(`Version\s*[.\d]+`)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	element := doc.Find("b").FilterFunction(func(i int, selection *goquery.Selection) bool {
		return re.MatchString(selection.Text())
	}).First().Text()
	return numVersion(element), nil
}

func parseBeyondCompare(resp *resty.Response) (string, error) {
	re, _ := regexp.Compile("Current Version.+")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	element := doc.Find("p").FilterFunction(func(i int, selection *goquery.Selection) bool {
		return re.MatchString(selection.Text())
	}).First().Text()
	return numVersion(element), err

}

func parsePython(resp *resty.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "", err
	}
	element := doc.Find("p.download-buttons>a").First().Text()
	return numVersion(element), nil
}

func numVersion(str string) string {
	return numRe.FindString(str)
}
