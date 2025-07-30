package client

import "fmt"

type projectInfoSlice []projectInfo
type robotInfoSlice []robotInfo

func (s projectInfoSlice) GetProjectByName(projectName string) (*projectInfo, error) {
	var foundProject *projectInfo

	for _, p := range s {
		if p.Name == projectName {
			foundProject = &p
			break
		}
	}
	if foundProject == nil {
		return foundProject, fmt.Errorf("Project '%v' not found", projectName)
	}
	return foundProject, nil
}

func (s robotInfoSlice) GetRobotByName(robotName string) (*robotInfo, error) {
	var foundProject *robotInfo

	for _, r := range s {
		if r.Name == robotName {
			foundProject = &r
			break
		}
	}
	if foundProject == nil {
		return foundProject, fmt.Errorf("Robot '%v' not found", robotName)
	}
	return foundProject, nil
}
