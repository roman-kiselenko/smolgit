package config

import (
	"testing"

	"smolgit/pkg/model"
)

func TestFindUserByKey(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		key       string
		expected  model.User
		expectErr bool
	}{
		{
			name: "User found by exact key match",
			config: &Config{
				Users: []model.User{
					{
						Name: "alice",
						Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... alice@mail.com"},
					},
				},
			},
			key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... alice@mail.com",
			expected: model.User{
				Name: "alice",
				Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... alice@mail.com"},
			},
			expectErr: false,
		},
		{
			name: "User found by partial key match",
			config: &Config{
				Users: []model.User{
					{
						Name: "bob",
						Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... bob@mail.com"},
					},
				},
			},
			key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ",
			expected: model.User{
				Name: "bob",
				Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... bob@mail.com"},
			},
			expectErr: false,
		},
		{
			name: "User not found",
			config: &Config{
				Users: []model.User{
					{
						Name: "charlie",
						Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... charlie@mail.com"},
					},
				},
			},
			key:       "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... unknown@mail.com",
			expected:  model.User{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := tt.config.FindUserByKey(tt.key)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if user.Name != tt.expected.Name {
					t.Errorf("Expected user %s, got %s", tt.expected.Name, user.Name)
				}

				if len(user.Keys) != len(tt.expected.Keys) {
					t.Errorf("Expected %d keys, got %d", len(tt.expected.Keys), len(user.Keys))
				}
			}
		})
	}
}

func TestFindUserByName(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		userName  string
		expected  model.User
		expectErr bool
	}{
		{
			name: "User found by name",
			config: &Config{
				Users: []model.User{
					{
						Name: "alice",
						Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... alice@mail.com"},
					},
				},
			},
			userName: "alice",
			expected: model.User{
				Name: "alice",
				Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... alice@mail.com"},
			},
			expectErr: false,
		},
		{
			name: "User not found",
			config: &Config{
				Users: []model.User{
					{
						Name: "bob",
						Keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ... bob@mail.com"},
					},
				},
			},
			userName:  "charlie",
			expected:  model.User{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := tt.config.FindUserByName(tt.userName)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if user.Name != tt.expected.Name {
					t.Errorf("Expected user %s, got %s", tt.expected.Name, user.Name)
				}

				if len(user.Keys) != len(tt.expected.Keys) {
					t.Errorf("Expected %d keys, got %d", len(tt.expected.Keys), len(user.Keys))
				}
			}
		})
	}
}
