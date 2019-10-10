// This file is automatically generated. DO NOT MODIFY!

package lumina

func (*PullMdPacket) getType() PacketType {
	return PKT_PULL_MD
}

func (*PullMdPacket) getResponseType() PacketType {
	return PKT_PULL_MD_RESULT
}

func (this *PullMdPacket) readFrom(r Reader) (err error) {
	// Field this.Flags
	// Basic uint32
	// Typed MdKeyFlag
	var v1 uint32
	if v1, err = readUint32(r); err != nil {
		return
	}
	this.Flags = MdKeyFlag(v1)
	// Field this.Keys
	// Slice []MdKey
	var v2 uint32
	if v2, err = readUint32(r); err != nil {
		return
	}
	this.Keys = make([]MdKey, v2)
	for i := uint32(0); i < v2; i++ {
		// Field this.Keys[i]
		// Basic uint32
		// Typed MdKey
		var v3 uint32
		if v3, err = readUint32(r); err != nil {
			return
		}
		this.Keys[i] = MdKey(v3)
	}
	// Field this.PatternIds
	// Slice []PatternId
	var v4 uint32
	if v4, err = readUint32(r); err != nil {
		return
	}
	this.PatternIds = make([]PatternId, v4)
	for i := uint32(0); i < v4; i++ {
		// Field this.PatternIds[i]
		// Struct PatternId
		if err = this.PatternIds[i].readFrom(r); err != nil {
			return
		}
	}
	return
}

func (this *PullMdPacket) writeTo(w Writer) (err error) {
	// Field this.Flags
	// Basic uint32
	// Typed MdKeyFlag
	if err = writeUint32(w, uint32(this.Flags)); err != nil {
		return
	}
	// Field this.Keys
	// Slice []MdKey
	if len(this.Keys) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v2 = uint32(len(this.Keys))
	if err = writeUint32(w, v2); err != nil {
		return
	}
	for i := uint32(0); i < v2; i++ {
		// Field this.Keys[i]
		// Basic uint32
		// Typed MdKey
		if err = writeUint32(w, uint32(this.Keys[i])); err != nil {
			return
		}
	}
	// Field this.PatternIds
	// Slice []PatternId
	if len(this.PatternIds) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v4 = uint32(len(this.PatternIds))
	if err = writeUint32(w, v4); err != nil {
		return
	}
	for i := uint32(0); i < v4; i++ {
		// Field this.PatternIds[i]
		// Struct PatternId
		if err = this.PatternIds[i].writeTo(w); err != nil {
			return
		}
	}
	return
}
