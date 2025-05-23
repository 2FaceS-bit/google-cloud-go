/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	pb "cloud.google.com/go/spanner/testdata/protos"
	"github.com/google/uuid"
	proto3 "google.golang.org/protobuf/types/known/structpb"
)

type customKeyToString string

func (k customKeyToString) EncodeSpanner() (interface{}, error) {
	return string(k), nil
}

type customKeyToInt int

func (k customKeyToInt) EncodeSpanner() (interface{}, error) {
	return int(k), nil
}

type customKeyToError struct{}

func (k customKeyToError) EncodeSpanner() (interface{}, error) {
	return nil, errors.New("always error")
}

// Test Key.String() and Key.proto().
func TestKey(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339Nano, "2016-11-15T15:04:05.999999999Z")
	dt, _ := civil.ParseDate("2016-11-15")
	for _, test := range []struct {
		k         Key
		wantProto *proto3.ListValue
		wantStr   string
	}{
		{
			k:         Key{int(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{int8(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{int16(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{int32(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{int64(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{uint8(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{uint16(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{uint32(1)},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{true},
			wantProto: listValueProto(boolProto(true)),
			wantStr:   "(true)",
		},
		{
			k:         Key{float32(1.5)},
			wantProto: listValueProto(floatProto(1.5)),
			wantStr:   "(1.5)",
		},
		{
			k:         Key{float64(1.5)},
			wantProto: listValueProto(floatProto(1.5)),
			wantStr:   "(1.5)",
		},
		{
			k:         Key{"value"},
			wantProto: listValueProto(stringProto("value")),
			wantStr:   `("value")`,
		},
		{
			k:         Key{[]byte(nil)},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{[]byte{}},
			wantProto: listValueProto(stringProto("")),
			wantStr:   `("")`,
		},
		{
			k:         Key{tm},
			wantProto: listValueProto(stringProto("2016-11-15T15:04:05.999999999Z")),
			wantStr:   `("2016-11-15T15:04:05.999999999Z")`,
		},
		{k: Key{dt},
			wantProto: listValueProto(stringProto("2016-11-15")),
			wantStr:   `("2016-11-15")`,
		},
		{
			k:         Key{*big.NewRat(1, 1)},
			wantProto: listValueProto(stringProto("1.000000000")),
			wantStr:   `(1.000000000)`,
		},
		{
			k:         Key{uuid1},
			wantProto: listValueProto(uuidProto(uuid1)),
			wantStr:   fmt.Sprintf("(%s)", uuid1.String()),
		},
		{
			k:         Key{uuid.NullUUID{UUID: uuid1, Valid: true}},
			wantProto: listValueProto(uuidProto(uuid1)),
			wantStr:   fmt.Sprintf("(%s)", uuid1.String()),
		},
		{
			k:         Key{uuid.NullUUID{UUID: uuid1, Valid: false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   `(<null>)`,
		},
		{
			k:         Key{NullUUID{uuid1, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   `(<null>)`,
		},
		{
			k:         Key{NullUUID{uuid1, true}},
			wantProto: listValueProto(uuidProto(uuid1)),
			wantStr:   fmt.Sprintf("(%s)", uuid1.String()),
		},
		{
			k:         Key{NullUUID{uuid1, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   `(<null>)`,
		},
		{
			k:         Key{[]byte("value")},
			wantProto: listValueProto(bytesProto([]byte("value"))),
			wantStr:   `("value")`,
		},
		{
			k:         Key{NullInt64{1, true}},
			wantProto: listValueProto(stringProto("1")),
			wantStr:   "(1)",
		},
		{
			k:         Key{NullInt64{2, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullFloat64{1.5, true}},
			wantProto: listValueProto(floatProto(1.5)),
			wantStr:   "(1.5)",
		},
		{
			k:         Key{NullFloat64{2.0, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullFloat32{3.14, true}},
			wantProto: listValueProto(floatProto(float64(float32(3.14)))),
			wantStr:   "(3.14)",
		},
		{
			k:         Key{NullFloat32{2.0, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullBool{true, true}},
			wantProto: listValueProto(boolProto(true)),
			wantStr:   "(true)",
		},
		{
			k:         Key{NullBool{true, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullString{"value", true}},
			wantProto: listValueProto(stringProto("value")),
			wantStr:   `("value")`,
		},
		{
			k:         Key{NullString{"value", false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullTime{tm, true}},
			wantProto: listValueProto(timeProto(tm)),
			wantStr:   `("2016-11-15T15:04:05.999999999Z")`,
		},

		{
			k:         Key{NullTime{time.Now(), false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{NullDate{dt, true}},
			wantProto: listValueProto(dateProto(dt)),
			wantStr:   `("2016-11-15")`,
		},
		{
			k:         Key{NullDate{civil.Date{}, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{int(1), NullString{"value", false}, "value", 1.5, true},
			wantProto: listValueProto(stringProto("1"), nullProto(), stringProto("value"), floatProto(1.5), boolProto(true)),
			wantStr:   `(1,<null>,"value",1.5,true)`,
		},
		{
			k:         Key{NullNumeric{*big.NewRat(2, 3), true}},
			wantProto: listValueProto(stringProto("0.666666667")),
			wantStr:   "(0.666666667)",
		},
		{
			k:         Key{NullNumeric{big.Rat{}, false}},
			wantProto: listValueProto(nullProto()),
			wantStr:   "(<null>)",
		},
		{
			k:         Key{customKeyToString("value")},
			wantProto: listValueProto(stringProto("value")),
			wantStr:   `("value")`,
		},
		{
			k:         Key{customKeyToInt(1)},
			wantProto: listValueProto(intProto(1)),
			wantStr:   `(1)`,
		},
		{
			k:         Key{customKeyToError{}},
			wantProto: nil,
			wantStr:   `(error)`,
		},
		{
			k:         Key{pb.Genre_ROCK},
			wantProto: listValueProto(stringProto("3")),
			wantStr:   "(ROCK)",
		},
		{
			k:         Key{NullProtoEnum{pb.Genre_FOLK, true}},
			wantProto: listValueProto(stringProto("2")),
			wantStr:   "(FOLK)",
		},
	} {
		if got := test.k.String(); got != test.wantStr {
			t.Errorf("%v.String() = %v, want %v", test.k, got, test.wantStr)
		}
		gotProto, err := test.k.proto()
		if test.wantProto != nil && err != nil {
			t.Errorf("%v.proto() returns error %v; want nil error", test.k, err)
		}
		if !testEqual(gotProto, test.wantProto) {
			t.Errorf("%v.proto() = \n%v\nwant:\n%v", test.k, gotProto, test.wantProto)
		}
	}
}

// Test KeyRange.String() and KeyRange.proto().
func TestKeyRange(t *testing.T) {
	for _, test := range []struct {
		kr        KeyRange
		wantProto *sppb.KeyRange
		wantStr   string
	}{
		{
			kr: KeyRange{Key{"A"}, Key{"D"}, OpenOpen},
			wantProto: &sppb.KeyRange{
				StartKeyType: &sppb.KeyRange_StartOpen{StartOpen: listValueProto(stringProto("A"))},
				EndKeyType:   &sppb.KeyRange_EndOpen{EndOpen: listValueProto(stringProto("D"))},
			},
			wantStr: `(("A"),("D"))`,
		},
		{
			kr: KeyRange{Key{1}, Key{10}, OpenClosed},
			wantProto: &sppb.KeyRange{
				StartKeyType: &sppb.KeyRange_StartOpen{StartOpen: listValueProto(stringProto("1"))},
				EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(stringProto("10"))},
			},
			wantStr: "((1),(10)]",
		},
		{
			kr: KeyRange{Key{1.5, 2.1, 0.2}, Key{1.9, 0.7}, ClosedOpen},
			wantProto: &sppb.KeyRange{
				StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(floatProto(1.5), floatProto(2.1), floatProto(0.2))},
				EndKeyType:   &sppb.KeyRange_EndOpen{EndOpen: listValueProto(floatProto(1.9), floatProto(0.7))},
			},
			wantStr: "[(1.5,2.1,0.2),(1.9,0.7))",
		},
		{
			kr: KeyRange{Key{NullInt64{1, true}}, Key{10}, ClosedClosed},
			wantProto: &sppb.KeyRange{
				StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(stringProto("1"))},
				EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(stringProto("10"))},
			},
			wantStr: "[(1),(10)]",
		},
		{
			kr: KeyRange{Key{customKeyToString("A")}, Key{customKeyToString("D")}, OpenOpen},
			wantProto: &sppb.KeyRange{
				StartKeyType: &sppb.KeyRange_StartOpen{StartOpen: listValueProto(stringProto("A"))},
				EndKeyType:   &sppb.KeyRange_EndOpen{EndOpen: listValueProto(stringProto("D"))},
			},
			wantStr: `(("A"),("D"))`,
		},
	} {
		if got := test.kr.String(); got != test.wantStr {
			t.Errorf("%v.String() = %v, want %v", test.kr, got, test.wantStr)
		}
		gotProto, err := test.kr.proto()
		if err != nil {
			t.Errorf("%v.proto() returns error %v; want nil error", test.kr, err)
		}
		if !testEqual(gotProto, test.wantProto) {
			t.Errorf("%v.proto() = \n%v\nwant:\n%v", test.kr, gotProto.String(), test.wantProto.String())
		}
	}
}

func TestPrefixRange(t *testing.T) {
	got := Key{1}.AsPrefix()
	want := KeyRange{Start: Key{1}, End: Key{1}, Kind: ClosedClosed}
	if !testEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestKeySetFromKeys(t *testing.T) {
	for i, test := range []struct {
		ks        KeySet
		wantProto *sppb.KeySet
	}{
		{
			KeySetFromKeys(),
			&sppb.KeySet{},
		},
		{
			KeySetFromKeys(Key{1}),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(intProto(1)),
				},
			},
		},
		{
			KeySetFromKeys(Key{1}, Key{2}),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(intProto(1)),
					listValueProto(intProto(2)),
				},
			},
		},
		{
			KeySetFromKeys(Key{1, "one"}, Key{2, "two"}),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(intProto(1), stringProto("one")),
					listValueProto(intProto(2), stringProto("two")),
				},
			},
		},
	} {
		gotProto, err := test.ks.keySetProto()
		if err != nil {
			t.Errorf("#%d: %v.proto() returns error %v; want nil error", i, test.ks, err)
		}
		if !testEqual(gotProto, test.wantProto) {
			t.Errorf("#%d: %v.proto() = \n%v\nwant:\n%v", i, test.ks, gotProto.String(), test.wantProto.String())
		}
	}
}

func TestKeySets(t *testing.T) {
	int1 := intProto(1)
	int2 := intProto(2)
	int3 := intProto(3)
	int4 := intProto(4)
	for i, test := range []struct {
		ks        KeySet
		wantProto *sppb.KeySet
	}{
		{
			KeySets(),
			&sppb.KeySet{},
		},
		{
			Key{4},
			&sppb.KeySet{
				Keys: []*proto3.ListValue{listValueProto(int4)},
			},
		},
		{
			AllKeys(),
			&sppb.KeySet{All: true},
		},
		{
			KeySets(Key{1, 2}, Key{3, 4}),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(int1, int2),
					listValueProto(int3, int4),
				},
			},
		},
		{
			KeyRange{Key{1}, Key{2}, ClosedOpen},
			&sppb.KeySet{Ranges: []*sppb.KeyRange{
				{
					StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(int1)},
					EndKeyType:   &sppb.KeyRange_EndOpen{EndOpen: listValueProto(int2)},
				},
			}},
		},
		{
			Key{2}.AsPrefix(),
			&sppb.KeySet{Ranges: []*sppb.KeyRange{
				{
					StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(int2)},
					EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(int2)},
				},
			}},
		},
		{
			KeySets(
				KeyRange{Key{1}, Key{2}, ClosedClosed},
				KeyRange{Key{3}, Key{4}, OpenClosed},
			),
			&sppb.KeySet{
				Ranges: []*sppb.KeyRange{
					{
						StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(int1)},
						EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(int2)},
					},
					{
						StartKeyType: &sppb.KeyRange_StartOpen{StartOpen: listValueProto(int3)},
						EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(int4)},
					},
				},
			},
		},
		{
			KeySets(
				Key{1},
				KeyRange{Key{2}, Key{3}, ClosedClosed},
				KeyRange{Key{4}, Key{5}, OpenClosed},
				KeySets(),
				Key{6}),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(int1),
					listValueProto(intProto(6)),
				},
				Ranges: []*sppb.KeyRange{
					{
						StartKeyType: &sppb.KeyRange_StartClosed{StartClosed: listValueProto(int2)},
						EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(int3)},
					},
					{
						StartKeyType: &sppb.KeyRange_StartOpen{StartOpen: listValueProto(int4)},
						EndKeyType:   &sppb.KeyRange_EndClosed{EndClosed: listValueProto(intProto(5))},
					},
				},
			},
		},
		{
			KeySets(
				Key{1},
				KeyRange{Key{2}, Key{3}, ClosedClosed},
				AllKeys(),
				KeyRange{Key{4}, Key{5}, OpenClosed},
				Key{6}),
			&sppb.KeySet{All: true},
		},
		{
			KeySets(
				Key{customKeyToInt(1), customKeyToInt(2)},
				Key{customKeyToInt(3), customKeyToInt(4)},
			),
			&sppb.KeySet{
				Keys: []*proto3.ListValue{
					listValueProto(int1, int2),
					listValueProto(int3, int4),
				},
			},
		},
	} {
		gotProto, err := test.ks.keySetProto()
		if err != nil {
			t.Errorf("#%d: %v.proto() returns error %v; want nil error", i, test.ks, err)
		}
		if !testEqual(gotProto, test.wantProto) {
			t.Errorf("#%d: %v.proto() = \n%v\nwant:\n%v", i, test.ks, gotProto.String(), test.wantProto.String())
		}
	}
}
