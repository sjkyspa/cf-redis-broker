package redisinstance

import "net/http"

type InstanceFinder interface {
	IDForHost(string) string
}

func NewHandler(instanceFinder InstanceFinder) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		instanceID := instanceFinder.IDForHost(req.URL.Query()["host"][0])

		res.Write([]byte(`{"instance_id":"` + instanceID + `"}`))
	}
}

// func NewHandler(repo *redis.RemoteRepository) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		res.Header().Add("Content-Type", "application/json")

// 		debugInfoBytes, err := buildDebugInfoBytes(repo)

// 		if err != nil {
// 			res.Write([]byte(http.StatusText(http.StatusInternalServerError)))
// 			res.WriteHeader(http.StatusInternalServerError)
// 		} else {
// 			res.Write(debugInfoBytes)
// 		}
// 	}
// }
