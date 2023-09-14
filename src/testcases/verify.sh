#!/bin/bash

../../utils/linux-amd64/bin/showsort INPUT | sort > REF_OUTPUT
../../utils/linux-amd64/bin/showsort OUTPUT > MY_OUTPUT
diff REF_OUTPUT MY_OUTPUT