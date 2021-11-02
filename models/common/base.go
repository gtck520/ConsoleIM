package common

import(
	"fmt"

	"github.com/gtck520/ConsoleIM/models/basemodel"
)

func TableName(name string) string {
	return fmt.Sprintf("%s%s%s", basemodel.GetTablePrefix(),"com_", name)
}