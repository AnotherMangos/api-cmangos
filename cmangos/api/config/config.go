package config

type ApiConfig struct {
  RequireInvite bool `json:"requireInvite,omitempty"`
  RealmdAddress string `json:"address,omitempty"`
  RealmdPort int `json:"port,omitempty"`
}
