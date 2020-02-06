// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package testing

import (
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf

func (t *SignedArray) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{129}); err != nil {
		return err
	}

	// t.Signed ([]uint64) (slice)
	if len(t.Signed) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Signed was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.Signed)))); err != nil {
		return err
	}
	for _, v := range t.Signed {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, v); err != nil {
			return err
		}
	}
	return nil
}

func (t *SignedArray) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Signed ([]uint64) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Signed: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.Signed = make([]uint64, extra)
	}
	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeader(br)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Signed slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Signed was not a uint, instead got %d", maj)
		}

		t.Signed[i] = val
	}

	return nil
}

func (t *SimpleTypeOne) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{132}); err != nil {
		return err
	}

	// t.Foo (string) (string)
	if len(t.Foo) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Foo was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Foo)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Foo)); err != nil {
		return err
	}

	// t.Value (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Value))); err != nil {
		return err
	}

	// t.Binary ([]uint8) (slice)
	if len(t.Binary) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Binary was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(t.Binary)))); err != nil {
		return err
	}
	if _, err := w.Write(t.Binary); err != nil {
		return err
	}

	// t.Signed (int64) (int64)
	if t.Signed >= 0 {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Signed))); err != nil {
			return err
		}
	} else {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(-t.Signed)-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *SimpleTypeOne) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 4 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Foo (string) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Foo = string(sval)
	}
	// t.Value (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.Value = uint64(extra)
	// t.Binary ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Binary: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}
	t.Binary = make([]byte, extra)
	if _, err := io.ReadFull(br, t.Binary); err != nil {
		return err
	}
	// t.Signed (int64) (int64)
	{
		maj, extra, err := cbg.CborReadHeader(br)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Signed = int64(extraI)
	}
	return nil
}

func (t *SimpleTypeTwo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{133}); err != nil {
		return err
	}

	// t.Stuff (testing.SimpleTypeTwo) (struct)
	if err := t.Stuff.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Others ([]uint64) (slice)
	if len(t.Others) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Others was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.Others)))); err != nil {
		return err
	}
	for _, v := range t.Others {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, v); err != nil {
			return err
		}
	}

	// t.SignedOthers ([]int64) (slice)
	if len(t.SignedOthers) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.SignedOthers was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.SignedOthers)))); err != nil {
		return err
	}
	for _, v := range t.SignedOthers {
		if v >= 0 {
			if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(v))); err != nil {
				return err
			}
		} else {
			if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(-v)-1)); err != nil {
				return err
			}
		}
	}

	// t.Test ([][]uint8) (slice)
	if len(t.Test) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Test was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.Test)))); err != nil {
		return err
	}
	for _, v := range t.Test {
		if len(v) > cbg.ByteArrayMaxLen {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(v)))); err != nil {
			return err
		}
		if _, err := w.Write(v); err != nil {
			return err
		}
	}

	// t.Dog (string) (string)
	if len(t.Dog) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Dog was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Dog)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Dog)); err != nil {
		return err
	}
	return nil
}

func (t *SimpleTypeTwo) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 5 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Stuff (testing.SimpleTypeTwo) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Stuff = new(SimpleTypeTwo)
			if err := t.Stuff.UnmarshalCBOR(br); err != nil {
				return err
			}
		}

	}
	// t.Others ([]uint64) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Others: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.Others = make([]uint64, extra)
	}
	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeader(br)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Others slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Others was not a uint, instead got %d", maj)
		}

		t.Others[i] = val
	}

	// t.SignedOthers ([]int64) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.SignedOthers: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.SignedOthers = make([]int64, extra)
	}
	for i := 0; i < int(extra); i++ {
		{
			maj, extra, err := cbg.CborReadHeader(br)
			var extraI int64
			if err != nil {
				return err
			}
			switch maj {
			case cbg.MajUnsignedInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 positive overflow")
				}
			case cbg.MajNegativeInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 negative oveflow")
				}
				extraI = -1 - extraI
			default:
				return fmt.Errorf("wrong type for int64 field: %d", maj)
			}

			t.SignedOthers[i] = int64(extraI)
		}
	}

	// t.Test ([][]uint8) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Test: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.Test = make([][]uint8, extra)
	}
	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error

			maj, extra, err = cbg.CborReadHeader(br)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Test[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}
			t.Test[i] = make([]byte, extra)
			if _, err := io.ReadFull(br, t.Test[i]); err != nil {
				return err
			}
		}
	}

	// t.Dog (string) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Dog = string(sval)
	}
	return nil
}