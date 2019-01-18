// This file is automatically generated. DO NOT MODIFY!

package lumina

func (*PushMdResultPacket) getType() PacketType {
    return PKT_PUSH_MD_RESULT
}

func (this *PushMdResultPacket) ReadFrom(r Reader) (err error) {
    // Field this.Codes
    // Slice []OpResult
    var v1 uint32
    if v1, err = readUint32(r); err != nil { return }
    this.Codes = make([]OpResult, v1)
    for i := uint32(0); i < v1; i++ {
        // Field this.Codes[i]
        // Basic int32
        // Typed OpResult
        var v2 int32
        if v2, err = readInt32(r); err != nil { return }
        this.Codes[i] = OpResult(v2)
    }
    return
}

func (this *PushMdResultPacket) WriteTo(w Writer) (err error) {
    // Field this.Codes
    // Slice []OpResult
    if len(this.Codes) > 0x7FFFFFFF { err = errTooLong; return }
    var v1 = uint32(len(this.Codes))
    if err = writeUint32(w, v1); err != nil { return }
    for i := uint32(0); i < v1; i++ {
        // Field this.Codes[i]
        // Basic int32
        // Typed OpResult
        if err = writeInt32(w, int32(this.Codes[i])); err != nil { return }
    }
    return
}