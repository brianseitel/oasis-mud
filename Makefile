server:
	go build -o ./bin/gomud *.go

serve:
	go build -o ./bin/gomud *.go
	./bin/gomud
