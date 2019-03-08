Telegraf is an agent written in Go for collecting, processing, aggregating,
and writing metrics.

Design goals are to have a minimal memory footprint with a plugin system so
that developers in the community can easily add support for collecting metrics.  
For an example configuration referencet from local or remote services.

Telegraf is plugin-driven and has the concept of 4 distinct plugins:

1. [Input Plugins](#input-plugins) collect metrics from the system, services, or 3rd party APIs
2. [Processor Plugins](#processor-plugins) transform, decorate, and/or filter metrics
3. [Aggregator Plugins](#aggregator-plugins) create aggregate metrics (e.g. mean, min, max, quantiles, etc.)
4. [Output Plugins](#output-plugins) write metrics to various destinations

For more information on Processor and Aggregator plugins please [read this](./docs/AGGREGATORS_AND_PROCESSORS.md).

It is a forked Telegraf repository adding a DDS input plugin (dds_consumer). 

## Installation:

Telegraf requires golang version 1.9 or newer, the Makefile requires GNU make.

1. [Install Go](https://golang.org/doc/install) >=1.9
   
2. Install dep   

   ```
   cd $GOPATH
   mkdir -p src/github.com/golang
   cd $GOPATH/src/github.com/golang
   git clone https://github.com/pohly/dep.git
   cd dep
   git checkout export-git-submodule
   go install ./cmd/dep
   ```

Current dep does not support git submodule (Go Connector uses git submodule for Connector C library).   
So you need to build dep including the fix. (See this [Issue Link](https://github.com/golang/dep/pull/1909) for details).   

You can find dep at $GOPATH/bin.   
Please add $GOPATH/bin to your $PATH if you haven't.   
   
3. Download Telegraf source:
   ```
   cd $GOPATH
   mkdir -p src/github.com/influxdata
   cd $GOPATH/src/github.com/influxdata
   git clone https://github.com/rticommunity/telegraf.git
   ```
   
4. Run make from the source directory
   ```
   cd "$GOPATH/src/github.com/influxdata/telegraf"
   make
   ```

## How to use it:

#### Include RTI Connector library to the library path (e.g. LD_LIBRARY_PATH)

Currently, RTI Go Connector dynamically links to connector library, so the connectory library should be included in the environment variable for library path (e.g. LD_LIBRARY_PATH). After you built Telegraf, Connector Git repository was checked out under YOUR_TELEGRAF_PATH/vendor. You can include the connector library like the following.

``` 
$ export LD_LIBRARY_PATH=$GOPATH/src/github.com/influxdata/telegraf/vendor/github.com/rticommunity/rticonnextdds-connector-go/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5:$LD_LIBRARY_PATH
```

See usage with:

```
./telegraf --help
```

#### Generate a telegraf config file:

```
./telegraf config > telegraf.conf
```

#### Generate config with DDS input & influxdb output plugins defined:

```
./telegraf --input-filter dds_consumer --output-filter influxdb config
```

#### Run a single telegraf collection, outputing metrics to stdout:

```
./telegraf --config telegraf.conf --test
```

When you run with a DDS input plugin, please make sure that a configuration file for XML Application Creation is at the location configured in your Telegraf configuration (e.g. telegraf.conf).


## Configuration

See the [configuration guide](docs/CONFIGURATION.md) for a rundown of the more advanced
configuration options.

## Input Plugins

* [activemq](./plugins/inputs/activemq)
* [aerospike](./plugins/inputs/aerospike)
* [amqp_consumer](./plugins/inputs/amqp_consumer) (rabbitmq)
* [apache](./plugins/inputs/apache)
* [aurora](./plugins/inputs/aurora)
* [aws cloudwatch](./plugins/inputs/cloudwatch)
* [bcache](./plugins/inputs/bcache)
* [bond](./plugins/inputs/bond)
* [burrow](./plugins/inputs/burrow)
* [cassandra](./plugins/inputs/cassandra) (deprecated, use [jolokia2](./plugins/inputs/jolokia2))
* [ceph](./plugins/inputs/ceph)
* [cgroup](./plugins/inputs/cgroup)
* [chrony](./plugins/inputs/chrony)
* [conntrack](./plugins/inputs/conntrack)
* [consul](./plugins/inputs/consul)
* [couchbase](./plugins/inputs/couchbase)
* [couchdb](./plugins/inputs/couchdb)
* [cpu](./plugins/inputs/cpu)
* [DC/OS](./plugins/inputs/dcos)
* [dds_consumer](./plugins/inputs/dds_consumer)
* [diskio](./plugins/inputs/diskio)
* [disk](./plugins/inputs/disk)
* [disque](./plugins/inputs/disque)
* [dmcache](./plugins/inputs/dmcache)
* [dns query time](./plugins/inputs/dns_query)
* [docker](./plugins/inputs/docker)
* [dovecot](./plugins/inputs/dovecot)
* [elasticsearch](./plugins/inputs/elasticsearch)
* [exec](./plugins/inputs/exec) (generic executable plugin, support JSON, influx, graphite and nagios)
* [fail2ban](./plugins/inputs/fail2ban)
* [fibaro](./plugins/inputs/fibaro)
* [file](./plugins/inputs/file)
* [filestat](./plugins/inputs/filestat)
* [filecount](./plugins/inputs/filecount)
* [fluentd](./plugins/inputs/fluentd)
* [graylog](./plugins/inputs/graylog)
* [haproxy](./plugins/inputs/haproxy)
* [hddtemp](./plugins/inputs/hddtemp)
* [httpjson](./plugins/inputs/httpjson) (generic JSON-emitting http service plugin)
* [http_listener](./plugins/inputs/http_listener)
* [http](./plugins/inputs/http) (generic HTTP plugin, supports using input data formats)
* [http_response](./plugins/inputs/http_response)
* [icinga2](./plugins/inputs/icinga2)
* [influxdb](./plugins/inputs/influxdb)
* [internal](./plugins/inputs/internal)
* [interrupts](./plugins/inputs/interrupts)
* [ipmi_sensor](./plugins/inputs/ipmi_sensor)
* [ipset](./plugins/inputs/ipset)
* [iptables](./plugins/inputs/iptables)
* [jolokia2](./plugins/inputs/jolokia2) (java, cassandra, kafka)
* [jolokia](./plugins/inputs/jolokia) (deprecated, use [jolokia2](./plugins/inputs/jolokia2))
* [jti_openconfig_telemetry](./plugins/inputs/jti_openconfig_telemetry)
* [kafka_consumer](./plugins/inputs/kafka_consumer)
* [kapacitor](./plugins/inputs/kapacitor)
* [kernel](./plugins/inputs/kernel)
* [kernel_vmstat](./plugins/inputs/kernel_vmstat)
* [kubernetes](./plugins/inputs/kubernetes)
* [leofs](./plugins/inputs/leofs)
* [linux_sysctl_fs](./plugins/inputs/linux_sysctl_fs)
* [logparser](./plugins/inputs/logparser)
* [lustre2](./plugins/inputs/lustre2)
* [mailchimp](./plugins/inputs/mailchimp)
* [mcrouter](./plugins/inputs/mcrouter)
* [memcached](./plugins/inputs/memcached)
* [mem](./plugins/inputs/mem)
* [mesos](./plugins/inputs/mesos)
* [minecraft](./plugins/inputs/minecraft)
* [mongodb](./plugins/inputs/mongodb)
* [mqtt_consumer](./plugins/inputs/mqtt_consumer)
* [mysql](./plugins/inputs/mysql)
* [nats_consumer](./plugins/inputs/nats_consumer)
* [nats](./plugins/inputs/nats)
* [net](./plugins/inputs/net)
* [net_response](./plugins/inputs/net_response)
* [netstat](./plugins/inputs/net)
* [nginx](./plugins/inputs/nginx)
* [nginx_plus](./plugins/inputs/nginx_plus)
* [nsq_consumer](./plugins/inputs/nsq_consumer)
* [nsq](./plugins/inputs/nsq)
* [nstat](./plugins/inputs/nstat)
* [ntpq](./plugins/inputs/ntpq)
* [nvidia_smi](./plugins/inputs/nvidia_smi)
* [openldap](./plugins/inputs/openldap)
* [opensmtpd](./plugins/inputs/opensmtpd)
* [pf](./plugins/inputs/pf)
* [pgbouncer](./plugins/inputs/pgbouncer)
* [phpfpm](./plugins/inputs/phpfpm)
* [phusion passenger](./plugins/inputs/passenger)
* [ping](./plugins/inputs/ping)
* [postfix](./plugins/inputs/postfix)
* [postgresql_extensible](./plugins/inputs/postgresql_extensible)
* [postgresql](./plugins/inputs/postgresql)
* [powerdns](./plugins/inputs/powerdns)
* [processes](./plugins/inputs/processes)
* [procstat](./plugins/inputs/procstat)
* [prometheus](./plugins/inputs/prometheus) (can be used for [Caddy server](./plugins/inputs/prometheus/README.md#usage-for-caddy-http-server))
* [puppetagent](./plugins/inputs/puppetagent)
* [rabbitmq](./plugins/inputs/rabbitmq)
* [raindrops](./plugins/inputs/raindrops)
* [redis](./plugins/inputs/redis)
* [rethinkdb](./plugins/inputs/rethinkdb)
* [riak](./plugins/inputs/riak)
* [salesforce](./plugins/inputs/salesforce)
* [sensors](./plugins/inputs/sensors)
* [smart](./plugins/inputs/smart)
* [snmp_legacy](./plugins/inputs/snmp_legacy)
* [snmp](./plugins/inputs/snmp)
* [socket_listener](./plugins/inputs/socket_listener)
* [solr](./plugins/inputs/solr)
* [sql server](./plugins/inputs/sqlserver) (microsoft)
* [statsd](./plugins/inputs/statsd)
* [swap](./plugins/inputs/swap)
* [syslog](./plugins/inputs/syslog)
* [sysstat](./plugins/inputs/sysstat)
* [system](./plugins/inputs/system)
* [tail](./plugins/inputs/tail)
* [tcp_listener](./plugins/inputs/socket_listener)
* [teamspeak](./plugins/inputs/teamspeak)
* [tengine](./plugins/inputs/tengine)
* [tomcat](./plugins/inputs/tomcat)
* [twemproxy](./plugins/inputs/twemproxy)
* [udp_listener](./plugins/inputs/socket_listener)
* [unbound](./plugins/inputs/unbound)
* [varnish](./plugins/inputs/varnish)
* [webhooks](./plugins/inputs/webhooks)
  * [filestack](./plugins/inputs/webhooks/filestack)
  * [github](./plugins/inputs/webhooks/github)
  * [mandrill](./plugins/inputs/webhooks/mandrill)
  * [papertrail](./plugins/inputs/webhooks/papertrail)
  * [particle](./plugins/inputs/webhooks/particle)
  * [rollbar](./plugins/inputs/webhooks/rollbar)
* [win_perf_counters](./plugins/inputs/win_perf_counters) (windows performance counters)
* [win_services](./plugins/inputs/win_services)
* [zfs](./plugins/inputs/zfs)
* [zipkin](./plugins/inputs/zipkin)
* [zookeeper](./plugins/inputs/zookeeper)
* [dds_consumer](./plugins/inputs/dds_consumer)

Telegraf is able to parse the following input data formats into metrics, these
formats may be used with input plugins supporting the `data_format` option:

* [InfluxDB Line Protocol](./docs/DATA_FORMATS_INPUT.md#influx)
* [JSON](./docs/DATA_FORMATS_INPUT.md#json)
* [Graphite](./docs/DATA_FORMATS_INPUT.md#graphite)
* [Value](./docs/DATA_FORMATS_INPUT.md#value)
* [Nagios](./docs/DATA_FORMATS_INPUT.md#nagios)
* [Collectd](./docs/DATA_FORMATS_INPUT.md#collectd)
* [Dropwizard](./docs/DATA_FORMATS_INPUT.md#dropwizard)

## Processor Plugins

* [converter](./plugins/processors/converter)
* [override](./plugins/processors/override)
* [printer](./plugins/processors/printer)
* [regex](./plugins/processors/regex)
* [rename](./plugins/processors/rename)
* [topk](./plugins/processors/topk)

## Aggregator Plugins

* [basicstats](./plugins/aggregators/basicstats)
* [minmax](./plugins/aggregators/minmax)
* [histogram](./plugins/aggregators/histogram)
* [valuecounter](./plugins/aggregators/valuecounter)

## Output Plugins

* [influxdb](./plugins/outputs/influxdb)
* [amon](./plugins/outputs/amon)
* [amqp](./plugins/outputs/amqp) (rabbitmq)
* [application_insights](./plugins/outputs/application_insights)
* [aws kinesis](./plugins/outputs/kinesis)
* [aws cloudwatch](./plugins/outputs/cloudwatch)
* [cratedb](./plugins/outputs/cratedb)
* [datadog](./plugins/outputs/datadog)
* [discard](./plugins/outputs/discard)
* [elasticsearch](./plugins/outputs/elasticsearch)
* [file](./plugins/outputs/file)
* [graphite](./plugins/outputs/graphite)
* [graylog](./plugins/outputs/graylog)
* [http](./plugins/outputs/http)
* [instrumental](./plugins/outputs/instrumental)
* [kafka](./plugins/outputs/kafka)
* [librato](./plugins/outputs/librato)
* [mqtt](./plugins/outputs/mqtt)
* [nats](./plugins/outputs/nats)
* [nsq](./plugins/outputs/nsq)
* [opentsdb](./plugins/outputs/opentsdb)
* [prometheus](./plugins/outputs/prometheus_client)
* [riemann](./plugins/outputs/riemann)
* [riemann_legacy](./plugins/outputs/riemann_legacy)
* [socket_writer](./plugins/outputs/socket_writer)
* [tcp](./plugins/outputs/socket_writer)
* [udp](./plugins/outputs/socket_writer)
* [wavefront](./plugins/outputs/wavefront)
