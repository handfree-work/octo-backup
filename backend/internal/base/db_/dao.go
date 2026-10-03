package db_

import (
	"context"
	"database/sql"
	"fmt"
	"handfree-work/octo-backup/internal/base/conv_"
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/error_/code_"
	"handfree-work/octo-backup/internal/base/log_"
	"handfree-work/octo-backup/internal/base/reflect_"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/olekukonko/errors"
	"gorm.io/gorm"
)

type OrderBy struct {
	Name string `json:"name"`
	Asc  bool   `json:"asc"`
}

type Page struct {
	Total   int64      `json:"total,omitempty"`
	Limit   int64      `json:"limit"`
	Start   int64      `json:"start"`
	OrderBy []*OrderBy `json:"orderBy,omitempty"`
}

type FilterItem struct {
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	StringValue *string `json:"stringValue,omitempty"`
	Int64Value  *int64  `json:"int64Value,omitempty"`
	BoolValue   *bool   `json:"boolValue,omitempty"`
}

type PageReq[T any] struct {
	Query   *T
	Page    *Page
	Filters []*FilterItem `json:"filter,omitempty"`
}

type PageReply[T any] struct {
	List []*T
	Page *Page
}

type DbCtx struct {
	Ctx context.Context
	Db  *gorm.DB
}

func NewCtx(ctx context.Context, db *gorm.DB) *DbCtx {
	return &DbCtx{
		Ctx: ctx,
		Db:  db,
	}
}

func (t *DbCtx) GetTx() *gorm.DB {
	tx := t.Ctx.Value("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return nil
}

func (t *DbCtx) GetDb() *gorm.DB {
	tx := t.GetTx()
	if tx != nil {
		return tx
	}
	return t.Db
}

func (t *DbCtx) SetTx(tx *gorm.DB) {
	t.Ctx = context.WithValue(t.Ctx, "tx", tx)
}

type Dao[T any] struct {
	//当前dao使用的db
	currentDb *gorm.DB
	//原始tx
	originTx *gorm.DB
}

type WhereItem struct {
	Where string
	Value interface{}
}

func Transaction(dbCtx *DbCtx, f func(tx *DbCtx) error, opts ...*sql.TxOptions) error {
	tx := dbCtx.GetTx()
	if tx != nil {
		return f(dbCtx)
	}

	return dbCtx.GetDb().Transaction(func(tx *gorm.DB) error {
		transactionCtx := NewCtx(context.WithValue(dbCtx.Ctx, "tx", tx), dbCtx.Db)
		return f(transactionCtx)
	}, opts...)
}

func New[T any](t *DbCtx) *Dao[T] {
	tx := t.GetDb()
	var m T
	return &Dao[T]{
		currentDb: tx.Model(&m),
		originTx:  tx,
	}
}

func NewFromCtx[T any](ctx context.Context, db *gorm.DB) *Dao[T] {
	dbCtx := NewCtx(ctx, db)
	return New[T](dbCtx)
}

func (d *Dao[T]) Reset() *Dao[T] {
	var m T
	d.currentDb = d.originTx.Model(&m)
	return d
}
func (d *Dao[T]) GetTx() *gorm.DB {
	return d.originTx
}

func (d *Dao[T]) Model() *gorm.DB {
	var m T
	return d.originTx.Model(&m)
}

func (d *Dao[T]) BuildPageReq(query interface{}, page interface{}) *PageReq[T] {
	var page_ Page
	var query_ T
	conv_.Convert(page, &page_)
	conv_.Convert(query, &query_)
	return &PageReq[T]{
		Query: &query_,
		Page:  &page_,
	}
}

func (d *Dao[T]) BuildFilterPageReq(query interface{}, page interface{}, filters interface{}) *PageReq[T] {
	var page_ Page
	var query_ T
	var filters_ []*FilterItem
	conv_.Convert(page, &page_)
	conv_.Convert(query, &query_)
	if filters != nil {
		conv_.ConvertList(filters, &filters_)
	}
	return &PageReq[T]{
		Query:   &query_,
		Page:    &page_,
		Filters: filters_,
	}
}

var SafeSymbols = []string{"=", "<>", "<", "<=", ">", ">=", "like"}

func (d *Dao[T]) FindPage(in *PageReq[T], more ...func(*gorm.DB)) (*[]T, error) {
	d.Reset()
	var entities []T
	var query = in.Query
	db := d.currentDb.Model(&query)
	if query != nil {
		db.Where(query)
	}
	if more != nil {
		for _, f := range more {
			f(db)
		}
	}
	if in.Filters != nil {
		for _, item := range in.Filters {
			//检查symbol
			if item.Symbol == "" || item.Name == "" {
				return nil, error_.NewTextError("symbol或name为空")
			}
			if !slices.Contains(SafeSymbols, item.Symbol) {
				return nil, error_.NewTextError("symbol不安全:" + item.Symbol)
			}

			if item.StringValue != nil {
				if item.Symbol == "like" {
					db.Where(fmt.Sprintf("`%s` like ?", item.Name), "%"+*item.StringValue+"%")
				} else {
					db.Where(fmt.Sprintf("`%s` %s ?", item.Name, item.Symbol), item.StringValue)
				}
			} else if item.Int64Value != nil {
				db.Where(fmt.Sprintf("`%s` %s ?", item.Name, item.Symbol), *item.Int64Value)
			} else if item.BoolValue != nil {
				db.Where(fmt.Sprintf("`%s` %s ?", item.Name, item.Symbol), *item.BoolValue)
			}
		}
	}
	if in.Page.OrderBy != nil {
		for _, orderBy := range in.Page.OrderBy {
			if orderBy.Asc {
				db.Order(fmt.Sprintf("`%s` asc", orderBy.Name))
			} else {
				db.Order(fmt.Sprintf("`%s` desc", orderBy.Name))
			}
		}
	} else {
		db.Order("id desc")
	}
	if in.Page != nil {
		var total int64
		var err = db.Count(&total).Error
		if err != nil {
			log_.Error("查询count失败", err)
			return nil, error_.NewTextError("查询count失败")
		}
		in.Page.Total = total
	}
	var err = db.Scopes(d.Paginate(in.Page)).Find(&entities).Error
	if err != nil {
		return nil, error_.NewTextError("查询page失败")
	}
	return &entities, err
}
func (d *Dao[T]) Paginate(info interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := reflect_.GetAttrValue(info, "Limit", reflect.Int).Int()
		start := reflect_.GetAttrValue(info, "Start", reflect.Int).Int()
		switch {
		case limit > 1000:
			limit = 1000
			break
		case limit < 0:
			limit = 10
			break
		}
		return db.Offset(int(start)).Limit(int(limit))
	}
}

