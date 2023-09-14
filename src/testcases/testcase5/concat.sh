#!/bin/bash

rm INPUT
touch INPUT

for i in $(seq 0 15)
do
	cat 'input-'$i'.dat' >> INPUT
done

rm OUTPUT
touch OUTPUT

for i in $(seq 0 15)
do
	cat 'output-'$i'.dat' >> OUTPUT
done