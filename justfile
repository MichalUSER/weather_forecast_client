rasp:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o bin ./src
	scp -rp bin pi@192.168.0.110:~/weather-forecast/server
