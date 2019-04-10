package core

type RegisterMessage struct {
	IP string `json:"ip"`
}

type ClientStats struct {
	Token          uint64 `json:"token,string,omitempty"`
	BW             int    `json:"bw,string,omitempty"`
	Timestamp      uint64 `json:"timestamp,string,omitempty"`
	Latency        int    `json:"latency,string,omitempty"`
	Bitrate        int    `json:"bitrate,string,omitempty"`
	BufferDuration int    `json:"buffered_duration,string,omitempty"`
}

type NetworkStats struct {
	BW      int `json:"bw"`
	Latency int `json:"latency"`
	Quality int `json:"quality"`
}

type UEGetMME struct {
	Message string `json:"message"`
	UeList  []struct {
		Imsi                    string `json:"imsi"`
		Sqn                     string `json:"sqn"`
		Imeisv                  string `json:"imeisv"`
		MTmsi                   int64  `json:"m_tmsi"`
		Registered              bool   `json:"registered"`
		UeAggregateMaxBitrateDl int    `json:"ue_aggregate_max_bitrate_dl"`
		UeAggregateMaxBitrateUl int    `json:"ue_aggregate_max_bitrate_ul"`
		Tac                     int    `json:"tac"`
		TacPlmn                 string `json:"tac_plmn"`
		Bearers                 []struct {
			ErabID       int    `json:"erab_id"`
			IP           string `json:"ip"`
			DlTotalBytes int    `json:"dl_total_bytes"`
			UlTotalBytes int    `json:"ul_total_bytes"`
			Apn          string `json:"apn"`
		} `json:"bearers"`
		EnbUeID int `json:"enb_ue_id"`
		MmeUeID int `json:"mme_ue_id"`
	} `json:"ue_list"`
	MessageID string `json:"message_id"`
}

type UEGeteNB struct {
	Message string `json:"message"`
	UeList  []struct {
		EnbUeID int `json:"enb_ue_id"`
		MmeUeID int `json:"mme_ue_id"`
		Rnti    int `json:"rnti"`
		Cells   []struct {
			CellID int `json:"cell_id"`
		} `json:"cells"`
	} `json:"ue_list"`
	MessageID string `json:"message_id"`
}

type MMEStats struct {
	Message string `json:"message"`
	CPU     struct {
		Global float64 `json:"global"`
	} `json:"cpu"`
	Counters struct {
		Messages struct {
			S1InitialContextSetupRequest  int `json:"s1_initial_context_setup_request"`
			S1InitialContextSetupResponse int `json:"s1_initial_context_setup_response"`
			S1DownlinkNasTransport        int `json:"s1_downlink_nas_transport"`
			S1InitialUeMessage            int `json:"s1_initial_ue_message"`
			S1UplinkNasTransport          int `json:"s1_uplink_nas_transport"`
			S1SetupRequest                int `json:"s1_setup_request"`
			S1SetupResponse               int `json:"s1_setup_response"`
			S1UeCapabilityInfoIndication  int `json:"s1_ue_capability_info_indication"`
			NasAttachRequest              int `json:"nas_attach_request"`
			NasAttachAccept               int `json:"nas_attach_accept"`
			NasAttachComplete             int `json:"nas_attach_complete"`
			NasAuthenticationRequest      int `json:"nas_authentication_request"`
			NasAuthenticationResponse     int `json:"nas_authentication_response"`
			NasSecurityModeCommand        int `json:"nas_security_mode_command"`
			NasSecurityModeComplete       int `json:"nas_security_mode_complete"`
			NasEsmInformationRequest      int `json:"nas_esm_information_request"`
			NasEsmInformationResponse     int `json:"nas_esm_information_response"`
			NasEmmInformation             int `json:"nas_emm_information"`
			DiamCer                       int `json:"diam_cer"`
			DiamCea                       int `json:"diam_cea"`
			DiamAir                       int `json:"diam_air"`
			DiamAia                       int `json:"diam_aia"`
		} `json:"messages"`
		Errors struct {
		} `json:"errors"`
	} `json:"counters"`
	EmmRegisteredUeCount int `json:"emm_registered_ue_count"`
	S1Connections        []struct {
		Plmn   string `json:"plmn"`
		EnbID  int    `json:"enb_id"`
		IPAddr string `json:"ip_addr"`
		TaList []struct {
			Plmn string `json:"plmn"`
			Tac  int    `json:"tac"`
		} `json:"ta_list"`
		EmmConnectedUeCount int `json:"emm_connected_ue_count"`
	} `json:"s1_connections"`
	InstanceID string `json:"instance_id"`
	MessageID  string `json:"message_id"`
}

type CommandResultCode struct {
	Message   string `json:"message"`
	Error     string `json:"error"`
	MessageID string `json:"message_id"`
}