func (d *Dao[T]) FindList(query *T, orderBy []*OrderBy, more ...func(*gorm.DB)) ([]*T, error) {
	d.Reset()
	var entities []*T
	db := d.currentDb.Model(query)
	if query != nil {
		db.Where(query)
	}
	if len(more) > 0 {
		for _, f := range more {
			if f != nil {
				f(db)
			}
		}
	}
	if orderBy != nil {
		for _, item := range orderBy {
			if item.Asc {
				db.Order("`" + item.Name + "` asc")
			} else {
				db.Order("`" + item.Name + "` desc")
			}
		}
	} else {
		db.Order("id desc")
	}

	var err = db.Debug().Find(&entities).Error
	if err != nil {
		return nil, error_.NewTextError("查询list失败")
	}
	return entities, err
}

func (d *Dao[T]) FindOne(query *T, more ...func(*gorm.DB)) (*T, error) {
	d.Reset()
	var entity T

	tx := d.currentDb.Where(query)
	if len(more) > 0 {
		for _, f := range more {
			f(tx)
		}
	}
	if err := tx.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log_.Error("数据库查询失败", err)
		return nil, error_.NewTextError("查询失败")
	}
	return &entity, nil
}

func (d *Dao[T]) GetForUpdate(id int64, more ...func(*gorm.DB)) (*T, error) {
	forUpdate := func(db *gorm.DB) {
		db.Set("gorm:query_option", "FOR UPDATE")
	}
	more = append(more, forUpdate)
	return d.GetById(id, more...)
}

func (d *Dao[T]) GetById(id int64, more ...func(*gorm.DB)) (*T, error) {
	d.Reset()
	if id <= 0 {
		return nil, error_.NewCodeTextError(code_.ParamError_, "id")
	}

	var query T
	idCondition := func(db *gorm.DB) {
		db.Where("id = ?", id)
	}
	more = append(more, idCondition)
	return d.FindOne(&query, more...)
}

