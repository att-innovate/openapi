package lteAdapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	m "openapi/models"
	"os"
	"os/exec"
	"strings"
)

/*
Specific iplementation for Amarisoft eNB stack.
*/
func GetAmarisoftMMEStats(ipaddr string) m.MMEStats {
	cmd := "node"
	args := []string{"<SYSTEM PATH TO WEBSOCKET>/ws.js", ipaddr, "{\"message\": \"stats\"}"}
	process := exec.Command(cmd, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		log.Println("ERROR in GetAmarisoftMMEStats: ", err)
	}
	defer stdin.Close()
	buf := new(bytes.Buffer)
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		log.Println("An error occured: ", err)
	}

	process.Wait()

	//fmt.Println("Generated string:", buf)

	var mmeStats m.MMEStats
	buf2 := new(bytes.Buffer)
	//log.Printf("read from buffer.")
	buf.ReadFrom(buf2)
	b := buf.Bytes()

	uncleanJSON := string(b)
	i := strings.Index(uncleanJSON, "{")
	cleanJSON := uncleanJSON
	if i > -1 {
		//temp := strings.SplitAfterN(uncleanJSON, "{", 1)
		//log.Printf("len %v and temp %v position i %v", len(temp), temp, i)
		cleanJSON = uncleanJSON[i:len(uncleanJSON)]
	}
	log.Printf("start unmarsheling b: %s", cleanJSON)
	err2 := json.Unmarshal([]byte(cleanJSON), &mmeStats)
	if err2 != nil {
		log.Printf("Error while unmarshaling: %v", err2)
	}
	//log.Printf("stats.CPU.global: %f", mmeStats.CPU.Global)
	return mmeStats
}

func GetAmarisoftUEGeteNB(ipaddr string) m.UEGeteNB {
	cmd := "node"
	args := []string{"<SYSTEM PATH TO WEBSOCKET>/ws.js", ipaddr, "{\"message\": \"ue_get\"}"}
	process := exec.Command(cmd, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		log.Println("Error in GetAmarisoftUEGeteNB", err)
	}
	defer stdin.Close()
	buf := new(bytes.Buffer) // THIS STORES THE NODEJS OUTPUT
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		log.Println("An error occured: ", err)
	}

	process.Wait()

	//fmt.Println("Generated string:", buf)

	var ueGeteNB m.UEGeteNB
	buf2 := new(bytes.Buffer)
	//log.Printf("read from buffer.")
	buf.ReadFrom(buf2)
	b := buf.Bytes()

	uncleanJSON := string(b)
	i := strings.Index(uncleanJSON, "{")
	cleanJSON := uncleanJSON
	if i > -1 {
		//temp := strings.SplitAfterN(uncleanJSON, "{", 1)
		//log.Printf("len %v and temp %v position i %v", len(temp), temp, i)
		cleanJSON = uncleanJSON[i:len(uncleanJSON)]
	}
	//log.Printf("start unmarsheling b: %s", cleanJSON)
	err2 := json.Unmarshal([]byte(cleanJSON), &ueGeteNB)
	if err2 != nil {
		log.Printf("Error while unmarshaling: %v", err2)
	}
	return ueGeteNB
}

func AmarisoftHandover(ipaddr string, enbueid, pci, earfcn int) bool {
	cmd := "node"
	handoverCommand := fmt.Sprintf("{\"message\": \"handover\", \"enb_ue_id\": %v, \"pci\": %v, \"dl_earfcn\":%v}", enbueid, pci, earfcn)
	args := []string{"<SYSTEM PATH TO WEBSOCKET>/ws.js", ipaddr, handoverCommand}
	process := exec.Command(cmd, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()
	buf := new(bytes.Buffer)
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		log.Println("An error occured: ", err)
	}

	process.Wait()
	//log.Println("Generated string: ", buf)

	var crc m.CommandResultCode
	buf2 := new(bytes.Buffer)
	//log.Printf("read from buffer.")
	buf.ReadFrom(buf2)
	b := buf.Bytes()

	uncleanJSON := string(b)
	i := strings.Index(uncleanJSON, "{")
	cleanJSON := uncleanJSON
	if i > -1 {
		//temp := strings.SplitAfterN(uncleanJSON, "{", 1)
		cleanJSON = uncleanJSON[i:len(uncleanJSON)]
	}
	//log.Printf("start unmarsheling b: %s", cleanJSON)
	err2 := json.Unmarshal([]byte(cleanJSON), &crc)
	if err2 != nil {
		log.Printf("Error while unmarshaling: %v", err2)
	}

	flag := validateCommandResultCode(crc)

	return flag
}

func validateCommandResultCode(crc m.CommandResultCode) bool {
	if crc.Error != "" {
		log.Println("ERROR found in CommandResultCode: ", crc.Error)
		return false
	} else {
		return true
	}
}

func GetOptiomalAmarisoftENB1TargetCell() (int, int) {
	cellID := rand.Intn(2) + 1
	earfcn := 0
	if cellID == 1 {
		earfcn = 3350
	} else if cellID == 2 {
		earfcn = 3248
	}
	return cellID, earfcn
}
func GetOptiomalAmarisoftENB2TargetCell() (int, int) {
	cellID := rand.Intn(2) + 3
	earfcn := 0
	if cellID == 3 {
		earfcn = 2849
	} else if cellID == 4 {
		earfcn = 3050
	}
	return cellID, earfcn
}
