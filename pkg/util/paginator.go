package util

import (
    "fmt"
    "math"
    "net/http"
    "net/url"
    "reflect"
    "strconv"
)

/**
   @author yhy
   @since 2024/5/21
   @desc 分页器
**/

type Paginator struct {
    Request     *http.Request
    PerPageNums int
    MaxPages    int
    
    nums      int64
    pageRange []int
    pageNums  int
    page      int
}

func (p *Paginator) PageNums() int {
    if p.pageNums != 0 {
        return p.pageNums
    }
    pageNums := math.Ceil(float64(p.nums) / float64(p.PerPageNums))
    if p.MaxPages > 0 {
        pageNums = math.Min(pageNums, float64(p.MaxPages))
    }
    p.pageNums = int(pageNums)
    return p.pageNums
}

func (p *Paginator) Nums() int64 {
    return p.nums
}

func (p *Paginator) SetNums(nums interface{}) {
    p.nums, _ = ToInt64(nums)
}

func (p *Paginator) Page() int {
    if p.page != 0 {
        return p.page
    }
    if p.Request.Form == nil {
        p.Request.ParseForm()
    }
    p.page, _ = strconv.Atoi(p.Request.Form.Get("current"))
    if p.page > p.PageNums() {
        p.page = p.PageNums()
    }
    if p.page <= 0 {
        p.page = 1
    }
    return p.page
}

func (p *Paginator) Pages() []int {
    if p.pageRange == nil && p.nums > 0 {
        var pages []int
        pageNums := p.PageNums()
        page := p.Page()
        switch {
        case page >= pageNums-4 && pageNums > 9:
            start := pageNums - 9 + 1
            pages = make([]int, 9)
            for i, _ := range pages {
                pages[i] = start + i
            }
        case page >= 5 && pageNums > 9:
            start := page - 5 + 1
            pages = make([]int, int(math.Min(9, float64(page+4+1))))
            for i, _ := range pages {
                pages[i] = start + i
            }
        default:
            pages = make([]int, int(math.Min(9, float64(pageNums))))
            for i, _ := range pages {
                pages[i] = i + 1
            }
        }
        p.pageRange = pages
    }
    return p.pageRange
}

func (p *Paginator) PageLink(page int) string {
    link, _ := url.ParseRequestURI(p.Request.RequestURI)
    values := link.Query()
    if page == 1 {
        values.Del("current")
    } else {
        values.Set("current", strconv.Itoa(page))
    }
    link.RawQuery = values.Encode()
    return link.String()
}

func (p *Paginator) PageLinkPrev() (link string) {
    if p.HasPrev() {
        link = p.PageLink(p.Page() - 1)
    }
    return
}

func (p *Paginator) PageLinkNext() (link string) {
    if p.HasNext() {
        link = p.PageLink(p.Page() + 1)
    }
    return
}

func (p *Paginator) PageLinkFirst() (link string) {
    return p.PageLink(1)
}

func (p *Paginator) PageLinkLast() (link string) {
    return p.PageLink(p.PageNums())
}

func (p *Paginator) HasPrev() bool {
    return p.Page() > 1
}

func (p *Paginator) HasNext() bool {
    return p.Page() < p.PageNums()
}

func (p *Paginator) IsActive(page int) bool {
    return p.Page() == page
}

func (p *Paginator) Offset() int {
    return (p.Page() - 1) * p.PerPageNums
}

func (p *Paginator) HasPages() bool {
    return p.PageNums() > 1
}

func NewPaginator(req *http.Request, per int, nums int64) *Paginator {
    p := Paginator{}
    p.Request = req
    p.PerPageNums = per
    p.nums = nums
    return &p
}

// ToInt64 convert any numeric value to int64
func ToInt64(value interface{}) (d int64, err error) {
    val := reflect.ValueOf(value)
    switch value.(type) {
    case int, int8, int16, int32, int64:
        d = val.Int()
    case uint, uint8, uint16, uint32, uint64:
        d = int64(val.Uint())
    default:
        err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
    }
    return
}
