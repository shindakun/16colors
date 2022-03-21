package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"
)

var pagesize = 500
var apiBase = "https://api.16colo.rs"
var urlBase = "https://16colo.rs"

// var raw = "/raw/"
var version = "v1"

type colorsPacks struct {
	Page struct {
		Total    int    `json:"total"`
		Sort     string `json:"sort"`
		Order    string `json:"order"`
		Pagesize int    `json:"pagesize"`
		Page     int    `json:"page"`
		Pages    int    `json:"pages"`
		Offset   int    `json:"offset"`
	} `json:"page"`
	Results []struct {
		Gallery string `json:"gallery"`
	} `json:"results"`
}

type colorsPack struct {
	Page struct {
	} `json:"page"`
	Results []struct {
		Files interface{} `json:"files"`
	} `json:"results"`
}

func getRaws(cp colorsPack) ([]string, error) {
	val := reflect.ValueOf(cp.Results[0].Files)

	// fmt.Println("VALUE = ", val)
	// fmt.Println("KIND = ", val.Kind())

	var raws []string

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			switch t := v.Interface().(type) {
			case map[string]interface{}:
				f := t["file"]
				r, _ := regexp.Compile(".asc|.ASC|.ans|.ANS")
				for a, b := range f.(map[string]interface{}) {
					if a == "raw" {

						match := r.MatchString(b.(string))
						if match {
							raws = append(raws, b.(string))
						}
					}
				}
			default:
				fmt.Println("not found")
			}
		}
	}
	if len(raws) == 0 {
		return nil, errors.New("use fallback.ans")
	}

	return raws, nil
}

func getPacks() {

}

func main() {
	fmt.Println("16colorsapi")

	url := fmt.Sprintf("%s/%s/pack?pagesize=%d", apiBase, version, pagesize)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	r, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var c colorsPacks
	json.Unmarshal(r, &c)

	rand.Seed(time.Now().UnixMicro())
	ran := rand.Intn(pagesize)

	url2 := fmt.Sprintf("%s/%s%s?pagesize=%d", apiBase, version, c.Results[ran].Gallery, pagesize)

	req2, err := http.NewRequest(http.MethodGet, url2, nil)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	defer res.Body.Close()
	r2, err := io.ReadAll(res2.Body)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	var cp colorsPack
	json.Unmarshal(r2, &cp)

	raws, err := getRaws(cp)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	ran1 := rand.Intn(len(raws))

	url3 := fmt.Sprintf("%s%s/%s/%s", urlBase, c.Results[ran].Gallery, "raw", raws[ran1])

	fmt.Println(url3)

	req3, err := http.NewRequest(http.MethodGet, url3, nil)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	res3, err := http.DefaultClient.Do(req3)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	defer res3.Body.Close()

	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}

	body, err := io.ReadAll(res3.Body)
	if err != nil {
		fmt.Println("use fallback.ans")
		os.Exit(0)
	}
	os.WriteFile("ansi.ans", body, 0777)
}
