package core

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Database struct {
	Config *Config `Inject:""`
	Gorm   *gorm.DB
}

const (
	dsnFormat = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func createDatabase() *Database {
	return &Database{}
}

func (db *Database) Connect() {
	Log.WithFields(logrus.Fields{"name": db.Config.Database.Name, "location": db.Config.Database.Location}).Info("Connecting to database...")
	dsn := fmt.Sprintf(dsnFormat, db.Config.Database.User, db.Config.Database.Password, db.Config.Database.Location, db.Config.Database.Name)
	gorm, err := gorm.Open(mysql.Open(dsn), db.getGormConfig())

	if err != nil {
		Log.Fatal(err)
	}

	db.Gorm = gorm
}

func (db *Database) getGormConfig() *gorm.Config {
	var logMode logger.LogLevel

	if db.Config.Debug == true {
		logMode = logger.Info
	} else {
		logMode = logger.Error
	}

	return &gorm.Config{Logger: NewLogger(Log).LogMode(logMode)}
}

type gormLogger struct {
	log   *logrus.Logger
	debug bool
}

func NewLogger(l *logrus.Logger) *gormLogger {
	return &gormLogger{
		log: l,
	}
}

// Implementation of the gorm logger.Interface methods
func (l *gormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	if logLevel == logger.Info {
		l.debug = true
	}
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	fields := logrus.Fields{}
	fields["loc"] = utils.FileWithLineNum()
	fields["rows"] = rows
	fields["ms"] = time.Since(begin)

	if err != nil {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf(sql)
		return
	}

	if l.debug {
		l.log.WithContext(ctx).WithFields(fields).Debugf(sql)
	}
}
