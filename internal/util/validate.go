package util

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/ceph/ceph-csi/internal/util/k8s"
)

// ValidateNodeStageVolumeRequest validates the node stage request.
func ValidateNodeStageVolumeRequest(req *csi.NodeStageVolumeRequest) error {
	if req.GetVolumeCapability() == nil {
		return status.Error(codes.InvalidArgument, "volume capability missing in request")
	}

	if req.GetVolumeId() == "" {
		return status.Error(codes.InvalidArgument, "volume ID missing in request")
	}

	if req.GetStagingTargetPath() == "" {
		return status.Error(codes.InvalidArgument, "staging target path missing in request")
	}

	if req.GetSecrets() == nil || len(req.GetSecrets()) == 0 {
		return status.Error(codes.InvalidArgument, "stage secrets cannot be nil or empty")
	}

	// Allow auto-creation of staging path when not running on Kubernetes
	if !k8s.RunsOnKubernetes() {
		if err := os.MkdirAll(req.StagingTargetPath, 0755); err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to create staging target path: %v", err))
		}
	} else {
		// validate staging path exists
		ok := checkDirExists(req.GetStagingTargetPath())
		if !ok {
			return status.Errorf(
				codes.InvalidArgument,
				"staging path %s does not exist on node",
				req.GetStagingTargetPath())
		}
	}

	return nil
}

// ValidateNodeUnstageVolumeRequest validates the node unstage request.
func ValidateNodeUnstageVolumeRequest(req *csi.NodeUnstageVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return status.Error(codes.InvalidArgument, "volume ID missing in request")
	}

	if req.GetStagingTargetPath() == "" {
		return status.Error(codes.InvalidArgument, "staging target path missing in request")
	}

	return nil
}

// ValidateNodePublishVolumeRequest validates the node publish request.
func ValidateNodePublishVolumeRequest(req *csi.NodePublishVolumeRequest) error {
	if req.GetVolumeCapability() == nil {
		return status.Error(codes.InvalidArgument, "volume capability missing in request")
	}

	if req.GetVolumeId() == "" {
		return status.Error(codes.InvalidArgument, "volume ID missing in request")
	}

	if req.GetTargetPath() == "" {
		return status.Error(codes.InvalidArgument, "target path missing in request")
	}

	if req.GetStagingTargetPath() == "" {
		return status.Error(codes.InvalidArgument, "staging target path missing in request")
	}

	return nil
}

// ValidateNodeUnpublishVolumeRequest validates the node unpublish request.
func ValidateNodeUnpublishVolumeRequest(req *csi.NodeUnpublishVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return status.Error(codes.InvalidArgument, "volume ID missing in request")
	}

	if req.GetTargetPath() == "" {
		return status.Error(codes.InvalidArgument, "target path missing in request")
	}

	return nil
}

// CheckReadOnlyManyIsSupported checks the request is to create ReadOnlyMany
// volume is from source as empty ReadOnlyMany is not supported.
func CheckReadOnlyManyIsSupported(req *csi.CreateVolumeRequest) error {
	for _, capability := range req.GetVolumeCapabilities() {
		if m := capability.GetAccessMode().GetMode(); m == csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY ||
			m == csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY {
			if req.GetVolumeContentSource() == nil {
				return status.Error(codes.InvalidArgument, "readOnly accessMode is supported only with content source")
			}
		}
	}

	return nil
}
