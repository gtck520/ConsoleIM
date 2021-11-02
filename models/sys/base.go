package sys

import(
	"fmt"

	"github.com/gtck520/ConsoleIM/models/basemodel"
)

func TableName(name string) string {
	return fmt.Sprintf("%s%s%s", basemodel.GetTablePrefix(),"sys_", name)
}