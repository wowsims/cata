"use strict";
var commonjsGlobal = typeof globalThis !== "undefined" ? globalThis : typeof window !== "undefined" ? window : typeof global !== "undefined" ? global : typeof self !== "undefined" ? self : {};
var freeGlobal$1 = typeof commonjsGlobal == "object" && commonjsGlobal && commonjsGlobal.Object === Object && commonjsGlobal;
var _freeGlobal = freeGlobal$1;
var freeGlobal = _freeGlobal;
var freeSelf = typeof self == "object" && self && self.Object === Object && self;
var root$8 = freeGlobal || freeSelf || Function("return this")();
var _root = root$8;
var root$7 = _root;
var Symbol$4 = root$7.Symbol;
var _Symbol = Symbol$4;
var Symbol$3 = _Symbol;
var objectProto$3 = Object.prototype;
var hasOwnProperty$2 = objectProto$3.hasOwnProperty;
var nativeObjectToString$1 = objectProto$3.toString;
var symToStringTag$1 = Symbol$3 ? Symbol$3.toStringTag : void 0;
function getRawTag$1(value) {
  var isOwn = hasOwnProperty$2.call(value, symToStringTag$1), tag = value[symToStringTag$1];
  try {
    value[symToStringTag$1] = void 0;
    var unmasked = true;
  } catch (e) {
  }
  var result = nativeObjectToString$1.call(value);
  if (unmasked) {
    if (isOwn) {
      value[symToStringTag$1] = tag;
    } else {
      delete value[symToStringTag$1];
    }
  }
  return result;
}
var _getRawTag = getRawTag$1;
var objectProto$2 = Object.prototype;
var nativeObjectToString = objectProto$2.toString;
function objectToString$1(value) {
  return nativeObjectToString.call(value);
}
var _objectToString = objectToString$1;
var Symbol$2 = _Symbol, getRawTag = _getRawTag, objectToString = _objectToString;
var nullTag = "[object Null]", undefinedTag = "[object Undefined]";
var symToStringTag = Symbol$2 ? Symbol$2.toStringTag : void 0;
function baseGetTag$3(value) {
  if (value == null) {
    return value === void 0 ? undefinedTag : nullTag;
  }
  return symToStringTag && symToStringTag in Object(value) ? getRawTag(value) : objectToString(value);
}
var _baseGetTag = baseGetTag$3;
function isObject$2(value) {
  var type = typeof value;
  return value != null && (type == "object" || type == "function");
}
var isObject_1 = isObject$2;
var baseGetTag$2 = _baseGetTag, isObject$1 = isObject_1;
var asyncTag = "[object AsyncFunction]", funcTag = "[object Function]", genTag = "[object GeneratorFunction]", proxyTag = "[object Proxy]";
function isFunction$1(value) {
  if (!isObject$1(value)) {
    return false;
  }
  var tag = baseGetTag$2(value);
  return tag == funcTag || tag == genTag || tag == asyncTag || tag == proxyTag;
}
var isFunction_1 = isFunction$1;
var root$6 = _root;
var coreJsData$1 = root$6["__core-js_shared__"];
var _coreJsData = coreJsData$1;
var coreJsData = _coreJsData;
var maskSrcKey = function() {
  var uid = /[^.]+$/.exec(coreJsData && coreJsData.keys && coreJsData.keys.IE_PROTO || "");
  return uid ? "Symbol(src)_1." + uid : "";
}();
function isMasked$1(func) {
  return !!maskSrcKey && maskSrcKey in func;
}
var _isMasked = isMasked$1;
var funcProto$1 = Function.prototype;
var funcToString$1 = funcProto$1.toString;
function toSource$2(func) {
  if (func != null) {
    try {
      return funcToString$1.call(func);
    } catch (e) {
    }
    try {
      return func + "";
    } catch (e) {
    }
  }
  return "";
}
var _toSource = toSource$2;
var isFunction = isFunction_1, isMasked = _isMasked, isObject = isObject_1, toSource$1 = _toSource;
var reRegExpChar = /[\\^$.*+?()[\]{}|]/g;
var reIsHostCtor = /^\[object .+?Constructor\]$/;
var funcProto = Function.prototype, objectProto$1 = Object.prototype;
var funcToString = funcProto.toString;
var hasOwnProperty$1 = objectProto$1.hasOwnProperty;
var reIsNative = RegExp(
  "^" + funcToString.call(hasOwnProperty$1).replace(reRegExpChar, "\\$&").replace(/hasOwnProperty|(function).*?(?=\\\()| for .+?(?=\\\])/g, "$1.*?") + "$"
);
function baseIsNative$1(value) {
  if (!isObject(value) || isMasked(value)) {
    return false;
  }
  var pattern = isFunction(value) ? reIsNative : reIsHostCtor;
  return pattern.test(toSource$1(value));
}
var _baseIsNative = baseIsNative$1;
function getValue$1(object, key) {
  return object == null ? void 0 : object[key];
}
var _getValue = getValue$1;
var baseIsNative = _baseIsNative, getValue = _getValue;
function getNative$7(object, key) {
  var value = getValue(object, key);
  return baseIsNative(value) ? value : void 0;
}
var _getNative = getNative$7;
var getNative$6 = _getNative, root$5 = _root;
var Map$1 = getNative$6(root$5, "Map");
var _Map = Map$1;
var getNative$5 = _getNative;
getNative$5(Object, "create");
var getNative$4 = _getNative;
(function() {
  try {
    var func = getNative$4(Object, "defineProperty");
    func({}, "", {});
    return func;
  } catch (e) {
  }
})();
function isObjectLike$2(value) {
  return value != null && typeof value == "object";
}
var isObjectLike_1 = isObjectLike$2;
var baseGetTag$1 = _baseGetTag, isObjectLike$1 = isObjectLike_1;
var argsTag = "[object Arguments]";
function baseIsArguments$1(value) {
  return isObjectLike$1(value) && baseGetTag$1(value) == argsTag;
}
var _baseIsArguments = baseIsArguments$1;
var baseIsArguments = _baseIsArguments, isObjectLike = isObjectLike_1;
var objectProto = Object.prototype;
var hasOwnProperty = objectProto.hasOwnProperty;
var propertyIsEnumerable = objectProto.propertyIsEnumerable;
baseIsArguments(/* @__PURE__ */ function() {
  return arguments;
}()) ? baseIsArguments : function(value) {
  return isObjectLike(value) && hasOwnProperty.call(value, "callee") && !propertyIsEnumerable.call(value, "callee");
};
var isBuffer = { exports: {} };
function stubFalse() {
  return false;
}
var stubFalse_1 = stubFalse;
isBuffer.exports;
(function(module, exports) {
  var root2 = _root, stubFalse2 = stubFalse_1;
  var freeExports = exports && !exports.nodeType && exports;
  var freeModule = freeExports && true && module && !module.nodeType && module;
  var moduleExports = freeModule && freeModule.exports === freeExports;
  var Buffer = moduleExports ? root2.Buffer : void 0;
  var nativeIsBuffer = Buffer ? Buffer.isBuffer : void 0;
  var isBuffer2 = nativeIsBuffer || stubFalse2;
  module.exports = isBuffer2;
})(isBuffer, isBuffer.exports);
isBuffer.exports;
var _nodeUtil = { exports: {} };
_nodeUtil.exports;
(function(module, exports) {
  var freeGlobal2 = _freeGlobal;
  var freeExports = exports && !exports.nodeType && exports;
  var freeModule = freeExports && true && module && !module.nodeType && module;
  var moduleExports = freeModule && freeModule.exports === freeExports;
  var freeProcess = moduleExports && freeGlobal2.process;
  var nodeUtil2 = function() {
    try {
      var types = freeModule && freeModule.require && freeModule.require("util").types;
      if (types) {
        return types;
      }
      return freeProcess && freeProcess.binding && freeProcess.binding("util");
    } catch (e) {
    }
  }();
  module.exports = nodeUtil2;
})(_nodeUtil, _nodeUtil.exports);
var _nodeUtilExports = _nodeUtil.exports;
var nodeUtil$2 = _nodeUtilExports;
nodeUtil$2 && nodeUtil$2.isTypedArray;
var _cloneBuffer = { exports: {} };
_cloneBuffer.exports;
(function(module, exports) {
  var root2 = _root;
  var freeExports = exports && !exports.nodeType && exports;
  var freeModule = freeExports && true && module && !module.nodeType && module;
  var moduleExports = freeModule && freeModule.exports === freeExports;
  var Buffer = moduleExports ? root2.Buffer : void 0, allocUnsafe = Buffer ? Buffer.allocUnsafe : void 0;
  function cloneBuffer(buffer, isDeep) {
    if (isDeep) {
      return buffer.slice();
    }
    var length = buffer.length, result = allocUnsafe ? allocUnsafe(length) : new buffer.constructor(length);
    buffer.copy(result);
    return result;
  }
  module.exports = cloneBuffer;
})(_cloneBuffer, _cloneBuffer.exports);
_cloneBuffer.exports;
var getNative$3 = _getNative, root$4 = _root;
var DataView$1 = getNative$3(root$4, "DataView");
var _DataView = DataView$1;
var getNative$2 = _getNative, root$3 = _root;
var Promise$2 = getNative$2(root$3, "Promise");
var _Promise = Promise$2;
var getNative$1 = _getNative, root$2 = _root;
var Set$1 = getNative$1(root$2, "Set");
var _Set = Set$1;
var getNative = _getNative, root$1 = _root;
var WeakMap$1 = getNative(root$1, "WeakMap");
var _WeakMap = WeakMap$1;
var DataView = _DataView, Map = _Map, Promise$1 = _Promise, Set = _Set, WeakMap = _WeakMap, baseGetTag = _baseGetTag, toSource = _toSource;
var mapTag = "[object Map]", objectTag = "[object Object]", promiseTag = "[object Promise]", setTag = "[object Set]", weakMapTag = "[object WeakMap]";
var dataViewTag = "[object DataView]";
var dataViewCtorString = toSource(DataView), mapCtorString = toSource(Map), promiseCtorString = toSource(Promise$1), setCtorString = toSource(Set), weakMapCtorString = toSource(WeakMap);
var getTag = baseGetTag;
if (DataView && getTag(new DataView(new ArrayBuffer(1))) != dataViewTag || Map && getTag(new Map()) != mapTag || Promise$1 && getTag(Promise$1.resolve()) != promiseTag || Set && getTag(new Set()) != setTag || WeakMap && getTag(new WeakMap()) != weakMapTag) {
  getTag = function(value) {
    var result = baseGetTag(value), Ctor = result == objectTag ? value.constructor : void 0, ctorString = Ctor ? toSource(Ctor) : "";
    if (ctorString) {
      switch (ctorString) {
        case dataViewCtorString:
          return dataViewTag;
        case mapCtorString:
          return mapTag;
        case promiseCtorString:
          return promiseTag;
        case setCtorString:
          return setTag;
        case weakMapCtorString:
          return weakMapTag;
      }
    }
    return result;
  };
}
var root = _root;
root.Uint8Array;
var Symbol$1 = _Symbol;
var symbolProto = Symbol$1 ? Symbol$1.prototype : void 0;
symbolProto ? symbolProto.valueOf : void 0;
var nodeUtil$1 = _nodeUtilExports;
nodeUtil$1 && nodeUtil$1.isMap;
var nodeUtil = _nodeUtilExports;
nodeUtil && nodeUtil.isSet;
const noop = () => {
};
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
class WorkerInterface {
  constructor(handlers) {
    this._workerId = "";
    this.handlers = handlers;
    addEventListener("message", async ({ data }) => {
      const { id, msg, inputData } = data;
      if (msg === "setID") {
        this._workerId = id;
        this.postMessage({ msg: "idConfirm" });
        return;
      }
      const handlerFunc = this.handlers?.[msg];
      if (!handlerFunc) {
        console.error(`Request msg: ${msg}, id: ${id}, is not handled!`);
        return;
      }
      const progressCallback = (prog) => {
        this.postMessage({
          msg: "progress",
          id: `${id}progress`,
          outputData: prog
        });
      };
      const outputData = await handlerFunc(inputData, progressCallback, id, msg);
      this.postMessage({ msg, id, outputData });
    });
  }
  postMessage(m) {
    postMessage(m);
  }
  get workerId() {
    return this._workerId;
  }
  /**
   * Tell UI that the worker is ready.
   * @param isWasm true if worker is using wasm.
   */
  ready(isWasm) {
    this.postMessage({ msg: "ready", outputData: new Uint8Array([+isWasm]) });
  }
}
const defaultRequestOptions = {
  method: "POST",
  headers: {
    "Content-Type": "application/x-protobuf"
  }
};
const setupHttpWorker = (baseURL) => {
  const makeHttpApiRequest = (endPoint, inputData, requestId) => fetch(`${baseURL}/${endPoint}?requestId=${requestId}`, {
    ...defaultRequestOptions,
    body: inputData
  });
  const syncHandler = async (inputData, _, id, msg) => {
    const response = await makeHttpApiRequest(msg, inputData, id);
    const ab = await response.arrayBuffer();
    return new Uint8Array(ab);
  };
  const asyncHandler = async (inputData, progress, id, msg) => {
    const asyncApiResult = await syncHandler(inputData, noop, id, msg);
    let outputData = new Uint8Array();
    while (true) {
      const progressResponse = await makeHttpApiRequest("asyncProgress", asyncApiResult, id);
      if ([204, 404].includes(progressResponse.status)) {
        break;
      }
      const ab = await progressResponse.arrayBuffer();
      outputData = new Uint8Array(ab);
      progress(outputData);
      await sleep(500);
    }
    return outputData;
  };
  const noWasmConcurrency = (inputData, progress, msg) => {
    const errmsg = `Tried to use ${msg} while using a http worker! This is only supported for wasm!`;
    console.error(errmsg);
    return new Uint8Array();
  };
  new WorkerInterface({
    bulkSimAsync: asyncHandler,
    bulkSimCombos: syncHandler,
    computeStats: syncHandler,
    computeStatsJson: syncHandler,
    raidSim: syncHandler,
    raidSimJson: syncHandler,
    raidSimAsync: asyncHandler,
    statWeights: syncHandler,
    statWeightsAsync: asyncHandler,
    statWeightRequests: syncHandler,
    statWeightCompute: syncHandler,
    raidSimRequestSplit: noWasmConcurrency,
    raidSimResultCombination: noWasmConcurrency,
    abortById: syncHandler
  }).ready(false);
};
setupHttpWorker("http://localhost:3333");
