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

var raw = "/raw/"
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
		return nil, errors.New("use fallback ansi")
	}

	return raws, nil
}

func getPacks(page int) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/pack?pagesize=%d&page=%d", apiBase, version, pagesize, page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	r, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func getPack(ran int, c *colorsPacks) ([]byte, error) {
	url := fmt.Sprintf("%s/%s%s?pagesize=%d", apiBase, version, c.Results[ran].Gallery, pagesize)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	r, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	fmt.Println("Loading random ANSi from 16colo.rs")

	// TODO: add support for ansi fallback
	// TODO: add support for getting multiple pages of results
	page := 0
	r, err := getPacks(page)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	var c colorsPacks
	json.Unmarshal(r, &c)

	rand.Seed(time.Now().UnixMicro())
	ran := rand.Intn(len(c.Results))

	r2, err := getPack(ran, &c)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	var cp colorsPack
	json.Unmarshal(r2, &cp)

	raws, err := getRaws(cp)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	ran1 := rand.Intn(len(raws))

	url := fmt.Sprintf("%s%s%s%s", urlBase, c.Results[ran].Gallery, raw, raws[ran1])

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	defer res.Body.Close()

	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}
	err = os.WriteFile("ansi.ans", body, 0777)
	if err != nil {
		fmt.Println("use fallback ansi")
		os.Exit(0)
	}
}
