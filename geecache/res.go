package geecache

type ResponseObj struct {
	data []int64
}

type RequestObj struct {
	key     string
	count   int32
	startNo int32
	endNo   int32
}

func NewRequestObj(key string, count, startNo, endNo int32) *RequestObj {
	return &RequestObj{
		key: key, count: count, startNo: startNo, endNo: endNo,
	}
}

func NewResponseObj(data []int64) *ResponseObj {
	return &ResponseObj{
		data: data,
	}
}
