export enum Phase {
	Phase1 = 1,
	Phase2,
	Phase3,
	Phase4,
	Phase5,
}

export const CURRENT_PHASE = Phase.Phase1;

// Github pages serves our site under the /cata directory (because the repo name is cata)
export const REPO_NAME = 'cata';

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];

export const LOCAL_STORAGE_PREFIX = '__cata';

export enum SortDirection {
	ASC,
	DESC,
}
