#!/bin/bash
# check all parents and childs have context (cancelable)
go vet
go get muxator/tor
go get muxator/proxy
