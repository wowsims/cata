import { readMessageOption } from '@protobuf-ts/runtime';

import { ProtoVersion } from '../proto/common';

export enum Phase {
	Phase1 = 1,
	Phase2,
	Phase3,
	Phase4,
}

export const CURRENT_PHASE = Phase.Phase4;

export const CURRENT_API_VERSION: number = readMessageOption(ProtoVersion, 'proto.current_version_number')! as number;

// Github pages serves our site under the /cata directory (because the repo name is cata)
export const REPO_NAME = 'cata';
export const REPO_URL = `https://github.com/wowsims/${REPO_NAME}`;
export const REPO_RELEASES_URL = `${REPO_URL}/releases`;
export const REPO_NEW_ISSUE_URL = `${REPO_URL}/issues/new`;
export const REPO_CHOOSE_NEW_ISSUE_URL = `${REPO_NEW_ISSUE_URL}/choose`;

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window?.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];

export const LOCAL_STORAGE_PREFIX = '__cata';

export enum SortDirection {
	ASC,
	DESC,
}
