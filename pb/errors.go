package pb

import "errors"

var (
	ErrUnknown            = errors.New("iknore: ErrUnknown")            // 未知的錯誤，沒有足夠詳細的資訊可以解釋。
	ErrCancelled          = errors.New("iknore: ErrCancelled")          // 動作被取消了。例如：使用者已經中斷連線，剩餘的動作直接被 Context 中止。
	ErrInvalidArgument    = errors.New("iknore: ErrInvalidArgument")    // 錯誤的請求參數。例如：字元長度沒有達到條件。
	ErrDeadlineExceeded   = errors.New("iknore: ErrDeadlineExceeded")   // 超過最後期限。例如：請求逾時沒有完成、已經過期。
	ErrNotFound           = errors.New("iknore: ErrNotFound")           // 找不到指定資源。注意：如果是權限不符而無法存取，應該使用 ErrForbidden。
	ErrAlreadyExists      = errors.New("iknore: ErrAlreadyExists")      // 建立的資源早已存在、重複。例如：使用相同的使用者名稱。
	ErrPermissionDenied   = errors.New("iknore: ErrPermissionDenied")   // 權限不足。例如：擁有的權利不足已執行這件事，應該提昇權限。存取一個不屬於自己的資源。
	ErrResourceExhausted  = errors.New("iknore: ErrResourceExhausted")  // 資源或配額耗盡。例如：硬碟空間已滿、聊天室人數已滿、沒有剩餘次數。
	ErrFailedPrecondition = errors.New("iknore: ErrFailedPrecondition") // 執行此動作的前置條件沒有達成。例如：資料夾必須先清空才能移除。
	ErrAborted            = errors.New("iknore: ErrFailedPrecondition") // 執行手續中離，例如：某個異步執行緒發生異常，應該重試。
	ErrOutOfRange         = errors.New("iknore: ErrOutOfRange")         // 在範圍之外，比起 ErrFailedPrecondition 這更能敘述問題。例如：超過日期區間、指定的載入區塊超過範圍。
	ErrUnimplemented      = errors.New("iknore: ErrUnimplemented")      // 尚未實作的功能。
	ErrInternal           = errors.New("iknore: ErrInternal")           // 系統內部發生非預期的錯誤。
	ErrUnavailable        = errors.New("iknore: ErrUnavailable")        // 服務目前暫時不可用，應該重試。
	ErrDataLoss           = errors.New("iknore: ErrDataLoss")           // 相依的資料遺失且無法恢復。
	ErrUnauthenticated    = errors.New("iknore: ErrUnauthenticated")    // 沒有任何憑證或身份，例如：尚未登入，若登入之後沒有權限請用 ErrPermissionDenied。
)
