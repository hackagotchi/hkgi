build:
	go build -o hkgi main.go

run: build
	./hkgi

watch:
	reflex -s -r '\.go$$' make run
