package middleware

import (
	"context"
	utilsGoServer "github.com/kintohub/utils-go/server"
	"github.com/kintoproj/kinto-core/pkg/types"
)

// This is a wrapper that allows adapter implementations of Middleware be able to be passed in as middleware reference
// Composition in go does not allow composition to be = to the base type you are composing. This interface bridges that gap.
type Interface interface {
	GetEnvironment(id string) (*types.Environment, *utilsGoServer.Error)
	GetEnvironments() (*types.Environments, *utilsGoServer.Error)
	CreateEnvironment(name string) (*types.Environment, *utilsGoServer.Error)
	UpdateEnvironment(id string, name string) (*types.Environment, *utilsGoServer.Error)
	DeleteEnvironment(id string) *utilsGoServer.Error

	CreateBlock(
		envId, name string,
		buildConfig *types.BuildConfig,
		runConfig *types.RunConfig) (string, string, *utilsGoServer.Error)
	GetBlock(name, envId string) (*types.Block, *utilsGoServer.Error)
	GetBlocks(envId string) (*types.Blocks, *utilsGoServer.Error)
	DeployBlockUpdate(
		name, envId, baseReleaseId string,
		buildConfig *types.BuildConfig,
		runConfig *types.RunConfig) (string, string, *utilsGoServer.Error)
	TriggerDeploy(
		name, envId string) (string, string, *utilsGoServer.Error)
	RollbackBlock(name, envId, releaseId string) (string, string, *utilsGoServer.Error)
	DeleteBlock(name, envId string) *utilsGoServer.Error
	WatchReleasesStatus(blockName, envId string, ctx context.Context, statusChan chan *types.ReleasesStatus) *utilsGoServer.Error
	GetBlocksHealthStatus(envId string) (*types.BlockStatuses, *utilsGoServer.Error)
	WatchJobsStatus(
		blockName, envId string, ctx context.Context, sendClientLogs func(jobStatus *types.JobStatus) error) *utilsGoServer.Error
	GetBlocksMetrics(name, envId string) (*types.BlocksMetrics, *utilsGoServer.Error)
	KillBlockInstance(id, envId string) *utilsGoServer.Error
	// Scale down all the resources of a block down to 0
	// return blockName, releaseId, error if any
	SuspendBlock(blockName, envId string) (string, string, *utilsGoServer.Error)

	WatchBuildLogs(releaseId, blockName, envId string, ctx context.Context, logsChan chan *types.Logs) *utilsGoServer.Error
	UpdateBuildStatus(releaseId, blockName, envId string, buildState types.BuildStatus_State) (*types.Release, *utilsGoServer.Error)
	UpdateBuildCommitSha(releaseId, blockName, envId, commitSha string) *utilsGoServer.Error

	WatchConsoleLogs(blockName, envId string, ctx context.Context, logsChan chan *types.ConsoleLog) *utilsGoServer.Error

	GetKintoConfiguration() (*types.KintoConfiguration, error)

	AbortRelease(ctx context.Context, blockName, releaseId, envId string) *utilsGoServer.Error

	EnableExternalURL(name, envId, releaseId string) *utilsGoServer.Error
	DisableExternalURL(name, envId string) *utilsGoServer.Error

	CreateCustomDomainName(blockName, envId, domainName string, protocol types.RunConfig_Protocol) *utilsGoServer.Error
	DeleteCustomDomainName(blockName, envId, domainName string, protocol types.RunConfig_Protocol) *utilsGoServer.Error
	CheckCertificateReadiness(blockName, envId string) bool

	StartTeleport(
		ctx context.Context, envId, blockNameToTeleport string) (*types.TeleportServiceData, *utilsGoServer.Error)
	StopTeleport(envId, blockNameTeleported string) *utilsGoServer.Error

	TagRelease(tag, blockName, envId, releaseId string) *utilsGoServer.Error
	PromoteRelease(tag, releaseId, blockName, envId, targetEnvId string) *utilsGoServer.Error

	GenReleaseConfigFromKintoFile(
		org, repo, branch, envId, githubUserToken string, blockType types.Block_Type) (*types.ReleaseConfig, *utilsGoServer.Error)

	// this is a private func used for internal middleware package only. Adapter impls should not have access to this
	add(m Interface)
}

// Implements middleware pattern + Implements controller.ControllerInterface for Adapter Pattern
// All functions implemented for the Controller interface by default forward the request to `next`
// To override the function as middleware, create a struct and insert Middleware as composition and override each
// function as you please. see test file for example
type Middleware struct {
	// Note - it should never be possible for this to be nil
	next Interface
}

func (m *Middleware) add(middleware Interface) {
	m.next = middleware
}

func NewControllerMiddlewareOrDie(middlewares ...Interface) Interface {
	// Build linked list of middleware
	var last Interface
	for _, middleware := range middlewares {
		if last == nil {
			last = middleware
		} else {
			last.add(middleware)
			last = middleware
		}
	}

	// Return first middleware as the top of the chain
	return middlewares[0]
}

func (m *Middleware) CreateEnvironment(name string) (*types.Environment, *utilsGoServer.Error) {
	return m.next.CreateEnvironment(name)
}

func (m *Middleware) GetEnvironment(id string) (*types.Environment, *utilsGoServer.Error) {
	return m.next.GetEnvironment(id)
}

func (m *Middleware) DeleteEnvironment(id string) *utilsGoServer.Error {
	return m.next.DeleteEnvironment(id)
}

func (m *Middleware) CreateBlock(
	envId, displayName string,
	buildConfig *types.BuildConfig,
	runConfig *types.RunConfig) (string, string, *utilsGoServer.Error) {
	return m.next.CreateBlock(envId, displayName, buildConfig, runConfig)
}

