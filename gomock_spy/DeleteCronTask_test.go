package gomock_spy

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCronTask_Start(t *testing.T) {
	tests := []struct {
		name                string
		expectedDeleteIndex string
	}{
		{
			name:                "successful",
			expectedDeleteIndex: "test_index",
		},
		{
			name:                "fail",
			expectedDeleteIndex: "dev_index_failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var actualDeleteIndex string
			spyDeleteService := newSpyDeleteService(t, &actualDeleteIndex)

			cron := NewDeleteCronTask(spyDeleteService)
			cron.Start()

			assert.Equal(t, tt.expectedDeleteIndex, actualDeleteIndex)
		})
	}
}

func newSpyDeleteService(t *testing.T, actualDeleteIndex *string) *MockDeleteService {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish() // 註解此行, 就不會判斷 Times(1) 的條件限制
	// ctl.Finish()：进行 mock 用例的期望值断言，一般会使用 defer 延迟执行，以防止我们忘记这一操作
	svc := NewMockDeleteService(ctrl)

	svc.EXPECT().
		Delete(gomock.Not(gomock.Eq("test_index"))). // 永遠不會進入這個條件, 因為 DeleteCronTask 只會輸入 "test_index"
		Do(func(indexName string) {
			*actualDeleteIndex = indexName + "_failed"
		}).Times(1)

	svc.EXPECT().
		Delete(gomock.Eq("test_index")).
		Do(func(indexName string) {
			*actualDeleteIndex = indexName
		}).Times(1)

	return svc
}
