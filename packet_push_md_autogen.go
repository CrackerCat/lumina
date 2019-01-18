// This file is automatically generated. DO NOT MODIFY!

package lumina

func (this *InputFile) ReadFrom(r Reader) (err error) {
    // Field this.Path
    // Basic string
    if this.Path, err = readString(r); err != nil { return }
    // Field this.MD5
    // Array [16]byte
    if err = readBytes(r, this.MD5[:]); err != nil { return }
    return
}

func (this *InputFile) WriteTo(w Writer) (err error) {
    // Field this.Path
    // Basic string
    if err = writeString(w, this.Path); err != nil { return }
    // Field this.MD5
    // Array [16]byte
    if err = writeBytes(w, this.MD5[:]); err != nil { return }
    return
}

func (*PushMdPacket) getType() PacketType {
    return PKT_PUSH_MD
}

func (*PushMdPacket) getResponseType() PacketType {
    return PKT_PUSH_MD_RESULT
}

func (this *PushMdPacket) ReadFrom(r Reader) (err error) {
    // Field this.Flags
    // Basic uint32
    // Typed PushMdFlag
    var v1 uint32
    if v1, err = readUint32(r); err != nil { return }
    this.Flags = PushMdFlag(v1)
    // Field this.Idb
    // Basic string
    if this.Idb, err = readString(r); err != nil { return }
    // Field this.Input
    // Struct InputFile
    if err = this.Input.ReadFrom(r); err != nil { return }
    // Field this.Hostname
    // Basic string
    if this.Hostname, err = readString(r); err != nil { return }
    // Field this.Contents
    // Slice []FuncInfoAndPattern
    var v2 uint32
    if v2, err = readUint32(r); err != nil { return }
    this.Contents = make([]FuncInfoAndPattern, v2)
    for i := uint32(0); i < v2; i++ {
        // Field this.Contents[i]
        // Struct FuncInfoAndPattern
        if err = this.Contents[i].ReadFrom(r); err != nil { return }
    }
    // Field this.EAs
    // Slice []uint64
    var v3 uint32
    if v3, err = readUint32(r); err != nil { return }
    this.EAs = make([]uint64, v3)
    for i := uint32(0); i < v3; i++ {
        // Field this.EAs[i]
        // Basic uint64
        if this.EAs[i], err = readUint64(r); err != nil { return }
    }
    return
}

func (this *PushMdPacket) WriteTo(w Writer) (err error) {
    // Field this.Flags
    // Basic uint32
    // Typed PushMdFlag
    if err = writeUint32(w, uint32(this.Flags)); err != nil { return }
    // Field this.Idb
    // Basic string
    if err = writeString(w, this.Idb); err != nil { return }
    // Field this.Input
    // Struct InputFile
    if err = this.Input.WriteTo(w); err != nil { return }
    // Field this.Hostname
    // Basic string
    if err = writeString(w, this.Hostname); err != nil { return }
    // Field this.Contents
    // Slice []FuncInfoAndPattern
    if len(this.Contents) > 0x7FFFFFFF { err = errTooLong; return }
    var v2 = uint32(len(this.Contents))
    if err = writeUint32(w, v2); err != nil { return }
    for i := uint32(0); i < v2; i++ {
        // Field this.Contents[i]
        // Struct FuncInfoAndPattern
        if err = this.Contents[i].WriteTo(w); err != nil { return }
    }
    // Field this.EAs
    // Slice []uint64
    if len(this.EAs) > 0x7FFFFFFF { err = errTooLong; return }
    var v3 = uint32(len(this.EAs))
    if err = writeUint32(w, v3); err != nil { return }
    for i := uint32(0); i < v3; i++ {
        // Field this.EAs[i]
        // Basic uint64
        if err = writeUint64(w, this.EAs[i]); err != nil { return }
    }
    return
}