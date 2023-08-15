package api

import (
        "net/http"
)

var (
        sessions = map[string]int{
                "123": 1,
                "543": 2,
                "675": 3,
        }
)

func authorized(r *http.Request) bool {
        sessionKey := r.Header.Get("X-Session-Key")
        if sessionKey == "" {
                return false
        }
        
        userId := sessions[sessionKey]
        if userId == 0 {
                return false // 0 means key is not in map
        } else {
                return true
        }
}

func getUserId(r *http.Request) int {
        sessionKey := r.Header.Get("X-Session-Key")
        if sessionKey == "" {
                panic("there is no session key: authorization check must be before rest of logic")
        }
        // Mock. NOT PRODUCTION CODE
        userId := sessions[sessionKey]
        if userId == 0 {
                panic("there is no such user: authorization check must be before rest of logic")
        } else {
                return userId
        }
}
