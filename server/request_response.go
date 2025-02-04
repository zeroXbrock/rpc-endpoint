package server

import (
	"encoding/json"
	"net/http"

	"github.com/flashbots/rpc-endpoint/types"
)

func (r *RpcRequestHandler) writeHeaderContentTypeJson() {
	(*r.respw).Header().Set("Content-Type", "application/json")
}

func (r *RpcRequestHandler) _writeRpcResponse(res *types.JsonRpcResponse) {

	// If the request is single and not batch
	// Write content type
	r.writeHeaderContentTypeJson() // Set content type to json

	// Choose httpStatusCode based on json-rpc error code
	statusCode := http.StatusOK
	if res.Error != nil {
		// TODO(Note): http.StatusUnauthorized is not mapped
		switch res.Error.Code {
		case types.JsonRpcInvalidRequest, types.JsonRpcInvalidParams:
			statusCode = http.StatusBadRequest
		case types.JsonRpcMethodNotFound:
			statusCode = http.StatusNotFound
		case types.JsonRpcInternalError, types.JsonRpcParseError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	}
	(*r.respw).WriteHeader(statusCode)

	// Write response
	if err := json.NewEncoder(*r.respw).Encode(res); err != nil {
		r.logger.logError("failed writing rpc response: %v", err)
		(*r.respw).WriteHeader(http.StatusInternalServerError)
	}
}

func (r *RpcRequestHandler) _writeRpcBatchResponse(res []*types.JsonRpcResponse) {
	r.writeHeaderContentTypeJson() // Set content type to json
	(*r.respw).WriteHeader(http.StatusOK)
	// Write response
	if err := json.NewEncoder(*r.respw).Encode(res); err != nil {
		r.logger.logError("failed writing rpc response: %v", err)
		(*r.respw).WriteHeader(http.StatusInternalServerError)
	}
}
