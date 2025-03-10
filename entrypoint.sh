#!/bin/bash

set -e

function fetch_mrtdump() {
  if [ -f "$mrtDumpPath" ]; then
    echo "File $mrtDumpPath is already exist"
    if test ! `find $mrtDumpPath -mmin +"$intvMins"`; then
      echo "File $mrtDumpPath is still fresh"
      return
    fi

    echo "File $mrtDumpPath exist but outdated, now updating it"
  fi

  echo "[$(date --rfc-3339=seconds)]" "Getting MRTDump data: $mrtBzUrl -> $mrtDumpPath"
  curl -L -o - $mrtBzUrl | bzip2 -d > $mrtDumpPath
}

pwd=$PWD
echo "[$(date --rfc-3339=seconds)]" "Present working directory: $pwd"

if ! [ -d "$pwd/venv" ]; then
  venvDir=$pwd/venv
  echo "[$(date --rfc-3339=seconds)]" "Creating venv dir at $venvDir"
  python3 -m venv $venvDir
fi

echo "[$(date --rfc-3339=seconds)]" "Activating python venv..."
source $pwd/venv/bin/activate

echo "[$(date --rfc-3339=seconds)]" "Installing dependencies..."
python3 -m pip install -r $pwd/requirements.txt

echo "[$(date --rfc-3339=seconds)]" "Building program..."
go build -o $pwd/main $pwd/main.go

intv=${INTV_SECS:-"1800"}
intvMins=$((intv/60))
iter=0
maxIter=${MAX_ITER:-"inf"}
while /bin/true; do
  echo "[$(date --rfc-3339=seconds)]" iter: $iter

  mrtBzUrl=$(cat $pwd/resources/dn42-mrtdump6.mrt.bz2.txt)
  mrtDumpPath=$pwd/data/master6.mrt
  mkdir -p $(dirname $mrtDumpPath)
  fetch_mrtdump

  asPathsDataFile=$pwd/data/aspaths.txt
  echo "[$(date --rfc-3339=seconds)]" "Parsing MRTDump data: $mrtDumpPath -> $asPathsDataFile"
  python3 $pwd/aspaths.py $mrtDumpPath > $asPathsDataFile

  linksDataFile=$pwd/data/links.txt
  echo "[$(date --rfc-3339=seconds)]" "Getting links: $asPathsDataFile -> $linksDataFile"
  cat $asPathsDataFile | python3 tolinks.py > $linksDataFile

  pagerankJSONFile=$pwd/data/pagerank.json
  echo "[$(date --rfc-3339=seconds)]" "Calculating PageRank: $linksDataFile -> $pagerankJSONFile"
  cat $linksDataFile | $pwd/main --json pagerank > $pagerankJSONFile

  echo "[$(date --rfc-3339=seconds)]" "Done. Waiting next round."

  iter=$((iter+1))
  if [ "$iter" = "$maxIter" ]; then
    break
  fi
  nextRun=$(date --rfc-3339=seconds -d "+${intv} seconds")
  echo "[$(date --rfc-3339=seconds)]" "Sleeping for ${intv} seconds, next run: $nextRun"
  sleep $intv
done
