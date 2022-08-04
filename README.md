# Introduction
不限功能小工具，現有小鬧鐘、列舉員工等功能

# Dependencies
- [gin](https://pkg.go.dev/github.com/gin-gonic/gin)
- [cron](https://pkg.go.dev/github.com/robfig/cron/v3) (standard style)
- [gorm](https://pkg.go.dev/gorm.io/gorm)
- [shlex](https://pkg.go.dev/github.com/google/shlex)
- [go-arg](https://pkg.go.dev/github.com/alexflint/go-arg)

# How to use
1. 建立鬧鐘: alarm create "0 15 * * 1-5" "message!"
2. 檢視鬧鐘: alarm list
3. 刪除鬧鐘: alarm delete $ID
