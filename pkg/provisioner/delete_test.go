// +build integration

package provisioner

import (
	storagev1 "k8s.io/api/storage/v1"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

func TestDelete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	p, _ := NewZFSProvisioner(logger)
	options := controller.ProvisionOptions{
		StorageClass: &storagev1.StorageClass{
			Parameters: map[string]string{
				"parentDataset": "test/volumes",
				"shareSubnet":   "10.0.0.0/8",
				"hostname":      "test",
			},
		},
		PVName: "pv-testdelete",
		PVC:    newClaim(resource.MustParse("1G"), []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce, v1.ReadOnlyMany}, nil),
	}
	pv, _ := p.Provision(options) // Already covered by TestProvision

	err := p.Delete(pv)
	assert.NoError(t, err, "Delete should not return an error")

	_, err = os.Stat(pv.Spec.PersistentVolumeSource.NFS.Path)
	assert.Error(t, err, "The volume should not exist on disk")
}
