# exec-ttynvt

The ttynvt executer is used to observe _ttynvt._tcp mdns services of io4edge devices. It starts a new ttynvt instance when a new service shows up and terminates the corresponding instance again when the service disappears.
The ttynvt creates for each _ttynvt._tcp mdns service a simulated tty `/dev/tty<mdns-instance-name>`.

For more information about ttynvt see https://gitlab.com/ci4rail/ttynvt.

# Usage
```
$ sudo ./exec-ttynvt <ttynvt-program-path> <major-driver-number>
```
