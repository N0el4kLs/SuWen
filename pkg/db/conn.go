package db

import (
    "fmt"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/logging"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    dblogger "gorm.io/gorm/logger"
    "time"
)

var GlobalDB *gorm.DB

func Init() {
    var err error
    // 配置mysql数据源
    if conf.GlobalConfig == nil || conf.GlobalConfig.DbConfig.User == "" ||
        conf.GlobalConfig.DbConfig.Password == "" ||
        conf.GlobalConfig.DbConfig.Host == "" ||
        conf.GlobalConfig.DbConfig.Port == "" ||
        conf.GlobalConfig.DbConfig.Database == "" {
        logging.Logger.Fatalf("db.Setup err: '%s' mysql config not set", conf.ConfigFileName)
    }
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
        conf.GlobalConfig.DbConfig.User,
        conf.GlobalConfig.DbConfig.Password,
        conf.GlobalConfig.DbConfig.Host,
        conf.GlobalConfig.DbConfig.Port,
        conf.GlobalConfig.DbConfig.Database)
    
    // 创建 Gorm 的 logger 实现
    GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: &DBLogger{
            Logger: logging.Logger,
            // 可以设置日志级别
            GormLogger: dblogger.New(
                logging.Logger, // io.Writer
                dblogger.Config{
                    SlowThreshold: 3 * time.Second,
                    Colorful:      true,
                    LogLevel:      dblogger.Error,
                },
            ),
        },
    })
    
    if err != nil {
        logging.Logger.Fatalf("db.Setup err: %v", err)
    }
    
    sqlDB, _ := GlobalDB.DB()
    
    // SetMaxIdleConns 设置空闲连接池中连接的最大数量, 该函数的作用就是保持等待连接操作状态的连接数，这个主要就是避免操作过程中频繁的获取连接，释放连接。
    sqlDB.SetMaxIdleConns(100)
    
    // SetMaxOpenConns 设置打开数据库连接的最大数量。 这个我们不设置默认就是不限制，可以无限创建连接，问题就在数据库本身有瓶颈，无限创建，会损耗性能。所以我们要根据我们自己的数据库瓶颈情况来进行相关的设置。当出现连接数超出了我们设定的数量时候，后面的用户等待超时时间之前，有连接释放就会自动获得操作的权限，否则返回连接超时。
    sqlDB.SetMaxOpenConns(100)
    
    // 5秒内连接没有活跃的话则自动关闭连接
    sqlDB.SetConnMaxLifetime(time.Second * 10)
    
    if GlobalDB == nil {
        logging.Logger.Fatalf("db.Setup err: db connect failed")
    }
    
    err = GlobalDB.AutoMigrate(&PressRelease{}, &Advisory{}, &Poc{}, &Pusher{}, &IPCounts{}, &PathCounts{})
    
    if err != nil {
        logging.Logger.Fatalf("db.Setup err: %v", err)
    }
}
