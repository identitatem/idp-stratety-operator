module github.com/identitatem/idp-strategy-operator

go 1.16

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v0.4.0
	github.com/identitatem/idp-mgmt-operator v0.0.0-20210816175333-8714ccda628f
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/code-generator v0.20.4
	k8s.io/klog/v2 v2.10.0
	open-cluster-management.io/api v0.0.0-20210804091127-340467ff6239 // indirect
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/controller-tools v0.5.0
)

replace (
	k8s.io/api => k8s.io/api v0.20.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.4
	k8s.io/client-go => k8s.io/client-go v0.20.4
)
