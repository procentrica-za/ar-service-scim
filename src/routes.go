package main

func (s *Server) routes() {
	s.router.HandleFunc("/verifycred", s.verifycredentials()).Methods("POST")
}
