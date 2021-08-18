package routes

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mitchellh/mapstructure"
	"github.com/shellhub-io/shellhub/api/apicontext"
	svc "github.com/shellhub-io/shellhub/api/services"
	api "github.com/shellhub-io/shellhub/pkg/api/client"
	"github.com/shellhub-io/shellhub/pkg/models"
)

const (
	AuthRequestURL   = "/auth"
	AuthDeviceURL    = "/devices/auth"
	AuthDeviceURLV2  = "/auth/device"
	AuthUserURL      = "/login"
	AuthUserURLV2    = "/auth/user"
	AuthUserTokenURL = "/auth/token/:tenant" //nolint:gosec
	AuthPublicKeyURL = "/auth/ssh"
)

func (h *Handler) AuthRequest(c apicontext.Context) error {
	token := c.Get("user").(*jwt.Token)
	rawClaims := token.Claims.(*jwt.MapClaims)

	switch claims := (*rawClaims)["claims"]; claims {
	case "user":
		var claims models.UserAuthClaims

		if err := DecodeMap(rawClaims, &claims); err != nil {
			return err
		}

		// Extract tenant and username from JWT
		c.Response().Header().Set("X-Tenant-ID", claims.Tenant)
		c.Response().Header().Set("X-Username", claims.Username)
		c.Response().Header().Set("X-ID", claims.ID)

		return nil
	case "device":
		var claims models.DeviceAuthClaims

		if err := DecodeMap(rawClaims, &claims); err != nil {
			return err
		}

		// Extract device UID from JWT
		c.Response().Header().Set(api.DeviceUIDHeader, claims.UID)

		return nil
	case "apiToken":
		var claims models.APITokenAuthClaims

		if err := DecodeMap(rawClaims, &claims); err != nil {
			return err
		}

		c.Response().Header().Set("X-ID", claims.ID)
		c.Response().Header().Set("X-Tenant-ID", claims.TenantID)

		return nil
	}

	return echo.ErrUnauthorized
}

func (h *Handler) AuthDevice(c apicontext.Context) error {
	var req models.DeviceAuthRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.service.AuthDevice(c.Ctx(), &req, c.Request().Header.Get("X-Real-IP"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthUser(c apicontext.Context) error {
	var req models.UserAuthRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.service.AuthUser(c.Ctx(), req)
	if err != nil {
		switch err {
		case svc.ErrForbidden:
			return c.JSON(http.StatusForbidden, nil)
		default:
			return echo.ErrUnauthorized
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthUserInfo(c apicontext.Context) error {
	username := c.Request().Header.Get("X-Username")
	tenant := c.Request().Header.Get("X-Tenant-ID")
	token := c.Request().Header.Get(echo.HeaderAuthorization)

	res, err := h.service.AuthUserInfo(c.Ctx(), username, tenant, token)
	if err != nil {
		switch err {
		case svc.ErrUnauthorized:
			return echo.ErrUnauthorized
		default:
			return err
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthGetToken(c apicontext.Context) error {
	res, err := h.service.AuthGetToken(c.Ctx(), c.Param("tenant"))
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthSwapToken(c apicontext.Context) error {
	id := ""
	if v := c.ID(); v != nil {
		id = v.ID
	}

	res, err := h.service.AuthSwapToken(c.Ctx(), id, c.Param("tenant"))
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthPublicKey(c apicontext.Context) error {
	var req models.PublicKeyAuthRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.service.AuthPublicKey(c.Ctx(), &req)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, res)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Get("ctx").(*apicontext.Context)

		jwt := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:        &jwt.MapClaims{},
			SigningKey:    ctx.Service().(svc.Service).PublicKey(),
			SigningMethod: "RS256",
		})

		return jwt(next)(c)
	}
}

func DecodeMap(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:  "json",
		Metadata: nil,
		Result:   output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
