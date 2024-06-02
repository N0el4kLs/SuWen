package util

import (
    "fmt"
    "sort"
    "time"
)

/**
   @author yhy
   @since 2024/5/30
   @desc //TODO
**/

func TimeNow() string {
    loc, _ := time.LoadLocation("Asia/Shanghai")
    return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func TimeSub(t *time.Time) string {
    loc, _ := time.LoadLocation("Asia/Shanghai")
    // 计算时间差
    duration := time.Now().In(loc).Sub(t.In(loc))
    
    hours := duration.Hours()
    minutes := duration.Minutes()
    
    // 获取完整小时数
    wholeHours := int(hours)
    // 获取剩余的分钟数
    remainingMinutes := int(minutes) % 60
    day := wholeHours / 24
    if day > 0 {
        wholeHours -= day * 24
        return fmt.Sprintf("%d day %d hours %d minutes", day, wholeHours, remainingMinutes)
    } else {
        return fmt.Sprintf("%d hours %d minutes", wholeHours, remainingMinutes)
    }
}

func CurrentlyYear(dateStr string) bool {
    givenDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return true
    }
    
    if givenDate.Year() == time.Now().Year() {
        return true
    }
    return false
}

func MonthlyCalculation(dateStr string) bool {
    givenDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    months := int(time.Now().Sub(givenDate).Hours() / 24 / 30)
    
    if months <= 3 {
        return true
    }
    
    return false
}

func CurrentlyMonth(dateStr string) bool {
    givenDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return true
    }
    
    currentTime := time.Now()
    isSameYear := givenDate.Year() == currentTime.Year()
    isSameMonth := givenDate.Month() == currentTime.Month()
    isLastMonth := givenDate.Month() == currentTime.Month()-1
    
    // 判断是否是当前两个月内的
    if isSameYear && (isSameMonth || isLastMonth) {
        return true
    }
    return false
}

func IsTodayBeijing(t time.Time) (bool, time.Time) {
    loc, _ := time.LoadLocation("Asia/Shanghai")
    t = t.In(loc)
    now := time.Now().In(loc)
    
    hours := int(now.Sub(t).Hours()) // 计算两个时间的小时差，只要小时差小于 24 , 这两天更新的
    
    return t.Year() == now.Year() && t.Month() == now.Month() && (t.Day() == now.Day() || hours <= 24), t
}

func SortTimeSlice(dateStrings []string) []string {
    // 解析日期字符串并将它们存储在一个新的 time.Time 切片中
    var dates []time.Time
    for _, dateString := range dateStrings {
        date, err := time.Parse("2006-01-02", dateString)
        if err != nil {
            fmt.Printf("Error parsing date: %v\n", err)
            continue
        }
        dates = append(dates, date)
    }
    
    // 对 time.Time 切片进行排序
    sort.Slice(dates, func(i, j int) bool {
        return dates[i].Before(dates[j])
    })
    
    // 将排序后的 time.Time 切片转换回字符串格式
    var sortedDateStrings []string
    for _, date := range dates {
        sortedDateStrings = append(sortedDateStrings, date.Format("2006-01-02"))
    }
    
    return sortedDateStrings
}
