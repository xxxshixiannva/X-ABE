# X-ABE

### Golang Version Tutorial

##### Step1 Install golang 1.21.0
##### Step2 Install pbc, https://pkg.go.dev/github.com/Nik-U/pbc#section-readme
##### Step3 Download go src file on each node, and modify ip address and ip list in 'ipinfo.go'
##### Step4 Get library "github.com/Nik-U/pbc" according to file 'go.mod' and 'go.sum'
##### Step5 On Each Node, execute command 'go run ipinfo.go main.go messagesender.go messagereceiver.go pedersenvss.go feldmanvss.go'. If all nodes are online, key-generating operation will be automatically done
