package stonenode

import (
	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	ptypes "github.com/bnb-chain/greenfield-storage-provider/pkg/types/v1"
	"github.com/bnb-chain/greenfield-storage-provider/util/log"
)

// dispatchSecondarySP convert pieceDataBySegment to ec dimensional slice, first dimensional is ec number such as ec1, contains [][]byte data
// pieceDataBySegment is a three-dimensional slice, first dimensional is segment index, second is [][]byte data
func (node *StoneNodeService) dispatchSecondarySP(pieceDataBySegment [][][]byte, redundancyType ptypes.RedundancyType,
	secondarySPs []string, targetIdx []uint32) ([][][]byte, error) {
	if len(pieceDataBySegment) == 0 {
		return nil, merrors.ErrInvalidPieceData
	}
	if len(secondarySPs) == 0 {
		return nil, merrors.ErrSecondarySPNumber
	}
	var pieceDataBySecondary [][][]byte
	var err error
	switch redundancyType {
	case ptypes.RedundancyType_REDUNDANCY_TYPE_REPLICA_TYPE, ptypes.RedundancyType_REDUNDANCY_TYPE_INLINE_TYPE:
		pieceDataBySecondary, err = dispatchReplicaData(pieceDataBySegment, secondarySPs, targetIdx)
	default: // ec type
		pieceDataBySecondary, err = dispatchECData(pieceDataBySegment, secondarySPs, targetIdx)
	}
	if err != nil {
		log.Errorw("dispatch piece data by secondary error", "error", err)
		return nil, err
	}
	return pieceDataBySecondary, nil
}

// dispatchReplicaData dispatches replica or inline data into different sp, each sp should store all segments data of an object
// if an object uses replica type, it's split into 10 segments and there are 6 sp, each sp should store 10 segments data
// if an object uses inline type, there is only one segment and there are 6 sp, each sp should store 1 segment data
func dispatchReplicaData(pieceDataBySegment [][][]byte, secondarySPs []string, targetIdx []uint32) ([][][]byte, error) {
	if len(secondarySPs) < len(targetIdx) {
		return nil, merrors.ErrSecondarySPNumber
	}

	segmentLength := len(pieceDataBySegment[0])
	if segmentLength != 1 {
		return nil, merrors.ErrInvalidSegmentData
	}

	data := convertSlice(pieceDataBySegment, segmentLength)
	segmentPieceSlice := make([][][]byte, len(targetIdx))
	for i := 0; i < len(targetIdx); i++ {
		segmentPieceSlice[i] = data[0]
	}
	return segmentPieceSlice, nil
}

// dispatchECData dispatches ec data into different sp
// one sp stores same ec column data: sp1 stores all ec1 data, sp2 stores all ec2 data, etc
func dispatchECData(pieceDataBySegment [][][]byte, secondarySPs []string, targetIdx []uint32) ([][][]byte, error) {
	segmentLength := len(pieceDataBySegment[0])
	if segmentLength < 6 {
		return nil, merrors.ErrInvalidECData
	}
	if segmentLength > len(secondarySPs) {
		return nil, merrors.ErrSecondarySPNumber
	}

	data := convertSlice(pieceDataBySegment, segmentLength)
	ecPieceSlice := make([][][]byte, len(targetIdx))
	for index, value := range targetIdx {
		ecPieceSlice[index] = data[value]
	}
	return ecPieceSlice, nil
}

func convertSlice(data [][][]byte, length int) [][][]byte {
	tempSlice := make([][][]byte, length)
	for i := 0; i < length; i++ {
		tempSlice[i] = make([][]byte, 0)
		for j := 0; j < len(data); j++ {
			tempSlice[i] = append(tempSlice[i], data[j][i])
		}
	}
	return tempSlice
}