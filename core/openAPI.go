package core

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	conf "openapi/conf"
	adpt "openapi/lteAdapter"
	m "openapi/models"
	u "openapi/user"
	"strconv"
	"strings"
	"time"
)

// Optional paramter for wifi setup to indicate the limitand identify
// the accepted user with IP ending on 192.168.1. and wiFiCLientIPOctet
var wiFiCLientIPOctet = 206

var counterStreaming = make(map[uint64]int)
var counterNormal = make(map[uint64]int)
var ueWLANmapping = make(map[uint64]int)

type OpenAPI struct {
	threshold         int
	Config            conf.Configuration
	wiFiCLientIPOctet int
}

func (oa OpenAPI) InitConfiguration(c conf.Configuration) {
	oa.Config = c
	oa.threshold = c.HandoverThreshold
	oa.wiFiCLientIPOctet = c.WiFiCLientIPOctet
}

func (oa OpenAPI) DetermineAction(token uint64, mode string, action string) {
	oa.Config = conf.LoadConfiguration(mode)

	log.Printf("param [%v], token [%v], mode [%s].", mode, token, action)
	switch oa.Config.Env {
	case conf.WLAN:
		//log.Printf("INFO: action defined for [%v] mode.", conf.WLAN)
		oa.executeWlanAction(token, mode)
	case conf.LTE:
		log.Printf("INFO: Action defined for [%v] mode.", conf.LTE)
		oa.executeLTEAction(token, action)
	case conf.TEST:
		log.Printf("INFO: no action defined for [%v] mode.", conf.TEST)
	default:
		log.Printf("ERROR in DetermineAction: no default action mode definded.")
	}
}

func (oa OpenAPI) executeLTEAction(token uint64, mode string) {
	log.Printf("starting executeLTEAction")
	switch mode {
	case conf.STREAMING:
		counter := counterStreaming[token]
		counter++
		counterStreaming[token] = counter
		//log.Printf("'streaming mode' counter = [%v], 'normal mode' counter = [0], threshold=[%v].", counter, oa.Config.HandoverThreshold)

		ueMap := make(map[int][]int)
		ueMap = oa.gatherStatistics(ueMap)

		if counter < oa.Config.HandoverThreshold {
			counterNormal[token] = 0
			return
		} else if ueMap[1] != nil && len(ueMap) >= -1 && counter == oa.Config.HandoverThreshold {
			ueList := ueMap[1]
			cellID, earfcn := adpt.GetOptiomalENB2TargetCell()
			log.Printf("ACTION: \tchange network from eNB1 (rate limited) to cellID (%v) / earfcn (%v) onto eNB2 (open) for UE_ID: %v", cellID, earfcn, ueList)
			adpt.Handover(conf.GeteNBonPosition(1), ueList[0], cellID, earfcn)
			return
		} else {
			//log.Printf("NO ACTION: \tUE already connected to best cell: cellID on eNB2")
			return
		}

	case conf.NORMAL:
		counterStreaming[token] = 0
		counter := counterNormal[token]
		counter++
		counterNormal[token] = counter
		log.Printf("Indicated 'normal mode', nulled 'streaming mode' counter and return to default network momentarily with counter %v.", counter)

		ueMap := make(map[int][]int)
		ueMap = oa.gatherStatistics(ueMap)

		log.Printf("debug %v and %v.", ueMap, ueMap[1])

		if ueMap[1] != nil && len(ueMap[1]) >= -1 {
			ueList := ueMap[1]
			cellID, earfcn := adpt.GetOptiomalENB1TargetCell()
			log.Printf("ACTION: \tchange network from eNB2 (open) to cellID (%v) / earfcn (%v) toon eNB1 (rate limited) for UE_ID: %v", cellID, earfcn, ueList)
			adpt.Handover(conf.GeteNBonPosition(0), ueList[0], cellID, earfcn)
			return
		}

	default:
		log.Printf("ERROR: no mode selection definded.")
	}
}

