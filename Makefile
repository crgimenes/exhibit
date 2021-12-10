@all:
	go build -o exhibit cmd/exhibit/main.go
	go build -o keycap cmd/keycap/main.go
	go build -o mdtoansi cmd/mdtoansi/main.go

clean:
	rm -f exhibit keycap mdtoansi


