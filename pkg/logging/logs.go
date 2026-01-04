package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
)

type udpDatagramWriter struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func (w *udpDatagramWriter) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func InitLogs(service string, vectorAddr string) (*zap.Logger, error) {
	addr, err := net.ResolveUDPAddr("udp", vectorAddr)
	if err != nil {
		return nil, err
	}

	// Dial вместо Listen
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	writer := &udpDatagramWriter{
		conn: conn,
		addr: addr,
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderCfg.EncodeCaller = nil
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(writer),
		zap.InfoLevel,
	)

	logger := zap.New(
		core,
		zap.Fields(
			zap.String("service", service),
			zap.String("source", "tiles-backend"),
		),
	)

	return logger, nil
}
