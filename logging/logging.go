package logging

const loggerConfigString = "LOGGER"

type Logger interface {
	Info(msg string, kvs map[string]string)
}
