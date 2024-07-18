package internal

type InternalSuperUserPermissions struct{}

func (InternalSuperUserPermissions) IsPredefinedScope() {}

func (InternalSuperUserPermissions) IsPermissionScope() {}
