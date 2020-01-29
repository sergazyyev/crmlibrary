package ocrmtypes

type LogLevel string
type SimpleResponseCode string
type LeadType string
type LeadStatus string

const (
	LogLevelTrace           LogLevel           = "trace"
	LogLevelDebug           LogLevel           = "debug"
	LogLevelInfo            LogLevel           = "info"
	LogLevelWarn            LogLevel           = "warn"
	LogLevelErr             LogLevel           = "error"
	SimpleErrCode           SimpleResponseCode = "ERROR"
	SimpleWarnCode          SimpleResponseCode = "WARN"
	SimpleEmptyCode         SimpleResponseCode = "EMPTY"
	SimpleSuccessCode       SimpleResponseCode = "SUCCESS"
	RecommendedLead         LeadType           = "RECOMMENDED_LEAD"
	TmOnlineLead            LeadType           = "TM_ONLINE_LEAD"
	PartnerLead             LeadType           = "PARTNER_LEAD"
	OfficeOnlineLead        LeadType           = "OFFICE_ONLINE_LEAD"
	OfficeOnlineAccOpenLead LeadType           = "OFFICE_ONLINE_ACC_OPEN_LEAD"
	LeadStatusNew           LeadStatus         = "NEW"
	LeadStatusEnqueue       LeadStatus         = "ENQUEUE"
	LeadStatusDequeue       LeadStatus         = "DEQUEUE"
	LeadStatusOnVerifier    LeadStatus         = "ONVERIFIER"
	LeadStatusSuccess       LeadStatus         = "SUCCESS"
	LeadStatusNotProcess    LeadStatus         = "NOTPROCESSED"
	LeadStatusIntProcess    LeadStatus         = "INPROCESSING"
	LeadStatusException     LeadStatus         = "EXCEPTION"
	LeadStatusReProcess     LeadStatus         = "TOREPROCESS"
)
