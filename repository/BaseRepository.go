package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/gtck520/ConsoleIM/common/datasource"
	"github.com/gtck520/ConsoleIM/common/logger"
)

//BaseRepository 注入IDb,Logger
type BaseRepository struct {
	Source datasource.IDb `inject:""`
	Log    logger.ILogger `inject:""`
}

type Where struct {
	Op  string                 //操作符
	Wh1 map[string]interface{} //条件1
	Wh2 map[string]interface{} //条件2
}

// Create 创建实体
func (b *BaseRepository) Create(value interface{}) error {
	return b.Source.DB().Create(value).Error
}

// Save 保存实体
func (b *BaseRepository) Save(value interface{}) error {
	return b.Source.DB().Save(value).Error
}

// // Updates 更新实体
// func (b *BaseRepository) Updates(model interface{}, value interface{}) error {
// 	return b.Source.DB().Model(model).Updates(value).Error
// }

// DeleteByWhere 根据条件删除实体
func (b *BaseRepository) DeleteByWhere(model, where interface{}) (count int64, err error) {
	db := b.Source.DB().Where(where).Delete(model)
	err = db.Error
	if err != nil {
		b.Log.Errorf("删除实体出错", err)
		return
	}
	count = db.RowsAffected
	return
}

// DeleteByID 根据id删除实体
func (b *BaseRepository) DeleteByID(model interface{}, id int) error {
	return b.Source.DB().Where("id=?", id).Delete(model).Error
}

// DeleteByIDS 根据多个id删除多个实体
func (b *BaseRepository) DeleteByIDS(model interface{}, ids []int) (count int64, err error) {
	db := b.Source.DB().Where("id in (?)", ids).Delete(model)
	err = db.Error
	if err != nil {
		b.Log.Errorf("删除多个实体出错", err)
		return
	}
	count = db.RowsAffected
	return
}

// First 根据条件获取一个实体
func (b *BaseRepository) First(where interface{}, out interface{}, selects ...string) error {
	db := b.Source.DB().Where(where)
	if len(selects) > 0 {
		for _, sel := range selects {
			db = db.Select(sel)
		}
	}
	return db.First(out).Error
}

// FirstByID 根据条件获取一个实体
func (b *BaseRepository) FirstByID(out interface{}, id uint) error {
	return b.Source.DB().First(out, id).Error
}

// Find 根据条件返回数据
func (b *BaseRepository) Find(where interface{}, out interface{}, sel string, orders ...string) error {
	var where1 interface{}
	var where2 interface{}
	value, ok := where.(*Where)
	if ok {
		where1 = value.Wh1
		where2 = value.Wh2
	} else {
		where1 = where
	}

	db := b.Source.DB().Where(where1)
	if sel != "" {
		db = db.Select(sel)
	}
	if value.Op == "or" {
		db = db.Or(where2)
	}
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// GetPages 分页返回数据
func (b *BaseRepository) GetPages(model interface{}, out interface{}, pageIndex, pageSize int, totalCount *uint64, where interface{}, orders ...string) error {
	db := b.Source.DB().Model(model).Where(model)
	db = db.Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	err := db.Count(totalCount).Error
	if err != nil {
		b.Log.Errorf("查询总数出错", err)
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

// PluckList 查询 model 中的一个列作为切片
func (b *BaseRepository) PluckList(model, where interface{}, out interface{}, fieldName string) error {
	return b.Source.DB().Model(model).Where(where).Pluck(fieldName, out).Error
}

//GetTransaction 获取事务
func (b *BaseRepository) GetTransaction() *gorm.DB {
	return b.Source.DB().Begin()
}
