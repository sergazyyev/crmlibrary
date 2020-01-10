package ocrmmodel

const (
	//Response codes
	SimpleErrCode     SimpleResponseCode = "ERROR"
	SimpleWarnCode    SimpleResponseCode = "WARN"
	SimpleEmptyCode   SimpleResponseCode = "EMPTY"
	SimpleSuccessCode SimpleResponseCode = "SUCCESS"
)

//Code of simple response
type SimpleResponseCode string

//Simple response struct
type SimpleResponse struct {
	Code    SimpleResponseCode `json:"code"`
	Message string             `json:"message"`
}

type Pageable struct {
	TotalElement    int  `json:"totalElement"`
	NumberOfElement int  `json:"numberOfElement"`
	TotalPages      int  `json:"totalPages"`
	First           bool `json:"first"`
	Last            bool `json:"last"`
	Page            int  `json:"page"`
	PageSize        int  `json:"-"`
}
