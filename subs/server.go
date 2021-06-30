package subs

import (
	//"fmt"
)

type Server struct {

        Cfg *Config
}

func NewServer(config *Config) (*Server, error) {

        if config == nil {
                config = DefaultConfig
        }
        s := &Server{Cfg: config}

        return s, nil

}

func (s* Server) Start() {

        /*if s.cfg.Categories != nil {

               dflt_ctgrs = s.cfg.Categories
        }*/
}
