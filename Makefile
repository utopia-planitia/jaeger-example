
run:
	CompileDaemon -build "go build ./cmd/midas-phalanx" -command "./midas-phalanx" -run-dir . -command-stop -graceful-kill -graceful-timeout 10 -color=true -log-prefix=false

test-loop:
	watch --color ./hack/test.sh

curl:
	curl --silent -v 127.0.0.1:8080
