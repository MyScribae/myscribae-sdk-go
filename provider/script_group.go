package provider

import "context"

type ScriptGroup struct {
	Uuid *string
	ScriptGroupInput
	Scripts []*Script
}

func (sg *ScriptGroup) Sync(ctx context.Context) {
	
}
