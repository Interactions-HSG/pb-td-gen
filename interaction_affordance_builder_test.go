package grpcwot

import (
	"errors"
	"github.com/emicklei/proto"
	"github.com/linksmart/thing-directory/wot"
	"testing"
)

var inMessages = []map[string]*wot.DataSchema{
	{
		"TestM1": {},
		"TestM2": {},
	},
	{
		"TestM1":         {},
		"TestM1.TestIM1": {},
		"TestM2":         {},
	},
}

var conformRPCTest = []struct {
	inMessages map[string]*wot.DataSchema
	inRPCs     []*proto.RPC
	out        map[string]affs
	err        error
}{
	{
		inMessages[0],
		[]*proto.RPC{
			{
				Name:        "TestRPC1",
				RequestType: "TestM1",
				ReturnsType: "TestM2",
			},
		},
		map[string]affs{
			"TestRPC1": {
				req: inMessages[0]["TestM1"],
				res: inMessages[0]["TestM2"],
			},
		},
		nil,
	},
	{
		inMessages[1],
		[]*proto.RPC{
			{
				Name:        "TestRPC2",
				RequestType: "TestM1.TestIM1",
				ReturnsType: "TestM2",
			},
		},
		map[string]affs{
			"TestRPC2": {
				req: inMessages[1]["TestM1.TestIM1"],
				res: inMessages[1]["TestM2"],
			},
		},
		nil,
	},
	{
		inMessages[0],
		[]*proto.RPC{
			{
				Name:        "TestRPC3",
				RequestType: "TestM1",
				ReturnsType: "TestM2",
			},
			{
				Name:        "TestRPC4",
				RequestType: "TestM2",
				ReturnsType: "TestM1",
			},
		},
		map[string]affs{
			"TestRPC3": {
				req: inMessages[0]["TestM1"],
				res: inMessages[0]["TestM2"],
			},
			"TestRPC4": {
				req: inMessages[0]["TestM2"],
				res: inMessages[0]["TestM1"],
			},
		},
		nil,
	},
	{
		map[string]*wot.DataSchema{
			"T": {},
		},
		[]*proto.RPC{
			{
				Name:        "TestRPCError1",
				RequestType: "T",
				ReturnsType: "T",
			},
			{
				Name:        "TestRPCError1",
				RequestType: "T",
				ReturnsType: "T",
			},
		},
		map[string]affs{},
		errors.New("Duplicate RPC name found in proto file for RPC name TestRPCError1"),
	},
	{
		map[string]*wot.DataSchema{
			"T": {},
		},
		[]*proto.RPC{
			{
				Name:        "TestRPCError2",
				RequestType: "E",
				ReturnsType: "T",
			},
		},
		map[string]affs{},
		errors.New("Not able to determine message for request type E in RPC TestRPCError2"),
	},
	{
		map[string]*wot.DataSchema{
			"T": {},
		},
		[]*proto.RPC{
			{
				Name:        "TestRPCError3",
				RequestType: "T",
				ReturnsType: "E",
			},
		},
		map[string]affs{},
		errors.New("Not able to determine message for return type E in RPC TestRPCError3"),
	},
}

func TestConformRPC(t *testing.T) {
	for _, v := range conformRPCTest {
		dsb := &dataSchemaBuilder{
			ds: v.inMessages,
		}
		iab := &interactionAffordanceBuilder{
			rpcs: v.inRPCs,
			dsb:  dsb,
		}
		err := iab.conformRPCs()
		if v.err != nil {
			if err == nil {
				t.Errorf("Expected the error %v,\n but nothing was raised", v.err.Error())
			} else if err.Error() != v.err.Error() {
				t.Errorf("Wrong error message:\n Expected: %v\n but actual is: %v\n", v.err.Error(), err.Error())
			}
		} else {
			for k, v := range v.out {
				if iab.affs[k].req != v.req {
					t.Errorf("Expected the request type %v,\n but got %v\n for RPC %v", v.req, iab.affs[k].req, k)
				}
				if iab.affs[k].res != v.res {
					t.Errorf("Expected the return type %v,\n but got %v\n for RPC %v", v.res, iab.affs[k].res, k)
				}
			}
		}
	}
}

