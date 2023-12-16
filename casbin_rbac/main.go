package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gin-gonic/gin"
)

func main() {
	m, err := model.NewModelFromFile("casbin_model.conf")
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	adapter := fileadapter.NewAdapter("policy.csv")

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	e.AddRoleForUser("admin", "data_admin")
	e.AddRoleForUser("user", "data_user")

	e.AddPermissionForUser("data_admin", "/api/admin", "GET")

	engine := gin.Default()

	engine.GET("api/admin", func(c *gin.Context) {
		req := make(map[string]string)
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		obj := c.Request.RequestURI
		fmt.Println("obj:", obj)
		sub := req["name"]
		act := c.Request.Method
		b, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		if b {
			c.JSON(200, gin.H{
				"name": sub,
				"pass": req["pass"],
			})
			return
		} else {
			c.JSON(500, gin.H{
				"msg": "dont have perrmisson",
			})
			return
		}

	})

	engine.Run(":8080")

}
