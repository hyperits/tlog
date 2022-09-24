package log_test

import (
	"errors"
	log "github.com/hyperits/tlog"
	"github.com/hyperits/tlog/plugin"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var defaultConfig = []log.OutputConfig{
	{
		Writer:    "console",
		Level:     "debug",
		Formatter: "console",
		FormatConfig: log.FormatConfig{
			TimeFmt: "2006.01.02 15:04:05",
		},
	},
	{
		Writer:    "file",
		Level:     "info",
		Formatter: "json",
		WriteConfig: log.WriteConfig{
			Filename:   "tlog_size.log",
			RollType:   "size",
			MaxAge:     7,
			MaxBackups: 10,
			MaxSize:    100,
		},
		FormatConfig: log.FormatConfig{
			TimeFmt: "2006.01.02 15:04:05",
		},
	},
	{
		Writer:    "file",
		Level:     "info",
		Formatter: "json",
		WriteConfig: log.WriteConfig{
			Filename:   "tlog_time.log",
			RollType:   "time",
			MaxAge:     7,
			MaxBackups: 10,
			MaxSize:    100,
			TimeUnit:   log.Day,
		},
		FormatConfig: log.FormatConfig{
			TimeFmt: "2006-01-02 15:04:05",
		},
	},
}

func TestNewZapLog(t *testing.T) {
	logger := log.NewZapLog(defaultConfig)
	assert.NotNil(t, logger)

	logger.SetLevel("0", log.LevelInfo)
	lvl := logger.GetLevel("0")
	assert.Equal(t, lvl, log.LevelInfo)

	l := logger.WithFields("test", "a")
	if tmp, ok := l.(*log.ZapLogWrapper); ok {
		tmp.GetLogger()
		tmp.Sync()
	}
	l.SetLevel("output", log.LevelDebug)
	assert.Equal(t, log.LevelDebug, l.GetLevel("output"))
}

func TestZapLogWithLevel(t *testing.T) {
	logger := log.NewZapLog(defaultConfig)
	assert.NotNil(t, logger)

	l := logger.WithFields("field1")
	l.SetLevel("0", log.LevelFatal)
	assert.Equal(t, log.LevelFatal, l.GetLevel("0"))

	l = l.With(log.Field{Key: "key1", Value: "val1"})
	l.SetLevel("0", log.LevelError)
	assert.Equal(t, log.LevelError, l.GetLevel("0"))
}

func BenchmarkDefaultTimeFormat(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		log.DefaultTimeFormat(t)
	}
}

func BenchmarkCustomTimeFormat(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		log.CustomTimeFormat(t, "2006-01-02 15:04:05.000")
	}
}

func TestCustomTimeFormat(t *testing.T) {
	date := time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)
	dateStr := log.CustomTimeFormat(date, "2006-01-02 15:04:05.000")
	assert.Equal(t, dateStr, "2006-01-02 15:04:05.000")
}

func TestDefaultTimeFormat(t *testing.T) {
	date := time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)
	dateStr := string(log.DefaultTimeFormat(date))
	assert.Equal(t, dateStr, "2006-01-02 15:04:05.000")
}

func TestGetLogEncoderKey(t *testing.T) {
	tests := []struct {
		name   string
		defKey string
		key    string
		want   string
	}{
		{"custom", "T", "Time", "Time"},
		{"default", "T", "", "T"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := log.GetLogEncoderKey(tt.defKey, tt.key); got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestNewTimeEncoder(t *testing.T) {
	encoder := log.NewTimeEncoder("")
	assert.NotNil(t, encoder)

	encoder = log.NewTimeEncoder("2006-01-02 15:04:05")
	assert.NotNil(t, encoder)

	tests := []struct {
		name string
		fmt  string
	}{
		{"seconds timestamp", "seconds"},
		{"milliseconds timestamp", "milliseconds"},
		{"nanoseconds timestamp", "nanoseconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := log.NewTimeEncoder(tt.fmt)
			assert.NotNil(t, got)
		})
	}
}

type withFieldsWriter struct {
	core zapcore.Core
}

func (f *withFieldsWriter) Type() string { return "log" }

func (f *withFieldsWriter) Setup(name string, dec plugin.Decoder) error {
	if dec == nil {
		return errors.New("empty decoder")
	}
	decoder, ok := dec.(*log.Decoder)
	if !ok {
		return errors.New("invalid decoder")
	}
	decoder.Core = f.core
	decoder.ZapLevel = zap.NewAtomicLevel()
	return nil
}

func TestWithFields(t *testing.T) {
	// register Writer.
	// use zap observer to support test.
	core, ob := observer.New(zap.InfoLevel)
	log.RegisterWriter("withfields", &withFieldsWriter{core: core})

	// config is configuration.
	cfg := []log.OutputConfig{
		{
			Writer: "withfields",
		},
	}

	// create a zap logger.
	zl := log.NewZapLog(cfg)
	assert.NotNil(t, zl)

	// test With.
	field := log.Field{Key: "abc", Value: int32(123)}
	logger := zl.With(field)
	assert.NotNil(t, logger)
	log.SetLogger(logger)
	log.Warn("with fields warning")
	assert.Equal(t, 1, ob.Len())
	entry := ob.All()[0]
	assert.Equal(t, zap.WarnLevel, entry.Level)
	assert.Equal(t, "with fields warning", entry.Message)
	assert.Equal(t, []zapcore.Field{{Key: "abc", Type: zapcore.Int32Type, Integer: 123}}, entry.Context)
}
