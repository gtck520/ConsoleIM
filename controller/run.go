package controller

import (
	"github.com/gtck520/ConsoleIM/cview"
	"github.com/spf13/cobra"
)

type RunController struct {
}

func (r *RunController) Run(cmd *cobra.Command, args []string, globals ...string) {
	//初始化显示参数
	cvi := cview.NewCView()
	cvi.Index()
}
