package main

func (s *Server) routes() {
	s.router.HandleFunc("/scim", s.verifycredentials()).Methods("POST")
}
