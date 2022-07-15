package cache

import (
	"fmt"
	"sync"

	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	guildDeleteStatusMap map[int64]struct{} = make(map[int64]struct{}, 10000)
	guildDeleteStatusMu  *sync.RWMutex      = new(sync.RWMutex)
)

func IsGuildDelete(guildId int64) (bool, error) {
	var (
		errout error
		isHit  bool = false
	)
	for {
		guildDeleteStatusMu.RLock()
		_, ok := guildDeleteStatusMap[guildId]
		guildDeleteStatusMu.RUnlock()
		if ok {
			isHit = true
			errout = nil
			break
		}

		gds, err := service.GetGuildInfosWithField([]int64{guildId}, []string{"deleted_at"})
		if err != nil {
			isHit = false
			errout = err
			break
		}

		if len(gds) != 1 {
			isHit = false
			errout = fmt.Errorf("not found guild delete info guild_id:%d", guildId)
			break
		}

		if gds[0].DeletedAt.Valid {
			guildDeleteStatusMu.Lock()
			guildDeleteStatusMap[guildId] = struct{}{}
			guildDeleteStatusMu.Unlock()

			isHit = true
			errout = nil

			break
		}

		break
	}

	return isHit, errout
}

func UpdateGuildDeleteStatus(datas map[int64]struct{}) int {
	guildDeleteStatusMu.Lock()
	defer guildDeleteStatusMu.Unlock()
	for k, _ := range datas {
		guildDeleteStatusMap[k] = struct{}{}
	}
	length := len(guildDeleteStatusMap)
	return length
}
