package binaryIoTypeA

import (
	"testing"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	fb "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestFunctionControlSetProto(t *testing.T) {
	assert := assert.New(t)
	cmd := binio.FunctionControlSet{
		Type: &binio.FunctionControlSet_Single{
			Single: &binio.SetSingle{
				Channel: 1,
				State:   true,
			},
		},
	}
	envelopeCmd, err := functionblock.FunctionControlSet(&cmd, string(cmd.Type.(*binio.FunctionControlSet_Single).Single.ProtoReflect().Descriptor().FullName()))
	assert.Nil(err)
	serialized, err := proto.Marshal(envelopeCmd)
	assert.Nil(err)
	newEnvelopeCmd := &fb.Command{}
	err = proto.Unmarshal(serialized, newEnvelopeCmd)
	assert.Nil(err)
	var newSet binio.FunctionControlSet
	err = newEnvelopeCmd.GetFunctionControl().GetFunctionSpecificFunctionControlSet().UnmarshalTo(&newSet)
	assert.Nil(err)
	assert.Equal(uint32(1), newSet.Type.(*binio.FunctionControlSet_Single).Single.Channel)
	assert.Equal(true, newSet.Type.(*binio.FunctionControlSet_Single).Single.State)
}
