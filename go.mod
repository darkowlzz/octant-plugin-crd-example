module github.com/darkowlzz/octant-plugin-crd-example

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/vmware-tanzu/octant v0.0.0-00010101000000-000000000000
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)

// Use a fork of octant with k8s v1.20.2 dependencies.
replace github.com/vmware-tanzu/octant => github.com/darkowlzz/octant v0.20.1-0.20210613125003-3f821bf41488
