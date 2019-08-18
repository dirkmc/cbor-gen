package typegen

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type ByteReader interface {
	io.ByteReader
	io.Reader
}

func CborReadHeader(br ByteReader) (byte, uint64, error) {
	first, err := br.ReadByte()
	if err != nil {
		return 0, 0, err
	}

	maj := (first & 0xe0) >> 5
	low := first & 0x1f

	switch {
	case low < 24:
		return maj, uint64(low), nil
	case low == 24:
		next, err := br.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		return maj, uint64(next), nil
	case low == 25:
		buf := make([]byte, 2)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, uint64(binary.BigEndian.Uint16(buf)), nil
	case low == 26:
		buf := make([]byte, 4)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, uint64(binary.BigEndian.Uint32(buf)), nil
	case low == 27:
		buf := make([]byte, 8)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, binary.BigEndian.Uint64(buf), nil
	default:
		return 0, 0, fmt.Errorf("invalid header: (%x)", first)
	}
}

func CborEncodeMajorType(t byte, l uint64) []byte {
	var b [9]byte
	switch {
	case l < 24:
		b[0] = (t << 5) | byte(l)
		return b[:1]
	case l < (1 << 8):
		b[0] = (t << 4) | 24
		b[1] = byte(l)
		return b[:2]
	case l < (1 << 16):
		b[0] = (t << 4) | 25
		binary.BigEndian.PutUint16(b[1:3], uint16(l))
		return b[:3]
	case l < (1 << 32):
		b[0] = (t << 4) | 26
		binary.BigEndian.PutUint32(b[1:5], uint32(l))
		return b[:5]
	default:
		b[0] = (t << 4) | 27
		binary.BigEndian.PutUint64(b[1:], uint64(l))
		return b[:]
	}
}

func PrintHeaderAndUtilityMethods(w io.Writer, pkg string) error {
	fmt.Fprintf(w, "package %s\n\n", pkg)

	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\t\"encoding/binary\"\n")
	fmt.Fprintf(w, "\t\"fmt\"\n")
	fmt.Fprintf(w, "\t\"io\"\n")
	fmt.Fprintf(w, "\n\tcbg \"github.com/whyrusleeping/cbor-gen\"\n")
	fmt.Fprintf(w, ")\n\n")

	return nil
}

