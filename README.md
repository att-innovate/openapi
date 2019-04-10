# Pro-Active Wireless Network Resource Management and Control

## Compile and Run

    $ go build openapi

    $ ./openapi --ip=<IP ADDRESS POINT> -port=<PORT> -mode=<RAN TYPE>

    $ ./openapi --ip=0.0.0.0 -port=8099 -mode=lte

    $ go build openapi && ./openapi --ip=0.0.0.0 -port=8099 -mode=lte

## Usage - Example API Calls and Responses

#### Hello message request
    $ curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:8099/hello?clientid=1.2.3.4

##### Example Hello message Response
    $ {"token":1234567890,"bw":9876543210,"latency":4}

##### Example Register Request
    $ curl -H "Content-Type: application/json" -X POST -d '{"ip": "1.2.3.4"}' http://0.0.0.0:8099/register

##### Example Streaming Mode Message Request
    $ curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:8099/streaming?client=1234

##### Example Streaming Mode Message Response
    $ {"bw":1000,"latency":20}

##### Example Upload Statistics
    $ curl -H "Content-Type: application/json" -X POST -d '[{"token": "16909062", "bw": "4000", "latency": "7", "bitrate": "1800", "timestamp": "1509138855"}]' http://0.0.0.0:8099/statistics

    $ curl -H "Content-Type: application/json" -X POST -d '[{"token": "16909060", "bw": "600", "latency": "7", "bitrate": "1000", "timestamp": "1509138855"}, {"token": "16909061", "bw": "6000", "latency": "10", "bitrate": "8000", "timestamp": "1509138855"}]' http://0.0.0.0:8099/statistics

##### Example Goodbye Message Request
    $ curl -H "Accept: application/xml" -H "Content-Type: application/xml" -X GET http://0.0.0.0:8099/goodbye?clientid=1234

##### Example Goodbye Message Response
    $ {"bw":1000,"latency":7}
