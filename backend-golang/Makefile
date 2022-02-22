PROJECT:=AdvSearchLan

.PHONY: build
build: build-yacc
	CGO_ENABLED=0 go build -o adv-search api.go
build-yacc:
	goyacc -o adv_parse/adv_parser.go -p Adv -v adv_parse/adv_parser.output adv_parse/adv.y
