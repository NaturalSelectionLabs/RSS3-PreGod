package cmd

type Interface interface {
	Initialize() error
	Run() error
}
