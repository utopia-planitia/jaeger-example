
run:
	CompileDaemon -build "go build ./cmd/midas-phalanx" -command "./midas-phalanx" -run-dir . -command-stop -graceful-kill -graceful-timeout 10 -color=true -log-prefix=false

test-loop:
	watch --color ./hack/test.sh

curl:
	curl --silent -v 127.0.0.1:8080

# https://medium.com/opentracing/take-opentracing-for-a-hotrod-ride-f6e3141f7941
jaeger:
	docker pull jaegertracing/all-in-one:latest
	docker run --name jaeger -d -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
