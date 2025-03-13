package main

import "github.com/volodymyrzuyev/goCsInspect/logger"

func main() {
	log := logger.NewLogger(nil)

	log.Error("Err test")
	log.Info("Info test")
	log.Debug("Debug test")
}
