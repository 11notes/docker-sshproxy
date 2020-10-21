#!/bin/ash
curl --max-time 5 -kILs --fail telnet://localhost:$SSHPIPERD_PORT