package errors

const (
	TypeExpenseNotFound         = "EXPENSE_NOT_FOUND"
	TypeExpenseAlreadyApproved  = "EXPENSE_ALREADY_APPROVED"
	TypeExpenseAlreadyRejected  = "EXPENSE_ALREADY_REJECTED"
	TypeExpenseInvalidStatus    = "EXPENSE_INVALID_STATUS"
	TypeExpenseCannotDelete     = "EXPENSE_CANNOT_DELETE"
	TypeExpenseAmountInvalid    = "EXPENSE_AMOUNT_INVALID"
	TypeExpenseNotPending       = "EXPENSE_NOT_PENDING"
	TypeExpenseNotOwnedByUser   = "EXPENSE_NOT_OWNED_BY_USER"
	TypeBatchNotFound           = "BATCH_NOT_FOUND"
	TypeBatchAlreadyCompleted   = "BATCH_ALREADY_COMPLETED"
	TypeBatchNoExpenses         = "BATCH_NO_EXPENSES"
	TypeBatchExpensesMixed      = "BATCH_EXPENSES_MIXED"
	TypeBatchCreateFailed       = "BATCH_CREATE_FAILED"
	TypeBatchNotOwnedByApprover = "BATCH_NOT_OWNED_BY_APPROVER"
)

var (
	ErrExpenseNotFound         = New(TypeExpenseNotFound, "报销单不存在")
	ErrExpenseAlreadyApproved  = New(TypeExpenseAlreadyApproved, "报销单已审批，无法修改")
	ErrExpenseAlreadyRejected  = New(TypeExpenseAlreadyRejected, "报销单已被驳回")
	ErrExpenseInvalidStatus    = New(TypeExpenseInvalidStatus, "报销单状态无效")
	ErrExpenseCannotDelete     = New(TypeExpenseCannotDelete, "报销单无法删除")
	ErrExpenseAmountInvalid    = New(TypeExpenseAmountInvalid, "报销金额无效")
	ErrExpenseNotPending       = New(TypeExpenseNotPending, "只能审批待审批状态的报销单")
	ErrExpenseNotOwnedByUser   = New(TypeExpenseNotOwnedByUser, "您无权操作该报销单")
	ErrBatchNotFound           = New(TypeBatchNotFound, "审批批次不存在")
	ErrBatchAlreadyCompleted   = New(TypeBatchAlreadyCompleted, "审批批次已完成")
	ErrBatchNoExpenses         = New(TypeBatchNoExpenses, "审批批次中没有报销单")
	ErrBatchExpensesMixed      = New(TypeBatchExpensesMixed, "批次中包含非待审批状态的报销单")
	ErrBatchCreateFailed       = New(TypeBatchCreateFailed, "创建审批批次失败")
	ErrBatchNotOwnedByApprover = New(TypeBatchNotOwnedByApprover, "您无权操作该审批批次")
)
