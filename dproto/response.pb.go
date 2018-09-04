// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dproto/response.proto

package dproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

//
// We will allways send same packet type, but it will differs by PacketType field.
// We will improve this when protobuf.Any will be supported on more platforms/languages.
type Response struct {
	Type          PacketType            `protobuf:"varint,1,opt,name=Type,proto3,enum=dproto.PacketType" json:"Type,omitempty"`
	Platform      *Platform             `protobuf:"bytes,2,opt,name=platform" json:"platform,omitempty"`
	Interfaces    map[string]*Interface `protobuf:"bytes,3,rep,name=Interfaces" json:"Interfaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	LldpNeighbors []*LldpNeighbor       `protobuf:"bytes,4,rep,name=LldpNeighbors" json:"LldpNeighbors,omitempty"`
	Vlans         []*Vlan               `protobuf:"bytes,5,rep,name=Vlans" json:"Vlans,omitempty"`
	Ipifs         []*Ipif               `protobuf:"bytes,6,rep,name=Ipifs" json:"Ipifs,omitempty"`
	Uplink        string                `protobuf:"bytes,7,opt,name=Uplink,proto3" json:"Uplink,omitempty"`
	Config        string                `protobuf:"bytes,8,opt,name=Config,proto3" json:"Config,omitempty"`
	// We need this for bulk responses like 'all'
	Errors               map[string]string `protobuf:"bytes,99,rep,name=Errors" json:"Errors,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Error                string            `protobuf:"bytes,100,opt,name=Error,proto3" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_4b9102458fcbcab3, []int{0}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Response.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return m.Size()
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetType() PacketType {
	if m != nil {
		return m.Type
	}
	return PacketType_ALL
}

func (m *Response) GetPlatform() *Platform {
	if m != nil {
		return m.Platform
	}
	return nil
}

func (m *Response) GetInterfaces() map[string]*Interface {
	if m != nil {
		return m.Interfaces
	}
	return nil
}

func (m *Response) GetLldpNeighbors() []*LldpNeighbor {
	if m != nil {
		return m.LldpNeighbors
	}
	return nil
}

func (m *Response) GetVlans() []*Vlan {
	if m != nil {
		return m.Vlans
	}
	return nil
}

func (m *Response) GetIpifs() []*Ipif {
	if m != nil {
		return m.Ipifs
	}
	return nil
}

func (m *Response) GetUplink() string {
	if m != nil {
		return m.Uplink
	}
	return ""
}

func (m *Response) GetConfig() string {
	if m != nil {
		return m.Config
	}
	return ""
}

func (m *Response) GetErrors() map[string]string {
	if m != nil {
		return m.Errors
	}
	return nil
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Response)(nil), "dproto.Response")
	proto.RegisterMapType((map[string]string)(nil), "dproto.Response.ErrorsEntry")
	proto.RegisterMapType((map[string]*Interface)(nil), "dproto.Response.InterfacesEntry")
}
func (m *Response) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Response) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Type != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintResponse(dAtA, i, uint64(m.Type))
	}
	if m.Platform != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintResponse(dAtA, i, uint64(m.Platform.Size()))
		n1, err := m.Platform.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.Interfaces) > 0 {
		for k, _ := range m.Interfaces {
			dAtA[i] = 0x1a
			i++
			v := m.Interfaces[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovResponse(uint64(msgSize))
			}
			mapSize := 1 + len(k) + sovResponse(uint64(len(k))) + msgSize
			i = encodeVarintResponse(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintResponse(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintResponse(dAtA, i, uint64(v.Size()))
				n2, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n2
			}
		}
	}
	if len(m.LldpNeighbors) > 0 {
		for _, msg := range m.LldpNeighbors {
			dAtA[i] = 0x22
			i++
			i = encodeVarintResponse(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Vlans) > 0 {
		for _, msg := range m.Vlans {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintResponse(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Ipifs) > 0 {
		for _, msg := range m.Ipifs {
			dAtA[i] = 0x32
			i++
			i = encodeVarintResponse(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Uplink) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintResponse(dAtA, i, uint64(len(m.Uplink)))
		i += copy(dAtA[i:], m.Uplink)
	}
	if len(m.Config) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintResponse(dAtA, i, uint64(len(m.Config)))
		i += copy(dAtA[i:], m.Config)
	}
	if len(m.Errors) > 0 {
		for k, _ := range m.Errors {
			dAtA[i] = 0x9a
			i++
			dAtA[i] = 0x6
			i++
			v := m.Errors[k]
			mapSize := 1 + len(k) + sovResponse(uint64(len(k))) + 1 + len(v) + sovResponse(uint64(len(v)))
			i = encodeVarintResponse(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintResponse(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintResponse(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.Error) > 0 {
		dAtA[i] = 0xa2
		i++
		dAtA[i] = 0x6
		i++
		i = encodeVarintResponse(dAtA, i, uint64(len(m.Error)))
		i += copy(dAtA[i:], m.Error)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintResponse(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Response) Size() (n int) {
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovResponse(uint64(m.Type))
	}
	if m.Platform != nil {
		l = m.Platform.Size()
		n += 1 + l + sovResponse(uint64(l))
	}
	if len(m.Interfaces) > 0 {
		for k, v := range m.Interfaces {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovResponse(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovResponse(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovResponse(uint64(mapEntrySize))
		}
	}
	if len(m.LldpNeighbors) > 0 {
		for _, e := range m.LldpNeighbors {
			l = e.Size()
			n += 1 + l + sovResponse(uint64(l))
		}
	}
	if len(m.Vlans) > 0 {
		for _, e := range m.Vlans {
			l = e.Size()
			n += 1 + l + sovResponse(uint64(l))
		}
	}
	if len(m.Ipifs) > 0 {
		for _, e := range m.Ipifs {
			l = e.Size()
			n += 1 + l + sovResponse(uint64(l))
		}
	}
	l = len(m.Uplink)
	if l > 0 {
		n += 1 + l + sovResponse(uint64(l))
	}
	l = len(m.Config)
	if l > 0 {
		n += 1 + l + sovResponse(uint64(l))
	}
	if len(m.Errors) > 0 {
		for k, v := range m.Errors {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovResponse(uint64(len(k))) + 1 + len(v) + sovResponse(uint64(len(v)))
			n += mapEntrySize + 2 + sovResponse(uint64(mapEntrySize))
		}
	}
	l = len(m.Error)
	if l > 0 {
		n += 2 + l + sovResponse(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovResponse(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozResponse(x uint64) (n int) {
	return sovResponse(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Response) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowResponse
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Response: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Response: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= (PacketType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Platform", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Platform == nil {
				m.Platform = &Platform{}
			}
			if err := m.Platform.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interfaces", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Interfaces == nil {
				m.Interfaces = make(map[string]*Interface)
			}
			var mapkey string
			var mapvalue *Interface
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowResponse
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowResponse
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthResponse
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowResponse
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= (int(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthResponse
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthResponse
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &Interface{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipResponse(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthResponse
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Interfaces[mapkey] = mapvalue
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LldpNeighbors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LldpNeighbors = append(m.LldpNeighbors, &LldpNeighbor{})
			if err := m.LldpNeighbors[len(m.LldpNeighbors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vlans", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vlans = append(m.Vlans, &Vlan{})
			if err := m.Vlans[len(m.Vlans)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ipifs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ipifs = append(m.Ipifs, &Ipif{})
			if err := m.Ipifs[len(m.Ipifs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uplink", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uplink = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Config = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 99:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Errors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Errors == nil {
				m.Errors = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowResponse
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowResponse
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthResponse
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowResponse
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthResponse
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipResponse(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthResponse
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Errors[mapkey] = mapvalue
			iNdEx = postIndex
		case 100:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipResponse(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthResponse
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipResponse(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowResponse
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthResponse
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowResponse
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipResponse(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthResponse = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowResponse   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("dproto/response.proto", fileDescriptor_response_4b9102458fcbcab3) }

var fileDescriptor_response_4b9102458fcbcab3 = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x4e, 0xab, 0x40,
	0x18, 0xc5, 0xef, 0xb4, 0xc0, 0x85, 0xe9, 0xfd, 0xc3, 0x9d, 0x70, 0xcd, 0x84, 0x18, 0x42, 0xba,
	0x50, 0x16, 0x06, 0x93, 0xea, 0x42, 0xbb, 0x32, 0x9a, 0x2e, 0x6a, 0x8c, 0x69, 0x26, 0xea, 0x9e,
	0xd2, 0xa1, 0x92, 0x22, 0x4c, 0x06, 0x6a, 0xc2, 0x9b, 0xf8, 0x22, 0xbe, 0x83, 0x4b, 0x1f, 0xc1,
	0xd4, 0x17, 0x31, 0xcc, 0x40, 0x4b, 0xd5, 0xdd, 0x7c, 0xe7, 0x77, 0xce, 0x19, 0x3e, 0x06, 0xfe,
	0x9f, 0x31, 0x9e, 0x15, 0xd9, 0x21, 0xa7, 0x39, 0xcb, 0xd2, 0x9c, 0xfa, 0x62, 0x44, 0x9a, 0x94,
	0x6d, 0x54, 0xe3, 0xa2, 0x64, 0x34, 0x97, 0xcc, 0xb6, 0x36, 0x91, 0x65, 0x52, 0xd4, 0x6a, 0xff,
	0x59, 0x81, 0x3a, 0xa9, 0x4b, 0xd0, 0x1e, 0x54, 0x6e, 0x4a, 0x46, 0x31, 0x70, 0x81, 0xf7, 0x67,
	0x80, 0x7c, 0x99, 0xf0, 0x27, 0x41, 0xb8, 0xa0, 0x45, 0x45, 0x88, 0xe0, 0xe8, 0x00, 0xea, 0x2c,
	0x09, 0x8a, 0x28, 0xe3, 0x0f, 0xb8, 0xe3, 0x02, 0xaf, 0x37, 0x30, 0xd7, 0xde, 0x5a, 0x27, 0x6b,
	0x07, 0x3a, 0x83, 0x70, 0x9c, 0x16, 0x94, 0x47, 0x41, 0x48, 0x73, 0xdc, 0x75, 0xbb, 0x5e, 0x6f,
	0xe0, 0x36, 0xfe, 0xe6, 0x6e, 0x7f, 0x63, 0x19, 0xa5, 0x05, 0x2f, 0x49, 0x2b, 0x83, 0x86, 0xf0,
	0xf7, 0x55, 0x32, 0x63, 0xd7, 0x34, 0x9e, 0xdf, 0x4f, 0x33, 0x9e, 0x63, 0x45, 0x94, 0x58, 0x4d,
	0x49, 0x1b, 0x92, 0x6d, 0x2b, 0xea, 0x43, 0xf5, 0x2e, 0x09, 0xd2, 0x1c, 0xab, 0x22, 0xf3, 0xab,
	0xc9, 0x54, 0x22, 0x91, 0xa8, 0xf2, 0x8c, 0x59, 0x1c, 0xe5, 0x58, 0xdb, 0xf6, 0x54, 0x22, 0x91,
	0x08, 0xed, 0x40, 0xed, 0x96, 0x25, 0x71, 0xba, 0xc0, 0x3f, 0x5d, 0xe0, 0x19, 0xa4, 0x9e, 0x2a,
	0xfd, 0x22, 0x4b, 0xa3, 0x78, 0x8e, 0x75, 0xa9, 0xcb, 0x09, 0x1d, 0x43, 0x6d, 0xc4, 0x79, 0xf5,
	0xb1, 0xa1, 0x28, 0xdd, 0xfd, 0xb2, 0xb1, 0xc4, 0x72, 0xdb, 0xda, 0x8b, 0x2c, 0xa8, 0x8a, 0x13,
	0x9e, 0x89, 0x32, 0x39, 0xd8, 0x13, 0xf8, 0xf7, 0xd3, 0xef, 0x41, 0x26, 0xec, 0x2e, 0x68, 0x29,
	0x5e, 0xca, 0x20, 0xd5, 0x11, 0xed, 0x43, 0xf5, 0x31, 0x48, 0x96, 0xb4, 0x7e, 0x91, 0x7f, 0xeb,
	0x25, 0x9a, 0x24, 0x91, 0x7c, 0xd8, 0x39, 0x01, 0xf6, 0x29, 0xec, 0xb5, 0xae, 0xff, 0xa6, 0xcd,
	0x6a, 0xb7, 0x19, 0xad, 0xe8, 0xa5, 0xa2, 0x1b, 0x66, 0x78, 0x6e, 0xbe, 0xac, 0x1c, 0xf0, 0xba,
	0x72, 0xc0, 0xdb, 0xca, 0x01, 0x4f, 0xef, 0xce, 0x8f, 0xa9, 0x26, 0xae, 0x3b, 0xfa, 0x08, 0x00,
	0x00, 0xff, 0xff, 0x7d, 0x71, 0xe4, 0xe9, 0x9b, 0x02, 0x00, 0x00,
}