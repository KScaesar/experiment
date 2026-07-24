package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	records, err := LoadFile("/Users/caesar.tsai/asdv0402.log")
	if err != nil {
		panic(err)
	}

	adIds := map[string]bool{
		"7471": true,
		"7369": true,
	}

	adUnits := map[string]bool{
		"adsv_litv_tv_livech_ftvn_ad1":       true,
		"adsv_litv_tv_livech_tvbsn_ad1":      true,
		"adsv_litv_appleTV_livech_tvbsn_ad1": true,
		"adsv_litv_webtv_livech_tvbsn_ad1":   true,
	}

	_, hit, views := SearchTarget(records, adIds, adUnits)
	for _, view := range views {
		fmt.Println(view)
	}

	report := `
ad_unit: 
%v

調查數量:
%v

第一順位委刊項目標數量:
%v


`
	fmt.Printf(report, ViewAdUnit(adUnits), len(records), hit)
}

func ViewAdUnit(data map[string]bool) string {
	r := ""
	for k := range data {
		r += k
		r += "\n"
	}
	return fmt.Sprintf(r[:len(r)-1])
}

func SearchTarget(records []Record, adIds map[string]bool, adUnits map[string]bool) (total, hit int, views []View) {
	for _, record := range records {
		if !adUnits[record.AdUnit] {
			continue
		}

		view := View{
			Datatime: record.Datatime,
			Puid:     record.PuId,
			AdUnit:   record.AdUnit,
		}

		if len(record.ShowList) > 0 {
			show := record.ShowList[0]
			view.FirstShow = show
			if adIds[show.AdId] {
				hit++
			}
		} else {
			view.FirstShow = Show{
				AdId:       "empty",
				CreativeId: "empty",
			}
		}

		for _, non := range record.NonList {
			if adIds[non.AdId] {
				view.NonList = append(view.NonList, non)
			}
		}
		// sort.Slice(view.NonList, func(i, j int) bool {
		// 	return view.NonList[i].AdId < view.NonList[j].AdId
		// })

		views = append(views, view)
	}
	return len(views), hit, views
}

type View struct {
	AdUnit    string
	Puid      string
	Datatime  string
	FirstShow Show
	NonList   []NonShow
}

func (v View) String() string {
	var nonList []string
	for _, non := range v.NonList {
		nonList = append(nonList, fmt.Sprintf("%v:%v", non.AdId, non.Reason))
	}
	return fmt.Sprintf(`%v	%v	%v	%v	%v	%v`,
		v.AdUnit,
		v.Puid,
		v.Datatime,
		v.FirstShow.AdId,
		v.FirstShow.CreativeId,
		strings.Join(nonList, "\t"),
	)
}

//

func LoadFile(filename string) (records []Record, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		record := ParseRecord(line)
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func ParseRecord(rawData []byte) Record {
	parts := bytes.Split(rawData, []byte("\t"))

	var record Record
	err := json.Unmarshal(parts[3], &record)
	if err != nil {
		panic(err)
	}

	record.Datatime = string(parts[0])
	return record
}

type Record struct {
	Datatime string
	Issue    string    `json:"issue"`
	PuId     string    `json:"puid"`
	AdUnit   string    `json:"ad_unit"`
	ShowList []Show    `json:"show"`
	NonList  []NonShow `json:"non"`
}

type Show struct {
	AdId       string
	CreativeId string
}

func (s *Show) UnmarshalJSON(data []byte) error {
	raw := bytes.Trim(data, `"`)
	parts := bytes.Split(raw, []byte(":"))
	if len(parts) != 2 {
		return errors.New("invalid Show format, expected 'AdId:CreativeId'")
	}
	s.AdId = string(parts[0])
	s.CreativeId = string(parts[1])
	return nil
}

type NonShow struct {
	AdId   string
	Reason string
}

func (s *NonShow) UnmarshalJSON(data []byte) error {
	raw := bytes.Trim(data, `"`)
	parts := bytes.Split(raw, []byte(":"))
	if len(parts) != 2 {
		return errors.New("invalid NonShow format, expected 'AdId:Reason'")
	}
	s.AdId = string(parts[0])
	s.Reason = string(parts[1])
	return nil
}
