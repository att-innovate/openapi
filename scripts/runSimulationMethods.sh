# OpenAPI project code, MIT license, author: Julius Mueller
# Script to auomatically run through the OpenAPI demo steps
# wireshark filter: (ip.src==192.168.1.131 || ip.dst==192.168.1.131) && http

echo exec simulation with interval: $1
IP=192.168.1.206
TOKEN=3232235982
PORT=8099
THRESHOLD=6


sleep $1
echo
echo register client
curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:$PORT/hello?clientid=$IP


# agreement to start demo in the normal/default network
sleep $1
echo
echo Force device into normal mode as default network to start experiemnt with
curl -sS 'http://<CLIENT IP ADDR>:8080?target_ssid=<SSID>&password=<PASSWORD>'



sleep $1
echo
echo activate streaming mode #1

counter=0
while [ $counter -le $THRESHOLD ]
do
    echo $counter
    sleep $1
    ((counter++))
    curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:$PORT/streaming?client=$TOKEN
done



sleep $1
echo
echo activate normal mode
counter=0
while [ $counter -le $THRESHOLD ]
do
    echo $counter
    sleep $1
    ((counter++))
    curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:$PORT/normal?client=$TOKEN
done

sleep $1
echo
echo activate streaming mode #2
counter=0
while [ $counter -le $THRESHOLD ]
do
    echo $counter
    sleep $1
    ((counter++))
    curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:$PORT/streaming?client=$TOKEN
done



sleep $1
echo
echo De-register client with separate api call
curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:$PORT/goodbye?clientid=$TOKEN
