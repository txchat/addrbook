package version

var (
	// The full version string
	Version = "2.0.2.0"

	GitCommit string
)

func GetVersion() string {
	if GitCommit != "" {
		return Version + "-" + GitCommit
	}
	return Version
}
