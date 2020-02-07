FROM golang:1.13.1-buster as builder

# Example commands for running the image:
# $ docker run -tdi -v $PWD/telegraf.conf:/app/telegraf.conf telegraf
RUN apt-get update && apt-get install -y curl git-lfs && \
	git lfs install && \
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
	cd $GOPATH && mkdir -p src/github.com/influxdata && cd $GOPATH/src/github.com/influxdata && git clone https://github.com/rticommunity/telegraf.git && \
	cd "$GOPATH/src/github.com/influxdata/telegraf" && make && \
	mkdir /app && cp $GOPATH/src/github.com/influxdata/telegraf/telegraf /app/ && \
	cp -rf $GOPATH/src/github.com/influxdata/telegraf/vendor/github.com/rticommunity/rticonnextdds-connector-go/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5 /app/lib && \
	cp $GOPATH/src/github.com/influxdata/telegraf/example_configs/system_to_ddslp.conf /app/telegraf.conf 

FROM ubuntu:18.04

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/lib

COPY --from=builder /app /app

WORKDIR /app
CMD ./telegraf --config telegraf.conf
