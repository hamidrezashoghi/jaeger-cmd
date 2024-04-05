# jaeger-cmd
Some related jobs to jaeger

### 1. restart-jaeger-collector

##### First create D-Bus connection

The `dbus` package connects to the [systemd D-Bus API](http://www.freedesktop.org/wiki/Software/systemd/dbus/) and lets you start, stop and introspect systemd units.
[API documentation][dbus-doc] is available online.

[dbus-doc]: https://pkg.go.dev/github.com/coreos/go-systemd/v22/dbus?tab=doc

Create `/etc/dbus-1/system-local.conf` that looks like this:

```
<!DOCTYPE busconfig PUBLIC
"-//freedesktop//DTD D-Bus Bus Configuration 1.0//EN"
"http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd">
<busconfig>
    <policy user="root">
        <allow eavesdrop="true"/>
        <allow eavesdrop="true" send_destination="*"/>
    </policy>
</busconfig>
```

##### Create jaeger-collector.service unitd

Create `/etc/systemd/system/jaeger-collector.service` that looks like this:
```
[Unit]
Description=Jaeger Collector
Requires=docker.service
Before=multi-user.target
After=docker.service
Wants=network-online.target

[Service]
Type=oneshot
WorkingDirectory=/opt/collector/
ExecStart=/usr/local/bin/docker-compose -f /opt/collector/collector.yaml up -d --remove-orphans
ExecStop=/usr/local/bin/docker-compose -f /opt/collector/collector.yaml down
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
```

### 2. consumers-status

##### Add kafka servers to /etc/hosts
```
192.168.1.2 kafka-01.local1.server
192.168.1.3 kafka-02.local2.server
192.168.1.4 kafka-03.local3.server
```

* You should be able to connect to the kafka brokers via port 9092

##### Download kafka and extract it in /usr/local/

```
sudo -i
wget -c wget -c https://downloads.apache.org/kafka/3.6.0/kafka_2.13-3.6.0.tgz
tar -xzf kafka_2.13-3.6.0.tgz -C /usr/local
mv /usr/local/kafka_2.13-3.6.0 /usr/local/kafka
```

##### Why node_exporter?
If you put a prom file in `/var/lib/node_exporter/textfile_collector/` path your metrics expose along other

### 3. architecture
Generate `System Architecture` in Jaeger UI
https://github.com/jaegertracing/spark-dependencies/tree/main

### 4. restart-cassandra
Restart Cassandra when leaves the cluster
