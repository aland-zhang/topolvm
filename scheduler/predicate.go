package scheduler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cybozu-go/topolvm"
	corev1 "k8s.io/api/core/v1"
)

func filterNodes(nodes corev1.NodeList, requested int64, spare uint64) ExtenderFilterResult {
	if requested <= 0 {
		return ExtenderFilterResult{
			Nodes: &nodes,
		}
	}

	required := spare + uint64(requested)
	filtered := corev1.NodeList{}
	failed := FailedNodesMap{}

	for _, node := range nodes.Items {
		val, ok := node.Annotations[topolvm.CapacityKey]
		if !ok {
			failed[node.Name] = "no capacity annotation"
			continue
		}
		capacity, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			failed[node.Name] = "bad capacity annotation: " + val
			continue
		}

		if capacity < required {
			failed[node.Name] = "out of VG free space"
			continue
		}

		filtered.Items = append(filtered.Items, node)
	}
	return ExtenderFilterResult{
		Nodes:       &filtered,
		FailedNodes: failed,
	}
}

func extractRequestedSize(pod *corev1.Pod) int64 {
	for _, container := range pod.Spec.Containers {
		for k, v := range container.Resources.Limits {
			if k == topolvm.CapacityResource {
				return v.Value()
			}
		}
		for k, v := range container.Resources.Requests {
			if k == topolvm.CapacityResource {
				return v.Value()
			}
		}
	}

	return 0
}

func (s scheduler) predicate(w http.ResponseWriter, r *http.Request) {
	var input ExtenderArgs

	reader := http.MaxBytesReader(w, r.Body, 10<<20)
	err := json.NewDecoder(reader).Decode(&input)
	if err != nil || input.Nodes == nil || input.Pod == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	requested := extractRequestedSize(input.Pod)
	result := filterNodes(*input.Nodes, requested, s.spareGB<<30)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}
