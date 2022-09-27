package main

import (
	"errors"
	"regexp"
	"strings"
)

type HeaderProcess struct {
	namespace string
	tables    []string
}

func (p *HeaderProcess) GetNamespace() string {
	if len(*packagename) > 0 {
		return *packagename
	}
	return p.namespace
}

func (p *HeaderProcess) GetTables() []string {
	return p.tables
}

func (p *HeaderProcess) GenExternConf() []*StreamDef {
	return nil
}

func (p *HeaderProcess) Read(content string) error {
	regNs := regexp.MustCompile(currProtoConfig.namespaceRegexp)
	ns := regNs.FindAllStringSubmatch(content, -1)
	if len(ns) < 1 {
		return errors.New("not found namespace")
	}
	//fmt.Println("ns:", ns)

	namespace := strings.TrimSpace(ns[0][1])
	if p.namespace != "" && p.namespace != namespace {
		return errors.New("must same namespace")
	}
	p.namespace = namespace

	regTab := regexp.MustCompile(currProtoConfig.messageRegexp)
	tab := regTab.FindAllStringSubmatch(content, -1)
	//fmt.Println("tab:", tab)
	for _, t := range tab {
		s := strings.TrimSpace(t[1])
		p.tables = append(p.tables, s)
	}
	return nil
}
