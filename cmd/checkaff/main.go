package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AffData struct {
	Intervals []struct {
		Count     string `json:"count"`
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
		Thornames []struct {
			Count     string `json:"count"`
			Thorname  string `json:"thorname"`
			Volume    string `json:"volume"`
			VolumeUSD string `json:"volumeUSD"`
		} `json:"thornames"`
		Volume    string `json:"volume"`
		VolumeUSD string `json:"volumeUSD"`
	} `json:"intervals"`
	Meta struct {
		Count     string `json:"count"`
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
		Volume    string `json:"volume"`
		VolumeUSD string `json:"volumeUSD"`
	} `json:"meta"`
}

func main() {
	// split affResult by each liuns
	affResults := strings.Split(affResult, "\n")
	for _, v := range affResults {
		v = strings.Split(v, ", ")[0]
		v = strings.ReplaceAll(v, "Affiliate:", "")
		v = strings.TrimSpace(v)
		v = fmt.Sprintf("    - \"%s\"", v)
		fmt.Println(v)
	}
	return
	// send http request to
	bt := []byte(affJson)
	var affData AffData
	if err := json.Unmarshal(bt, &affData); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	vultisigAffilaites := []string{
		"va",
		"vi",
		"v0",
	}
	uniqAffilaites := make(map[string]struct{})
	for _, interval := range affData.Intervals {
		for _, thorname := range interval.Thornames {
			uniqAffilaites[thorname.Thorname] = struct{}{}
		}
	}
	fmt.Println(len(uniqAffilaites))
	for k := range uniqAffilaites {
		for _, v := range vultisigAffilaites {
			//time.Sleep(1 * time.Second)
			newAff := fmt.Sprintf("%s/%s", k, v)
			count, err := getCount(newAff)
			if err != nil {
				fmt.Printf("Error getting count for %s: %v\n", newAff, err)
				continue
			}
			if count > 0 {
				fmt.Printf("Affiliate: %s, Count: %d\n", newAff, count)
			}

		}
	}
}

func getCount(affiliate string) (int, error) {
	url := fmt.Sprintf("https://vanaheimex.com/actions?affiliate=%s&limit=1&type=swap", affiliate)
	httpResp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data for %s: %v\n", affiliate, err)
		return 0, err
	}
	defer httpResp.Body.Close()

	//print body
	bt, err := io.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for %s: %v\n", affiliate, err)
		return 0, err
	}
	var result ActionData
	if err := json.Unmarshal(bt, &result); err != nil {
		fmt.Printf("Error decoding response for %s: %v\n", affiliate, err)
		fmt.Println("Response body:", string(bt))
		return 0, err
	}
	return result.Count, nil
}

/*
- "va"
- "vi"
- "v0"
*/
type ActionData struct {
	Count int `json:"count,string"`
}

var affResult string = `
Affiliate: GRIF/vi, Count: 2
Affiliate: jp/vi, Count: 3
Affiliate: TARD/vi, Count: 3
Affiliate: DOOM/vi, Count: 10
Affiliate: HOGA/vi, Count: 217
Affiliate: QBIT/vi, Count: 2
Affiliate: VALT/vi, Count: 8
Affiliate: DROP/vi, Count: 5
Affiliate: ZENO/vi, Count: 11
`

