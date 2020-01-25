package sh

func init() {
	initSshClientCache()
}

type SshResponse struct {
	Address string
	Stdout  string
	Stderr  string
	Err     error
	Status  []string
}
