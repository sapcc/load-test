#!/usr/bin/env bash
for file in results/*.bin
do
  file_without_extension=${file%.*}
  rate_duration=${file_without_extension#*-}
  # echo "processing file $file_without_extension with rate_duration $rate_duration"
  vegeta report $file > $file_without_extension.txt
  vegeta plot  --title "Rate-Duration: $rate_duration" $file > $file_without_extension.html
done
