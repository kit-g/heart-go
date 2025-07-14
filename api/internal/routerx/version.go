package routerx

var (
	Commit     = "dev"
	DeployedAt = "unknown"
)

func Info() map[string]string {
	return map[string]string{
		"commit":     Commit,
		"deployedAt": DeployedAt,
	}
}

func String() string {
	return Commit + " @ " + DeployedAt
}
