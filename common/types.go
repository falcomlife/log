package common

import (
	v1 "k8s.io/api/apps/v1"
)

type HttpClient struct {
	Name     string
	Protocol string
	Host     string
	Port     string
}
type ReplicaSetSlice []v1.ReplicaSet

func (r ReplicaSetSlice) Len() int {
	return len(r)
}
func (r ReplicaSetSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ReplicaSetSlice) Less(i, j int) bool {
	return !r[i].CreationTimestamp.Before(&r[j].CreationTimestamp)
}
