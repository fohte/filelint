#!/bin/bash

go generate
goxc -d=./dist -tasks=clean-destination,xc,archive,rmbin -bc="windows linux darwin"
