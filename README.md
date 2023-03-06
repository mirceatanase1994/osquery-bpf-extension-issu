# osquery-bpf-extension-issue

This repo contains the necessary tools for reproducing an osquery extension bpf events issue.

The repo contains:
* The source code for the osquery table plugin, which is very similar to the osquery-go extension example, the main difference being it creates and populates 3 tables instead of one
* The source code for the osquery logger plugin, which posts all the logs to a web server
* The osquery conf and flags files necessary for reproducing the issue
* Scripts for building the extensions and preparing the issue reproduction setup

## Build the extensions from source

### Prerequisites: go 1.18

The pre-built extensions are available in the repo, but they can be built from source using the build.sh script:

`./build.sh`

## Prepare the issue reproduction setup

### Prerequisites: osquery 5.7.0

The prepare.sh script adds the extensions to osquery and adds the conf and flags files necessary for reproducing the issue:

`sudo ./prepare.sh`


## Reproduce the issue

Start osqueryctl

`sudo osqueryctl start`

Watch the osquery service logs and wait...

`sudo journalct -u osqueryd.service -f`
