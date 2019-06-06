#!/bin/bash

# WHERE time> (now() -"$UPPER"d) AND time <= (now()-"$LOWER"d) GROUP BY

table=$1

for UPPER in {200..1}
do
  LOWER=$(expr $UPPER - 1)
  
  outql="SELECT \"duration\",processingDurationInMilliseconds,profile,state,under100,under10k,under120k,under1k,under20k,under2k,under30k,under3k,under40k,under500,under50k,under5k,under60k,under8k,under90k   INTO PREET_ProfileSyncIndividualProfileData FROM ProfileSyncIndividualProfileData WHERE time> (now() -"$UPPER"d) AND time <= (now()-"$LOWER"d) GROUP BY client,environment,host,profileId"

  echo
  echo $outql
done

echo 
echo

for UPPER in {200..1}
do
  LOWER=$(expr $UPPER - 1)
  

  outql="SELECT \"duration\",processingDurationInMilliseconds,profile,state,under100,under10k,under120k,under1k,under20k,under2k,under30k,under3k,under40k,under500,under50k,under5k,under60k,under8k,under90k   INTO ProfileSyncIndividualProfileData FROM PREET_ProfileSyncIndividualProfileData WHERE time> (now() -"$UPPER"d) AND time <= (now()-"$LOWER"d) GROUP BY client,environment,host,profileId"

  echo
  echo $outql
done