func (oa OpenAPI) gatherStatistics(ueMap map[int][]int) map[int][]int {
	//eNB1stats := adpt.GetUEGeteNB("192.168.1.207:9001")
	eNBs := oa.Config.Enbs

	for eNBn := 0; eNBn < len(eNBs); eNBn++ {
		eNBaddr := eNBs[eNBn]
		eNBstats := adpt.GetUEGeteNB(eNBaddr)

		for ue := 0; ue < len(eNBstats.UeList); ue++ {
			ueList := ueMap[1]
			ueList = append(ueList, eNBstats.UeList[ue].EnbUeID)
			ueMap[1] = ueList
			//log.Printf("\t\teNBstats.UeList.EnbUeID: %v", eNBstats.UeList[ue].EnbUeID)
			for c := 0; c < len(eNBstats.UeList[ue].Cells); c++ {
				//log.Printf("\t\teNBstats.UeList.cellID: %v", eNBstats.UeList[ue].Cells[c].CellID)
			}
		}
	}

	// Print and log statistics
	/*
		eNB1stats := adpt.GetUEGeteNB("192.168.1.207:9001")
		eNB2stats := adpt.GetUEGeteNB("192.168.1.244:9001")
		mmeStats := adpt.GetMMEStats("192.168.1.244:9000")

		//log.Printf("stats.CPU.global: %f", mmeStats.CPU.Global)
		for i := 0; i < len(mmeStats.S1Connections); i++ {
			//log.Printf("stats.S1Connections.IPAddr: %v", mmeStats.S1Connections[i].IPAddr)
			//log.Printf("stats.S1Connections.EnbID: %v", mmeStats.S1Connections[i].EnbID)
		}

		for ue := 0; ue < len(eNB1stats.UeList); ue++ {
			ueList := ueMap[1]
			ueList = append(ueList, eNB1stats.UeList[ue].EnbUeID)
			ueMap[1] = ueList
			log.Printf("\t\teNB1stats.UeList.EnbUeID: %v", eNB1stats.UeList[ue].EnbUeID)
			for c := 0; c < len(eNB1stats.UeList[ue].Cells); c++ {
				log.Printf("\t\teNB1stats.UeList.cellID: %v", eNB1stats.UeList[ue].Cells[c].CellID)
			}
		}

		for ue := 0; ue < len(eNB2stats.UeList); ue++ {
			ueList := ueMap[2]
			ueList = append(ueList, eNB2stats.UeList[ue].EnbUeID)
			ueMap[2] = ueList
			log.Printf("\\tteNB2stats.UeList.EnbUeID: %v", eNB2stats.UeList[ue].EnbUeID)
			for c := 0; c < len(eNB2stats.UeList[ue].Cells); c++ {
				log.Printf("\t\teNB2stats.UeList.cellID: %v", eNB2stats.UeList[ue].Cells[c].CellID)
			}
		}
	*/
	return ueMap
}

func (oa OpenAPI) executeWlanAction(token uint64, mode string) {
	url := ""
	switch mode {
	case conf.STREAMING:
		counter := counterStreaming[token]
		counter++
		counterStreaming[token] = counter
		//log.Printf("increased 'streaming mode' counter to %v and nulled 'normal mode' counter - return without handover.", counter)
		if counter <= oa.threshold {
			counterNormal[token] = 0
			return
		} else if ueWLANmapping[token] != 16 {
			//log.Printf("change network to 16")
			url = fmt.Sprintf("http://192.168.1.%v:8080?target_ssid=OpenWRT_16&password=foundry123", wiFiCLientIPOctet)
			ueWLANmapping[token] = 16
			//counterWLANstreaming[token] = 0
		} else if ueWLANmapping[token] == 16 {
			//log.Printf("Remain on network 16")
			return
		}

	case conf.NORMAL:
		url = fmt.Sprintf("http://192.168.1.%v:8080?target_ssid=OpenWRT_23&password=foundry123", wiFiCLientIPOctet)
		ueWLANmapping[token] = 23
		counterStreaming[token] = 0
		counter := counterNormal[token]
		counter++
		counterNormal[token] = counter
		//log.Printf("Indicated 'normal mode', nulled 'streaming mode' counter and return to default network.", counter)
		/*
			if counter <= threshold {
				counterWLANstreaming[token] = 0
				return
			} else if ueWLANmapping[token] != 23 {
				log.Printf("change network to 23")
				url = fmt.Sprintf("http://192.168.1.%v:8080?target_ssid=OpenWRT_23&password=foundry123", wiFiCLientIPOctet)
				ueWLANmapping[token] = 23
			} else if ueWLANmapping[token] == 23 {
				log.Printf("Remain on network 23")
				return
			}
		*/
	default:
		log.Printf("ERROR: no mode selection definded.")
	}
	resp, err := http.PostForm(url, nil)
	if resp != nil {
		log.Printf("StreamingModeHandler resp: %v.", resp)
	}
	if err != nil {
		log.Printf("StreamingModeHandler network action failed with ERROR: %v.", err)
	}
}

func (oa OpenAPI) RegisterClient(rm m.RegisterMessage) m.ClientStats {
	log.Printf("register client under ip addr: '%v'.", rm.IP)

	var token = calculateuint64(rm.IP)
	ns := GetNetworkStats(token)
	cs := m.ClientStats{Token: token, BW: ns.BW, Latency: ns.Latency}

	newUser := u.User{Token: token, Registrationtime: time.Now().UTC().UnixNano(), IP: rm.IP}
	u.AddUser(newUser)

	return cs
}

func (oa OpenAPI) DeregisterClient(token uint64) {
	log.Printf("de-register client under token:[%v].", token)

	u.RemoveUser(token)
}

func (oa OpenAPI) GetNetworkStats(token uint64) m.NetworkStats {
	ns := GetNetworkStats(token)

	return ns
}

/*
Interface to retreive network statistics
Return: network parameters for client
*/
func GetNetworkStats(token uint64) m.NetworkStats {
	var latency int = -1
	var bandwidth int = -1
	var quality int = -1

	/*
		Overwrite default values with RAN and/or core specific parameters.
		Either query network, service and cache/database synchronous or asynchronous.
	*/

	ns := m.NetworkStats{Latency: latency, BW: bandwidth, Quality: quality}

	return ns
}

func calculateuint64(ip string) (number uint64) {
	bits := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return uint64(sum)
}

func calculateHash(input string) (hash string) {
	s := input
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	result := fmt.Sprintf("%x\n", bs)
	//log.Printf("calculated %v for %v.\n", result, input)

	return result
}
