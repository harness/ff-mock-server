package internal

const (
	// ServerKey is a randomly generated UUID, it can be used only in
	// server SDKs
	ServerKey = "2e182b14-9944-4bd4-9c9f-3e859e2a2954"
	// ClientKey is a randomly generated UUID, it can be used only in
	// client SDKs
	ClientKey = "2e2ecf62-ce53-4e9e-8006-b4db0386688c"
	// DefaultAuthSecret is used only if there is no value in env variable
	DefaultAuthSecret = "mock-server"
	// DefaultClusterIdentifier is used only if there is no value in env variable
	DefaultClusterIdentifier = "cluster"
	// Project mocked value
	Project = "demo"
	// Environment mocked value
	Environment = "dev"
	// EnvironmentUUID mocked value
	EnvironmentUUID = "265597ad-516c-4575-a16f-b3d17adffc44"
	// JWTKey mocked value
	JWTKey = "jwt"
)
