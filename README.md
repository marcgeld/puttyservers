# Puttyservers
Dumps Putty session (servers) values from Windows registry to text/json
A small utility thing that exports [Putty](https://www.putty.org/) server list from Windows Registry to json or text.

## Build
`go build -o puttyservers.exe main.go`

## Run
Write a json file with filename 'servers.json'
```shell
puttyservers.exe -json -filename=servers.json
```

Write json output to screen
```shell
puttyservers.exe -json
```

Write text output to screen
```shell
puttyservers.exe
```

Show help text
```shell
puttyservers.exe -h
```

## Clean
`go clean`