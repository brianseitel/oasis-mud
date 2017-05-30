server:
	colorgo build -o ./bin/gomud *.go

serve:
	colorgo build -o ./bin/gomud *.go
	./bin/gomud
