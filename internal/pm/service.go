package pm

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"nmsappsrv/pkg/redis"
)

// Service contains PM business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new PM service.
func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

// ---------- PerformanceKpi ----------

func (s *Service) ListKPIs(tenancyId int) ([]PerformanceKpi, error) {
	return s.repo.FindKPIs(tenancyId)
}

func (s *Service) GetKPI(id string) (*PerformanceKpi, error) {
	return s.repo.FindKPIByID(id)
}

func (s *Service) CreateKPI(k *PerformanceKpi) error {
	return s.repo.CreateKPI(k)
}

func (s *Service) UpdateKPI(k *PerformanceKpi) error {
	return s.repo.UpdateKPI(k)
}

func (s *Service) DeleteKPI(id string) error {
	return s.repo.DeleteKPI(id)
}

// ---------- PerformanceKpiSet ----------

func (s *Service) ListKPISets(tenancyId int) ([]PerformanceKpiSet, error) {
	return s.repo.FindKPISets(tenancyId)
}

func (s *Service) CreateKPISet(set *PerformanceKpiSet) error {
	return s.repo.CreateKPISet(set)
}

// ---------- PerformanceKpiTemplate ----------

func (s *Service) ListKPITemplates(tenancyId int) ([]PerformanceKpiTemplate, error) {
	return s.repo.FindKPITemplates(tenancyId)
}

func (s *Service) CreateKPITemplate(t *PerformanceKpiTemplate) error {
	return s.repo.CreateKPITemplate(t)
}

func (s *Service) UpdateKPITemplate(t *PerformanceKpiTemplate) error {
	return s.repo.UpdateKPITemplate(t)
}

func (s *Service) DeleteKPITemplate(id int) error {
	return s.repo.DeleteKPITemplate(id)
}

// ---------- PMFileLog ----------

func (s *Service) ListPMFileLogs(tenancyId int, page, pageSize int) ([]PMFileLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindPMFileLogs(tenancyId, offset, pageSize)
}

// ---------- KpiAlarmTemplate ----------

func (s *Service) ListKPIAlarmTemplates(tenancyId int) ([]KpiAlarmTemplate, error) {
	return s.repo.FindKPIAlarmTemplates(tenancyId)
}

func (s *Service) CreateKPIAlarmTemplate(t *KpiAlarmTemplate) error {
	return s.repo.CreateKPIAlarmTemplate(t)
}

func (s *Service) UpdateKPIAlarmTemplate(t *KpiAlarmTemplate) error {
	return s.repo.UpdateKPIAlarmTemplate(t)
}

func (s *Service) DeleteKPIAlarmTemplate(id int) error {
	return s.repo.DeleteKPIAlarmTemplate(id)
}

// ---------- Dashboard ----------

func (s *Service) GetDashboardData(tenancyId int, startTime, endTime string) ([]DashboardPmStatisticData, error) {
	st, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, err
	}
	et, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, err
	}
	return s.repo.FindDashboardData(tenancyId, st, et)
}

// ---------- PDCPTraffic ----------

func (s *Service) GetPDCPTraffic(tenancyId int, startTime, endTime string) ([]PDCPTraffic, error) {
	st, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, err
	}
	et, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, err
	}
	return s.repo.FindPDCPTraffic(tenancyId, st, et)
}

// ---------- Dashboard: Device Online Info ----------

// GetDeviceOnlineInfo 统计 gNB/eNB/CPE 各自在线/离线设备数
func (s *Service) GetDeviceOnlineInfo(tenancyId int) (*DeviceOnlineInfoVO, error) {
	rows, err := s.repo.FindAllActiveElements(tenancyId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	vo := &DeviceOnlineInfoVO{}

	for _, row := range rows {
		online := redis.Exists(ctx, fmt.Sprintf("online_%d", row.NeNeid))
		dt := strVal(row.DeviceType)
		gen := strVal(row.Generation)

		switch {
		case dt == "enb" && gen == "NR":
			vo.GnbTotal++
			if online {
				vo.GnbOnline++
			} else {
				vo.GnbOffline++
			}
		case dt == "enb":
			vo.EnbTotal++
			if online {
				vo.EnbOnline++
			} else {
				vo.EnbOffline++
			}
		default:
			vo.CpeTotal++
			if online {
				vo.CpeOnline++
			} else {
				vo.CpeOffline++
			}
		}
	}
	return vo, nil
}

// GetProductTypeAndDeviceCount 按产品型号统计设备数量及在线情况
// mode: "all" 查全部租户, 否则按 tenancyId 过滤
func (s *Service) GetProductTypeAndDeviceCount(tenancyId int, mode string) ([]ProductTypeAndCount, error) {
	var rows []elementRow
	var err error
	if mode == "all" {
		rows, err = s.repo.FindAllActiveElementsAllTenants()
	} else {
		rows, err = s.repo.FindAllActiveElements(tenancyId)
	}
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// group by model_name
	type agg struct {
		count        int64
		onlineCount  int64
	}
	grouped := make(map[string]*agg)

	for _, row := range rows {
		modelName := strVal(row.ModelName)
		if modelName == "" {
			modelName = "Unknown"
		}
		a, ok := grouped[modelName]
		if !ok {
			a = &agg{}
			grouped[modelName] = a
		}
		a.count++
		if redis.Exists(ctx, fmt.Sprintf("online_%d", row.NeNeid)) {
			a.onlineCount++
		}
	}

	result := make([]ProductTypeAndCount, 0, len(grouped))
	for pt, a := range grouped {
		result = append(result, ProductTypeAndCount{
			ProductType:  pt,
			Count:        a.count,
			OnlineCount:  a.onlineCount,
			OfflineCount: a.count - a.onlineCount,
		})
	}
	return result, nil
}

// strVal safely dereference a *string, returning "" for nil.
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
