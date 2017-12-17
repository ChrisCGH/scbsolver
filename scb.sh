#!/bin/bash

file=$1
data_file=$2
crib=$3
limit=$4
verbose=$5
if [ "${limit}" = "" ]
then
    limit=2000000
fi

ct=$(tr -d '\r' < $file | tr '\n' ' ')

if [ "${crib}" = "" ]
then
    post_data='{"Verbose" : '${verbose}', "Ciphertext" : "'${ct}'", "Limit" : "'${limit}'", "Data_file" : "'${data_file}'"}'
else
    post_data='{"Verbose" : '${verbose}', "Ciphertext" : "'${ct}'", "Limit" : "'${limit}'", "Data_file" : "'${data_file}'", "Crib" : "'${crib}'"}'
fi
echo ${post_data}
curl -s -H 'Content-type: application/json' http://localhost:8080/function/scbsolver -X POST --data-binary @- <<EOF
${post_data}
EOF
