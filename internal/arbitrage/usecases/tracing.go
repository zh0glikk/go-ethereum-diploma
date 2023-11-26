package usecases

import "github.com/ethereum/go-ethereum/eth/tracers"

func setTraceCallConfigDefaultTracer(config *tracers.TraceCallConfig) *tracers.TraceCallConfig {
	if config == nil {
		config = &tracers.TraceCallConfig{}
	}

	if config.Tracer == nil {
		tracer := "callTracerParity"
		config.Tracer = &tracer
	}

	return config
}

func getTraceConfigFromTraceCallConfig(config *tracers.TraceCallConfig) *tracers.TraceConfig {
	var traceConfig *tracers.TraceConfig
	if config != nil {
		traceConfig = &tracers.TraceConfig{
			Config:  config.Config,
			Tracer:  config.Tracer,
			Timeout: config.Timeout,
			Reexec:  config.Reexec,
		}
	}
	return traceConfig
}

func decorateNestedTraceResponse(res interface{}, tracer string) interface{} {
	out := map[string]interface{}{}
	if tracer == "callTracerParity" {
		out["trace"] = res
	} else if tracer == "stateDiffTracer" {
		out["stateDiff"] = res
	} else {
		return res
	}
	return out
}

func decorateResponse(res interface{}, config *tracers.TraceConfig) (interface{}, error) {
	if config != nil && config.Tracer != nil {
		return decorateNestedTraceResponse(res, *config.Tracer), nil
	}
	return res, nil
}
