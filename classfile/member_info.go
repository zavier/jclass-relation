package classfile

// 字段和方法信息
type MemberInfo struct {
	cp               ConstantPool
	accessFlags      uint16
	nameIndex        uint16
	descriptionIndex uint16
	attributes       []AttributeInfo
}

func readMembers(reader *ClassReader, cp ConstantPool) []*MemberInfo {
	memberCount := reader.readUint16()
	members := make([]*MemberInfo, memberCount)
	for i := range members {
		members[i] = readMember(reader, cp)
	}
	return members
}

func readMember(reader *ClassReader, cp ConstantPool) *MemberInfo {
	return &MemberInfo{
		cp:               cp,
		accessFlags:      reader.readUint16(),
		nameIndex:        reader.readUint16(),
		descriptionIndex: reader.readUint16(),
		attributes:       readAttributes(reader, cp),
	}
}

func (self *MemberInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}

func (self *MemberInfo) Description() string {
	return self.cp.getUtf8(self.descriptionIndex)
}
