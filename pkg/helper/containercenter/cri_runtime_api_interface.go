// Copyright 2025 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package containercenter

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba/ilogtail/pkg/logger"

	"google.golang.org/grpc"
)

type CriContainerState int32

const (
	ContainerStateContainerCreated CriContainerState = 0
	ContainerStateContainerRunning CriContainerState = 1
	ContainerStateContainerExited  CriContainerState = 2
	ContainerStateContainerUnknown CriContainerState = 3
	// v2 新增状态
	ContainerStateContainerCreated2 CriContainerState = 4
	ContainerStateContainerRunning2 CriContainerState = 5
	ContainerStateContainerExited2  CriContainerState = 6
	ContainerStateContainerUnknown2 CriContainerState = 7
)

type CriMountPropagation int32

const (
	MountPropagationPropagationPrivate         CriMountPropagation = 0
	MountPropagationPropagationHostToContainer CriMountPropagation = 1
	MountPropagationPropagationBidirectional   CriMountPropagation = 2
)

type CriVersionInfo struct {
	Version           string
	RuntimeName       string
	RuntimeVersion    string
	RuntimeAPIVersion string
}

type CriVersionResponse struct {
	Version           string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	RuntimeName       string `protobuf:"bytes,2,opt,name=runtime_name,proto3" json:"runtime_name,omitempty"`
	RuntimeVersion    string `protobuf:"bytes,3,opt,name=runtime_version,proto3" json:"runtime_version,omitempty"`
	RuntimeAPIVersion string `protobuf:"bytes,4,opt,name=runtime_api_version,proto3" json:"runtime_api_version,omitempty"`
}

type CriContainerMetadata struct {
	Name    string
	Attempt uint32
}

type CriImageSpec struct {
	Image       string
	Annotations map[string]string
}

// v2 新增：资源限制
type CriLinuxContainerResources struct {
	CPUPeriod              int64             `json:"cpu_period,omitempty"`
	CPUQuota               int64             `json:"cpu_quota,omitempty"`
	CPUShares              int64             `json:"cpu_shares,omitempty"`
	MemoryLimitInBytes     int64             `json:"memory_limit_in_bytes,omitempty"`
	OomScoreAdj            int64             `json:"oom_score_adj,omitempty"`
	CpusetCpus             string            `json:"cpuset_cpus,omitempty"`
	CpusetMems             string            `json:"cpuset_mems,omitempty"`
	HugepageLimits         map[string]uint64 `json:"hugepage_limits,omitempty"`
	Unified                map[string]string `json:"unified,omitempty"`
	MemorySwapLimitInBytes int64             `json:"memory_swap_limit_in_bytes,omitempty"`
}

// v2 新增：容器安全上下文
type CriLinuxContainerSecurityContext struct {
	Capabilities       *CriCapability      `json:"capabilities,omitempty"`
	Privileged         bool                `json:"privileged,omitempty"`
	NamespaceOptions   *CriNamespaceOption `json:"namespace_options,omitempty"`
	SelinuxOptions     *CriSELinuxOption   `json:"selinux_options,omitempty"`
	RunAsUser          *CriInt64Value      `json:"run_as_user,omitempty"`
	RunAsGroup         *CriInt64Value      `json:"run_as_group,omitempty"`
	RunAsUsername      string              `json:"run_as_username,omitempty"`
	ReadonlyRootfs     bool                `json:"readonly_rootfs,omitempty"`
	SupplementalGroups []int64             `json:"supplemental_groups,omitempty"`
	ApparmorProfile    string              `json:"apparmor_profile,omitempty"`
	SeccompProfilePath string              `json:"seccomp_profile_path,omitempty"`
	NoNewPrivs         bool                `json:"no_new_privs,omitempty"`
	MaskedPaths        []string            `json:"masked_paths,omitempty"`
	ReadonlyPaths      []string            `json:"readonly_paths,omitempty"`
}

type CriCapability struct {
	AddCapabilities  []string `json:"add_capabilities,omitempty"`
	DropCapabilities []string `json:"drop_capabilities,omitempty"`
}

type CriNamespaceOption struct {
	Network  CriNamespaceMode `json:"network,omitempty"`
	Pid      CriNamespaceMode `json:"pid,omitempty"`
	Ipc      CriNamespaceMode `json:"ipc,omitempty"`
	TargetId string           `json:"target_id,omitempty"`
}

type CriNamespaceMode int32

const (
	NamespaceModePod       CriNamespaceMode = 0
	NamespaceModeContainer CriNamespaceMode = 1
	NamespaceModeNode      CriNamespaceMode = 2
	NamespaceModeTarget    CriNamespaceMode = 3
)

type CriSELinuxOption struct {
	User  string `json:"user,omitempty"`
	Role  string `json:"role,omitempty"`
	Type  string `json:"type,omitempty"`
	Level string `json:"level,omitempty"`
}

type CriInt64Value struct {
	Value int64 `json:"value,omitempty"`
}

type CriContainer struct {
	ID           string
	PodSandboxID string
	Metadata     *CriContainerMetadata
	Image        *CriImageSpec
	ImageRef     string
	State        CriContainerState
	CreatedAt    int64
	Labels       map[string]string
	Annotations  map[string]string
	// v2 新增字段
	Resources       *CriLinuxContainerResources
	SecurityContext *CriLinuxContainerSecurityContext
}

type CriListContainersResponse struct {
	Containers []*CriContainer
}

