package topolvm

import corev1 "k8s.io/api/core/v1"

// CapacityKey is a key of Node annotation that represents VG free space.
const CapacityKey = "topolvm.cybozu.com/capacity"

// CapacityResource is the resource name of topolvm capacity.
const CapacityResource = corev1.ResourceName("topolvm.cybozu.com/capacity")
