package mserv

var log Logger

// Logger interface for package things
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// SetLogger for package usage
func SetLogger(l Logger) {
	log = l
}
