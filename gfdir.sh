#!/usr/bin/env bash

gofind_output=$(/usr/local/bin/gofind $@)
if [[ $gofind_output != "Update!" ]];
then
    cd "$gofind_output"
else
    echo "You have successfully updated!"
fi