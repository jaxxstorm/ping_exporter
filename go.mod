module github.com/czerwonk/ping_exporter

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/digineo/go-logwrap v0.0.0-20181106161722-a178c58ea3f0 // indirect
	github.com/digineo/go-ping v0.0.0-20181106162602-34ca5077449a
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/montanaflynn/stats v0.0.0-20180911141734-db72e6cae808 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910 // indirect
	github.com/prometheus/common v0.0.0-20180801064454-c7de2306084e
	github.com/prometheus/procfs v0.0.0-20180920065004-418d78d0b9a7 // indirect
	github.com/sirupsen/logrus v1.0.6 // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20180910181607-0e37d006457b // indirect
	golang.org/x/net v0.0.0-20180921000356-2f5d2388922f // indirect
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	golang.org/x/sys v0.0.0-20180921163948-d47a0f339242 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/yaml.v2 v2.2.1
)

replace github.com/digineo/go-ping => github.com/jaxxstorm/go-ping v0.0.0-20181211164210-7bc8a4db2ba1

replace github.com/digineo/go-ping/monitor => github.com/jaxxstorm/go-ping v0.0.0-20181211164210-7bc8a4db2ba1
