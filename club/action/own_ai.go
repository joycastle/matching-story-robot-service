package action

import "time"

func cycleTimeHandlerOwnAi(job *Job) (int64, error) {

	return time.Now().Unix() + 3600, nil
}
