@all:
	go build -o exhibit cmd/exhibit/main.go
	go build -o keycap cmd/keycap/main.go

clean:
	rm -f exhibit keycap


