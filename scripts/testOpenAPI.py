import requests, random, time

url = 'http://0.0.0.0:8099/statistics'


curtime = 0
tput_mean = [15, 18, 14, 16, 20, 11]
tput_var = [0.1, 0.1, 0.1, 0.1, 0.1, 0.1]

lat_mean = [10, 30, 25, 12, 33, 23]
lat_var  = [1, 1, 1, 1, 1, 1]

bw_mean = [30, 35, 25, 37, 42, 23]
bw_var  = [0.5, 0.5, 0.5, 0.5, 0.5, 0.5]

token = [3232235977, 3232235982, 16909060]
#token = ["ue1", "ue2"]
#token = [16309060]


while True:

	payload = "["

	for i in range(len(token)):

		tput = int(random.normalvariate(9000,100))
		latency = int(random.normalvariate(lat_mean[i],lat_var[i])-5)
		bw = int(random.normalvariate(10000,150))
		ts = int(int(time.time()*100000))
		tput2 = int(random.normalvariate(1000,100))
		latency2 = int(random.normalvariate(lat_mean[i],lat_var[i]))
		bw2 = int(random.normalvariate(3000,100))
		ts2 = int(int(time.time()*100000))
		tput3 = int(random.normalvariate(2000,10))
		latency3 = int(random.normalvariate(lat_mean[i],lat_var[i])+10)
		bw3 = int(random.normalvariate(5000,100))
		ts3 = int(int(time.time()*100000))
		
		# single entry
		#datapoint = {"token": str(token[0]), "bw": str(bw), "latency": str(latency), "bitrate": str(tput), "timestamp": str(ts)}
		
		# double entry
		#datapoint = {"token": str(token[0]), "bw": str(bw), "latency": str(latency), "bitrate": str(tput), "timestamp": str(ts)}, {"token": str(token[1]), "bw": str(bw2), "latency": str(latency2), "bitrate": str(tput2), "timestamp": str(ts2)}
		
		#tripple entry
		datapoint = {"token": str(token[0]), "bw": str(bw), "latency": str(latency), "bitrate": str(tput), "timestamp": str(ts)}, {"token": str(token[1]), "bw": str(bw2), "latency": str(latency2), "bitrate": str(tput2), "timestamp": str(ts2)}, {"token": str(token[2]), "bw": str(bw3), "latency": str(latency3), "bitrate": str(tput3), "timestamp": str(ts3)}
		
		payload = payload + str(datapoint)

		if i < len(token)-1:
			payload = payload + ", "
		
	payload = payload + "]"

	print datapoint

	r = requests.post(url=url, json=datapoint)

	print r

	curtime += 1

	if (curtime%50) == 0:
		tput_mean[0] = 15
		lat_mean[0] = 10


	if (curtime%150) == 0:
		tput_mean[0] = 1.5
		lat_mean[0] = 20

	time.sleep(1)