func (d *Dao[T]) GetByPtrId(id *int64, more ...func(*gorm.DB)) (*T, error) {
	if id == nil || *id <= 0 {
		return nil, error_.NewTextError("id不能为空")
	}
	return d.GetById(*id, more...)
}

func (d *Dao[T]) NewModel(id *int64) *T {
	var bean T
	name := reflect.ValueOf(&bean).Elem().FieldByName("BaseModel")
	name.Set(reflect.ValueOf(BaseModel{Id: id}))
	return &bean
}

func (d *Dao[T]) NewPtrModel(id *int64) *T {
	var bean T
	name := reflect.ValueOf(&bean).Elem().FieldByName("PtrBaseModel")
	name.Set(reflect.ValueOf(BaseModel{Id: id}))
	return &bean
}

func (d *Dao[T]) SetId(entity *T, id *int64) *T {
	name := reflect.ValueOf(entity).Elem().FieldByName("BaseModel")
	name.Set(reflect.ValueOf(BaseModel{Id: id}))
	return entity
}

func (d *Dao[T]) Create(entity *T) error {
	d.Reset()
	tx := d.currentDb.Create(entity)
	if tx.Error != nil {
		log_.Error("insert error", tx.Error)
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return error_.NewTextError("记录已存在")
		}
		return error_.NewCodeTextError(code_.DbAddError_, tx.Error.Error())
	}
	return tx.Error
}

func (d *Dao[T]) Update(entity *T) (int64, error) {
	d.Reset()
	var query T
	queryBM := reflect.ValueOf(&query).Elem().FieldByName("BaseModel")
	entityBM := reflect.ValueOf(entity).Elem().FieldByName("BaseModel")
	queryBM.Set(reflect.ValueOf(BaseModel{
		Id: tea.Int64(entityBM.FieldByName("Id").Int()),
	}))
	//TODO 这里会把update at ，create	at 也作为查询条件
	res := d.currentDb.Where(&query).UpdateColumns(entity)

	if res.Error != nil {
		log_.Error(res.Error)
	}
	return res.RowsAffected, res.Error
}

func (d *Dao[T]) UpdateByPtrId(entity *T) (int64, error) {
	d.Reset()
	entityBM := reflect.ValueOf(entity).Elem().FieldByName("PtrBaseModel")
	//TODO 这里会把update at ，create	at 也作为查询条件
	queryBM := d.NewPtrModel(tea.Int64(entityBM.FieldByName("Id").Int()))
	res := d.currentDb.Where(queryBM).UpdateColumns(entity)

	if res.Error != nil {
		log_.Error(res.Error)
	}
	return res.RowsAffected, res.Error
}

func (d *Dao[T]) UpdateById(id int64, entity *T) (int64, error) {
	d.Reset()
	if id <= 0 {
		return 0, error_.NewFormatError(code_.ParamError_, "id")
	}
	res := d.currentDb.Where("id = ?", id).UpdateColumns(entity)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return res.RowsAffected, res.Error
}

func (d *Dao[T]) UpdateWhere(query *T, entity *T, where ...func(*gorm.DB)) (int64, error) {
	d.Reset()
	tx := d.currentDb.Where(query)
	if where != nil {
		for _, f := range where {
			f(tx)
		}
	}
	var db = tx
	if entity != nil {
		db = tx.UpdateColumns(entity)
	}

	d.Reset()
	if tx.Error != nil {
		log_.Error(tx.Error)
	}
	return db.RowsAffected, db.Error
}

func (d *Dao[T]) UpdateMap(query *T, target *map[string]any, where ...func(*gorm.DB)) (int64, error) {
	d.Reset()
	tx := d.currentDb.Where(query)
	if where != nil {
		for _, f := range where {
			f(tx)
		}
	}
	if target != nil {
		tx = tx.Updates(target)
	}

	d.Reset()
	if tx.Error != nil {
		log_.Error(tx.Error)
	}
	return tx.RowsAffected, tx.Error
}

