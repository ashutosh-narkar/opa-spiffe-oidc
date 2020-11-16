#!/bin/sh
backend-app -log /tmp/web-server.log &
/usr/local/bin/envoy -l debug -c /etc/envoy/envoy.yaml --log-path /tmp/envoy.log
