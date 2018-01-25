package cgminer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
	"bytes"
)

type CGMiner struct {
	server string
	timeout time.Duration
}

/**
  "STATUS": [
    {
      "STATUS": "S",
      "When": 1516692550,
      "Code": 11,
      "Msg": "Summary",
      "Description": "cpuminer 2.3.2"
    }
  ]
 */
type status struct {
	Code        int
	Description string
	Status      string `json:"STATUS"`
	When        int64
}

//type Summary struct {
//	TotalPushWork				int64
//	Elapsed						int64
//	MHSav						float64 `json:"MHS av"`
//	MHS5s						float64 `json:"MHS 5s"`
//	MHS1m						float64 `json:"MHS 1m"`
//	MHS5m						float64 `json:"MHS 5m"`
//	MHS15m						float64 `json:"MHS 15m"`
//	Temper1						float64 `json:"temper1"`
//	Temper2						float64 `json:"temper2"`
//	Temper3						float64 `json:"temper3"`
//	//temper1 43.00
//	//temper2 43.00
//	//temper3 43.00
//	//inputVolt1 12.000
//	//inputVolt2 12.000
//	//inputVolt3 12.000
//	//inputCurrent1 30.120
//	//inputCurrent2 32.230
//	//inputCurrent3 32.100
//	//outputVolt1 5.600
//	//outputVolt2 5.600
//	//outputVolt3 5.600
//	//Found Blocks 0
//	//Getworks 0
//	//Accepted 1637997447479296
//	//Rejected 0
//	//Hardware Errors 0
//	//Utility 0.00
//	//Discarded 0
//	//Stale 0
//	//Get Failures 0
//	//Local Work 0
//	//Remote Failures 0
//	//Network Blocks 0
//	//Total MH 0.0000
//	//Work Utility 0.00
//	//Difficulty Accepted 0.00000000
//	//Difficulty Rejected 0.00000000
//	//Difficulty Stale 0.00000000
//	//Best Share 0
//	//Device Hardware% 1.0000
//	//Device Rejected% 1.0000
//	//Pool Rejected% 1.0000
//	//Pool Stale% 1.0000
//	//Last getwork 1516692550
//
//
//
//	//Accepted               int64
//	//BestShare              int64   `json:"Best Share"`
//	//DeviceHardwarePercent  float64 `json:"Device Hardware%"`
//	//DeviceRejectedPercent  float64 `json:"Device Rejected%"`
//	//DifficultyAccepted     float64 `json:"Difficulty Accepted"`
//	//DifficultyRejected     float64 `json:"Difficulty Rejected"`
//	//DifficultyStale        float64 `json:"Difficulty Stale"`
//	//Discarded              int64
//	//Elapsed                int64
//	//FoundBlocks            int64 `json:"Found Blocks"`
//	//GetFailures            int64 `json:"Get Failures"`
//	//Getworks               int64
//	//HardwareErrors         int64   `json:"Hardware Errors"`
//	//LocalWork              int64   `json:"Local Work"`
//	//MHS5s                  float64 `json:"MHS 5s"`
//	//MHSav                  float64 `json:"MHS av"`
//	//NetworkBlocks          int64   `json:"Network Blocks"`
//	//PoolRejectedPercentage float64 `json:"Pool Rejected%"`
//	//PoolStalePercentage    float64 `json:"Pool Stale%"`
//	//Rejected               int64
//	//RemoteFailures         int64 `json:"Remote Failures"`
//	//Stale                  int64
//	//TotalMH                float64 `json:"Total MH"`
//	//Utilty                 float64
//	//WorkUtility            float64 `json:"Work Utility"`
//}

type Devs struct {
	ASC				int
	Name			string
	ID				int
	Enabled			string
	Status			string
	MHSav			float64		`json:"MHS av"`
	MHS5s			float64		`json:"MHS 5s"`
	MHS1m			float64		`json:"MHS 1m"`
	MHS5m			float64		`json:"MHS 5m"`
	MHS15m			float64		`json:"MHS 15m"`
	Accepted		int64
	Rejected		int64
	HardwareErrors	int64		`json:"Hardware Errors"`
	DeviceElapsed	int64		`json:"Device Elapsed"`
	FansSpeed		int64		`json:"Fans Speed"`
	Temperature		float64		`json:"temperature"`
}

