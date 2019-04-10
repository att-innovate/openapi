package lteAdapter

import (
	m "openapi/models"
)

/*
Interface layer to abstract technology and Amarisoft-specific implementation.
*/

func GetMMEStats(ipaddr string) m.MMEStats {
	return GetAmarisoftMMEStats(ipaddr)
}

func GetUEGeteNB(ipaddr string) m.UEGeteNB {
	return GetAmarisoftUEGeteNB(ipaddr)
}

func GetOptiomalENB1TargetCell() (int, int) {
	return GetOptiomalAmarisoftENB1TargetCell()
}

func GetOptiomalENB2TargetCell() (int, int) {
	return GetOptiomalAmarisoftENB2TargetCell()
}

func Handover(ipaddr string, enbueid, pci, earfcn int) bool {
	return AmarisoftHandover(ipaddr, enbueid, pci, earfcn)
}
