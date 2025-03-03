#!/bin/bash

for i in {1..1000}
do
  curl -s http://172.18.0.2:30000 >> responses.txt
done

grep -o 'Hello from [^!]*' responses.txt | sort | uniq -c