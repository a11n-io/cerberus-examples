package common

const (

	// Resource Types
	Account_RT = "Account"
	Project_RT = "Project"
	Sprint_RT  = "Sprint"
	Story_RT   = "Story"

	// Account Actions
	CreateProject_A = "CreateProject"
	AddUser_A       = "AddUser"

	// Project Actions
	ReadProject_A   = "ReadProject"
	DeleteProject_A = "DeleteProject"
	CreateSprint_A  = "CreateSprint"

	// Sprint Actions
	ReadSprint_A  = "ReadSprint"
	StartSprint_A = "StartSprint"
	EndSprint_A   = "EndSprint"
	CreateStory_A = "CreateStory"

	// Story Actions
	ReadStory_A           = "ReadStory"
	EstimateStory_A       = "EstimateStory"
	ChangeStoryStatus_A   = "ChangeStoryStatus"
	ChangeStoryAssignee_A = "ChangeStoryAssignee"

	// Policies
	CanManageAccount_P = "CanManageAccount"
	CanManageProject_P = "CanManageProject"
	CanManageSprint_P  = "CanManageSprint"

	// Roles
	AccountAdministrator_R = "AccountAdministrator"
)
