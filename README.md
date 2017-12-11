# caddy-licserver
License server plugin for Caddy server (base on https://github.com/SaturnsVoid/HWID-Based-License-System)

licensemanager.exe - utilities for manage license db, commands:
list, add, add bulk, remove, exit


Make this steps for compilation caddy with plugin caddy-licserver:
- add directive in var section: "licserver", // github.com/linkonoid/caddy-licserver 
(in file github.com\mholt\caddy\caddyhttp\httpserver\plugin.go)
- add in import section _ "github.com/linkonoid/caddy-licserver" (in caddymain/run.go)
- add directives in Caddyfile ()

