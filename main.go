package main

import (
	"io/ioutil"
	"net/http"
  "strings"
  "strconv"
  "fmt"
)
type Party int
const (
  CDU = Party(0)
  SPD = Party(1)
  FDP = Party(2)
  B90 = Party(3)
  Linke = Party(4)
)
var parties = []string{"CDU/CSU","SPD","FDP","B90/Grüne","DieLinke"}

func main() {
	url := "https://www.bundestag.de/parlament/plenum/abstimmung/abstimmung?id="
	
  //max := 683
  for id := 1; id < 3; id++ {
    noSpace := GetHTML(fmt.Sprintf("%s%v",url,id))
    fmt.Printf("Date: %v\n", GetDateDDMMJJJJ(GetDate(noSpace)))
    for _,pt := range(parties) {
      fmt.Printf("%s Stats: %v\n", pt, GetPartyStatsInt(noSpace, pt))
    }
    fmt.Println()
  }
}
//Ja,Nein,Enthalten,Nicht Abgegeben
func GetPartyStatsInt(html, party string) [4]int {
  party_html := GetParty(party, html)
  stats_str := GetPartyStats(party_html)
  stats := strings.Split(stats_str, ",")
  ints := [4]int{}
  for i,s := range(stats) {
    ints[i],_ = strconv.Atoi(s)
  }
  return ints
}
const statsBgn = "data-chart-values=\""
const statsEnd = "\"data-chart-type"
func GetPartyStats(party_html string) string {
  bgn := strings.Index(party_html, statsBgn)
  end := strings.Index(party_html, statsEnd)
  return party_html[bgn+len(statsBgn):end]
}
const PartyBgn = "\"bt-chart-fraktion\">"
func GetParty(party, html string) string {
  idx := strings.Index(html, fmt.Sprintf("%s%s", PartyBgn, party))
  return html[idx:idx+1000]
}

func GetDateDDMMJJJJ(date string) (d [3]int) {
  dl := strings.Split(date, ".")
  d[0],_ = strconv.Atoi(dl[0])
  d[1],_ = strconv.Atoi(dl[1])
  d[2],_ = strconv.Atoi(dl[2])
  return
}

const DateBgn = "Sitzungvom"
func GetDate(html string) string {
  bgn := strings.Index(html, DateBgn)+len(DateBgn)
  return html[bgn:bgn+10]
}
const InfoBgn = "bt-dachzeile\">"
func GetInfoBody(html string) string {
  didx := strings.Index(html, InfoBgn)
  return html[didx:didx+1000]
}


func GetHTML(url string) string {
  resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
  return strings.ReplaceAll(strings.ReplaceAll(string(html), " ", ""), "\n", "")
}
