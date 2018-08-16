package exchange

import (
	"reflect"

	"github.com/aghape/aghape"
	"github.com/aghape/aghape/resource"
	"github.com/aghape/roles"
)

// Meta defines importable/exportable fields
type Meta struct {
	base       *Resource
	resource.Meta
	Header     string
	Valuer     func(interface{}, *qor.Context) interface{}
	Setter     func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) error
	Permission *roles.Permission
}

func NewMeta(name string) *Meta {
	m := &Meta{}
	m.MetaName = &resource.MetaName{Name:name}
	return m
}

// GetMetas get defined sub metas
func (meta *Meta) GetMetas() []resource.Metaor {
	return []resource.Metaor{}
}

func (meta *Meta) GetContextMetas(recorde interface{}, context *qor.Context) []resource.Metaor {
	return meta.GetMetas()
}

// GetResource get its resource
func (meta *Meta) GetResource() resource.Resourcer {
	return nil
}

func (meta *Meta) updateMeta() {
	if meta.MetaName == nil {
		meta.MetaName = &resource.MetaName{}
	}
	meta.Meta = resource.Meta{
		MetaName:     meta.MetaName,
		Alias:        meta.Alias,
		FieldName:    meta.FieldName,
		Setter:       meta.Setter,
		Valuer:       meta.Valuer,
		Permission:   meta.Permission,
		BaseResource: meta.base,
	}

	meta.PreInitialize()
	if meta.FieldStruct != nil {
		if injector, ok := reflect.New(meta.FieldStruct.Struct.Type).Interface().(resource.ConfigureMetaBeforeInitializeInterface); ok {
			injector.ConfigureQorMetaBeforeInitialize(meta)
		}
	}

	meta.Initialize()

	if meta.FieldStruct != nil {
		if injector, ok := reflect.New(meta.FieldStruct.Struct.Type).Interface().(resource.ConfigureMetaInterface); ok {
			injector.ConfigureQorMeta(meta)
		}
	}

	meta.SetFormattedValuer(func(record interface{}, context *qor.Context) interface{} {
		if valuer := meta.GetValuer(); valuer != nil {
			result := valuer(record, context)

			if reflectValue := reflect.ValueOf(result); reflectValue.IsValid() {
				if reflectValue.Kind() == reflect.Ptr {
					if reflectValue.IsNil() || !reflectValue.Elem().IsValid() {
						return ""
					}

					result = reflectValue.Elem().Interface()
				}

				return result
			}
		}
		return ""
	})
}