//type Pool struct {
//	Pool				int		`json:"POOL"`
//	URL					string
//	Status				string
//	Priority               int64
//	Quota                  int64
//	Accepted               int64
//	Rejected               int64
//	User                   string
//	LastShareTime          int64   `json:"Last Share Time"`
//
//	//"POOL": 0,
//	//"URL": "stratum+tcp://us2.litecoinpool.org:3333",
//	//"Status": "Alive",
//	//"Priority": 16842752,
//	//"Quota": 1,
//	//"Long Poll": "Y",
//	//"Getworks": 0,
//	//"Accepted": 1273,
//	//"Rejected": 76,
//	//"Works": 0,
//	//"Discarded": 0,
//	//"Stale": 0,
//	//"Get Failures": 0,
//	//"Remote Failures": 0,
//	//"User": "hodl4now.10011",
//	//"Passwd": "123123",
//	//"Last Share Time": 1516693404,
//	//"Diff1 Shares": 2048,
//	//"Proxy Type": "",
//	//"Proxy": "",
//	//"Difficulty Accepted": 0.00000000,
//	//"Difficulty Rejected": 0.00000000,
//	//"Difficulty Stale": 0.00000000,
//	//"Last Share Difficulty": 0.00000000,
//	//"Has Stratum": true,
//	//"Stratum Active": true,
//	//"Stratum URL": "stratum+tcp://us2.litecoinpool.org:3333",
//	//"Has GBT": false,
//	//"Best Share": 0,
//	//"Pool Rejected%": 0.0000,
//	//"Pool Stale%": 0.0000,
//	//"Bad Work": 0
//
//
//
//	//Accepted               int64
//	//BestShare              int64   `json:"Best Share"`
//	//Diff1Shares            int64   `json:"Diff1 Shares"`
//	//DifficultyAccepted     float64 `json:"Difficulty Accepted"`
//	//DifficultyRejected     float64 `json:"Difficulty Rejected"`
//	//DifficultyStale        float64 `json:"Difficulty Stale"`
//	//Discarded              int64
//	//GetFailures            int64 `json:"Get Failures"`
//	//Getworks               int64
//	//HasGBT                 bool    `json:"Has GBT"`
//	//HasStratum             bool    `json:"Has Stratum"`
//	//LastShareDifficulty    float64 `json:"Last Share Difficulty"`
//	//LastShareTime          int64   `json:"Last Share Time"`
//	//LongPoll               string  `json:"Long Poll"`
//	//Pool                   int64   `json:"POOL"`
//	//PoolRejectedPercentage float64 `json:"Pool Rejected%"`
//	//PoolStalePercentage    float64 `json:"Pool Stale%"`
//	//Priority               int64
//	//ProxyType              string `json:"Proxy Type"`
//	//Proxy                  string
//	//Quota                  int64
//	//Rejected               int64
//	//RemoteFailures         int64 `json:"Remote Failures"`
//	//Stale                  int64
//	//Status                 string
//	//StratumActive          bool   `json:"Stratum Active"`
//	//StratumURL             string `json:"Stratum URL"`
//	//URL                    string
//	//User                   string
//	//Works                  int64
//}

//type summaryResponse struct {
//	Status  []status  `json:"STATUS"`
//	Summary []Summary `json:"SUMMARY"`
//	Id      int64     `json:"id"`
//}

type devsResponse struct {
	Status  []status  `json:"STATUS"`
	Devs    []Devs    `json:"DEVS"`
	Id      int64     `json:"id"`
}

type ChipStat map[string]float64

type chipStatResponse struct {
	Status  []status  `json:"STATUS"`
	ChipStat ChipStat	`json:"CHIPSTAT"`
}

//type poolsResponse struct {
//	Status []status `json:"STATUS"`
//	Pools  []Pool   `json:"POOLS"`
//	Id     int64    `json:"id"`
//}
//
//type addPoolResponse struct {
//	Status []status `json:"STATUS"`
//	Id     int64    `json:"id"`
//}

// New returns a CGMiner pointer, which is used to communicate with a running
// CGMiner instance. Note that New does not attempt to connect to the miner.
func New(hostname string, port int64, timeout time.Duration) *CGMiner {
	miner := new(CGMiner)
	miner.server = fmt.Sprintf("%s:%d", hostname, port)
	miner.timeout = time.Second * timeout
	return miner
}

