build-rpi:
	GOOS=linux GOARCH=arm GOARM=6 go build -o rpi-ws .

run-test-app:
	go run . -fake-bme -term-screen -update-delay 1s
