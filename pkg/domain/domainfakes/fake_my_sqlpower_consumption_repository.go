// Code generated by counterfeiter. DO NOT EDIT.
package domainfakes

import (
	"sync"
	"time"

	"github.com/jeffleon1/consumption-ms/pkg/domain"
)

type FakeMySQLPowerConsumptionRepository struct {
	CreatePowerConsumptionRecordsStub        func([]*domain.UserConsumption) error
	createPowerConsumptionRecordsMutex       sync.RWMutex
	createPowerConsumptionRecordsArgsForCall []struct {
		arg1 []*domain.UserConsumption
	}
	createPowerConsumptionRecordsReturns struct {
		result1 error
	}
	createPowerConsumptionRecordsReturnsOnCall map[int]struct {
		result1 error
	}
	GetConsumptionByMeterIDAndWindowTimeStub        func(time.Time, time.Time, int) ([]domain.UserConsumption, error)
	getConsumptionByMeterIDAndWindowTimeMutex       sync.RWMutex
	getConsumptionByMeterIDAndWindowTimeArgsForCall []struct {
		arg1 time.Time
		arg2 time.Time
		arg3 int
	}
	getConsumptionByMeterIDAndWindowTimeReturns struct {
		result1 []domain.UserConsumption
		result2 error
	}
	getConsumptionByMeterIDAndWindowTimeReturnsOnCall map[int]struct {
		result1 []domain.UserConsumption
		result2 error
	}
	ModelMigrationStub        func() error
	modelMigrationMutex       sync.RWMutex
	modelMigrationArgsForCall []struct {
	}
	modelMigrationReturns struct {
		result1 error
	}
	modelMigrationReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecords(arg1 []*domain.UserConsumption) error {
	var arg1Copy []*domain.UserConsumption
	if arg1 != nil {
		arg1Copy = make([]*domain.UserConsumption, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.createPowerConsumptionRecordsMutex.Lock()
	ret, specificReturn := fake.createPowerConsumptionRecordsReturnsOnCall[len(fake.createPowerConsumptionRecordsArgsForCall)]
	fake.createPowerConsumptionRecordsArgsForCall = append(fake.createPowerConsumptionRecordsArgsForCall, struct {
		arg1 []*domain.UserConsumption
	}{arg1Copy})
	stub := fake.CreatePowerConsumptionRecordsStub
	fakeReturns := fake.createPowerConsumptionRecordsReturns
	fake.recordInvocation("CreatePowerConsumptionRecords", []interface{}{arg1Copy})
	fake.createPowerConsumptionRecordsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecordsCallCount() int {
	fake.createPowerConsumptionRecordsMutex.RLock()
	defer fake.createPowerConsumptionRecordsMutex.RUnlock()
	return len(fake.createPowerConsumptionRecordsArgsForCall)
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecordsCalls(stub func([]*domain.UserConsumption) error) {
	fake.createPowerConsumptionRecordsMutex.Lock()
	defer fake.createPowerConsumptionRecordsMutex.Unlock()
	fake.CreatePowerConsumptionRecordsStub = stub
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecordsArgsForCall(i int) []*domain.UserConsumption {
	fake.createPowerConsumptionRecordsMutex.RLock()
	defer fake.createPowerConsumptionRecordsMutex.RUnlock()
	argsForCall := fake.createPowerConsumptionRecordsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecordsReturns(result1 error) {
	fake.createPowerConsumptionRecordsMutex.Lock()
	defer fake.createPowerConsumptionRecordsMutex.Unlock()
	fake.CreatePowerConsumptionRecordsStub = nil
	fake.createPowerConsumptionRecordsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMySQLPowerConsumptionRepository) CreatePowerConsumptionRecordsReturnsOnCall(i int, result1 error) {
	fake.createPowerConsumptionRecordsMutex.Lock()
	defer fake.createPowerConsumptionRecordsMutex.Unlock()
	fake.CreatePowerConsumptionRecordsStub = nil
	if fake.createPowerConsumptionRecordsReturnsOnCall == nil {
		fake.createPowerConsumptionRecordsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createPowerConsumptionRecordsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTime(arg1 time.Time, arg2 time.Time, arg3 int) ([]domain.UserConsumption, error) {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.Lock()
	ret, specificReturn := fake.getConsumptionByMeterIDAndWindowTimeReturnsOnCall[len(fake.getConsumptionByMeterIDAndWindowTimeArgsForCall)]
	fake.getConsumptionByMeterIDAndWindowTimeArgsForCall = append(fake.getConsumptionByMeterIDAndWindowTimeArgsForCall, struct {
		arg1 time.Time
		arg2 time.Time
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.GetConsumptionByMeterIDAndWindowTimeStub
	fakeReturns := fake.getConsumptionByMeterIDAndWindowTimeReturns
	fake.recordInvocation("GetConsumptionByMeterIDAndWindowTime", []interface{}{arg1, arg2, arg3})
	fake.getConsumptionByMeterIDAndWindowTimeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTimeCallCount() int {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.RLock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.RUnlock()
	return len(fake.getConsumptionByMeterIDAndWindowTimeArgsForCall)
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTimeCalls(stub func(time.Time, time.Time, int) ([]domain.UserConsumption, error)) {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.Lock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.Unlock()
	fake.GetConsumptionByMeterIDAndWindowTimeStub = stub
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTimeArgsForCall(i int) (time.Time, time.Time, int) {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.RLock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.RUnlock()
	argsForCall := fake.getConsumptionByMeterIDAndWindowTimeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTimeReturns(result1 []domain.UserConsumption, result2 error) {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.Lock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.Unlock()
	fake.GetConsumptionByMeterIDAndWindowTimeStub = nil
	fake.getConsumptionByMeterIDAndWindowTimeReturns = struct {
		result1 []domain.UserConsumption
		result2 error
	}{result1, result2}
}

func (fake *FakeMySQLPowerConsumptionRepository) GetConsumptionByMeterIDAndWindowTimeReturnsOnCall(i int, result1 []domain.UserConsumption, result2 error) {
	fake.getConsumptionByMeterIDAndWindowTimeMutex.Lock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.Unlock()
	fake.GetConsumptionByMeterIDAndWindowTimeStub = nil
	if fake.getConsumptionByMeterIDAndWindowTimeReturnsOnCall == nil {
		fake.getConsumptionByMeterIDAndWindowTimeReturnsOnCall = make(map[int]struct {
			result1 []domain.UserConsumption
			result2 error
		})
	}
	fake.getConsumptionByMeterIDAndWindowTimeReturnsOnCall[i] = struct {
		result1 []domain.UserConsumption
		result2 error
	}{result1, result2}
}

func (fake *FakeMySQLPowerConsumptionRepository) ModelMigration() error {
	fake.modelMigrationMutex.Lock()
	ret, specificReturn := fake.modelMigrationReturnsOnCall[len(fake.modelMigrationArgsForCall)]
	fake.modelMigrationArgsForCall = append(fake.modelMigrationArgsForCall, struct {
	}{})
	stub := fake.ModelMigrationStub
	fakeReturns := fake.modelMigrationReturns
	fake.recordInvocation("ModelMigration", []interface{}{})
	fake.modelMigrationMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeMySQLPowerConsumptionRepository) ModelMigrationCallCount() int {
	fake.modelMigrationMutex.RLock()
	defer fake.modelMigrationMutex.RUnlock()
	return len(fake.modelMigrationArgsForCall)
}

func (fake *FakeMySQLPowerConsumptionRepository) ModelMigrationCalls(stub func() error) {
	fake.modelMigrationMutex.Lock()
	defer fake.modelMigrationMutex.Unlock()
	fake.ModelMigrationStub = stub
}

func (fake *FakeMySQLPowerConsumptionRepository) ModelMigrationReturns(result1 error) {
	fake.modelMigrationMutex.Lock()
	defer fake.modelMigrationMutex.Unlock()
	fake.ModelMigrationStub = nil
	fake.modelMigrationReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMySQLPowerConsumptionRepository) ModelMigrationReturnsOnCall(i int, result1 error) {
	fake.modelMigrationMutex.Lock()
	defer fake.modelMigrationMutex.Unlock()
	fake.ModelMigrationStub = nil
	if fake.modelMigrationReturnsOnCall == nil {
		fake.modelMigrationReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.modelMigrationReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMySQLPowerConsumptionRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createPowerConsumptionRecordsMutex.RLock()
	defer fake.createPowerConsumptionRecordsMutex.RUnlock()
	fake.getConsumptionByMeterIDAndWindowTimeMutex.RLock()
	defer fake.getConsumptionByMeterIDAndWindowTimeMutex.RUnlock()
	fake.modelMigrationMutex.RLock()
	defer fake.modelMigrationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMySQLPowerConsumptionRepository) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ domain.MySQLPowerConsumptionRepository = new(FakeMySQLPowerConsumptionRepository)