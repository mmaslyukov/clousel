package store

type IPortStoreAdapterBookRepo interface {
	StoreAddBookEntry(entry *BookEntry) IError
	StoreMarkBookEntryBySessionIdWithData(sessionId Session, status string, err *string) IError
	StoreReadBookEntryBySessionId(sessionId Session) (BookEntry, IError)
	StoreReadBookEntriesByCarosuelId(carId Carousel) ([]BookEntry, IError)
}
