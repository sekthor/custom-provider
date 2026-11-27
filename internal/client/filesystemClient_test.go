package client

import (
	"testing"
)

func isEqual(a Resource, b Resource) bool {
	return a.Id == b.Id && a.Name == b.Name
}

func Test_filesystemClient_GetResource(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		file      string
		want      Resource
		wantErr   bool
	}{
		{
			name:      "Read resource file from filesystem",
			directory: "/home/sekthor/code/apps/custom-provider/test",
			file:      "resource.yaml",
			want: Resource{
				Id:   "xyz",
				Name: "resourcename",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFilesystemClient(tt.directory)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := f.GetResource(t.Context(), tt.file)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetResource() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetResource() succeeded unexpectedly")
			}
			if !isEqual(tt.want, got) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
