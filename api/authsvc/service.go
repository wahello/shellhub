package authsvc

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/cnf/structhash"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/models"
	"gopkg.in/go-playground/validator.v9"
)

var ErrUnauthorized = errors.New("unauthorized")

type Service interface {
	AuthDevice(ctx context.Context, req *models.DeviceAuthRequest) (*models.DeviceAuthResponse, error)
	AuthUser(ctx context.Context, req models.UserAuthRequest) (*models.UserAuthResponse, error)
	AuthGetToken(ctx context.Context, tenant string) (*models.UserAuthResponse, error)
	AuthPublicKey(ctx context.Context, req *models.PublicKeyAuthRequest) (*models.PublicKeyAuthResponse, error)
	AuthSwapToken(ctx context.Context, ID, tenant string) (*models.UserAuthResponse, error)
	AuthUserInfo(ctx context.Context, username, tenant, token string) (*models.UserAuthResponse, error)
	PublicKey() *rsa.PublicKey
}

type service struct {
	store   store.Store
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	cache   *cache.Cache
}

type config struct {
	RedisUri string `envconfig:"redis_uri" default:"redis://redis:6379"`
}

var myCache *cache.Cache

func NewService(store store.Store, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) Service {
	if privKey == nil || pubKey == nil {
		var err error
		privKey, pubKey, err = loadKeys()
		if err != nil {
			panic(err)
		}
	}
	var cfg config
	if err := envconfig.Process("api", &cfg); err != nil {
		panic(err.Error())
	}

	opt, err := redis.ParseURL(cfg.RedisUri)
	if err != nil {
		panic(err)
	}

	if myCache == nil {
		myCache = cache.New(&cache.Options{
			Redis: redis.NewClient(opt),
		})
	}

	return &service{store, privKey, pubKey, myCache}
}

func (s *service) AuthDevice(ctx context.Context, req *models.DeviceAuthRequest) (*models.DeviceAuthResponse, error) {
	uid := sha256.Sum256(structhash.Dump(req.DeviceAuth, 1))
	key := hex.EncodeToString(uid[:])

	type Device struct {
		Name      string
		Namespace string
	}

	var value Device

	if err := s.cache.Get(ctx, key, &value); err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, models.DeviceAuthClaims{
			UID: hex.EncodeToString(uid[:]),
			AuthClaims: models.AuthClaims{
				Claims: "device",
			},
		})

		tokenStr, err := token.SignedString(s.privKey)
		if err != nil {
			return nil, err
		}

		return &models.DeviceAuthResponse{
			UID:       hex.EncodeToString(uid[:]),
			Token:     tokenStr,
			Name:      value.Name,
			Namespace: value.Namespace,
		}, nil

	}

	device := models.Device{
		UID:       hex.EncodeToString(uid[:]),
		Identity:  req.Identity,
		Info:      req.Info,
		PublicKey: req.PublicKey,
		TenantID:  req.TenantID,
		LastSeen:  time.Now(),
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	hostname := strings.ToLower(req.DeviceAuth.Hostname)

	if err := s.store.DeviceCreate(ctx, device, hostname); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, models.DeviceAuthClaims{
		UID: hex.EncodeToString(uid[:]),
		AuthClaims: models.AuthClaims{
			Claims: "device",
		},
	})

	tokenStr, err := token.SignedString(s.privKey)
	if err != nil {
		return nil, err
	}

	if err := s.store.DeviceSetOnline(ctx, models.UID(device.UID), true); err != nil {
		return nil, err
	}

	for _, uid := range req.Sessions {
		if err := s.store.SessionSetLastSeen(ctx, models.UID(uid)); err != nil {
			continue
		}
	}

	dev, err := s.store.DeviceGetByUID(ctx, models.UID(device.UID), device.TenantID)
	if err != nil {
		return nil, err
	}

	namespace, err := s.store.NamespaceGet(ctx, device.TenantID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &Device{Name: dev.Name, Namespace: namespace.Name},
		TTL:   time.Second * 30,
	}); err != nil {
		return nil, err
	}

	return &models.DeviceAuthResponse{
		UID:       hex.EncodeToString(uid[:]),
		Token:     tokenStr,
		Name:      dev.Name,
		Namespace: namespace.Name,
	}, nil
}

