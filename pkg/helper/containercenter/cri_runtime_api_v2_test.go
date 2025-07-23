// Copyright 2025 iLogtail Authors
// Licensed under Apache License, Version 2.0 (the "License")

package containercenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCriV2NewFields(t *testing.T) {
	// 测试 v2 新增字段的定义
	resources := &CriLinuxContainerResources{
		CPUPeriod:              100000,
		CPUQuota:               50000,
		CPUShares:              1024,
		MemoryLimitInBytes:     1073741824, // 1GB
		OomScoreAdj:            -1000,
		CpusetCpus:             "0-1",
		CpusetMems:             "0",
		HugepageLimits:         map[string]uint64{"2Mi": 10},
		Unified:                map[string]string{"cgroup.type": "threaded"},
		MemorySwapLimitInBytes: 2147483648, // 2GB
	}

	assert.Equal(t, int64(100000), resources.CPUPeriod)
	assert.Equal(t, int64(50000), resources.CPUQuota)
	assert.Equal(t, int64(1073741824), resources.MemoryLimitInBytes)
	assert.Equal(t, "0-1", resources.CpusetCpus)

	// 测试安全上下文
	securityContext := &CriLinuxContainerSecurityContext{
		Privileged:     true,
		ReadonlyRootfs: false,
		RunAsUser:      &CriInt64Value{Value: 1000},
		RunAsGroup:     &CriInt64Value{Value: 1000},
		Capabilities: &CriCapability{
			AddCapabilities:  []string{"NET_ADMIN", "SYS_ADMIN"},
			DropCapabilities: []string{"CHOWN", "SETUID"},
		},
		NamespaceOptions: &CriNamespaceOption{
			Network: NamespaceModePod,
			Pid:     NamespaceModeContainer,
			Ipc:     NamespaceModePod,
		},
		SelinuxOptions: &CriSELinuxOption{
			User:  "system_u",
			Role:  "system_r",
			Type:  "container_t",
			Level: "s0:c123,c456",
		},
	}

	assert.True(t, securityContext.Privileged)
	assert.Equal(t, int64(1000), securityContext.RunAsUser.Value)
	assert.Equal(t, 2, len(securityContext.Capabilities.AddCapabilities))
	assert.Equal(t, NamespaceModePod, securityContext.NamespaceOptions.Network)
}

func TestCriContainerStateCompatibility(t *testing.T) {
	// 测试容器状态兼容性
	states := []CriContainerState{
		ContainerStateContainerCreated,
		ContainerStateContainerRunning,
		ContainerStateContainerExited,
		ContainerStateContainerUnknown,
		// v2 新增状态
		ContainerStateContainerCreated2,
		ContainerStateContainerRunning2,
		ContainerStateContainerExited2,
		ContainerStateContainerUnknown2,
	}

	// 验证所有状态都有不同的值
	stateMap := make(map[CriContainerState]bool)
	for _, state := range states {
		assert.False(t, stateMap[state], "Duplicate state value: %d", state)
		stateMap[state] = true
	}
}

func TestCriNamespaceMode(t *testing.T) {
	// 测试命名空间模式
	assert.Equal(t, CriNamespaceMode(0), NamespaceModePod)
	assert.Equal(t, CriNamespaceMode(1), NamespaceModeContainer)
	assert.Equal(t, CriNamespaceMode(2), NamespaceModeNode)
	assert.Equal(t, CriNamespaceMode(3), NamespaceModeTarget)
}

func TestCriContainerWithV2Fields(t *testing.T) {
	// 测试包含 v2 字段的容器结构
	container := &CriContainer{
		ID:           "test-container-id",
		PodSandboxID: "test-sandbox-id",
		Metadata: &CriContainerMetadata{
			Name:    "test-container",
			Attempt: 1,
		},
		Image: &CriImageSpec{
			Image:       "nginx:latest",
			Annotations: map[string]string{"test": "value"},
		},
		ImageRef:    "docker.io/library/nginx:latest",
		State:       ContainerStateContainerRunning,
		CreatedAt:   1640995200000000000, // 2022-01-01 00:00:00
		Labels:      map[string]string{"app": "nginx"},
		Annotations: map[string]string{"io.kubernetes.container.name": "nginx"},
		Resources: &CriLinuxContainerResources{
			MemoryLimitInBytes: 1073741824,
			CPUQuota:           50000,
		},
		SecurityContext: &CriLinuxContainerSecurityContext{
			Privileged: true,
			RunAsUser:  &CriInt64Value{Value: 1000},
		},
	}

	assert.Equal(t, "test-container-id", container.ID)
	assert.Equal(t, ContainerStateContainerRunning, container.State)
	assert.NotNil(t, container.Resources)
	assert.Equal(t, int64(1073741824), container.Resources.MemoryLimitInBytes)
	assert.NotNil(t, container.SecurityContext)
	assert.True(t, container.SecurityContext.Privileged)
	assert.Equal(t, int64(1000), container.SecurityContext.RunAsUser.Value)
}

func TestCriContainerStatusWithV2Fields(t *testing.T) {
	// 测试包含 v2 字段的容器状态结构
	status := &CriContainerStatus{
		ID: "test-container-id",
		Metadata: &CriContainerMetadata{
			Name:    "test-container",
			Attempt: 1,
		},
		State:      ContainerStateContainerRunning,
		CreatedAt:  1640995200000000000,
		StartedAt:  1640995201000000000,
		FinishedAt: 0,
		ExitCode:   0,
		Image: &CriImageSpec{
			Image: "nginx:latest",
		},
		ImageRef:    "docker.io/library/nginx:latest",
		Reason:      "",
		Message:     "",
		Labels:      map[string]string{"app": "nginx"},
		Annotations: map[string]string{"io.kubernetes.container.name": "nginx"},
		Mounts: []*CriMount{
			{
				ContainerPath:  "/var/log",
				HostPath:       "/host/var/log",
				Readonly:       false,
				SelinuxRelabel: false,
				Propagation:    MountPropagationPropagationPrivate,
			},
		},
		LogPath: "/var/log/containers/test-container.log",
		Resources: &CriLinuxContainerResources{
			MemoryLimitInBytes: 1073741824,
			CPUQuota:           50000,
		},
		SecurityContext: &CriLinuxContainerSecurityContext{
			Privileged: false,
			RunAsUser:  &CriInt64Value{Value: 1000},
		},
	}

	assert.Equal(t, "test-container-id", status.ID)
	assert.Equal(t, ContainerStateContainerRunning, status.State)
	assert.Equal(t, 1, len(status.Mounts))
	assert.Equal(t, "/var/log", status.Mounts[0].ContainerPath)
	assert.NotNil(t, status.Resources)
	assert.NotNil(t, status.SecurityContext)
	assert.False(t, status.SecurityContext.Privileged)
}
