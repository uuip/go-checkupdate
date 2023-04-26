package rule

import (
	"os"
	"strings"

	"checkupdate/models"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

var TOKEN = os.Getenv("GITHUB_TOKEN")

const UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:112.0) Gecko/20100101 Firefox/112.0"

func FetchApp(app *models.VerModel) (string, error) {
	client := resty.New().SetHeader("user-agent", UA)
	if app.Name == "Fences" {
		resp, err := client.R().Head(app.Url)
		if err != nil {
			return "", err
		}
		return resp.Header().Get("Content-Length"), nil
	} else if app.Name == "EmEditor" {
		client.SetRedirectPolicy(resty.NoRedirectPolicy())
		resp, err := client.R().Get(app.Url)
		if strings.Contains(err.Error(), "auto redirect is disabled") {
			arg := resp.Header().Get("location")
			s := strings.Split(arg, "_")
			return numVersion(s[len(s)-1]), nil
		} else if err != nil {
			return "", err
		}
		return "", nil

	} else if app.Json == 1 {
		var resp *resty.Response
		var err error
		var rst gjson.Result

		if strings.HasPrefix(app.Url, "https://api.github.com") {
			resp, err = client.SetHeader("Authorization", "token "+TOKEN).R().Get(app.Url)
		} else {
			resp, err = client.R().Get(app.Url)
		}
		if err != nil {
			return "", err
		}

		switch app.Name {
		case "PyCharm":
			{
				rst = gjson.GetBytes(resp.Body(), "PCP.0.version")
			}
		case "Clash":
			{
				rst = gjson.GetBytes(resp.Body(), "name")
			}
		default:
			{
				rst = gjson.GetBytes(resp.Body(), "tag_name")
			}
		}
		return rst.String(), nil
	} else {
		resp, err := client.R().Get(app.Url)
		if err != nil {
			return "", err
		}
		return findRuleFn(app, resp)
	}
}

func findRuleFn(app *models.VerModel, resp *resty.Response) (string, error) {
	if fn, ok := fnRules[app.Name]; ok {
		return fn(resp)
	}
	return parseByCss(app, resp)
}