func (s *service) AuthUser(ctx context.Context, req models.UserAuthRequest) (*models.UserAuthResponse, error) {
	user, err := s.store.UserGetByUsername(ctx, strings.ToLower(req.Username))
	if err != nil {
		user, err = s.store.UserGetByEmail(ctx, strings.ToLower(req.Username))
		if err != nil {
			return nil, err
		}

	}

	namespace, err := s.store.NamespaceGetFirst(ctx, user.ID)
	if err != nil && err != store.ErrNamespaceNoDocuments {
		return nil, err
	}

	tenant := ""
	if namespace != nil {
		tenant = namespace.TenantID
	}

	password := sha256.Sum256([]byte(req.Password))
	if user.Password == hex.EncodeToString(password[:]) {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, models.UserAuthClaims{
			Username: user.Username,
			Admin:    true,
			Tenant:   tenant,
			ID:       user.ID,
			AuthClaims: models.AuthClaims{
				Claims: "user",
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		})

		tokenStr, err := token.SignedString(s.privKey)
		if err != nil {
			return nil, err
		}
		return &models.UserAuthResponse{
			Token:  tokenStr,
			Name:   user.Name,
			ID:     user.ID,
			User:   user.Username,
			Tenant: tenant,
			Email:  user.Email,
		}, nil
	}

	return nil, errors.New("unauthorized")
}

func (s *service) AuthGetToken(ctx context.Context, ID string) (*models.UserAuthResponse, error) {
	user, err := s.store.UserGetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	namespace, err := s.store.NamespaceGetFirst(ctx, user.ID)
	if err != nil && err != store.ErrNamespaceNoDocuments {
		return nil, err
	}

	tenant := ""
	if namespace != nil {
		tenant = namespace.TenantID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, models.UserAuthClaims{
		Username: user.Username,
		Admin:    true,
		Tenant:   tenant,
		ID:       user.ID,
		AuthClaims: models.AuthClaims{
			Claims: "user",
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	tokenStr, err := token.SignedString(s.privKey)
	if err != nil {
		return nil, err
	}
	return &models.UserAuthResponse{
		Token:  tokenStr,
		Name:   user.Name,
		ID:     user.ID,
		User:   user.Username,
		Tenant: tenant,
		Email:  user.Email,
	}, nil
}

func (s *service) AuthPublicKey(ctx context.Context, req *models.PublicKeyAuthRequest) (*models.PublicKeyAuthResponse, error) {
	privKey, err := s.store.PrivateKeyGet(ctx, req.Fingerprint)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privKey.Data)
	if block == nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	digest := sha256.Sum256([]byte(req.Data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		return nil, err
	}

	return &models.PublicKeyAuthResponse{
		Signature: base64.StdEncoding.EncodeToString(signature),
	}, nil
}

func (s *service) AuthSwapToken(ctx context.Context, id, tenant string) (*models.UserAuthResponse, error) {
	namespace, err := s.store.NamespaceGet(ctx, tenant)
	if err != nil {
		return nil, err
	}

	user, err := s.store.UserGetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, i := range namespace.Members {
		if user.ID == i.(string) {
			token := jwt.NewWithClaims(jwt.SigningMethodRS256, models.UserAuthClaims{
				Username: user.Username,
				Admin:    true,
				Tenant:   namespace.TenantID,
				ID:       user.ID,
				AuthClaims: models.AuthClaims{
					Claims: "user",
				},
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				},
			})

			tokenStr, err := token.SignedString(s.privKey)
			if err != nil {
				return nil, err
			}
			return &models.UserAuthResponse{
				Token:  tokenStr,
				Name:   user.Name,
				ID:     user.ID,
				User:   user.Username,
				Tenant: namespace.TenantID,
				Email:  user.Email}, nil
		}
	}

	return nil, nil
}

func (s *service) AuthUserInfo(ctx context.Context, username, tenant, token string) (*models.UserAuthResponse, error) {
	user, err := s.store.UserGetByUsername(ctx, username)
	if err != nil {
		if err == store.ErrUserNoDocuments {
			return nil, ErrUnauthorized
		}

		return nil, err
	}

	if _, err = s.store.NamespaceGet(ctx, tenant); err != nil && tenant != "" {
		if err == store.ErrNamespaceNoDocuments {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	return &models.UserAuthResponse{
		Token:  token,
		Name:   user.Name,
		User:   user.Username,
		Tenant: tenant,
		ID:     user.ID,
		Email:  user.Email,
	}, nil
}
func (s *service) PublicKey() *rsa.PublicKey {
	return s.pubKey
}

func loadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signBytes, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, nil, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, nil, err
	}

	verifyBytes, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
