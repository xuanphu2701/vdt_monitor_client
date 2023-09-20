package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
	"k8s.io/klog/v2"
)

type Token struct {
	Token string `json:"token"`
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func decryptToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRET_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		klog.Infof("Decrypt token success: %s", claims)
		return claims, nil
	} else {
		klog.Fatalf("Decrypt token failed: %s", err)
		return nil, err
	}
}

func main() {
	// client := redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName:    "mymaster",
	// 	SentinelAddrs: strings.Split(os.Getenv("SENTINEL_HOSTS_AND_PORTS"), " "),
	// 	Password:      os.Getenv("REDIS_PASSWORD"), // no password set
	// 	DB:            0,
	// })

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	klog.Info("Connected to Redis")

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	remote, err := url.Parse("http://mimir-1:8080")
	if err != nil {
		panic(err)
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
		return
	}

	reverseProxyHandler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			klog.Info("Redirect")
			tokenSplit := strings.Split(r.Header.Get("Authorization"), " ")
			// if len(tokenSplit) != 2 {
			// 	http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			// 	return
			// }
			tokenString := tokenSplit[1]
			tenantId := strings.Split(r.Header.Get("X-Scope-OrgID"), " ")[0]
			ctx := context.Background()
			value, err := client.Get(ctx, tokenString).Result()
			if err != nil {
				fmt.Println("Error compare token with Redis")
				panic(err)
			}

			if value != tenantId {
				fmt.Printf("Tenant %v is not the same with the one in the jwt key", tenantId)
				return
			}
			klog.Infof("Check tenancy of %v successful!", tenantId)
			// NO need to decrypt
			// tokenPayload, err := decryptToken(tokenString)
			// if err != nil {
			// 	panic(err)
			// }

			// r.URL.Path = "/"
			r.Host = remote.Host
			// r.Header.Del("Authorization")
			// Add http header from token payload
			// for k, v := range tokenPayload {
			// 	r.Header.Set(k, v.(string))
			// }
			p.ServeHTTP(w, r)
		}
	}

	jwtHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			klog.Info("Create token")
			tenantId := r.Header.Get("X-Scope-OrgID")
			klog.Info("Tenant Id: ", tenantId)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"X-Scope-OrgID": tenantId,
			})
			tokenString, err := token.SignedString([]byte(SECRET_KEY))
			ctx := context.Background()
			err = client.Set(ctx, tokenString, tenantId, 0).Err()
			if err != nil {
				panic(err)
			}
			data := Token{
				Token: tokenString,
			}
			w.Header().Set("Content-Type", "application/json")
			byteArray, err := json.Marshal(data)
			if err != nil {
				klog.Fatalf("JSON marshal data failed: %s", err)
				panic(err)
			}
			fmt.Fprint(w, string(byteArray))
			return
		}

		if r.Method == "DELETE" {
			klog.Info("Delete token")
			tokenSplit := strings.Split(r.Header.Get("Authorization"), " ")
			tokenString := tokenSplit[1]
			decryptToken(tokenString)
			ctx := context.Background()
			err := client.Del(ctx, tokenString).Err()
			if err != nil {
				panic(err)
			}
			klog.Info("Delete token success")
			return
		}
		panic("Method not allowed")

	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	serverMuxA := http.NewServeMux()
	serverMuxA.HandleFunc("/", indexHandler)
	serverMuxA.HandleFunc("/jwt", jwtHandler)

	serverMuxB := http.NewServeMux()
	serverMuxB.HandleFunc("/api/v1/push", reverseProxyHandler(proxy))

	klog.Info("serving metrics at: :9099")
	go http.ListenAndServe(":9099", promhttp.Handler())

	klog.Info("serving server at: :8081")
	go func() {
		err = http.ListenAndServe(":8081", middlewarestd.Handler("", mdlw, serverMuxA))
		if err != nil {
			panic(err)
		}
	}()

	klog.Info("serving reverse proxy at: :8082")
	err = http.ListenAndServe(":8082", middlewarestd.Handler("", mdlw, serverMuxB))
	if err != nil {
		panic(err)
	}
}
