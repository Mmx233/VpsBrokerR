package Controllers

import (
	"Mmx/Modules"
	"errors"
	"github.com/gin-gonic/gin"
)

func Main(c *gin.Context) {
	var form struct {
		Time    uint   `form:"time" binding:"required"`
		Name    string `form:"name" binding:"required"`
		Sign    string `form:"sign" binding:"required"`
		Backend string `form:"backend"`
	}
	if err := c.ShouldBind(&form); err != nil {
		Modules.Callback.Request.Error(c, err)
	}
	switch form.Sign {
	case "ip":
		form.Sign = c.ClientIP()
	case "name":
		form.Sign = form.Name
	default:
		Modules.Callback.Request.Error(c, errors.New("sign参数有误"))
		return
	}
	if Modules.Global.Online[form.Sign] == nil { //主机不在线
		Modules.Worker.Do(form.Sign, Modules.VpsInfo{ //开启协程
			Time:    form.Time,
			Name:    form.Name,
			Ip:      c.ClientIP(),
			Backend: form.Backend,
		})
		 //新主机
		 Modules.Callback.Backend.New(form.Sign)
		 Modules.Callback.Request.Info(c, "init")
	} else {
		Modules.Global.Online[form.Sign].Name = form.Name
		Modules.Global.Online[form.Sign].Time = form.Time
		Modules.Global.Online[form.Sign].Count = form.Time
		Modules.Global.Online[form.Sign].Backend = form.Backend
		if Modules.Global.Offline[form.Sign] != 0 { //恢复
			Modules.Callback.Backend.Up(form.Sign)
			Modules.Callback.Request.Info(c, "up")
			if Modules.Global.Online[form.Sign].Backend==""{
				Modules.Worker.Do(form.Sign, Modules.VpsInfo{ //未配置backend，重新开启协程
					Time:    form.Time,
					Name:    form.Name,
					Ip:      c.ClientIP(),
					Backend: form.Backend,
				})
			}
		}else { //续期
			Modules.Callback.Request.Info(c, "continue")
		}
	}
}
