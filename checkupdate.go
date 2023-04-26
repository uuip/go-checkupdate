package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"checkupdate/models"
	"checkupdate/rule"
	"github.com/fatih/color"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var green = color.HiGreenString
var red = color.HiRedString

func main() {
	t1 := time.Now()

	status := make(map[string][]string)
	status["success"] = []string{}
	status["failed"] = []string{}

	var dsn string
	if runtime.GOOS == "darwin" {
		dsn = "/Users/sharp/Downloads/ver_tab.db"
	} else {
		dsn = `c:/Users/sharp/AppData/Local/Programs/checkupdate/ver_tab.db`
	}
	db, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	var (
		app  models.VerModel
		apps []models.VerModel
	)
	db.Where("name=?", "fzf").Take(&app)
	aj, _ := json.Marshal(app)
	fmt.Println(string(aj))

	var wg sync.WaitGroup
	db.Find(&apps)
	wg.Add(len(apps))
	ch := make(chan [2]string, 10)
	for _, item := range apps {
		item := item
		go func() {
			defer wg.Done()
			newVer, err := rule.FetchApp(&item)
			if err != nil || newVer == "" {
				ch <- [2]string{item.Name, ""}
				if err == nil {
					fmt.Printf(red("%s failed\n"), item.Name)
				} else {
					fmt.Printf(red("%s failed\n %s\n"), item.Name, err)
				}
				return
			}
			if newVer != item.Ver {
				ch <- [2]string{item.Name, newVer}
				fmt.Println(item.Name, green(newVer))
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for val := range ch {
		appname, newver := val[0], val[1]
		if newver == "" {
			status["failed"] = append(status["failed"], appname)
		} else {
			status["success"] = append(status["success"], appname)
			db.Model(&models.VerModel{}).Where("name=?", appname).Update("Ver", newver)
		}
	}
	for k, v := range status {
		fmt.Println(k, ": ", strings.Join(v, ", "))
	}
	fmt.Printf("用时 %.2f 秒\n", time.Since(t1).Seconds())
	if runtime.GOOS == "windows" {
		_, _ = fmt.Scanln()
	}

	//fmt.Println(strconv.FormatInt(55, 10))

	//var arr = []int{1, 2, 3}
	//scoreMap := make(map[string]int)

	//s := "abcd"
	//for i, n := 0, len(s); i < n; i++ {
	//	fmt.Println(n)
	//	println(i, string(s[i]))
	//}
}
