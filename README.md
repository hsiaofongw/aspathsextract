# Guides of How To

Initialize Python virtual environment: (only do once per project)

```
python3 -m venv $PWD/venv
python3 -m pip install -r requirements.txt
```

Enter Python venv: (do this everytime after shell is restarted)

```
source venv/bin/activate
```

Download Route Map (in format of bzipped MRTDUMP)

```
curl -Lo - https://mrt.collector.dn42/master6_latest.mrt.bz2 | \
  bzip2 -d > master6_latest.mrt
```

Preview MRTDUMP data:

```
python3 mrt2json.py master6_latest.mrt 10 | jq
```

Get BGP ASPaths:

```
python3 aspaths.py master6_latest.mrt > aspaths.txt
```

Get Links:

```
cat aspaths.txt | python3 tolinks.py > links.txt
```

Analyze (seeing helps first):

```
go run main.go --help
```
