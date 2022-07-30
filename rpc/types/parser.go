package types

import "strings"

type NewTypedEventParser struct {
	protoPkg   string
	protoMsg   string
	protoField string
}

func (p *NewTypedEventParser) SetProtoPkg(protoPkgName string) {
	p.protoPkg = protoPkgName
}

func (p *NewTypedEventParser) SetProtoMsg(protoMsgName string) {
	p.protoMsg = protoMsgName
}

func (p *NewTypedEventParser) SetProtoField(protoFieldName string) {
	p.protoField = protoFieldName
}

func (p *NewTypedEventParser) Event() string {
	if p.protoField == "" {
		panic("You should add proto field name")
	}
	return p.EventWithFieldName(p.protoField)
}

func (p *NewTypedEventParser) EventWithFieldName(protoFieldName string) string {
	return strings.Join([]string{p.protoPkg, p.protoMsg, protoFieldName}, ".")
}
