package controllers

import (
	"context"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/errors"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/auth"
	authutils "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/utils"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/log"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	"net/http"

	qkm "github.com/lugondev/wallet-signer-manager/pkg/client"

	"github.com/gorilla/mux"
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
)

type Builder struct {
	accountsCtrl *AccountsController
	auth         Auth
}

type Auth struct {
	checker      auth.Checker
	multitenancy bool
}

func NewBuilder(
	multitenancy bool, ucs usecases.UseCases, keyManagerClient qkm.KeyManagerClient, qkmStoreID string,
	jwt, key auth.Checker) *Builder {
	return &Builder{
		auth: Auth{
			checker:      auth.NewCombineCheckers(key, jwt),
			multitenancy: multitenancy,
		},
		accountsCtrl: NewAccountsController(ucs, keyManagerClient, qkmStoreID),
	}
}

func (b *Builder) Build(_ context.Context, _ string, _ func(response *http.Response) error) (http.Handler, error) {
	router := mux.NewRouter()
	b.accountsCtrl.Append(router)

	return router, nil
}

func (b *Builder) BuildHandlerFunc(router *mux.Router, path string, f func(http.ResponseWriter, *http.Request)) *mux.Router {
	router.HandleFunc(path, f)

	return router
}

func (b *Builder) BuildHandle(router *mux.Router, path string, h http.Handler) *mux.Router {
	router.Handle(path, h)

	return router
}

func (b *Builder) BuildRouter(router *mux.Router, subPath string) *mux.Router {
	router.Use(b.AuthMiddlewareHandler)
	subRouter := router.PathPrefix(subPath).Subrouter()
	b.accountsCtrl.Append(subRouter)

	return router
}

func (b *Builder) AuthMiddlewareHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !b.auth.multitenancy {
			userInfo := multitenancy.DefaultUser()
			b.serveNext(rw, req.WithContext(multitenancy.WithUserInfo(req.Context(), userInfo)), h)
			return
		}

		// Extract Authorization credentials from HTTP headers
		authCtx := authutils.WithAuthorization(
			req.Context(),
			authutils.GetAuthorizationHeader(req),
		)

		// Extract API Key credentials from HTTP headers
		authCtx = authutils.WithAPIKey(
			authCtx,
			authutils.GetAPIKeyHeaderValue(req),
		)

		// Extract TenantID from HTTP headers
		authCtx = authutils.WithTenantID(
			authCtx,
			authutils.GetTenantIDHeaderValue(req),
		)

		// Extract Username from HTTP headers
		authCtx = authutils.WithUsername(
			authCtx,
			authutils.GetUsernameHeaderValue(req),
		)

		userInfo, err := b.auth.checker.Check(authCtx)
		if err != nil {
			log.FromContext(authCtx).WithError(err).Errorf("unauthorized request user info")
			b.writeUnauthorized(rw, err)
			return
		}

		if userInfo != nil {
			// Bypass JWT authentication
			log.FromContext(authCtx).
				WithField("tenant_id", userInfo.TenantID).
				WithField("username", userInfo.Username).
				WithField("allowed_tenants", userInfo.AllowedTenants).
				Debugf("authentication succeeded (%s)", userInfo.AuthMode)

			b.serveNext(rw, req.WithContext(multitenancy.WithUserInfo(authCtx, userInfo)), h)
			return
		}

		err = errors.UnauthorizedError("missing required credentials")
		log.FromContext(authCtx).WithError(err).Errorf("unauthorized request no auth info")
		b.writeUnauthorized(rw, err)
	})
}

func (b *Builder) writeUnauthorized(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusUnauthorized)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusUnauthorized, err.Error())))
}

func (b *Builder) serveNext(rw http.ResponseWriter, req *http.Request, h http.Handler) {
	// Remove authorization header
	// So possibly another Authorization will be set by Proxy
	authutils.DeleteAuthorizationHeaderValue(req)
	authutils.DeleteAPIKeyHeaderValue(req)

	// Execute next handlers
	h.ServeHTTP(rw, req)
}
