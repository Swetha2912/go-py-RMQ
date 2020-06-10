//+build wireinject

package main

import (
	"sample_rmq/common/infrastructures"
	"sample_rmq/common/utilities"
	"sample_rmq/gateway/controllers"

	"github.com/google/wire"
)

func InitializeContainer() Container {
	wire.Build(
		controllers.ProvideInferenceController,
		infrastructures.ProvideFiberRouter,
		infrastructures.ProvideRMQRouter,
		utilities.ProvideGoValidator,
		infrastructures.ProvideZmqBroker,
		ProvideContainer)
	return Container{}
}
