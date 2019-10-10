// This file is automatically generated. DO NOT MODIFY!

package lumina

func (*PullMdResultPacket) getType() PacketType {
	return PKT_PULL_MD_RESULT
}

func (this *PullMdResultPacket) readFrom(r Reader) (err error) {
	// Field this.Codes
	// Slice []OpResult
	var v1 uint32
	if v1, err = readUint32(r); err != nil {
		return
	}
	this.Codes = make([]OpResult, v1)
	for i := uint32(0); i < v1; i++ {
		// Field this.Codes[i]
		// Basic int32
		// Typed OpResult
		var v2 int32
		if v2, err = readInt32(r); err != nil {
			return
		}
		this.Codes[i] = OpResult(v2)
	}
	// Field this.Results
	// Slice []FuncInfoAndFrequency
	var v3 uint32
	if v3, err = readUint32(r); err != nil {
		return
	}
	this.Results = make([]FuncInfoAndFrequency, v3)
	for i := uint32(0); i < v3; i++ {
		// Field this.Results[i]
		// Struct FuncInfoAndFrequency
		if err = this.Results[i].readFrom(r); err != nil {
			return
		}
	}
	return
}

func (this *PullMdResultPacket) writeTo(w Writer) (err error) {
	// Field this.Codes
	// Slice []OpResult
	if len(this.Codes) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v1 = uint32(len(this.Codes))
	if err = writeUint32(w, v1); err != nil {
		return
	}
	for i := uint32(0); i < v1; i++ {
		// Field this.Codes[i]
		// Basic int32
		// Typed OpResult
		if err = writeInt32(w, int32(this.Codes[i])); err != nil {
			return
		}
	}
	// Field this.Results
	// Slice []FuncInfoAndFrequency
	if len(this.Results) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v3 = uint32(len(this.Results))
	if err = writeUint32(w, v3); err != nil {
		return
	}
	for i := uint32(0); i < v3; i++ {
		// Field this.Results[i]
		// Struct FuncInfoAndFrequency
		if err = this.Results[i].writeTo(w); err != nil {
			return
		}
	}
	return
}
