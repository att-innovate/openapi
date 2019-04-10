package comm

import (
	"encoding/json"
	"log"
	"net/http"
	conf "openapi/conf"
	"openapi/core"
	m "openapi/models"
	p "openapi/persistence"
	u "openapi/user"
	"strconv"
)

var openAPI = core.OpenAPI{}

func SetOpenAPI(oa core.OpenAPI) {
	openAPI = oa
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var registerMessage m.RegisterMessage
	err := decoder.Decode(&registerMessage)
	log.Printf("received %v", r)

	if err != nil {
		log.Printf("panic in unpacking JSON under RegisterHandler")
	}
	if &registerMessage != nil && registerMessage.IP != "" {
		//log.Printf("Received register message for client IP %v", registerMessage)

		clientStats := openAPI.RegisterClient(registerMessage)

		if err := json.NewEncoder(w).Encode(clientStats); err != nil {
			log.Printf("ERROR in RegisterHandler encoding interfaces as JSON: %v.", err)
		}
	} else {
		log.Printf("Incomplete POST request in RegisterHandler")
		w.WriteHeader(http.StatusNoContent)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	keys, ok := r.URL.Query()["clientid"]
	if !ok || len(keys) < 1 {
		log.Println("ERROR in HelloHandler: Url Param 'clientid' is missing")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	key := keys[0]
	log.Println("Url Param 'clientid' (ip addr for the start) is: " + string(key))

	var registerMessage m.RegisterMessage
	registerMessage.IP = string(key)

	clientStats := openAPI.RegisterClient(registerMessage)

	if err := json.NewEncoder(w).Encode(clientStats); err != nil {
		log.Printf("ERROR in HelloHandler while encoding interfaces as JSON: %v.", err)
		panic(err)
	}
}

func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)

	t, err := dec.Token()
	if err != nil {
		log.Printf("Error in StatisticsHandler for dec.Token %v.", err)
	} else if t == nil {
		log.Printf("Error in StatisticsHandler - Token is nil.")
	}

	if dec.More() == false {
		log.Printf("\tNo further action - return")
		return
	}
	for dec.More() {
		var cs m.ClientStats
		err := dec.Decode(&cs)
		if err != nil {
			log.Printf("Error in StatisticsHandler in dec.Decode(&cs) %v.", err)
		}
		if isTokenValid(cs.Token) {
			//log.Printf("db disabled - token: %v, \tbitrate: %v, \tbw: %v, \tlatency: %v, \ttimestamp: %v,
			//	\tbuffered_duration: %v.\n", cs.Token, cs.Bitrate, cs.BW, cs.Latency, cs.Timestamp, cs.BufferDuration)
			p.WriteDB(cs)
		}
	}

	t, err = dec.Token()
	if err != nil {
		log.Printf("ERROR in reading read closing bracket %v.", err)
	}
}

func StreamingModeHandler(w http.ResponseWriter, r *http.Request, env string) {
	defer r.Body.Close()

	handleMultiModeHandler(w, r, env, conf.STREAMING)
}

func LatencyModeHandler(w http.ResponseWriter, r *http.Request, env string) {
	defer r.Body.Close()

	handleMultiModeHandler(w, r, env, conf.LATENCY)
}

func NormalModeHandler(w http.ResponseWriter, r *http.Request, env string) {
	defer r.Body.Close()

	handleMultiModeHandler(w, r, env, conf.NORMAL)
}

func handleMultiModeHandler(w http.ResponseWriter, r *http.Request, mode string, action string) {
	keys, ok := r.URL.Query()["client"]
	log.Printf("keys: %s", keys)

	if !ok || len(keys) < 1 {
		log.Println("ERROR in handleMultiModeHandler: Url Param 'client' is missing")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	key := keys[0]
	token := stringToUint64(key)

	if isTokenValid(token) {
		log.Printf("token valid - continue")

		openAPI.DetermineAction(token, mode, action)
		networkStats := openAPI.GetNetworkStats(token)
		if err := json.NewEncoder(w).Encode(networkStats); err != nil {
			log.Printf("ERROR in encoding interfaces as JSON: %v.", err)
		}
	} else {
		log.Printf("token not valid - user unknown or not registered - discard request")
	}
	log.Printf("close connection")
}

func GoodbyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	keys, ok := r.URL.Query()["clientid"]

	if !ok || len(keys) < 1 {
		log.Println("ERROR in GoodbyHandler: Url Param 'clientid' is missing")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	key := keys[0]
	token := stringToUint64(key)

	openAPI.DeregisterClient(token)

	w.Header().Set("Server", "OpenAPI call")
	w.WriteHeader(200)
}

func stringToUint64(input string) (result uint64) {
	u, _ := strconv.ParseUint(string(input), 0, 64)

	return u
}

func isTokenValid(token uint64) bool {
	flag := u.LookupUser(token)

	return flag
}
