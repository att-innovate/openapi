package mobility

type Action struct {
	command string
}
type MobilityManager interface {
	mobilityAction() Action
}

type MobilityManagerWlan struct {
}

type MobilityManagerLTE struct {
}