func (d *Dao[T]) Delete(id *int64) (int64, error) {
	d.Reset()
	query := d.NewModel(id)
	res, err := d.DeleteWhere(query, nil)
	d.Reset()
	return res, err

}
func (d *Dao[T]) DeleteByIds(ids []int64) (int64, error) {
	d.Reset()
	var query T
	res := d.currentDb.Where("id in ?", ids).Delete(&query)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return res.RowsAffected, res.Error
}

func (d *Dao[T]) DeleteWhere(query *T, wheres []*WhereItem, more ...func(*gorm.DB)) (int64, error) {
	d.Reset()
	tx := d.currentDb
	if query != nil {
		tx = tx.Where(query)
	}

	if wheres != nil {
		for _, where := range wheres {
			tx = tx.Where(where.Where, where.Value)
		}
	}
	if more != nil {
		for _, f := range more {
			f(tx)
		}
	}
	var entity T
	res := tx.Delete(&entity)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return res.RowsAffected, res.Error
}

func (d *Dao[T]) Count(query *T, more ...func(*gorm.DB)) (int64, error) {
	d.Reset()
	var count = int64(-1)
	tx := d.currentDb.Where(query)
	if more != nil {
		for _, f := range more {
			f(tx)
		}
	}
	res := tx.Count(&count)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return count, res.Error
}

func (d *Dao[T]) Sum(query *T, field string, more ...func(*gorm.DB)) (int64, error) {
	d.Reset()
	var count = int64(-1)
	tx := d.currentDb.Where(query)
	if more != nil {
		for _, f := range more {
			f(tx)
		}
	}
	res := tx.Select(fmt.Sprintf("sum(%s)", field)).Scan(&count)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return count, res.Error
}

func (d *Dao[T]) Avg(query *T, field string, more ...func(*gorm.DB)) (float64, error) {
	d.Reset()
	var average float64
	tx := d.currentDb.Where(query)
	if more != nil {
		for _, f := range more {
			f(tx)
		}
	}
	res := tx.Select(fmt.Sprintf("COALESCE(AVG(%s), 0)", field)).Scan(&average)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	return average, res.Error
}

func (d *Dao[T]) SumMulti(query *T, where func(*gorm.DB), fields ...string) (map[string]int64, error) {
	d.Reset()
	tx := d.currentDb.Where(query)
	if where != nil {
		where(tx)
	}
	selects := make([]string, len(fields))
	for i, field := range fields {
		selects[i] = fmt.Sprintf("COALESCE(SUM(%s), 0) as %s", field, field)
	}
	tx.Select(strings.Join(selects, ", "))

	var results map[string]interface{}
	res := tx.Scan(&results)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	var ints = make(map[string]int64, len(fields))
	for _, field := range fields {
		parseInt, err := strconv.ParseInt(results[field].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		ints[field] = parseInt
	}
	return ints, res.Error
}

func (d *Dao[T]) SumAndCount(query *T, where func(*gorm.DB), fields ...string) (map[string]int64, error) {
	d.Reset()
	tx := d.currentDb.Where(query)
	if where != nil {
		where(tx)
	}
	selects := make([]string, len(fields)+1)
	for i, field := range fields {
		selects[i] = fmt.Sprintf("COALESCE(SUM(%s), 0) as %s", field, field)
	}
	selects[len(fields)] = fmt.Sprintf("COALESCE(Count(id), 0) as _count")
	tx.Select(strings.Join(selects, ", "))

	var results map[string]interface{}
	res := tx.Scan(&results)
	if res.Error != nil {
		log_.Error(res.Error)
	}
	allFields := append(fields, "_count")
	var ints = make(map[string]int64, len(allFields))
	for _, field := range fields {
		parseInt, err := strconv.ParseInt(results[field].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		ints[field] = parseInt
	}
	ints["_count"] = results["_count"].(int64)
	return ints, res.Error
}

func (d *Dao[T]) CheckPermission(id int64, userId int64, userIdName string) error {

	if userIdName == "" {
		userIdName = "user_id"
	}
	more := func(db *gorm.DB) {
		db.Where("id = ?", id)
		db.Where("? = ?", userIdName, userId)
	}
	var bean T
	count, err := d.Count(&bean, more)

	if err != nil {
		log_.Error(err)
		return err
	}
	if count == 0 {
		return error_.NewTextError("Permission denied")
	}

	return nil
}
