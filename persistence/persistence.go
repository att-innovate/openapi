package persistence

import (
	"fmt"
	"log"
	m "openapi/models"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

const (
	INFLUX_SERVER   = "http://0.0.0.0"
	INFLUX_PORT     = "8086"
	INFLUX_DBNAME   = ""
	INFLUX_TABLE    = ""
	INFLUX_USERNAME = ""
	INFLUX_PASSWORD = ""
)

func WriteDB(cs m.ClientStats) {
	//log.Printf("perform db write action")

	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("%s:%s", INFLUX_SERVER, INFLUX_PORT),
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})

	if err != nil {
		log.Printf("Error in instantiating inflixdb client %v.", err)
	}

	defer clnt.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  INFLUX_DBNAME,
		Precision: "ms",
	})

	tags := map[string]string{"token": fmt.Sprintf("%v", cs.Token)}
	fields := map[string]interface{}{
		"bandwidth": cs.BW,
		"latency":   cs.Latency,
		"bitrate":   cs.Bitrate,
		"timestamp": cs.Timestamp,
		"buffer":    cs.BufferDuration,
	}

	pt, err := client.NewPoint(INFLUX_TABLE, tags, fields, time.Now())
	if err != nil {
		log.Println("Error storing new point in InfluxDB: ", err.Error())
	}
	bp.AddPoint(pt)

	if err := clnt.Write(bp); err != nil {
		log.Fatalln("Error: ", err)
	} else {
		log.Println("write complete")
	}
}
