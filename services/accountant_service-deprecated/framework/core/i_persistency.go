package core

// import (
// 	"accountant_service/framework/utils"
// )
/**

Update(id).Where(func())
Delete()
*/

// type Table interface {
// 	Name() string
// }
// type IPersistency interface {
// 	Update(table Table) IPersistency
// 	SelectAll(table Table) IPersistency
// 	Select(table Table, name string) IPersistency
// 	Insert(table Table) IPersistency
// 	Delete(table Table) IPersistency
// }

// type IPersistencyUpdate interface {
// 	Set(name string, value any) IPersistencyUpdate
// 	Where(name string, value any) IPersistencyUpdate
// 	And() IPersistencyUpdate
// }

// type IPersistencySelect interface {
// 	Where(name string, value any) IPersistencyUpdate
// }

// type RecordInterface[T any, I any] interface {
// 	Create(record T) (T, error)
// 	Update(record T) (T, error)
// 	// UpdateBy(id ID, sets []Pair) error
// 	Delete(id I) error
// 	Read(id I) (utils.Optional[T], error)
// 	ReadOneBy(where func() string) (utils.Optional[T], error)
// 	ReadManyBy(where func() string) (utils.Optional[[]T], error)
// 	ReadAll() (utils.Optional[[]T], error)
// }
