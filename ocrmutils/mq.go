package ocrmutils

import "github.com/streadway/amqp"

func GetRetryCount(delivery amqp.Delivery) int {
	count := 1
	_, ok := delivery.Headers["retryCount"]
	if ok {
		count = int(delivery.Headers["retryCount"].(int32) + 1)
	}
	return count
}

func GetPublishingWithRetryCount(retryCount int, data []byte) amqp.Publishing {
	if retryCount == 0 {
		return amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		}
	} else {
		header := make(map[string]interface{})
		header["retryCount"] = retryCount
		return amqp.Publishing{
			Headers:      header,
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		}
	}
}

func GetPublishingWithOutHeader(data []byte) amqp.Publishing {
	return amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         data,
	}
}

func GetPublishingWithHeader(headers map[string]interface{}, data []byte) amqp.Publishing {
	return amqp.Publishing{
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         data,
	}
}

func GetHeaderForPublishing(header string, value interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result[header] = value
	return result
}
