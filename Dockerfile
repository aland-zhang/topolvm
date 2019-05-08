FROM quay.io/cybozu/ubuntu:18.04

COPY ./_output/hostpathplugin /hostpathplugin
ENTRYPOINT ["/hostpathplugin"]

