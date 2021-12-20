package binaryIoTypeA

import (
	"fmt"
	"testing"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestFunctionControlSetProto(t *testing.T) {
	assert := assert.New(t)
	// cmd := binio.FunctionControlSet{
	// 	Type: &binio.FunctionControlSet_Single{
	// 		Single: &binio.SetSingle{
	// 			Channel: 1,
	// 			State:   true,
	// 		},
	// 	},
	// }
	// envelopeCmd, err := functionblock.FunctionControlSet(&cmd)
	// assert.Nil(err)
	single := binio.SetSingle{
		Channel: 0,
		State:   false,
	}

	fmt.Println(single.ProtoReflect().Descriptor().FullName())
	assert.Nil(nil)
	// fmt.Printf("%+v\n", envelopeCmd)
}
