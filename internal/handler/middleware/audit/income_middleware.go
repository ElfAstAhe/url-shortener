package audit

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener/internal/handler/middleware"
	"github.com/ElfAstAhe/url-shortener/internal/service/audit"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
	"go.uber.org/zap"
)

type IncomeAuditPath struct {
	Method string
	Path   string
}

func NewIncomeAuditPath(method, path string) *IncomeAuditPath {
	return &IncomeAuditPath{
		Method: method,
		Path:   path,
	}
}

type IncomeAuditMiddleware struct {
	watchPaths []*IncomeAuditPath
	publisher  audit.IncomePublisher
	log        *zap.SugaredLogger
}

func NewIncomeAuditMiddleware(watchPaths []*IncomeAuditPath, log *zap.SugaredLogger, incomeObservers ...audit.IncomeObserver) *IncomeAuditMiddleware {
	middleware := &IncomeAuditMiddleware{
		watchPaths: watchPaths,
		publisher:  audit.NewIncomeEventService(log),
		log:        log,
	}

	//

	//    var observer audit.IncomeObserver =

	return middleware
}

func (ia *IncomeAuditMiddleware) Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		crw := middleware.NewCommonResponseWriter(rw)

		next.ServeHTTP(crw, r)

		if utils.IsSuccess(crw.Info.StatusCode) || utils.IsRedirection(crw.Info.StatusCode) {
			// ToDo: implement
		}
	})
}

func (ia *IncomeAuditMiddleware) isWatchable(method, path string) bool {
	for _, watchPath := range ia.watchPaths {
		// ToDo: need algorithm improvement
		if watchPath.Method == method && watchPath.Path == path {
			return true
		}
	}

	return false
}