func (miner *CGMiner) runCommand(command, argument string) (string, error) {
	conn, err := net.DialTimeout("tcp", miner.server, miner.timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	type commandRequest struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	request := &commandRequest{
		Command: command,
	}

	if argument != "" {
		request.Parameter = argument
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(conn, "%s", requestBody)
	result, err := bufio.NewReader(conn).ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(result, "\x00"), nil
}

// Devs returns basic information on the miner. See the Devs struct.
func (miner *CGMiner) Devs() (*[]Devs, error) {
	result, err := miner.runCommand("devs", "")
	if err != nil {
		return nil, err
	}

	var devsResponse devsResponse
	err = json.Unmarshal([]byte(result), &devsResponse)
	if err != nil {
		return nil, err
	}

	var devs = devsResponse.Devs
	return &devs, err
}

func (miner *CGMiner) ChipStat() (*ChipStat, error) {
	response, err := miner.runCommand("chipstat", "")
	if err != nil {
		return nil, err
	}

	result, err := processChipStat(response)
	return &result.ChipStat, err
}

func processChipStat(response string) (*chipStatResponse, error) {
	// BW kindly messed up the json on this one.
	fixResponse := bytes.Replace([]byte(response), []byte(",\"SUMMARY\":["), []byte(""), 1)
	fixResponse = bytes.Replace(fixResponse, []byte("]{"), []byte("],\"CHIPSTAT\":{"), 1)
	fixResponse = bytes.Replace(fixResponse, []byte("],\"id\":1"), []byte(""), 1)

	var chipStatResponse chipStatResponse
	err := json.Unmarshal(fixResponse, &chipStatResponse)
	if err != nil {
		return nil, err
	}

	return &chipStatResponse, err
}

//// Summary returns basic information on the miner. See the Summary struct.
//func (miner *CGMiner) Summary() (*Summary, error) {
//	result, err := miner.runCommand("summary", "")
//	if err != nil {
//		return nil, err
//	}
//
//	var summaryResponse summaryResponse
//	err = json.Unmarshal([]byte(result), &summaryResponse)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(summaryResponse.Summary) != 1 {
//		return nil, errors.New("Received multiple Summary objects")
//	}
//
//	var summary = summaryResponse.Summary[0]
//	return &summary, err
//}

//// Pools returns a slice of Pool structs, one per pool.
//func (miner *CGMiner) Pools() ([]Pool, error) {
//	result, err := miner.runCommand("pools", "")
//	if err != nil {
//		return nil, err
//	}
//
//	var poolsResponse poolsResponse
//	err = json.Unmarshal([]byte(result), &poolsResponse)
//	if err != nil {
//		return nil, err
//	}
//
//	var pools = poolsResponse.Pools
//	return pools, nil
//}

//// AddPool adds the given URL/username/password combination to the miner's
//// pool list.
//func (miner *CGMiner) AddPool(url, username, password string) error {
//	// TODO: Don't allow adding a pool that's already in the pool list
//	// TODO: Escape commas in the URL, username, and password
//	parameter := fmt.Sprintf("%s,%s,%s", url, username, password)
//	result, err := miner.runCommand("addpool", parameter)
//	if err != nil {
//		return err
//	}
//
//	var addPoolResponse addPoolResponse
//	err = json.Unmarshal([]byte(result), &addPoolResponse)
//	if err != nil {
//		// If there an error here, it's possible that the pool was actually added
//		return err
//	}
//
//	status := addPoolResponse.Status[0]
//
//	if status.Status != "S" {
//		return errors.New(fmt.Sprintf("%d: %s", status.Code, status.Description))
//	}
//
//	return nil
//}
//
//func (miner *CGMiner) Enable(pool *Pool) error {
//	parameter := fmt.Sprintf("%d", pool.Pool)
//	_, err := miner.runCommand("enablepool", parameter)
//	return err
//}
//
//func (miner *CGMiner) Disable(pool *Pool) error {
//	parameter := fmt.Sprintf("%d", pool.Pool)
//	_, err := miner.runCommand("disablepool", parameter)
//	return err
//}
//
//func (miner *CGMiner) Delete(pool *Pool) error {
//	parameter := fmt.Sprintf("%d", pool.Pool)
//	_, err := miner.runCommand("removepool", parameter)
//	return err
//}
//
//func (miner *CGMiner) SwitchPool(pool *Pool) error {
//	parameter := fmt.Sprintf("%d", pool.Pool)
//	_, err := miner.runCommand("switchpool", parameter)
//	return err
//}
//
//func (miner *CGMiner) Restart() error {
//	_, err := miner.runCommand("restart", "")
//	return err
//}
//
//func (miner *CGMiner) Quit() error {
//	_, err := miner.runCommand("quit", "")
//	return err
//}