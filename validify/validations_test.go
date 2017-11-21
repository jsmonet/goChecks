package validify

import (
	"fmt"
	"testing"
)

func TestPort(t *testing.T) {
	cases := []struct {
		port         int
		shouldBeBool bool
	}{
		{
			port:         65595,
			shouldBeBool: false,
		},
		{
			port:         22,
			shouldBeBool: true,
		},
	}

	for _, c := range cases {
		boolResult, err := Port(c.port)
		if err != nil {
			fmt.Println(err)
		}
		if boolResult != c.shouldBeBool {
			t.Errorf("Expected %v, but got %v", c.shouldBeBool, boolResult)
		}
	}
}

func TestNeorole(t *testing.T) {
	cases := []struct {
		testRole     string
		shouldBeBool bool
	}{
		{
			testRole:     "MaStEr",
			shouldBeBool: true,
		},
		{
			testRole:     "master",
			shouldBeBool: true,
		},
		{
			testRole:     "MASTER",
			shouldBeBool: true,
		},
		{
			testRole:     "SLavE",
			shouldBeBool: true,
		},
		{
			testRole:     "slave",
			shouldBeBool: true,
		},
		{
			testRole:     "SLAVE",
			shouldBeBool: true,
		},
		{
			testRole:     "skidding",
			shouldBeBool: false,
		},
		{
			testRole:     "S1av3",
			shouldBeBool: false,
		},
		{
			testRole:     "123e987435**",
			shouldBeBool: false,
		},
	}

	for _, c := range cases {
		boolResult, _ := Neorole(c.testRole)
		// bad roles are going to throw exceptions. We can safely discard this error
		if boolResult != c.shouldBeBool {
			t.Errorf("Expected %v, but got %v", c.shouldBeBool, boolResult)
		}
	}
}

func TestAuthb64(t *testing.T) {
	cases := []struct {
		base64EncodedString string
		shouldBeBool        bool
	}{
		{
			base64EncodedString: "dGhpcyBpcyBvYnZpb3VzbHkgYSBkZWNvZGVhYmxlIGJhc2U2NCBzdHJpbmcK",
			shouldBeBool:        true,
		},
		{
			base64EncodedString: "this is definitely not a base64 encoded string",
			shouldBeBool:        false,
		},
	}

	for _, c := range cases {
		boolResult, _ := Authb64(c.base64EncodedString)
		if boolResult != c.shouldBeBool {
			t.Errorf("Expected %v, but got %v", c.shouldBeBool, boolResult)
		}
	}
}

func TestPercentages(t *testing.T) {
	cases := []struct {
		testWarn int
		testCrit int
		testBool bool
	}{
		{
			testWarn: 85,
			testCrit: 95,
			testBool: true,
		},
		{
			testWarn: 85,
			testCrit: 75,
			testBool: false,
		},
		{
			testWarn: 75,
			testCrit: 100,
			testBool: false,
		},
	}

	for _, c := range cases {
		percentagesBool, _ := Percentages(c.testWarn, c.testCrit)
		if percentagesBool != c.testBool {
			t.Errorf("Expected %v, got %v", c.testBool, percentagesBool)
		}
	}
}
