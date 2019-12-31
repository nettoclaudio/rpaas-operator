// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func Test_Init(t *testing.T) {
	tests := []struct {
		config   string
		envs     map[string]string
		expected RpaasConfig
	}{
		{
			expected: RpaasConfig{
				ServiceName:               "rpaasv2",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
port-range-max: 31000
sync-interval: 2m
`,
			expected: RpaasConfig{
				ServiceName:               "rpaasv2",
				SyncInterval:              2 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              31000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
tls-certificate: /var/share/tls/mycert.pem
tls-key: /var/share/tls/key.pem
`,
			expected: RpaasConfig{
				ServiceName:               "rpaasv2",
				TLSCertificate:            "/var/share/tls/mycert.pem",
				TLSKey:                    "/var/share/tls/key.pem",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
api-username: u1
`,
			expected: RpaasConfig{
				APIUsername:               "u1",
				ServiceName:               "rpaasv2",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
api-username: ignored1
service-name: rpaasv2be
`,
			envs: map[string]string{
				"RPAASV2_API_USERNAME": "u1",
				"RPAASV2_API_PASSWORD": "p1",
			},
			expected: RpaasConfig{
				APIUsername:               "u1",
				APIPassword:               "p1",
				ServiceName:               "rpaasv2be",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
service-name: ignored-service-name
`,
			envs: map[string]string{
				"RPAASV2_SERVICE_NAME": "my-custom-service-name",
			},
			expected: RpaasConfig{
				ServiceName:               "my-custom-service-name",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
default-affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: pool
          operator: In
          values:
          - dev
team-affinity:
  team1:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: pool
            operator: NotIn
            values:
            - dev
`,
			expected: RpaasConfig{
				ServiceName:  "rpaasv2",
				SyncInterval: 5 * time.Minute,
				PortRangeMin: 20000,
				PortRangeMax: 30000,
				DefaultAffinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								{
									MatchExpressions: []corev1.NodeSelectorRequirement{
										{
											Key:      "pool",
											Operator: corev1.NodeSelectorOpIn,
											Values:   []string{"dev"},
										},
									},
								},
							},
						},
					},
				},
				TeamAffinity: map[string]corev1.Affinity{
					"team1": {
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "pool",
												Operator: corev1.NodeSelectorOpNotIn,
												Values:   []string{"dev"},
											},
										},
									},
								},
							},
						},
					},
				},
				WebsocketHandshakeTimeout: 5 * time.Second,
				WebsocketReadBufferSize:   1024,
				WebsocketWriteBufferSize:  4096,
				WebsocketPingInterval:     5 * time.Second,
				WebsocketMaxIdleTime:      10 * time.Second,
				WebsocketWriteWait:        time.Second,
			},
		},
		{
			config: `
websocket-handshake-timeout: 10s
websocket-read-buffer-size: 512
websocket-write-buffer-size: 512
websocket-ping-interval: 2s
websocket-max-idle-time: 30s
websocket-write-wait: 2s
`,
			expected: RpaasConfig{
				ServiceName:               "rpaasv2",
				SyncInterval:              5 * time.Minute,
				PortRangeMin:              20000,
				PortRangeMax:              30000,
				WebsocketHandshakeTimeout: 10 * time.Second,
				WebsocketReadBufferSize:   512,
				WebsocketWriteBufferSize:  512,
				WebsocketPingInterval:     2 * time.Second,
				WebsocketMaxIdleTime:      30 * time.Second,
				WebsocketWriteWait:        2 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			defer viper.Reset()
			for k, v := range tt.envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			dir, err := ioutil.TempDir("", "")
			require.NoError(t, err)
			name := filepath.Join(dir, "config.yaml")
			err = ioutil.WriteFile(name, []byte(tt.config), 0644)
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			os.Args = []string{"test", "--config", name}
			err = Init()
			require.NoError(t, err)
			config := Get()
			assert.Equal(t, tt.expected, config)
		})
	}
}
