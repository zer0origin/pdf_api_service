package eureka

import (
	"fmt"
	"github.com/ArthurHlt/go-eureka-client/eureka"
)

type Eureka struct {
}

func (Eureka) JoinEureka(hostname, ip, app string, port int) error {
	client := eureka.NewClient([]string{
		"http://127.0.0.1:8761/eureka", //From a spring boot based eureka server
		// add others servers here
	})
	instance := eureka.NewInstanceInfo(hostname, ip, app, port, 30, false) //Create a new instance to register
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	instance.Metadata.Map["foo"] = "bar" //add metadata for example
	err := client.RegisterInstance("go-backend", instance)
	if err != nil {
		return fmt.Errorf("failed to join eureka: %s", err)
	}

	return nil
}
