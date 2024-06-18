//go:build wasm
// +build wasm

package main

import (
	"log"
	"runtime/debug"
	"syscall/js"

	"github.com/wowsims/cata/sim"
	"github.com/wowsims/cata/sim/core"
	proto "github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	protojson "google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	core.SetRunningInWasm()
	sim.RegisterAll()
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("computeStats", js.FuncOf(computeStats))
	js.Global().Set("computeStatsJson", js.FuncOf(computeStatsJson))
	js.Global().Set("raidSim", js.FuncOf(raidSim))
	js.Global().Set("raidSimJson", js.FuncOf(raidSimJson))
	js.Global().Set("raidSimAsync", js.FuncOf(raidSimAsync))
	js.Global().Set("raidSimRequestSplit", js.FuncOf(raidSimRequestSplit))
	js.Global().Set("raidSimResultCombination", js.FuncOf(raidSimResultCombination))
	js.Global().Set("statWeights", js.FuncOf(statWeights))
	js.Global().Set("statWeightsAsync", js.FuncOf(statWeightsAsync))
	js.Global().Set("statWeightRequests", js.FuncOf(statWeightRequests))
	js.Global().Set("statWeightCompute", js.FuncOf(statWeightCompute))
	js.Global().Set("bulkSimAsync", js.FuncOf(bulkSimAsync))
	js.Global().Set("abortById", js.FuncOf(abortById))
	js.Global().Call("wasmready")
	<-c
}

func computeStats(this js.Value, args []js.Value) (response interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}

			errStr += "\nStack Trace:\n" + string(debug.Stack())
			result := &proto.ComputeStatsResult{
				ErrorResult: errStr,
			}
			outbytes, err := googleProto.Marshal(result)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal error (%s) result: %s", errStr, err.Error())
				return
			}
			outArray := js.Global().Get("Uint8Array").New(len(outbytes))
			js.CopyBytesToJS(outArray, outbytes)
			response = outArray
		}
	}()
	csr := &proto.ComputeStatsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), csr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.ComputeStats(csr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	response = outArray
	return response
}

func computeStatsJson(this js.Value, args []js.Value) (response interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}

			errStr += "\nStack Trace:\n" + string(debug.Stack())
			result := &proto.ComputeStatsResult{
				ErrorResult: errStr,
			}

			output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
			}
			response = js.ValueOf(string(output))
		}
	}()
	csr := &proto.ComputeStatsRequest{}
	log.Printf("Compute stats request: %s", getArgsJson(args[0]))
	if err := protojson.Unmarshal(getArgsJson(args[0]), csr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.ComputeStats(csr)

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	response = js.ValueOf(string(output))
	return response
}

func raidSim(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.RunRaidSim(rsr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func raidSimJson(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := protojson.Unmarshal(getArgsJson(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.RunRaidSim(rsr)

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	response := js.ValueOf(string(output))
	return response
}

func raidSimAsync(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	reporter := make(chan *proto.ProgressMetrics, 100)

	go core.RunRaidSimAsync(rsr, reporter)
	go processAsyncProgress(args[1], reporter)
	return js.Undefined()
}

func statWeights(this js.Value, args []js.Value) interface{} {
	swr := &proto.StatWeightsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), swr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.StatWeights(swr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func statWeightsAsync(this js.Value, args []js.Value) interface{} {
	rsr := &proto.StatWeightsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	reporter := make(chan *proto.ProgressMetrics, 100)
	go core.StatWeightsAsync(rsr, reporter)
	go processAsyncProgress(args[1], reporter)
	return js.Undefined()
}

func statWeightRequests(this js.Value, args []js.Value) interface{} {
	req := &proto.StatWeightsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), req); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}

	res := core.StatWeightRequests(req)

	outbytes, err := googleProto.Marshal(res)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func statWeightCompute(this js.Value, args []js.Value) interface{} {
	req := &proto.StatWeightsCalcRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), req); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}

	res := core.StatWeightCompute(req)

	outbytes, err := googleProto.Marshal(res)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func bulkSimAsync(this js.Value, args []js.Value) interface{} {
	rsr := &proto.BulkSimRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	reporter := make(chan *proto.ProgressMetrics, 100)
	go core.RunBulkSimAsync(rsr, reporter)
	go processAsyncProgress(args[1], reporter)
	return js.Undefined()
}

func raidSimRequestSplit(this js.Value, args []js.Value) interface{} {
	splitRequest := &proto.RaidSimRequestSplitRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), splitRequest); err != nil {
		log.Printf("Failed to parse RaidSimRequestSplitRequest: %s", err)
		return nil
	}

	splitRes := core.SplitSimRequestForConcurrency(splitRequest.Request, splitRequest.SplitCount)

	outbytes, err := googleProto.Marshal(splitRes)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal RaidSimRequestSplitResult: %s", err.Error())
		return nil
	}
	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func raidSimResultCombination(this js.Value, args []js.Value) interface{} {
	combRequest := &proto.RaidSimResultCombinationRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), combRequest); err != nil {
		log.Printf("Failed to parse RaidSimResultCombinationRequest: %s", err)
		return nil
	}

	combineRes := func() (res *proto.RaidSimResult) {
		defer func() {
			if err := recover(); err != nil {
				errStr := ""
				switch errt := err.(type) {
				case string:
					errStr = errt
				case error:
					errStr = errt.Error()
				}
				errStr += "\nStack Trace:\n" + string(debug.Stack())
				res = &proto.RaidSimResult{Error: &proto.ErrorOutcome{Message: errStr}}
			}
		}()
		return core.CombineConcurrentSimResults(combRequest.Results, false)
	}()

	outbytes, err := googleProto.Marshal(combineRes)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal RaidSimResult: %s", err.Error())
		return nil
	}
	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func abortById(this js.Value, args []js.Value) interface{} {
	abortRequest := &proto.AbortRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), abortRequest); err != nil {
		log.Printf("Failed to parse AbortRequest: %s", err)
		return nil
	}

	success := simsignals.AbortById(abortRequest.RequestId)

	outbytes, err := googleProto.Marshal(&proto.AbortResponse{
		RequestId:    abortRequest.RequestId,
		WasTriggered: success,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to marshal AbortResponse: %s", err.Error())
		return nil
	}
	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

// Assumes args[0] is a Uint8Array
func getArgsBinary(value js.Value) []byte {
	data := make([]byte, value.Get("length").Int())
	js.CopyBytesToGo(data, value)
	return data
}

func getArgsJson(value js.Value) []byte {
	str := value.String()
	return []byte(str)
}

func processAsyncProgress(progFunc js.Value, reporter chan *proto.ProgressMetrics) {
	for {
		select {
		case progMetric, ok := <-reporter:
			if !ok {
				return
			}
			outbytes, err := googleProto.Marshal(progMetric)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
				return
			}

			outArray := js.Global().Get("Uint8Array").New(len(outbytes))
			js.CopyBytesToJS(outArray, outbytes)
			progFunc.Invoke(outArray)

			if progMetric.FinalWeightResult != nil || progMetric.FinalRaidResult != nil || progMetric.FinalBulkResult != nil {
				return
			}
		}
	}
}
