module icode.baidu.com/k8s/k8s-test/tools/k8s

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	icode.baidu.com/k8s/k8s-test/tools/errorHelper v0.0.0
	k8s.io/api v0.0.0-20191010143144-fbf594f18f80
	k8s.io/apimachinery v0.0.0-20191014065749-fb3eea214746
	k8s.io/client-go v12.0.0+incompatible
)

replace icode.baidu.com/k8s/k8s-test/tools/errorHelper => ../errorHelper
