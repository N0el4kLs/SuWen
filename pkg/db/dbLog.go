package db

import (
    "context"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm/logger"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/23
  @desc: //TODO
**/

// DBLogger 实现了 Gorm 的 Logger 接口
type DBLogger struct {
    Logger     *logrus.Logger
    GormLogger logger.Interface
}

func (l *DBLogger) Info(ctx context.Context, s string, i ...interface{}) {
    l.Logger.Infof(s, "")
}

func (l *DBLogger) Warn(ctx context.Context, s string, i ...interface{}) {
    l.Logger.Warnf(s, "")
}

func (l *DBLogger) Error(ctx context.Context, s string, i ...interface{}) {
    l.Logger.Errorf(s, "")
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    // 调用 fc 函数获取 SQL 语句和影响的行数
    sql, rowsAffected := fc()
    
    // 计算执行时间
    duration := time.Since(begin)
    
    l.Logger.Tracef("rows_affected:%v, sql: %s , use: %v", rowsAffected, sql, duration)
}

func (l *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
    // 根据 Gorm 日志级别设置 logrus 日志级别
    switch level {
    case logger.Silent:
        l.Logger.SetLevel(logrus.PanicLevel)
    case logger.Error:
        l.Logger.SetLevel(logrus.ErrorLevel)
    case logger.Warn:
        l.Logger.SetLevel(logrus.WarnLevel)
    case logger.Info:
        l.Logger.SetLevel(logrus.InfoLevel)
    default:
        l.Logger.SetLevel(logrus.DebugLevel)
    }
    return l
}
