module github.com/lynkdb/kvgo-server

go 1.15

// replace github.com/lynkdb/kvgo v1.1.5 => /opt/workspace/src/github.com/lynkdb/kvgo
// replace github.com/hooto/hmetrics v0.0.1 => /opt/workspace/src/github.com/hooto/hmetrics

require (
	github.com/hooto/hlog4g v0.9.4
	github.com/hooto/hmetrics v0.0.1
	github.com/hooto/htoml4g v0.9.4
	github.com/lufia/plan9stats v0.0.0-20220517141722-cf486979b281 // indirect
	github.com/lynkdb/kvgo v1.1.6
	github.com/lynkdb/kvspec/v2 v2.0.4
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/shirou/gopsutil/v3 v3.23.2
	github.com/spf13/cobra v1.2.1
)
