package session

import (
	"sync"
)

var sessionMap sync.Map

type sessOpts func(s *Server)

type Server struct {
	// 客户端连接
	locker   sync.RWMutex
	sessions map[int64]*Session
}

func LoadServerStore(opts ...sessOpts) (s *Server) {
	var isUp bool
	var key = "user_session"
	if val, ok := sessionMap.Load(key); ok {
		s = val.(*Server)
	} else {
		s = &Server{sessions: make(map[int64]*Session)}
		isUp = true
	}

	for _, opt := range opts {
		opt(s)
		isUp = true
	}

	if isUp {
		sessionMap.Store(key, s)
	}
	return
}

func AddSession(uid int64, session *Session) sessOpts {
	return func(s *Server) {
		s.locker.RLock()
		defer s.locker.RUnlock()

		oldSession, ok := s.sessions[uid]

		if ok {
			oldSession.Quit()
		} else {
			s.sessions[uid] = session
		}
	}
}

func RemoveSession(uid int64) sessOpts {
	return func(s *Server) {
		s.locker.Lock()
		defer s.locker.Unlock()
		if oldSession, ok := s.sessions[uid]; ok {
			delete(s.sessions, uid)
			oldSession.Quit()
		}
	}
}