type CriMount struct {
	ContainerPath  string
	HostPath       string
	Readonly       bool
	SelinuxRelabel bool
	Propagation    CriMountPropagation
}

type CriContainerStatus struct {
	ID          string
	Metadata    *CriContainerMetadata
	State       CriContainerState
	CreatedAt   int64
	StartedAt   int64
	FinishedAt  int64
	ExitCode    int32
	Image       *CriImageSpec
	ImageRef    string
	Reason      string
	Message     string
	Labels      map[string]string
	Annotations map[string]string
	Mounts      []*CriMount
	LogPath     string
	// v2 新增字段
	Resources       *CriLinuxContainerResources
	SecurityContext *CriLinuxContainerSecurityContext
}

type CriContainerStatusResponse struct {
	Status *CriContainerStatus
	Info   map[string]string
}

type CriPodSandboxMetadata struct {
	Name      string
	UID       string
	Namespace string
	Attempt   uint32
}

type CriPodSandbox struct {
	ID             string
	Metadata       *CriPodSandboxMetadata
	State          CriContainerState
	CreatedAt      int64
	Labels         map[string]string
	Annotations    map[string]string
	RuntimeHandler string
}

type CriListPodSandboxResponse struct {
	Items []*CriPodSandbox
}

type CriPodSandboxStatus struct {
	ID             string
	Metadata       *CriPodSandboxMetadata
	State          CriContainerState
	CreatedAt      int64
	Labels         map[string]string
	Annotations    map[string]string
	RuntimeHandler string
}

type CriPodSandboxStatusResponse struct {
	Status *CriPodSandboxStatus
	Info   map[string]string
}

type RuntimeService interface {
	Version(ctx context.Context) (*CriVersionResponse, error)
	ListContainers(ctx context.Context) (*CriListContainersResponse, error)
	ContainerStatus(ctx context.Context, containerID string, verbose bool) (*CriContainerStatusResponse, error)
	ListPodSandbox(ctx context.Context) (*CriListPodSandboxResponse, error)
	PodSandboxStatus(ctx context.Context, sandboxID string, verbose bool) (*CriPodSandboxStatusResponse, error)
}

type RuntimeServiceClient struct {
	service RuntimeService
	info    CriVersionInfo
	conn    *grpc.ClientConn
}

var ( // for mock
	getAddressAndDialer = GetAddressAndDialer
	grpcDialContext     = grpc.DialContext
)

func NewRuntimeServiceClient(contextTimeout time.Duration, grpcMaxCallRecvMsgSize int) (*RuntimeServiceClient, error) {
	addr, dailer, err := getAddressAndDialer(containerdUnixSocket)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getContextWithTimeout(contextTimeout)
	defer cancel()

	conn, err := grpcDialContext(ctx, addr, grpc.WithInsecure(), grpc.WithDialer(dailer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallRecvMsgSize)))
	if err != nil {
		return nil, err
	}

	client := &RuntimeServiceClient{
		conn: conn,
	}
	// Try v1
	client.service = newCRIRuntimeServiceV1Adapter(conn)
	if client.getVersion(ctx) == nil {
		logger.Info(ctx, "Init cri client v1 success, cri info", client.info)
		return client, nil
	}

	// if create client failed, close the connection
	_ = conn.Close()
	return nil, fmt.Errorf("failed to initialize RuntimeServiceClient")
}

func (c *RuntimeServiceClient) Version(ctx context.Context) (*CriVersionResponse, error) {
	if c.service != nil {
		return c.service.Version(ctx)
	}
	return &CriVersionResponse{}, fmt.Errorf("invalid RuntimeServiceClient")
}

func (c *RuntimeServiceClient) ListContainers(ctx context.Context) (*CriListContainersResponse, error) {
	if c.service != nil {
		return c.service.ListContainers(ctx)
	}
	return &CriListContainersResponse{}, fmt.Errorf("invalid RuntimeServiceClient")
}

func (c *RuntimeServiceClient) ContainerStatus(ctx context.Context, containerID string, verbose bool) (*CriContainerStatusResponse, error) {
	if c.service != nil {
		return c.service.ContainerStatus(ctx, containerID, verbose)
	}
	return &CriContainerStatusResponse{}, fmt.Errorf("invalid RuntimeServiceClient")
}

func (c *RuntimeServiceClient) ListPodSandbox(ctx context.Context) (*CriListPodSandboxResponse, error) {
	if c.service != nil {
		return c.service.ListPodSandbox(ctx)
	}
	return &CriListPodSandboxResponse{}, fmt.Errorf("invalid RuntimeServiceClient")
}

func (c *RuntimeServiceClient) PodSandboxStatus(ctx context.Context, sandboxID string, verbose bool) (*CriPodSandboxStatusResponse, error) {
	if c.service != nil {
		return c.service.PodSandboxStatus(ctx, sandboxID, verbose)
	}
	return &CriPodSandboxStatusResponse{}, fmt.Errorf("invalid RuntimeServiceClient")
}

func (c *RuntimeServiceClient) getVersion(ctx context.Context) error {
	versionResp, err := c.service.Version(ctx)
	if err == nil {
		c.info = CriVersionInfo{
			versionResp.Version,
			versionResp.RuntimeName,
			versionResp.RuntimeVersion,
			versionResp.RuntimeAPIVersion,
		}
	}
	return err
}

func (c *RuntimeServiceClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
