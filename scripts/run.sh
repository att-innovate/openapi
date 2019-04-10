docker run -d --name influxdb \
  -p 8083:8083 -p 8086:8086 \
  -v /<PATH TO CONTAINER DATA>/influxSBGrafanainfluxdb:/var/lib/influxdb \
  influxdb:latest

docker run -d --name grafana \
  -p 3000:3000 \
  -v /<PATH TO CONTAINER DATA>/influxSBGrafana/grafana:/var/lib/grafana  \
  --link influxdb \
  grafana/grafana:3.1.1

sleep 2

curl -i -XPOST http://localhost:8086/query --data-urlencode "q=CREATE DATABASE openapi"

open -a Safari http://localhost:3000

echo " "
echo "Do you want to start random data generator? [Y,n]"
read input
if [[ $input == "Y" || $input == "y" ]]; then
        python p.py
fi
