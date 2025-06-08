// These utils are copied from core/utils.ts to keep workers completely isolated and prevent a circular dependency when building chunks.
// The worker code needs to be kept completely isolated.
// eslint-disable-next-line @typescript-eslint/no-empty-function
export const noop = () => {};
export const sleep = (ms: number) => new Promise(r => setTimeout(r, ms));
