#!/bin/bash

for i in $(seq 1 10); do
    curl -i "http://localhost:1718/?s=$i"
done
