package rabbit

import "go-practice/internal/common/config"

var rq []config.QueueConfig

func (c *Client) getRegisteredQueue() []config.QueueConfig {
	if rq != nil {
		return rq
	}
	rq = append(rq, c.queuesConfig.Order.OrderCreated)

	return rq
}
