"use strict";
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
      const outputData = await handlerFunc(inputData, progressCallback, msg);
      this.postMessage({ msg, id, outputData });
    });
  }
  postMessage(m) {
    postMessage(m);
  }
  get workerId() {
    return this._workerId;
  }
  /** Tell UI that the worker is ready. */
  ready() {
    this.postMessage({ msg: "ready" });
  }
}
const defaultRequestOptions = {
  method: "POST",
  headers: {
    "Content-Type": "application/x-protobuf"
  }
};
const setupHttpWorker = (baseURL) => {
  const makeHttpApiRequest = (endPoint, inputData) => fetch(`${baseURL}/${endPoint}`, {
    ...defaultRequestOptions,
    body: inputData
  });
  const syncHandler = async (inputData, _, msg) => {
    const response = await makeHttpApiRequest(msg, inputData);
    const ab = await response.arrayBuffer();
    return new Uint8Array(ab);
  };
  const asyncHandler = async (inputData, progress, msg) => {
    const asyncApiResult = await syncHandler(inputData, noop, msg);
    let outputData = new Uint8Array();
    while (true) {
      const progressResponse = await makeHttpApiRequest("asyncProgress", asyncApiResult);
      if ([204, 404].includes(progressResponse.status)) {
        break;
      }
      const ab = await progressResponse.arrayBuffer();
      outputData = new Uint8Array(ab);
      progress?.(outputData);
      await sleep(500);
    }
    return outputData;
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
    statWeightsAsync: asyncHandler
  }).ready();
};
setupHttpWorker("");
