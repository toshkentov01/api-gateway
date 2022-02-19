#!/bin/bash
CURRENT_DIR=$1
if [[ $CURRENT_DIR == "" ]]; then
    export CURRENT_DIR=`pwd`
fi

echo "Extracting module name from the GOPATH..."
go_path=$GOPATH/src/
size=${#go_path}
go_module=${CURRENT_DIR:size}
let sizes=${#go_path}+${#go_module}
mkdir -p genmock

echo "Generating mock files..."
for x in $(find $go_path$go_module/genproto/* -type f); do
    export file="$(echo ${x:sizes} | cut -d'/' -f 4)"
    export file_name="$(echo ${file} | cut -d'.' -f 1)"

    if [[ ($file == *"service"*) && ($file != *"follow"*) ]]; then
        mkdir -p "./genmock/"$file_name
        mockgen -source=$x > "./genmock/"$file_name"/mock_"$file_name".go"
    fi
done

echo "Sucessfully generated! You are ready to Go :)"
