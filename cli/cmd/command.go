package cmd

type Command interface {
	Initialize() error
	Run() error
}
