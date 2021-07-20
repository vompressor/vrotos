package async_rw

type RWCallback func(int, error)
type CopyCallback func(int64, error)
