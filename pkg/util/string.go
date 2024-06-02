package util

import (
    "github.com/russross/blackfriday/v2"
    "html/template"
    "sort"
    "strings"
)

/**
   @author yhy
   @since 2024/5/22
   @desc //TODO
**/

// Nl2br 将换行符转换为HTML的<br>标签
func Nl2br(text string) template.HTML {
    str := strings.ReplaceAll(template.HTMLEscapeString(text), "\r\n", "<br/>")
    str = strings.ReplaceAll(str, "\n", "<br/>")
    str = strings.ReplaceAll(str, "\r", "<br/>")
    return template.HTML(str)
}

// SplitString 将字符串按分隔符分割成数组
func SplitString(s, sep string) []string {
    return strings.Split(s, sep)
}

// Contains 检查字符串s是否包含子串substr
func Contains(s, substr string) bool {
    return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func Pointer(s *string, substr string) bool {
    return *s == substr
}

func ParseMarkdown(msg string) template.HTML {
    output := blackfriday.Run([]byte(msg))
    return template.HTML(output)
}

func UniqueStrings(input []string) []string {
    unique := make(map[string]bool)
    var result []string
    
    for _, value := range input {
        if _, exists := unique[value]; !exists {
            unique[value] = true
            result = append(result, value)
        }
    }
    
    return result
}

func Sort(data map[string]int) ([]string, []int) {
    // 创建一个切片来保存map的键值对
    type kv struct {
        Key   string
        Value int
    }
    var ss []kv
    for k, v := range data {
        ss = append(ss, kv{k, v})
    }
    
    // 按值排序
    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Value > ss[j].Value
    })
    
    var labels []string
    var series []int
    // 打印排序后的切片
    for _, k := range ss {
        labels = append(labels, k.Key)
        series = append(series, k.Value)
    }
    
    return labels, series
}
