module sigs.k8s.io/cluster-api

go 1.16

replace github.com/vmware-tanzu/cluster-api-provider-byoh => /Users/anushah/Documents/byoh/cluster-api-provider-byoh/

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/coredns/corefile-migration v1.0.11
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/drone/envsubst/v2 v2.0.0-20210305151453-490366e43a3c
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/fatih/color v1.10.0
	github.com/go-logr/logr v0.4.0
	github.com/gobuffalo/flect v0.2.2
	github.com/google/go-cmp v0.5.5
	github.com/google/go-github/v33 v33.0.0
	github.com/google/gofuzz v1.2.0
	github.com/gosuri/uitable v0.0.4
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/vmware-tanzu/cluster-api-provider-byoh v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.27.1
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/apiserver v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/cluster-bootstrap v0.21.1
	k8s.io/component-base v0.21.1
	k8s.io/klog/v2 v2.8.0
	k8s.io/kubectl v0.21.1
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10
	sigs.k8s.io/controller-runtime v0.9.0-beta.5
	sigs.k8s.io/kind v0.11.0
	sigs.k8s.io/yaml v1.2.0
)
