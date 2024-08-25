package services

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/common"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/constants"
	"github.com/ehsan-hosseiny/golang-web-api/data/db"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/service_errors"
	"gorm.io/gorm"
)

type preload struct {
	string
}

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Database *gorm.DB
	Logger   logging.Logger
	Preloads []preload
}

func NewBaseService[T any, Tc any, Tu any, Tr any](cfg *config.Config) *BaseService[T, Tc, Tu, Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		Database: db.GetDb(),
		Logger:   logging.NewLogger(cfg),
	}
}

func (s *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {

	model, _ := common.TypeConverter[T](req)
	tx := s.Database.WithContext(ctx).Begin()

	err := tx.Create(model).
		Error
	if err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, id int, req *Tu) (*Tr, error) {

	updateMap, _ := common.TypeConverter[map[string]interface{}](req)
	(*updateMap)["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constants.UserIdKey).(float64)), Valid: true}
	(*updateMap)["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}

	model := new(T)
	tx := s.Database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where("id = ? and deleted_by is null", id).
		Updates(*updateMap).
		Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return s.GetById(ctx, id)
}

func (s *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context, id int) error {
	tx := s.Database.WithContext(ctx).Begin()

	model := new(T)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constants.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	deleteMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constants.UserIdKey).(float64)), Valid: true}
	deleteMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}

	if ctx.Value(constants.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	if cnt := tx.Model(model).
		Where("id = ? and deleted_by is null", id).
		Updates(deleteMap).
		RowsAffected; cnt == 0 {
		s.Logger.Error(logging.Postgres, logging.Delete, service_errors.RecordNotFound, nil)
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int) (*Tr, error) {

	model := new(T)

	err := s.Database.
		Where("id = ? and deleted_by is null", id).
		First(model).Error
	if err != nil {
		return nil, err
	}
	return common.TypeConverter[Tr](model)

}

// Paginate

func getQuery[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	query := make([]string, 0)
	query = append(query, "deleted_by is null")
	if filter.Filter != nil {
		for name, filter := range filter.Filter {
			fld, ok := typeT.FieldByName(name)
			if ok {
				fld.Name = common.ToSnakeCase(fld.Name)
				switch filter.Type {
				case "contains":
					query = append(query, fmt.Sprintf("%s ILike '%%%s%%'", fld.Name, filter.From))
				case "notContains":
					query = append(query, fmt.Sprintf("%s not ILike '%%%s%%'", fld.Name, filter.From))
				case "startsWith":
					query = append(query, fmt.Sprintf("%s ILike '%s%%'", fld.Name, filter.From))
				case "endsWith":
					query = append(query, fmt.Sprintf("%s ILike '%%%s'", fld.Name, filter.From))
				case "equals":
					query = append(query, fmt.Sprintf("%s = '%s'", fld.Name, filter.From))
				case "notEqual":
					query = append(query, fmt.Sprintf("%s != '%s'", fld.Name, filter.From))
				case "lessThan":
					query = append(query, fmt.Sprintf("%s < %s", fld.Name, filter.From))
				case "lessThanOrEqual":
					query = append(query, fmt.Sprintf("%s <= %s", fld.Name, filter.From))
				case "greaterThan":
					query = append(query, fmt.Sprintf("%s > %s", fld.Name, filter.From))
				case "greaterThanOrEqual":
					query = append(query, fmt.Sprintf("%s >= %s", fld.Name, filter.From))
				case "inRange":
					if fld.Type.Kind() == reflect.String {
						query = append(query, fmt.Sprintf("%s >= '%s'", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= '%s'", fld.Name, filter.To))
					} else {
						query = append(query, fmt.Sprintf("%s >= %s", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= %s", fld.Name, filter.To))
					}
				}
			}
		}
	}
	return strings.Join(query, " AND ")
}

// getSort
func getSort[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	sort := make([]string, 0)
	if filter.Sort != nil {
		for _, tp := range *filter.Sort {
			fld, ok := typeT.FieldByName(tp.ColId)
			if ok && (tp.Sort == "asc" || tp.Sort == "desc") {
				fld.Name = common.ToSnakeCase(fld.Name)
				sort = append(sort, fmt.Sprintf("%s %s", fld.Name, tp.Sort))
			}
		}
	}
	return strings.Join(sort, ", ")
}

// Preload
func Preload(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}
