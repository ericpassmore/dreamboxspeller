#!/usr/bin/env bash
# PURPOSE: builds go scripts and any onetime setup

###################### PRESET VARS ##################
GO_SEARCH_PATH="/usr/local"

##################### FUNCTIONS #####################
##########################
# CHECK: IS GO INSTALLED #
##########################
check_for_go () {
  status_code=1 # failed by default
  version=$(go version)
  # 127 is command not found
  if [ "$?" -eq 127 ]
  then
    # maybe just bad path look in obvious /usr/local/go
    possible_go=$(ls -1d $GO_SEARCH_PATH/go*)
    if [ "$?" -eq 0 ]
    then
      printf "Found: "
      # return first match
      echo $possible_go | head -1
      echo set PATH and alias to run 'go'
    else # go lang not found user to fix
      echo "Go language not found"
      echo "please download and install see https://golang.org/dl/ "
      echo "please set PATH and alias to run 'go'"
      return 127
    fi
  else # all ok found go
    echo "Found $version"
    return 0
  fi
  return 1
}

####################### MAIN #########################
check_for_go
if [ $? -ne 0 ]
then
  # return status from func
  exit $?
fi
# set path for imports
GOPATH=$(pwd); export GOPATH
# build stuff
go build dreamboxspeller.go
echo "Build Complete"
echo "run '$ ./dreamboxspeller' and goto http://localhost:8080 "
