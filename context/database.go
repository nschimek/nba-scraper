package context

import (
	goctx "context"
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

type dbLogger struct {
	log   *logrus.Logger
	debug bool
}

func NewLogger(l *logrus.Logger) *dbLogger {
	return &dbLogger{
		log: l,
	}
}

func (l *dbLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	if logLevel == logger.Info {
		l.debug = true
	}
	return l
}

func (l *dbLogger) Info(ctx goctx.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args)
}

func (l *dbLogger) Warn(ctx goctx.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args)
}

func (l *dbLogger) Error(ctx goctx.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args)
}

func (l *dbLogger) Trace(ctx goctx.Context, begin time.Time, fc func() (string, int64), err error) {
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
