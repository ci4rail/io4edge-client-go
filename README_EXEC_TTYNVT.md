# exec_ttynvt

The ttynvt executer is used to observe _ttynvt._tcp mdns services of io4edge devices. It starts a new ttynvt instance when a new service shows up and terminates the corresponding instance again when the service disappears.

For more information about ttynvt see https://gitlab.com/ci4rail/ttynvt.

# Usage

$ sudo ./exec_ttynvt <major-driver-number>
