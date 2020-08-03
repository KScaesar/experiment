package gomock_spy

//go:generate mockgen -source=DeleteCronTask.go -destination=mock.go -package=gomock_spy -mock_names DeleteService=MockDeleteService
type DeleteService interface {
	Delete(indexName string)
}

type DeleteCronTask struct {
	svc DeleteService
}

func NewDeleteCronTask(svc DeleteService) *DeleteCronTask {
	return &DeleteCronTask{svc: svc}
}

func (cron *DeleteCronTask) Start() {
	cron.svc.Delete("test_index")
}
