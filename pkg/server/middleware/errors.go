package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tbd54566975/ssi-service/pkg/server/framework"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// Errors handles errors coming out of the call stack. It detects safe application
// errors (aka SafeError) that are used to respond to the requester in a
// normalized way. Unexpected errors (status >= 500) are logged.
func Errors() framework.Middleware {
	mw := func(handler framework.Handler) framework.Handler {
		wrapped := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("logger").Start(ctx, "business.middleware.errors")
			defer span.End()

			v, ok := ctx.Value(framework.KeyRequestState).(*framework.RequestState)
			if !ok {
				return framework.NewShutdownError("request state missing from context.")
			}

			if err := handler(ctx, w, r); err != nil {
				// log the error
				logrus.Printf("%s : ERROR : %v", v.TraceID, err)

				// send an error response back to the requester.
				if err := framework.RespondError(ctx, w, err); err != nil {
					return err
				}

				if ok := framework.IsShutdown(err); ok {
					return err
				}
			}

			return nil
		}

		return wrapped
	}

	return mw
}
