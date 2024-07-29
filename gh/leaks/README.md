Ghleaks parses GitHub events for repos that were made public, clones them and searches them for leaked credentials.

```
wget https://data.gharchive.org/2024-07-25-{15..17}.json.gz # events from 3pm to 5pm UTC
gunzip *.gz
cat 2024-07-25-{15..17}.json | go run main.go | tee ghleaks.json
```

Alternatively you can use something like https://github.com/WillAbides/gharchive-client to get GitHub events:

```
gharchive-client 2024-07-25 --type=PublicEvent | go run main.go
```

Analysing leaks:

```
jq -r '.Leaks[] | .RuleID' < ghleaks.json  | sort | uniq -c
```