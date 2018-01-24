package main

type Config struct {
	Address string `default:":7070"`
	Token   string `default:""`
	Mode    string `default:"docker"`
}

type Deployer interface {
	Update(target, image string) error
}

type Handler struct {
	deployer Deployer
}
