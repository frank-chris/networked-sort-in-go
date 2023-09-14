#!/bin/bash

FILES="input-*.dat"

rm INPUT
touch INPUT

for f in $FILES
do
	cat $f >> INPUT
done

FILES="output-*.dat"

rm OUTPUT
touch OUTPUT

for f in $FILES
do
	cat $f >> OUTPUT
done