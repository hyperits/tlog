package log_test

import (
	log "github.com/hyperits/tlog"
)

func Example() {
	l := log.WithFields("uid", "10012")

	l.Trace("helloworld")
	l.Debug("helloworld")
	l.Info("helloworld")
	l.Warn("helloworld")
	l.Error("helloworld")
	l.Tracef("helloworld")
	l.Debugf("helloworld")
	l.Infof("helloworld")
	l.Warnf("helloworld")
	l.Errorf("helloworld")
	// Output:
}
