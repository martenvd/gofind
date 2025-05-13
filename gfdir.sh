#!/usr/bin/env bash

gofind_output=$(/usr/local/bin/gofind $@)
if [[ $gofind_output == "Update!" ]];
then
    echo "You have successfully updated!"
elif [[ $gofind_output != "" ]];
then
    echo "cd $gofind_output"
    cd "$gofind_output"
fi
