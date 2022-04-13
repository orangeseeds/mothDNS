package zone

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

type ZoneList struct {
	List []ZoneEntry `json:"list"`
}

type ZoneEntry struct {
	QuestionSection   []core.DnsQuestion `json:"question_section"`
	AnswerSection     []core.DnsRecord   `json:"answer_section"`
	AdditionalSection []core.DnsRecord   `json:"additional_section"`
	AuthoritySection  []core.DnsRecord   `json:"authority_section"`
	ReturnCode        core.ResultCode    `json:"return_code"`
	ID                uint16             `json:"id"`
	AA                bool               `json:"authorative_answer"`
	TC                bool               `json:"truncation"`
	RD                bool               `json:"recursion_desired"`
	RA                bool               `json:"recursion_available"`
	AD                bool               `json:"authentic_data"`
	Query             map[string]string  `json:"query"`
}

var Zlist = ZoneList{List: []ZoneEntry{}}
var Queries = map[string]string{}

func (z *ZoneEntry) ToJson() string {
	val, _ := utils.PrettyStruct(z)
	return val
}

func init() {
	ReadJson("userfile.json")

	for _, entry := range Zlist.List {
		keys := []string{}

		i := 0
		for _, k := range entry.Query {
			keys = append(keys, k)
			i++
		}
		// fmt.Println(keys)
		Queries[keys[0]] = entry.Query[keys[0]]
	}

	// fmt.Println(Queries)
}

func NewZEntry() (ZoneEntry, error) {
	z := ZoneEntry{
		QuestionSection:   []core.DnsQuestion{},
		AnswerSection:     []core.DnsRecord{},
		AdditionalSection: []core.DnsRecord{},
		AuthoritySection:  []core.DnsRecord{},
		ReturnCode:        core.NOERROR,
		ID:                0,
		AA:                false,
		RD:                false,
		RA:                false,
		AD:                false,
		Query:             map[string]string{},
	}
	return z, nil
}

func (z *ZoneEntry) PacketToZEntry(d *core.DnsPacket) error {
	zPacket := ZoneEntry{
		QuestionSection:   d.Questions,
		AnswerSection:     d.Answers,
		AdditionalSection: d.Resources,
		AuthoritySection:  d.Authorities,
		ReturnCode:        d.Header.Rescode,
		ID:                d.Header.Id,
		AA:                d.Header.Authoritative_answer,
		RD:                d.Header.Recursion_desired,
		RA:                d.Header.Recursion_available,
		AD:                d.Header.Authed_data,
		Query: map[string]string{
			"server": "127.0.0.1",
		},
	}

	z = &zPacket
	Zlist.List = append(Zlist.List, *z)

	WriteJson()

	return nil
}

func WriteJson() {
	val, _ := utils.PrettyStruct(Zlist)
	err := ioutil.WriteFile("userfile.json", []byte(val), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadJson(path string) (*ZoneList, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	zList := ZoneList{}
	err = json.Unmarshal(content, &zList)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(newZEntry.ToJson())
	return &zList, nil
}

// func GetQueries(path string) {
// 	zones, _ := ReadJson(path)
// 	for _, zone := range zones {
// 	    if val,ok := Queries[zone.query]
// 	}
// }