var affJson string = `{
	"intervals": [
		{
			"count": "386716",
			"endTime": "1754338613",
			"startTime": "1647913096",
			"thornames": [
				{
					"count": "413",
					"thorname": "-",
					"volume": "195249662689",
					"volumeUSD": "300992"
				},
				{
					"count": "48707",
					"thorname": "-_",
					"volume": "65075225498652",
					"volumeUSD": "102639416"
				},
				{
					"count": "18",
					"thorname": "-t",
					"volume": "40717915489",
					"volumeUSD": "57646"
				},
				{
					"count": "2",
					"thorname": "0xGM",
					"volume": "16600479",
					"volumeUSD": "26"
				},
				{
					"count": "1",
					"thorname": "77tr777",
					"volume": "7256600784",
					"volumeUSD": "9960"
				},
				{
					"count": "1",
					"thorname": "A1",
					"volume": "272928",
					"volumeUSD": "0"
				},
				{
					"count": "1",
					"thorname": "B1",
					"volume": "272928",
					"volumeUSD": "0"
				},
				{
					"count": "135",
					"thorname": "CS",
					"volume": "317932143778",
					"volumeUSD": "476132"
				},
				{
					"count": "9",
					"thorname": "DOOM",
					"volume": "223784509",
					"volumeUSD": "351"
				},
				{
					"count": "6",
					"thorname": "DROP",
					"volume": "673579307",
					"volumeUSD": "959"
				},
				{
					"count": "2",
					"thorname": "GRIF",
					"volume": "184855718",
					"volumeUSD": "284"
				},
				{
					"count": "215",
					"thorname": "HOGA",
					"volume": "42500615533",
					"volumeUSD": "62820"
				},
				{
					"count": "51",
					"thorname": "KNS",
					"volume": "17585903468",
					"volumeUSD": "28092"
				},
				{
					"count": "535",
					"thorname": "OKY",
					"volume": "1999189674655",
					"volumeUSD": "3474916"
				},
				{
					"count": "137",
					"thorname": "Origamix",
					"volume": "3621359275",
					"volumeUSD": "4542"
				},
				{
					"count": "1",
					"thorname": "Pete",
					"volume": "197900552",
					"volumeUSD": "525"
				},
				{
					"count": "2",
					"thorname": "QBIT",
					"volume": "184062752",
					"volumeUSD": "301"
				},
				{
					"count": "2",
					"thorname": "QRYPTON",
					"volume": "782551317",
					"volumeUSD": "1010"
				},
				{
					"count": "2",
					"thorname": "RUST",
					"volume": "1076580678",
					"volumeUSD": "1523"
				},
				{
					"count": "7",
					"thorname": "SALE",
					"volume": "5265964779",
					"volumeUSD": "8025"
				},
				{
					"count": "5",
					"thorname": "TARD",
					"volume": "925577841",
					"volumeUSD": "1357"
				},
				{
					"count": "25",
					"thorname": "VALT",
					"volume": "2636180407",
					"volumeUSD": "3860"
				},
				{
					"count": "3",
					"thorname": "ViGay",
					"volume": "1499585604",
					"volumeUSD": "1963"
				},
				{
					"count": "2872",
					"thorname": "XDF",
					"volume": "6920233603981",
					"volumeUSD": "10810933"
				},
				{
					"count": "11",
					"thorname": "ZENO",
					"volume": "2862472145",
					"volumeUSD": "4489"
				},
				{
					"count": "17",
					"thorname": "_-",
					"volume": "3887730158",
					"volumeUSD": "6418"
				},
				{
					"count": "4",
					"thorname": "abtc",
					"volume": "1466201476",
					"volumeUSD": "2486"
				},
				{
					"count": "371",
					"thorname": "ahi",
					"volume": "67405432751",
					"volumeUSD": "99351"
				},
				{
					"count": "13",
					"thorname": "be",
					"volume": "245030184",
					"volumeUSD": "427"
				},
				{
					"count": "10304",
					"thorname": "bgw",
					"volume": "35364444588473",
					"volumeUSD": "57242771"
				},
				{
					"count": "3",
					"thorname": "bit",
					"volume": "41825616",
					"volumeUSD": "94"
				},
				{
					"count": "6",
					"thorname": "bqx",
					"volume": "1296851064",
					"volumeUSD": "2265"
				},
				{
					"count": "15",
					"thorname": "brt",
					"volume": "151401911",
					"volumeUSD": "231"
				},
				{
					"count": "3",
					"thorname": "burn",
					"volume": "156606804",
					"volumeUSD": "311"
				},
				{
					"count": "12",
					"thorname": "burrito",
					"volume": "26485972",
					"volumeUSD": "133"
				},
				{
					"count": "1",
					"thorname": "bwt",
					"volume": "1813551",
					"volumeUSD": "10"
				},
				{
					"count": "1",
					"thorname": "byb",
					"volume": "8075967",
					"volumeUSD": "11"
				},
				{
					"count": "28",
					"thorname": "cakewallet",
					"volume": "13192908945",
					"volumeUSD": "23303"
				},
				{
					"count": "52",
					"thorname": "commission",
					"volume": "868681206",
					"volumeUSD": "3166"
				},
				{
					"count": "840",
					"thorname": "ct",
					"volume": "453175350022",
					"volumeUSD": "868037"
				},
				{
					"count": "206",
					"thorname": "dcf",
					"volume": "9080222719",
					"volumeUSD": "20303"
				},
				{
					"count": "12",
					"thorname": "dh",
					"volume": "46188584",
					"volumeUSD": "82"
				},
				{
					"count": "4",
					"thorname": "dp",
					"volume": "311929198",
					"volumeUSD": "771"
				},
				{
					"count": "44",
					"thorname": "ds",
					"volume": "11513505462",
					"volumeUSD": "26842"
				},
				{
					"count": "7026",
					"thorname": "dx",
					"volume": "130969567902217",
					"volumeUSD": "181071636"
				},
				{
					"count": "13",
					"thorname": "dz",
					"volume": "355382068",
					"volumeUSD": "575"
				},
				{
					"count": "19",
					"thorname": "ecx",
					"volume": "16713407974",
					"volumeUSD": "88022"
				},
				{
					"count": "10656",
					"thorname": "ej",
					"volume": "8858746503940",
					"volumeUSD": "15064660"
				},
				{
					"count": "190",
					"thorname": "eld",
					"volume": "665470474674",
					"volumeUSD": "890171"
				},
				{
					"count": "1",
					"thorname": "esref",
					"volume": "6900963",
					"volumeUSD": "15"
				},
				{
					"count": "2",
					"thorname": "faizan",
					"volume": "74611546",
					"volumeUSD": "105"
				},
				{
					"count": "8",
					"thorname": "fs",
					"volume": "21156892293",
					"volumeUSD": "34391"
				},
				{
					"count": "22",
					"thorname": "fulwdvrt",
					"volume": "318011492",
					"volumeUSD": "605"
				},
				{
					"count": "782",
					"thorname": "g1",
					"volume": "5373392755514",
					"volumeUSD": "7464146"
				},
				{
					"count": "13",
					"thorname": "giddy",
					"volume": "29201763282",
					"volumeUSD": "39873"
				},
				{
					"count": "119",
					"thorname": "hrz",
					"volume": "28805167025",
					"volumeUSD": "42834"
				},
				{
					"count": "162",
					"thorname": "hrz_ios",
					"volume": "86130473149",
					"volumeUSD": "126954"
				},
				{
					"count": "20",
					"thorname": "hw",
					"volume": "283433140",
					"volumeUSD": "398"
				},
				{
					"count": "374",
					"thorname": "is",
					"volume": "806866059887",
					"volumeUSD": "1448731"
				},
				{
					"count": "1",
					"thorname": "iw",
					"volume": "19051098",
					"volumeUSD": "28"
				},
				{
					"count": "3",
					"thorname": "jp",
					"volume": "25119264872",
					"volumeUSD": "36544"
				},
				{
					"count": "24",
					"thorname": "jun",
					"volume": "348328486",
					"volumeUSD": "503"
				},
				{
					"count": "1",
					"thorname": "ke",
					"volume": "20218119",
					"volumeUSD": "26"
				},
				{
					"count": "17",
					"thorname": "krt",
					"volume": "2207902900",
					"volumeUSD": "3433"
				},
				{
					"count": "45",
					"thorname": "lends",
					"volume": "1279602951251",
					"volumeUSD": "2236712"
				},
				{
					"count": "1021",
					"thorname": "leo",
					"volume": "985304103539",
					"volumeUSD": "1476920"
				},
				{
					"count": "12639",
					"thorname": "ll",
					"volume": "173603847157175",
					"volumeUSD": "291438071"
				},
				{
					"count": "90",
					"thorname": "lref",
					"volume": "154680826298",
					"volumeUSD": "222263"
				},
				{
					"count": "4",
					"thorname": "mis",
					"volume": "57740927",
					"volumeUSD": "81"
				},
				{
					"count": "43",
					"thorname": "moca",
					"volume": "122587922400",
					"volumeUSD": "587688"
				},
				{
					"count": "3",
					"thorname": "odin",
					"volume": "359516132",
					"volumeUSD": "2401"
				},
				{
					"count": "8",
					"thorname": "omni",
					"volume": "7592100",
					"volumeUSD": "11"
				},
				{
					"count": "2",
					"thorname": "opzw",
					"volume": "25422306",
					"volumeUSD": "31"
				},
				{
					"count": "2",
					"thorname": "pb",
					"volume": "5123894",
					"volumeUSD": "6"
				},
				{
					"count": "1",
					"thorname": "piv",
					"volume": "2956622",
					"volumeUSD": "19"
				},
				{
					"count": "1",
					"thorname": "pivotal",
					"volume": "3298980",
					"volumeUSD": "21"
				},
				{
					"count": "1",
					"thorname": "pix",
					"volume": "67219777",
					"volumeUSD": "405"
				},
				{
					"count": "20",
					"thorname": "pl",
					"volume": "0",
					"volumeUSD": "0"
				},
				{
					"count": "27",
					"thorname": "pxc",
					"volume": "4641428360",
					"volumeUSD": "16648"
				},
				{
					"count": "1",
					"thorname": "qaz",
					"volume": "203934127",
					"volumeUSD": "427"
				},
				{
					"count": "15",
					"thorname": "qe",
					"volume": "4291645204",
					"volumeUSD": "6287"
				},
				{
					"count": "5",
					"thorname": "qn",
					"volume": "290983067",
					"volumeUSD": "1223"
				},
				{
					"count": "9",
					"thorname": "ref",
					"volume": "96559386",
					"volumeUSD": "131"
				},
				{
					"count": "9620",
					"thorname": "rj",
					"volume": "1877421551642",
					"volumeUSD": "2775551"
				},
				{
					"count": "1583",
					"thorname": "ro",
					"volume": "594234771429",
					"volumeUSD": "1041371"
				},
				{
					"count": "2",
					"thorname": "rujira",
					"volume": "78123",
					"volumeUSD": "0"
				},
				{
					"count": "2",
					"thorname": "sk",
					"volume": "1630411",
					"volumeUSD": "6"
				},
				{
					"count": "3291",
					"thorname": "ss",
					"volume": "16421012203388",
					"volumeUSD": "27956675"
				},
				{
					"count": "458",
					"thorname": "sy",
					"volume": "374380927225",
					"volumeUSD": "1110463"
				},
				{
					"count": "19864",
					"thorname": "t",
					"volume": "134862112552575",
					"volumeUSD": "226835017"
				},
				{
					"count": "14847",
					"thorname": "t1",
					"volume": "40846151410193",
					"volumeUSD": "60620024"
				},
				{
					"count": "5852",
					"thorname": "tb",
					"volume": "49185855108451",
					"volumeUSD": "69169236"
				},
				{
					"count": "260",
					"thorname": "tcb",
					"volume": "392072488443",
					"volumeUSD": "623115"
				},
				{
					"count": "1",
					"thorname": "tch",
					"volume": "27499572",
					"volumeUSD": "46"
				},
				{
					"count": "41",
					"thorname": "tchain",
					"volume": "147541132357",
					"volumeUSD": "226854"
				},
				{
					"count": "84480",
					"thorname": "td",
					"volume": "51153624686230",
					"volumeUSD": "93358183"
				},
				{
					"count": "23",
					"thorname": "thor13gymhg2atqujy5jhs6q4vg0h2mpcmxk485cvmp",
					"volume": "354666290",
					"volumeUSD": "493"
				},
				{
					"count": "1",
					"thorname": "thor1nyvskcndxfxne0hqwceselq55anayg7a9tes5h",
					"volume": "13660429",
					"volumeUSD": "24"
				},
				{
					"count": "5",
					"thorname": "thor1vtrlsuydxrst9ufgapgxfvldy2lp07pjayptwm",
					"volume": "2648318",
					"volumeUSD": "4"
				},
				{
					"count": "2",
					"thorname": "thor1zkkthewc4xmn0zfyn6qc9968cw73lfjdxa0ehx",
					"volume": "0",
					"volumeUSD": "0"
				},
				{
					"count": "122728",
					"thorname": "ti",
					"volume": "189658927216294",
					"volumeUSD": "345466958"
				},
				{
					"count": "28",
					"thorname": "tl",
					"volume": "779885508",
					"volumeUSD": "1257"
				},
				{
					"count": "1",
					"thorname": "to",
					"volume": "41045786",
					"volumeUSD": "68"
				},
				{
					"count": "4700",
					"thorname": "tps",
					"volume": "2551473680648",
					"volumeUSD": "6661031"
				},
				{
					"count": "8",
					"thorname": "tsw",
					"volume": "254037929",
					"volumeUSD": "439"
				},
				{
					"count": "1",
					"thorname": "ttx",
					"volume": "872846445",
					"volumeUSD": "1173"
				},
				{
					"count": "102",
					"thorname": "unizen-utxo",
					"volume": "70300892757",
					"volumeUSD": "100972"
				},
				{
					"count": "369",
					"thorname": "v0",
					"volume": "8047528345501",
					"volumeUSD": "9718147"
				},
				{
					"count": "911",
					"thorname": "va",
					"volume": "2120711823706",
					"volumeUSD": "3689603"
				},
				{
					"count": "1",
					"thorname": "valswappa",
					"volume": "1729013",
					"volumeUSD": "6"
				},
				{
					"count": "2023",
					"thorname": "vi",
					"volume": "10954291309197",
					"volumeUSD": "16598636"
				},
				{
					"count": "6",
					"thorname": "wehodl",
					"volume": "39015375",
					"volumeUSD": "50"
				},
				{
					"count": "13775",
					"thorname": "wr",
					"volume": "152124572725528",
					"volumeUSD": "228185641"
				},
				{
					"count": "36",
					"thorname": "x1",
					"volume": "9446535012",
					"volumeUSD": "12577"
				},
				{
					"count": "11",
					"thorname": "zak",
					"volume": "1258697130281",
					"volumeUSD": "1875177"
				},
				{
					"count": "3137",
					"thorname": "zengo",
					"volume": "5710273483841",
					"volumeUSD": "9941638"
				},
				{
					"count": "58",
					"thorname": "zengo-qa",
					"volume": "543176562",
					"volumeUSD": "861"
				}
			],
			"volume": "1101967806830954",
			"volumeUSD": "1784461425"
		}
	],
	"meta": {
		"count": "386716",
		"endTime": "1754338613",
		"startTime": "1647913096",
		"volume": "1101967806830954",
		"volumeUSD": "1784461425"
	}
}`
