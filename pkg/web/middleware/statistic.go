package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/qqwry"
    "sync"
    "time"
)

/**
   @author yhy
   @since 2024/6/2
   @desc 统计各个 api 访问次数
**/

// VisitInfo 包含访问次数和IP地址的访问次数映射
type VisitInfo struct {
    Count     int // 总访问次数
    LastCount int
    IP        map[string]int // 每个IP的访问次数
}

// VisitCounter 访问计数器结构体
type VisitCounter struct {
    mu     sync.Mutex
    Counts map[string]*VisitInfo
}

// NewVisitCounter 创建一个新的访问计数器
func NewVisitCounter() *VisitCounter {
    return &VisitCounter{
        Counts: make(map[string]*VisitInfo),
    }
}

// Increment 增加给定路径的访问次数
func (vc *VisitCounter) Increment(path string, clientIP string) {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    if _, ok := vc.Counts[path]; !ok {
        vc.Counts[path] = &VisitInfo{
            Count: 0,
            IP:    make(map[string]int),
        }
    }
    vc.Counts[path].Count++
    vc.Counts[path].LastCount = vc.Counts[path].Count - 1
    vc.Counts[path].IP[clientIP]++
}

// GetCounts 获取每个路径的访问次数和每个IP的访问次数
func (vc *VisitCounter) GetCounts() map[string]*VisitInfo {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    pathCounts := db.GetPathCounts()
    ipCounts := db.GetIPCounts()
    
    counts := make(map[string]*VisitInfo)
    for _, pc := range pathCounts {
        counts[pc.Path] = &VisitInfo{
            Count: pc.Count,
            IP:    make(map[string]int),
        }
    }
    
    for _, ipc := range ipCounts {
        if _, exists := counts[ipc.Path]; !exists {
            counts[ipc.Path] = &VisitInfo{
                Count: 0,
                IP:    make(map[string]int),
            }
        }
        counts[ipc.Path].IP[ipc.IP] = ipc.Count
    }
    return counts
}

// VisitCountMiddleware 访问计数中间件
func VisitCountMiddleware(vc *VisitCounter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取客户端IP地址
        clientIP := c.ClientIP()
        
        // 在处理请求之前增加计数
        vc.Increment(c.Request.URL.Path, clientIP)
        // 继续处理请求
        c.Next()
    }
}

// UpdateDatabase 每分钟 将统计信息更新到数据库中
func (vc *VisitCounter) UpdateDatabase() {
    // 设置定时任务，每分钟更新数据库
    ticker := time.NewTicker(1 * time.Minute)
    
    for {
        select {
        case <-ticker.C:
            for path, info := range vc.Counts {
                if info.LastCount == info.LastCount { // 说明这个 api 路径这个时间段就没有人访问，不用再去更新数据库了
                    continue
                }
                vc.Counts[path].LastCount = vc.Counts[path].Count
                count := db.AddOrUpdatePathCounts(path, &db.PathCounts{
                    Path:  path,
                    Count: info.Count,
                })
                if count != 0 {
                    vc.mu.Lock()
                    vc.Counts[path].Count = count
                    vc.mu.Unlock()
                }
                
                for ip, ipCount := range info.IP {
                    address := ""
                    if ip != "" && qqwry.DB != nil {
                        _address, _ := qqwry.DB.Find(ip)
                        if _address != nil {
                            address = _address.String()
                        }
                    }
                    
                    _count := db.AddOrUpdateIPCounts(path, ip, &db.IPCounts{
                        IP:      ip,
                        Path:    path,
                        Count:   ipCount,
                        Address: address,
                    })
                    
                    vc.mu.Lock()
                    vc.Counts[path].IP[ip] = _count
                    vc.mu.Unlock()
                }
            }
        }
    }
}