func (m *Middleware) GetBlock(name, envId string) (*types.Block, *utilsGoServer.Error) {
	return m.next.GetBlock(name, envId)
}

func (m *Middleware) GetBlocks(envId string) (*types.Blocks, *utilsGoServer.Error) {
	return m.next.GetBlocks(envId)
}

func (m *Middleware) DeployBlockUpdate(
	name, envId, baseReleaseId string,
	buildConfig *types.BuildConfig,
	runConfig *types.RunConfig) (string, string, *utilsGoServer.Error) {
	return m.next.DeployBlockUpdate(name, envId, baseReleaseId, buildConfig, runConfig)
}

func (m *Middleware) RollbackBlock(name, envId, releaseId string) (string, string, *utilsGoServer.Error) {
	return m.next.RollbackBlock(name, envId, releaseId)
}

func (m *Middleware) DeleteBlock(name, envId string) *utilsGoServer.Error {
	return m.next.DeleteBlock(name, envId)
}

func (m *Middleware) WatchReleasesStatus(blockName, envId string, ctx context.Context, statusChan chan *types.ReleasesStatus) *utilsGoServer.Error {
	return m.next.WatchReleasesStatus(blockName, envId, ctx, statusChan)
}

func (m *Middleware) GetBlocksHealthStatus(envId string) (*types.BlockStatuses, *utilsGoServer.Error) {
	return m.next.GetBlocksHealthStatus(envId)
}

func (m *Middleware) WatchJobsStatus(blockName, envId string, ctx context.Context, sendClientLogs func(jobStatus *types.JobStatus) error) *utilsGoServer.Error {
	return m.next.WatchJobsStatus(blockName, envId, ctx, sendClientLogs)
}

func (m *Middleware) GetBlocksMetrics(name, envId string) (*types.BlocksMetrics, *utilsGoServer.Error) {
	return m.next.GetBlocksMetrics(name, envId)
}

func (m *Middleware) KillBlockInstance(id, envId string) *utilsGoServer.Error {
	return m.next.KillBlockInstance(id, envId)
}

func (m *Middleware) SuspendBlock(blockName, envId string) (string, string, *utilsGoServer.Error) {
	return m.next.SuspendBlock(blockName, envId)
}

func (m *Middleware) WatchBuildLogs(releaseId, blockName, envId string, ctx context.Context, logsChan chan *types.Logs) *utilsGoServer.Error {
	return m.next.WatchBuildLogs(releaseId, blockName, envId, ctx, logsChan)
}

func (m *Middleware) UpdateBuildCommitSha(releaseId, blockName, envId, commitSha string) *utilsGoServer.Error {
	return m.next.UpdateBuildCommitSha(releaseId, blockName, envId, commitSha)
}

func (m *Middleware) UpdateBuildStatus(releaseId, blockName, envId string, buildState types.BuildStatus_State) (*types.Release, *utilsGoServer.Error) {
	return m.next.UpdateBuildStatus(releaseId, blockName, envId, buildState)
}

func (m *Middleware) WatchConsoleLogs(blockName, envId string, ctx context.Context, logsChan chan *types.ConsoleLog) *utilsGoServer.Error {
	return m.next.WatchConsoleLogs(blockName, envId, ctx, logsChan)
}

func (m *Middleware) GetKintoConfiguration() (*types.KintoConfiguration, error) {
	return m.next.GetKintoConfiguration()
}

func (m *Middleware) AbortRelease(ctx context.Context, blockName, releaseId, envId string) *utilsGoServer.Error {
	return m.next.AbortRelease(ctx, blockName, releaseId, envId)
}

func (m *Middleware) EnableExternalURL(name, envId, releaseId string) *utilsGoServer.Error {
	return m.next.EnableExternalURL(name, envId, releaseId)
}

func (m *Middleware) DisableExternalURL(name, envId string) *utilsGoServer.Error {
	return m.next.DisableExternalURL(name, envId)
}

func (m *Middleware) CreateCustomDomainName(
	blockName, envId, domainName string, protocol types.RunConfig_Protocol) *utilsGoServer.Error {

	return m.next.CreateCustomDomainName(blockName, envId, domainName, protocol)
}

func (m *Middleware) DeleteCustomDomainName(
	blockName, envId, domainName string, protocol types.RunConfig_Protocol) *utilsGoServer.Error {

	return m.next.DeleteCustomDomainName(blockName, envId, domainName, protocol)
}

func (m *Middleware) CheckCertificateReadiness(blockName, envId string) bool {
	return m.next.CheckCertificateReadiness(blockName, envId)
}

func (m *Middleware) StartTeleport(ctx context.Context, envId, blockName string) (*types.TeleportServiceData, *utilsGoServer.Error) {
	return m.next.StartTeleport(ctx, envId, blockName)
}

func (m *Middleware) StopTeleport(envId, blockName string) *utilsGoServer.Error {
	return m.next.StopTeleport(envId, blockName)
}

func (m *Middleware) TagRelease(tag, blockName, envId, releaseId string) *utilsGoServer.Error {
	return m.next.TagRelease(tag, blockName, envId, releaseId)
}

func (m *Middleware) PromoteRelease(tag, releaseId, blockName, envId, targetEnvId string) *utilsGoServer.Error {
	return m.next.PromoteRelease(tag, releaseId, blockName, envId, targetEnvId)
}

func (m *Middleware) GenReleaseConfigFromKintoFile(
	org, repo, branch, envId, githubUserToken string, blockType types.Block_Type) (*types.ReleaseConfig, *utilsGoServer.Error) {
	return m.next.GenReleaseConfigFromKintoFile(org, repo, branch, envId, githubUserToken, blockType)
}
