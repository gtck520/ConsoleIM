package repository

import (
	//"github.com/jinzhu/gorm"

	"github.com/gtck520/ConsoleIM/common/logger"
	//models "github.com/gtck520/ConsoleIM/models/common"
)

//GroupRepository 注入IDb
type UserGroupRepository struct {
	Log  logger.ILogger `inject:""`
	Base BaseRepository `inject:"inline"`
}
