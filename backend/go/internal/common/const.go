package common

const (

	// Resource Types
	AccountRt = "Account"
	ProjectRt = "Project"
	SprintRt  = "Sprint"
	StoryRt   = "Story"

	// Account Actions
	CreateProjectA = "CreateProject"
	AdduserA       = "AddUser"

	// Project Actions
	ReadProjectA   = "ReadProject"
	DeleteProjectA = "DeleteProject"
	CreateSprintA  = "CreateSprint"

	// Sprint Actions
	ReadSprintA  = "ReadSprint"
	StartSprintA = "StartSprint"
	EndSprintA   = "EndSprint"
	CreateStoryA = "CreateStory"

	// Story Actions
	ReadStoryA           = "ReadStory"
	EstimateStoryA       = "EstimateStory"
	ChangeStoryStatusA   = "ChangeStoryStatus"
	ChangeStoryAssigneeA = "ChangeStoryAssignee"

	// Policies
	CanManageAccountP = "CanManageAccount"

	// Roles
	AccountAdministratorR = "AccountAdministrator"
)
