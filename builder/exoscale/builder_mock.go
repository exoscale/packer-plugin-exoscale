package exoscale

import (
	"context"

	egoscale "github.com/exoscale/egoscale/v2"

	"github.com/stretchr/testify/mock"
)

type exoscaleClientMock struct {
	mock.Mock
}

func (m *exoscaleClientMock) AttachInstanceToPrivateNetwork(
	ctx context.Context,
	zone string,
	instance *egoscale.Instance,
	privateNetwork *egoscale.PrivateNetwork,
	opts ...egoscale.AttachInstanceToPrivateNetworkOpt,
) error {
	args := m.Called(ctx, zone, instance, privateNetwork, opts)
	return args.Error(0)
}

func (m *exoscaleClientMock) CopyTemplate(
	ctx context.Context,
	zone string,
	template *egoscale.Template,
	dst string,
) (*egoscale.Template, error) {
	args := m.Called(ctx, zone, template, dst)
	return args.Get(0).(*egoscale.Template), args.Error(1)
}

func (m *exoscaleClientMock) CreateInstance(
	ctx context.Context,
	zone string,
	instance *egoscale.Instance,
) (*egoscale.Instance, error) {
	args := m.Called(ctx, zone, instance)
	return args.Get(0).(*egoscale.Instance), args.Error(1)
}

func (m *exoscaleClientMock) CreateInstanceSnapshot(
	ctx context.Context,
	zone string,
	instance *egoscale.Instance,
) (*egoscale.Snapshot, error) {
	args := m.Called(ctx, zone, instance)
	return args.Get(0).(*egoscale.Snapshot), args.Error(1)
}

func (m *exoscaleClientMock) DeleteInstance(
	ctx context.Context,
	zone string,
	instance *egoscale.Instance,
) error {
	args := m.Called(ctx, zone, instance)
	return args.Error(0)
}

func (m *exoscaleClientMock) DeleteSSHKey(ctx context.Context, zone string, sshKey *egoscale.SSHKey) error {
	args := m.Called(ctx, zone, sshKey)
	return args.Error(0)
}

func (m *exoscaleClientMock) DeleteTemplate(ctx context.Context, zone string, template *egoscale.Template) error {
	args := m.Called(ctx, zone, template)
	return args.Error(0)
}

func (m *exoscaleClientMock) ExportSnapshot(
	ctx context.Context,
	zone string,
	snapshot *egoscale.Snapshot,
) (*egoscale.SnapshotExport, error) {
	args := m.Called(ctx, zone, snapshot)
	return args.Get(0).(*egoscale.SnapshotExport), args.Error(1)
}

func (m *exoscaleClientMock) FindInstanceType(
	ctx context.Context,
	zone string,
	x string,
) (*egoscale.InstanceType, error) {
	args := m.Called(ctx, zone, x)
	return args.Get(0).(*egoscale.InstanceType), args.Error(1)
}

func (m *exoscaleClientMock) FindPrivateNetwork(
	ctx context.Context,
	zone string,
	x string,
) (*egoscale.PrivateNetwork, error) {
	args := m.Called(ctx, zone, x)
	return args.Get(0).(*egoscale.PrivateNetwork), args.Error(1)
}

func (m *exoscaleClientMock) FindSecurityGroup(
	ctx context.Context,
	zone string,
	x string,
) (*egoscale.SecurityGroup, error) {
	args := m.Called(ctx, zone, x)
	return args.Get(0).(*egoscale.SecurityGroup), args.Error(1)
}

func (m *exoscaleClientMock) GetTemplate(ctx context.Context, zone string, id string) (*egoscale.Template, error) {
	args := m.Called(ctx, zone, id)
	return args.Get(0).(*egoscale.Template), args.Error(1)
}

func (m *exoscaleClientMock) ListTemplates(
	ctx context.Context,
	zone string,
	opts ...egoscale.ListTemplatesOpt,
) ([]*egoscale.Template, error) {
	args := m.Called(ctx, zone, opts)
	return args.Get(0).([]*egoscale.Template), args.Error(1)
}

func (m *exoscaleClientMock) RegisterSSHKey(
	ctx context.Context,
	zone string,
	name string,
	publicKey string,
) (*egoscale.SSHKey, error) {
	args := m.Called(ctx, zone, name, publicKey)
	return args.Get(0).(*egoscale.SSHKey), args.Error(1)
}

func (m *exoscaleClientMock) RegisterTemplate(
	ctx context.Context,
	zone string,
	template *egoscale.Template,
) (*egoscale.Template, error) {
	args := m.Called(ctx, zone, template)
	return args.Get(0).(*egoscale.Template), args.Error(1)
}

func (m *exoscaleClientMock) StopInstance(ctx context.Context, zone string, instance *egoscale.Instance) error {
	args := m.Called(ctx, zone, instance)
	return args.Error(0)
}