// Generates 'tuple representation' cbor encoders for the given type
func GenTupleEncodersForType(i interface{}, w io.Writer) error {
	t := reflect.TypeOf(i)

	fmt.Fprintf(w, "func (t *%s) MarshalCBOR(w io.Writer) error {\n", t.Name())

	if t.NumField() > 23 {
		return fmt.Errorf("lazy programmer doesnt want to handle field counters > 23")
	}

	firstByte := (4 << 5) | byte(t.NumField())

	fmt.Fprintf(w, "\tw.Write([]byte{0x%x})\n", firstByte)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Fprintf(w, "\t// t.%s (%s)\n", f.Name, f.Type)
		switch f.Type.Kind() {
		case reflect.String:
			return fmt.Errorf("strings arent handled yet")
		case reflect.Struct:
			fmt.Fprintf(w, "\tif err := t.%s.MarshalCBOR(w); err != nil {\n", f.Name)
			fmt.Fprintf(w, "\t\treturn err\n\t}\n\n")
		case reflect.Uint64:
			fmt.Fprintf(w, "\t w.Write(cbg.CborEncodeMajorType(0, t.%s))\n\n", f.Name)
		case reflect.Slice:
			e := f.Type.Elem()
			if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}

			if e.Kind() == reflect.Uint8 {
				fmt.Fprintf(w, "\tw.Write(cbg.CborEncodeMajorType(2, uint64(len(t.%s))))\n", f.Name)
				fmt.Fprintf(w, "\tw.Write(t.%s)\n\n", f.Name)
				continue
			}

			fmt.Fprintf(w, "\tw.Write(cbg.CborEncodeMajorType(4, uint64(len(t.%s))))\n", f.Name)
			fmt.Fprintf(w, "\tfor i, v := range t.%s {\n", f.Name)
			switch e.Kind() {
			case reflect.Struct:
				fmt.Fprintf(w, "\t\tif err := v.MarshalCBOR(w); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n")
			default:
				return fmt.Errorf("do not yet support slices of non-structs: %s %s", f.Type.Elem(), e.Kind())
			}
		default:
			return fmt.Errorf("field %q of %q has unsupported kind %q", f.Name, t.Name(), f.Type.Kind())
		}
	}

	fmt.Fprintf(w, "\treturn nil\n}\n\n")

	// Now for the unmarshal

	fmt.Fprintf(w, "func (t *%s) UnmarshalCBOR(br ByteReader) error {\n", t.Name())

	fmt.Fprintf(w, "\tfirst, err := br.ReadByte()\n")
	fmt.Fprintf(w, "\tif err != nil {\n\t\treturn err\n\t}\n\n")
	fmt.Fprintf(w, "\tif first != 0x%x {\n", firstByte)
	fmt.Fprintf(w, "\t\treturn fmt.Errorf(\"object had incorrect type or length\")\n\t}\n\n")

	fmt.Fprintf(w, "\tvar maj byte\n\tvar extra uint64\n\n")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Fprintf(w, "\t// t.%s (%s)\n", f.Name, f.Type)
		switch f.Type.Kind() {
		case reflect.String:
			return fmt.Errorf("strings arent handled yet")
		case reflect.Struct:
			fmt.Fprintf(w, "\tif err := t.%s.UnmarshalCBOR(br); err != nil {\n", f.Name)
			fmt.Fprintf(w, "\t\treturn err\n\t}\n\n")
		case reflect.Uint64:
			fmt.Fprintf(w, "\tmaj, extra, err = cbg.CborReadHeader(br)\n\tif err != nil {\n\t\treturn err\n\t}\n\n")
			fmt.Fprintf(w, "\tif maj != 0 {\n\t\treturn fmt.Errorf(\"wrong type for uint64 field\")\n\t}\n")
			fmt.Fprintf(w, "\tt.%s = extra\n\n", f.Name)
		case reflect.Slice:
			e := f.Type.Elem()
			if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}

			fmt.Fprintf(w, "\tmaj, extra, err = cbg.CborReadHeader(br)\n\tif err != nil {\n\t\treturn err\n\t}\n\n")
			fmt.Fprintf(w, "\tif extra > 8192 {\n\t\treturn fmt.Errorf(\"array too large\")\n\t}\n")
			if e.Kind() == reflect.Uint8 {
				fmt.Fprintf(w, "\tt.%s = make([]byte, extra)\n", f.Name)
				fmt.Fprintf(w, "\tif _, err := io.ReadFull(br, t.%s); err != nil {\n\t\treturn err\n\t}\n\n", f.Name)
				continue
			}

			fmt.Fprintf(w, "\tt.%s = make([]%s, 0, extra)\n", f.Name, f.Type)
			fmt.Fprintf(w, "\tfor i := 0; i < extra; i++ {\n")
			switch e.Kind() {
			case reflect.Struct:
				fmt.Fprintf(w, "\t\tvar v %s\n", f.Type)
				fmt.Fprintf(w, "\t\tif err := v.UnmarshalCBOR(br); err != nil {\n\t\t\treturn err\n\t\t}\n\n")
				fmt.Fprintf(w, "\t\tt.%s = append(t.%s, v)\n", f.Name, f.Name)
			default:
				return fmt.Errorf("do not yet support slices of non-structs: %s", f.Type.Elem())
			}
			fmt.Fprintf(w, "\t}\n\n")
		default:
			return fmt.Errorf("field %q of %q has unsupported kind %q", f.Name, t.Name(), f.Type.Kind())
		}
	}

	fmt.Fprintf(w, "\treturn nil\n}\n\n")

	return nil
}