#!/bin/sh
curl -XPUT "http://${API_CONTAINER}:8080/api/feedinfo"
echo Refresh done!