var categorizeRPCTestAffordances = map[string]affs{
	"SimpleTest": {
		name: "SimpleTest",
		req:  &wot.DataSchema{},
		res:  &wot.DataSchema{},
	},
	"GetTest": {
		name: "GetTest",
		req:  &wot.DataSchema{},
		res:  &wot.DataSchema{},
	},
	"GetTest2": {
		name: "GetTest2",
		req:  &wot.DataSchema{},
		res:  &wot.DataSchema{},
	},
	"SetTest": {
		name: "SetTest",
		req:  &wot.DataSchema{},
		res:  &wot.DataSchema{},
	},
	"TestWithReturn": {
		name: "TestWithReturn",
		req:  &wot.DataSchema{},
		res: &wot.DataSchema{
			DataType: "object",
			ObjectSchema: &wot.ObjectSchema{
				Properties: map[string]wot.DataSchema{
					"Test": {},
				},
			},
		},
	},
	"TestWithRequest": {
		name: "TestWithRequest",
		req: &wot.DataSchema{
			DataType: "object",
			ObjectSchema: &wot.ObjectSchema{
				Properties: map[string]wot.DataSchema{
					"Test": {},
				},
			},
		},
		res: &wot.DataSchema{},
	},
	"TestWithRequestAndReturn": {
		name: "TestWithRequestAndReturn",
		req: &wot.DataSchema{
			DataType: "object",
			ObjectSchema: &wot.ObjectSchema{
				Properties: map[string]wot.DataSchema{
					"Test": {},
				},
			},
		},
		res: &wot.DataSchema{
			DataType: "object",
			ObjectSchema: &wot.ObjectSchema{
				Properties: map[string]wot.DataSchema{
					"Test": {},
				},
			},
		},
	},
}

var categorizeRPCTest = []struct {
	inIab interactionAffordanceBuilder
	out   affClasses
}{
	{
		interactionAffordanceBuilder{
			affs: map[string]affs{
				"SimpleTest": categorizeRPCTestAffordances["SimpleTest"],
			},
			cats: catProps{
				prop:   func(a affs) bool { return true },
				action: func(a affs) bool { return false },
				event:  func(a affs) bool { return false },
			},
		},
		affClasses{
			prop:   []affs{categorizeRPCTestAffordances["SimpleTest"]},
			action: []affs{},
			event:  []affs{},
		},
	},
	{
		interactionAffordanceBuilder{
			affs: map[string]affs{
				"SimpleTest": categorizeRPCTestAffordances["SimpleTest"],
				"GetTest":    categorizeRPCTestAffordances["GetTest"],
				"SetTest":    categorizeRPCTestAffordances["SetTest"],
			},
			cats: catProps{
				prop:   or(startsWithGet, startsWithSet),
				action: func(a affs) bool { return true },
				event:  func(a affs) bool { return true },
			},
		},
		affClasses{
			prop:   []affs{categorizeRPCTestAffordances["GetTest"], categorizeRPCTestAffordances["SetTest"]},
			action: []affs{},
			event:  []affs{categorizeRPCTestAffordances["SimpleTest"]},
		},
	},
	{
		interactionAffordanceBuilder{
			affs: map[string]affs{
				"SimpleTest": categorizeRPCTestAffordances["SimpleTest"],
				"GetTest":    categorizeRPCTestAffordances["GetTest"],
				"SetTest":    categorizeRPCTestAffordances["SetTest"],
			},
			cats: catProps{
				prop:   or(startsWithGet, startsWithSet),
				action: func(a affs) bool { return true },
				event:  func(a affs) bool { return true },
			},
		},
		affClasses{
			prop:   []affs{categorizeRPCTestAffordances["GetTest"], categorizeRPCTestAffordances["SetTest"]},
			action: []affs{},
			event:  []affs{categorizeRPCTestAffordances["SimpleTest"]},
		},
	},
	{
		interactionAffordanceBuilder{
			affs: map[string]affs{
				"TestWithRequest":          categorizeRPCTestAffordances["TestWithRequest"],
				"TestWithReturn":           categorizeRPCTestAffordances["TestWithReturn"],
				"TestWithRequestAndReturn": categorizeRPCTestAffordances["TestWithRequestAndReturn"],
			},
			cats: catProps{
				prop:   and(hasRequestType, hasReturnType),
				action: hasReturnType,
				event:  not(hasRequestType),
			},
		},
		affClasses{
			prop:   []affs{categorizeRPCTestAffordances["TestWithRequestAndReturn"]},
			action: []affs{categorizeRPCTestAffordances["TestWithRequest"]},
			event:  []affs{categorizeRPCTestAffordances["TestWithReturn"]},
		},
	},
}

func TestCategorizeRPC(t *testing.T) {
	for _, v := range categorizeRPCTest {
		v.inIab.categorizeRPCs()
		equals(v.out.prop, v.inIab.affC.prop, t)
		equals(v.out.action, v.inIab.affC.action, t)
		equals(v.out.event, v.inIab.affC.event, t)
	}
}

func equals(a1 []affs, a2 []affs, t *testing.T) {
	if len(a1) != len(a2) {
		t.Errorf("The length differs for the provided affordances.\n Expected %v\n but got: %v\n", a1, a2)
	} else {
	l:
		for k, v := range a1 {
			if a2[k] != v {
				for _, a := range a2 {
					if v == a {
						continue l
					}
				}
				t.Errorf("One expected element was not found. \n Expected: %v\n but was not in: %v\n", v, a2)
			}
		}
	}
}