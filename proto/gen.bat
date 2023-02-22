@echo on
protoc --go_out=./vkit_example --vkit_out=./vkit_example/ --vkit_opt=--handlePath=../handler --swagger_out=./ --validate_out="lang=go:./vkit_example" .\vkit_example\vkit_example.proto

exit

