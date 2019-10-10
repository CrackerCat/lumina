// This file is automatically generated. DO NOT MODIFY!

package lumina

func (*HeloPacket) getType() PacketType {
	return PKT_HELO
}

func (*HeloPacket) getResponseType() PacketType {
	return PKT_RPC_OK
}

func (this *HeloPacket) readFrom(r Reader) (err error) {
	// Field this.ClientVersion
	// Basic int32
	if this.ClientVersion, err = readInt32(r); err != nil {
		return
	}
	// Field this.Key
	// Slice []byte
	var v1 uint32
	if v1, err = readUint32(r); err != nil {
		return
	}
	this.Key = make([]byte, v1)
	if err = readBytes(r, this.Key); err != nil {
		return
	}
	// Field this.LicenseId
	// Array [6]byte
	if err = readBytes(r, this.LicenseId[:]); err != nil {
		return
	}
	// Field this.RecordConv
	// Basic bool
	if this.RecordConv, err = readBool(r); err != nil {
		return
	}
	return
}

func (this *HeloPacket) writeTo(w Writer) (err error) {
	// Field this.ClientVersion
	// Basic int32
	if err = writeInt32(w, this.ClientVersion); err != nil {
		return
	}
	// Field this.Key
	// Slice []byte
	if len(this.Key) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v1 = uint32(len(this.Key))
	if err = writeUint32(w, v1); err != nil {
		return
	}
	if err = writeBytes(w, this.Key); err != nil {
		return
	}
	// Field this.LicenseId
	// Array [6]byte
	if err = writeBytes(w, this.LicenseId[:]); err != nil {
		return
	}
	// Field this.RecordConv
	// Basic bool
	if err = writeBool(w, this.RecordConv); err != nil {
		return
	}
	return
}
