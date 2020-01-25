package sh

func init() {
	initSshClientCache()
}

type sshResponse struct {
	Address string
	Stdout  string
	Stderr  string
	Err     error
	Status  []string
}
