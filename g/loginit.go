package g

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var (
	Logger *logrus.Logger
)
func init()  {
	Logger = logrus.New()
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
