// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: coolcat/catdrop/v1beta1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ClaimAuthorization struct {
	ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty" yaml:"contract_address"`
	Action          Action `protobuf:"varint,2,opt,name=action,proto3,enum=coolcat.catdrop.v1beta1.Action" json:"action,omitempty" yaml:"action"`
}

func (m *ClaimAuthorization) Reset()         { *m = ClaimAuthorization{} }
func (m *ClaimAuthorization) String() string { return proto.CompactTextString(m) }
func (*ClaimAuthorization) ProtoMessage()    {}
func (*ClaimAuthorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_e495f9a954f17123, []int{0}
}
func (m *ClaimAuthorization) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClaimAuthorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClaimAuthorization.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClaimAuthorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClaimAuthorization.Merge(m, src)
}
func (m *ClaimAuthorization) XXX_Size() int {
	return m.Size()
}
func (m *ClaimAuthorization) XXX_DiscardUnknown() {
	xxx_messageInfo_ClaimAuthorization.DiscardUnknown(m)
}

var xxx_messageInfo_ClaimAuthorization proto.InternalMessageInfo

func (m *ClaimAuthorization) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func (m *ClaimAuthorization) GetAction() Action {
	if m != nil {
		return m.Action
	}
	return ActionVote
}

// Params defines the claim module's parameters.
type Params struct {
	AirdropEnabled     bool          `protobuf:"varint,1,opt,name=airdrop_enabled,json=airdropEnabled,proto3" json:"airdrop_enabled,omitempty"`
	AirdropStartTime   time.Time     `protobuf:"bytes,2,opt,name=airdrop_start_time,json=airdropStartTime,proto3,stdtime" json:"airdrop_start_time" yaml:"airdrop_start_time"`
	DurationUntilDecay time.Duration `protobuf:"bytes,3,opt,name=duration_until_decay,json=durationUntilDecay,proto3,stdduration" json:"duration_until_decay,omitempty" yaml:"duration_until_decay"`
	DurationOfDecay    time.Duration `protobuf:"bytes,4,opt,name=duration_of_decay,json=durationOfDecay,proto3,stdduration" json:"duration_of_decay,omitempty" yaml:"duration_of_decay"`
	// denom of claimable asset
	ClaimDenom string `protobuf:"bytes,5,opt,name=claim_denom,json=claimDenom,proto3" json:"claim_denom,omitempty"`
	// list of contracts and their allowed claim actions
	AllowedClaimers []ClaimAuthorization `protobuf:"bytes,6,rep,name=allowed_claimers,json=allowedClaimers,proto3" json:"allowed_claimers" yaml:"allowed_claimers"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_e495f9a954f17123, []int{1}
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

func (m *Params) GetAirdropEnabled() bool {
	if m != nil {
		return m.AirdropEnabled
	}
	return false
}

func (m *Params) GetAirdropStartTime() time.Time {
	if m != nil {
		return m.AirdropStartTime
	}
	return time.Time{}
}

func (m *Params) GetDurationUntilDecay() time.Duration {
	if m != nil {
		return m.DurationUntilDecay
	}
	return 0
}

func (m *Params) GetDurationOfDecay() time.Duration {
	if m != nil {
		return m.DurationOfDecay
	}
	return 0
}

func (m *Params) GetClaimDenom() string {
	if m != nil {
		return m.ClaimDenom
	}
	return ""
}

func (m *Params) GetAllowedClaimers() []ClaimAuthorization {
	if m != nil {
		return m.AllowedClaimers
	}
	return nil
}

func init() {
	proto.RegisterType((*ClaimAuthorization)(nil), "coolcat.catdrop.v1beta1.ClaimAuthorization")
	proto.RegisterType((*Params)(nil), "coolcat.catdrop.v1beta1.Params")
}

func init() {
	proto.RegisterFile("coolcat/catdrop/v1beta1/params.proto", fileDescriptor_e495f9a954f17123)
}

var fileDescriptor_e495f9a954f17123 = []byte{
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x4d, 0x8b, 0xd3, 0x5e,
	0x14, 0xc6, 0x9b, 0xff, 0xcc, 0xbf, 0xe8, 0x2d, 0x4e, 0x3b, 0x61, 0x60, 0x32, 0x2d, 0x24, 0x25,
	0x28, 0x16, 0x95, 0x84, 0x76, 0x36, 0x32, 0xbb, 0xa6, 0x55, 0xf0, 0x05, 0x94, 0xa8, 0x1b, 0x37,
	0xe1, 0x26, 0xb9, 0x4d, 0x03, 0x49, 0x6e, 0xb8, 0xb9, 0x1d, 0xad, 0x1f, 0x41, 0x37, 0xb3, 0x92,
	0x59, 0xba, 0xf4, 0xa3, 0xcc, 0x72, 0x96, 0xae, 0xa2, 0xb4, 0x3b, 0x97, 0xdd, 0x0b, 0x72, 0x5f,
	0x32, 0x62, 0x3b, 0xc5, 0x5d, 0xfb, 0x9c, 0xe7, 0x9c, 0xf3, 0xe3, 0xc9, 0xb9, 0xe0, 0x76, 0x80,
	0x71, 0x12, 0x40, 0x6a, 0x07, 0x90, 0x86, 0x04, 0xe7, 0xf6, 0x69, 0xdf, 0x47, 0x14, 0xf6, 0xed,
	0x1c, 0x12, 0x98, 0x16, 0x56, 0x4e, 0x30, 0xc5, 0xea, 0xa1, 0x74, 0x59, 0xd2, 0x65, 0x49, 0x57,
	0xfb, 0x20, 0xc2, 0x11, 0xe6, 0x1e, 0x9b, 0xfd, 0x12, 0xf6, 0xb6, 0x1e, 0x61, 0x1c, 0x25, 0xc8,
	0xe6, 0xff, 0xfc, 0xd9, 0xc4, 0x0e, 0x67, 0x04, 0xd2, 0x18, 0x67, 0xb2, 0x6e, 0xac, 0xd7, 0x69,
	0x9c, 0xa2, 0x82, 0xc2, 0x34, 0x97, 0x86, 0x7b, 0xdb, 0xa8, 0x82, 0x04, 0xc6, 0xa9, 0x47, 0x50,
	0x80, 0x49, 0x28, 0xbc, 0xe6, 0x57, 0x05, 0xa8, 0x23, 0x26, 0x0f, 0x67, 0x74, 0x8a, 0x49, 0xfc,
	0x81, 0x6f, 0x52, 0x1f, 0x83, 0x56, 0x80, 0x33, 0x4a, 0x60, 0x40, 0x3d, 0x18, 0x86, 0x04, 0x15,
	0x85, 0xa6, 0x74, 0x95, 0xde, 0x4d, 0xa7, 0xb3, 0x2a, 0x8d, 0xc3, 0x39, 0x4c, 0x93, 0x13, 0x73,
	0xdd, 0x61, 0xba, 0xcd, 0x4a, 0x1a, 0x0a, 0x45, 0x7d, 0x0a, 0xea, 0x30, 0x60, 0x13, 0xb5, 0xff,
	0xba, 0x4a, 0x6f, 0x6f, 0x60, 0x58, 0x5b, 0xb2, 0xb0, 0x86, 0xdc, 0xe6, 0xec, 0xaf, 0x4a, 0xe3,
	0x96, 0x18, 0x2f, 0x1a, 0x4d, 0x57, 0x4e, 0x30, 0x7f, 0xed, 0x82, 0xfa, 0x4b, 0x9e, 0xab, 0x7a,
	0x17, 0x34, 0x61, 0x4c, 0x58, 0xbf, 0x87, 0x32, 0xe8, 0x27, 0x28, 0xe4, 0x74, 0x37, 0xdc, 0x3d,
	0x29, 0x3f, 0x12, 0xaa, 0x8a, 0x81, 0x5a, 0x19, 0x0b, 0x0a, 0x09, 0xf5, 0x58, 0x56, 0x9c, 0xa5,
	0x31, 0x68, 0x5b, 0x22, 0x48, 0xab, 0x0a, 0xd2, 0x7a, 0x5d, 0x05, 0xe9, 0xdc, 0xb9, 0x28, 0x8d,
	0xda, 0xaa, 0x34, 0x8e, 0x24, 0xca, 0xc6, 0x0c, 0xf3, 0xec, 0xbb, 0xa1, 0xb8, 0x2d, 0x59, 0x78,
	0xc5, 0x74, 0xd6, 0xad, 0x7e, 0x56, 0xc0, 0x41, 0xf5, 0xbd, 0xbc, 0x59, 0x46, 0xe3, 0xc4, 0x0b,
	0x51, 0x00, 0xe7, 0xda, 0x0e, 0xdf, 0x79, 0xb4, 0xb1, 0x73, 0x2c, 0xcd, 0xce, 0x13, 0xb6, 0xf2,
	0x67, 0x69, 0xe8, 0xd7, 0xb5, 0x3f, 0xc0, 0x69, 0x4c, 0x51, 0x9a, 0xd3, 0xf9, 0xaa, 0x34, 0x3a,
	0x02, 0xea, 0x3a, 0x9f, 0x79, 0xce, 0xb0, 0xd4, 0xaa, 0xf4, 0x86, 0x55, 0xc6, 0xac, 0xa0, 0x7e,
	0x54, 0xc0, 0xfe, 0x55, 0x07, 0x9e, 0x48, 0xaa, 0xdd, 0x7f, 0x51, 0x8d, 0x24, 0x55, 0x67, 0xa3,
	0xf7, 0x2f, 0x24, 0x6d, 0x0d, 0xa9, 0x32, 0x09, 0x9e, 0x66, 0xa5, 0xbf, 0x98, 0x08, 0x18, 0x03,
	0x34, 0xc4, 0x2d, 0x86, 0x28, 0xc3, 0xa9, 0xf6, 0x3f, 0xbb, 0x2c, 0x17, 0x70, 0x69, 0xcc, 0x14,
	0xf5, 0x93, 0x02, 0x5a, 0x30, 0x49, 0xf0, 0x3b, 0x14, 0x7a, 0x5c, 0x46, 0xa4, 0xd0, 0xea, 0xdd,
	0x9d, 0x5e, 0x63, 0x70, 0x7f, 0xeb, 0x09, 0x6d, 0xde, 0xb1, 0x73, 0x2c, 0xf1, 0x37, 0x86, 0xfd,
	0xb9, 0xe2, 0xf5, 0x8a, 0xe9, 0x36, 0xa5, 0x34, 0x92, 0xca, 0xc9, 0xee, 0xf9, 0x17, 0xa3, 0xe6,
	0xb8, 0x17, 0x0b, 0x5d, 0xb9, 0x5c, 0xe8, 0xca, 0x8f, 0x85, 0xae, 0x9c, 0x2d, 0xf5, 0xda, 0xe5,
	0x52, 0xaf, 0x7d, 0x5b, 0xea, 0xb5, 0xb7, 0x0f, 0xa3, 0x98, 0x4e, 0x67, 0xbe, 0x15, 0xe0, 0xd4,
	0x1e, 0xc7, 0x51, 0x4c, 0x61, 0xf2, 0x2c, 0xa6, 0xc1, 0x14, 0x65, 0xcf, 0xa1, 0x5f, 0xd8, 0xd5,
	0x73, 0x3c, 0xed, 0xdb, 0xef, 0xaf, 0xde, 0x24, 0x9d, 0xe7, 0xa8, 0xf0, 0xeb, 0x3c, 0xf1, 0xe3,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x24, 0x47, 0xeb, 0xff, 0x49, 0x04, 0x00, 0x00,
}

func (m *ClaimAuthorization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClaimAuthorization) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClaimAuthorization) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Action != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Action))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
	if len(m.AllowedClaimers) > 0 {
		for iNdEx := len(m.AllowedClaimers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AllowedClaimers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.ClaimDenom) > 0 {
		i -= len(m.ClaimDenom)
		copy(dAtA[i:], m.ClaimDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ClaimDenom)))
		i--
		dAtA[i] = 0x2a
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.DurationOfDecay, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.DurationOfDecay):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.DurationUntilDecay, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.DurationUntilDecay):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	n3, err3 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.AirdropStartTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.AirdropStartTime):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintParams(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x12
	if m.AirdropEnabled {
		i--
		if m.AirdropEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClaimAuthorization) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.Action != 0 {
		n += 1 + sovParams(uint64(m.Action))
	}
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AirdropEnabled {
		n += 2
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.AirdropStartTime)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.DurationUntilDecay)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.DurationOfDecay)
	n += 1 + l + sovParams(uint64(l))
	l = len(m.ClaimDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if len(m.AllowedClaimers) > 0 {
		for _, e := range m.AllowedClaimers {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClaimAuthorization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: ClaimAuthorization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClaimAuthorization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Action", wireType)
			}
			m.Action = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Action |= Action(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AirdropEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AirdropEnabled = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AirdropStartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.AirdropStartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DurationUntilDecay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.DurationUntilDecay, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DurationOfDecay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.DurationOfDecay, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClaimDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowedClaimers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AllowedClaimers = append(m.AllowedClaimers, ClaimAuthorization{})
			if err := m.AllowedClaimers[len(m.AllowedClaimers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
