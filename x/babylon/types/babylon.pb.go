// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: babylonchain/babylon/v1beta1/babylon.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the x/babylon module.
type Params struct {
	// MaxGasEndBlocker defines the maximum gas that can be spent in a contract
	// sudo callback
	MaxGasEndBlocker uint32 `protobuf:"varint,3,opt,name=max_gas_end_blocker,json=maxGasEndBlocker,proto3" json:"max_gas_end_blocker,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5add0b76ad5fde9, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "babylonchain.babylon.v1beta1.Params")
}

func init() {
	proto.RegisterFile("babylonchain/babylon/v1beta1/babylon.proto", fileDescriptor_b5add0b76ad5fde9)
}

var fileDescriptor_b5add0b76ad5fde9 = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4a, 0x4a, 0x4c, 0xaa,
	0xcc, 0xc9, 0xcf, 0x4b, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x87, 0x72, 0xf4, 0xcb, 0x0c, 0x93, 0x52,
	0x4b, 0x12, 0x0d, 0x61, 0x7c, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x19, 0x64, 0xb5, 0x7a,
	0x30, 0x39, 0xa8, 0x5a, 0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0, 0x42, 0x7d, 0x10, 0x0b, 0xa2,
	0x47, 0x4a, 0x30, 0x31, 0x37, 0x33, 0x2f, 0x5f, 0x1f, 0x4c, 0x42, 0x84, 0x94, 0x02, 0xb8, 0xd8,
	0x02, 0x12, 0x8b, 0x12, 0x73, 0x8b, 0x85, 0x74, 0xb9, 0x84, 0x73, 0x13, 0x2b, 0xe2, 0xd3, 0x13,
	0x8b, 0xe3, 0x53, 0xf3, 0x52, 0xe2, 0x93, 0x72, 0xf2, 0x93, 0xb3, 0x53, 0x8b, 0x24, 0x98, 0x15,
	0x18, 0x35, 0x78, 0x83, 0x04, 0x72, 0x13, 0x2b, 0xdc, 0x13, 0x8b, 0x5d, 0xf3, 0x52, 0x9c, 0x20,
	0xe2, 0x56, 0xe2, 0x2f, 0x16, 0xc8, 0x33, 0x76, 0x3d, 0xdf, 0xa0, 0xc5, 0x07, 0x73, 0x27, 0xc4,
	0x1c, 0xa7, 0xd0, 0x13, 0x0f, 0xe5, 0x18, 0x56, 0x3c, 0x92, 0x63, 0x38, 0xf1, 0x48, 0x8e, 0xf1,
	0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e,
	0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xe3, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc,
	0x5c, 0x7d, 0x6c, 0x3e, 0xd6, 0x2d, 0x4e, 0xc9, 0xd6, 0xaf, 0x80, 0xfb, 0xbf, 0xa4, 0xb2, 0x20,
	0xb5, 0x38, 0x89, 0x0d, 0xec, 0x5e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfd, 0xa8, 0x76,
	0x0e, 0x24, 0x01, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MaxGasEndBlocker != that1.MaxGasEndBlocker {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxGasEndBlocker != 0 {
		i = encodeVarintBabylon(dAtA, i, uint64(m.MaxGasEndBlocker))
		i--
		dAtA[i] = 0x18
	}
	return len(dAtA) - i, nil
}

func encodeVarintBabylon(dAtA []byte, offset int, v uint64) int {
	offset -= sovBabylon(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxGasEndBlocker != 0 {
		n += 1 + sovBabylon(uint64(m.MaxGasEndBlocker))
	}
	return n
}

func sovBabylon(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBabylon(x uint64) (n int) {
	return sovBabylon(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBabylon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxGasEndBlocker", wireType)
			}
			m.MaxGasEndBlocker = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBabylon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxGasEndBlocker |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBabylon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBabylon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBabylon(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBabylon
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
					return 0, ErrIntOverflowBabylon
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBabylon
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
			if length < 0 {
				return 0, ErrInvalidLengthBabylon
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBabylon
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBabylon
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBabylon        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBabylon          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBabylon = fmt.Errorf("proto: unexpected end of group")
)
