package client

import (
	"reflect"
	"testing"
)

func Test_projectInfoSlice_GetProjectByName(t *testing.T) {
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		s       projectInfoSlice
		args    args
		want    *projectInfo
		wantErr bool
	}{
		{
			name: "GetProjectByName - Found",
			s: []projectInfo{
				{ID: 1, Name: "name1", Description: "descr1"},
				{ID: 2, Name: "test", Description: "descr_test"},
				{ID: 3, Name: "name3", Description: "descr3"},
			},
			args:    args{projectName: "test"},
			want:    &projectInfo{ID: 2, Name: "test", Description: "descr_test"},
			wantErr: false,
		},
		{
			name: "GetProjectByName - Not Found",
			s: []projectInfo{
				{ID: 1, Name: "name1", Description: "descr1"},
				{ID: 2, Name: "name2", Description: "descr2"},
				{ID: 3, Name: "name3", Description: "descr3"},
			},
			args:    args{projectName: "test"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "GetProjectByName - Empty slice",
			s:       []projectInfo{},
			args:    args{projectName: "test"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetProjectByName(tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("projectInfoSlice.GetProjectByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("projectInfoSlice.GetProjectByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_robotInfoSlice_GetRobotByName(t *testing.T) {
	type args struct {
		robotName string
	}
	tests := []struct {
		name    string
		s       robotInfoSlice
		args    args
		want    *robotInfo
		wantErr bool
	}{
		{
			name: "GetProjectByName - Found",
			s: []robotInfo{
				{ID: 1, Name: "name1", DeploymentStatus: 1},
				{ID: 2, Name: "test", DeploymentStatus: 3},
				{ID: 3, Name: "name3", DeploymentStatus: 1},
			},
			args:    args{robotName: "test"},
			want:    &robotInfo{ID: 2, Name: "test", DeploymentStatus: 3},
			wantErr: false,
		},
		{
			name: "GetProjectByName - Not Found",
			s: []robotInfo{
				{ID: 1, Name: "name1", DeploymentStatus: 1},
				{ID: 2, Name: "name2", DeploymentStatus: 1},
				{ID: 3, Name: "name3", DeploymentStatus: 1},
			},
			args:    args{robotName: "test"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "GetProjectByName - Empty slice",
			s:       []robotInfo{},
			args:    args{robotName: "test"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetRobotByName(tt.args.robotName)
			if (err != nil) != tt.wantErr {
				t.Errorf("robotInfoSlice.GetRobotByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("robotInfoSlice.GetRobotByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
