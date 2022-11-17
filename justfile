set export
startfile := "cmd/gofeeder.go"
run:
 go run $startfile

install:
 go install $startfile

config:
  cp ./gofeeder_sample.json ~/.gofeeder.json
