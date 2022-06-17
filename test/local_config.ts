export const sleep = (ms: number) => new Promise(res => setTimeout(res, ms));
export const readEnv = (key: string, defaultValue = '') => {
	const value = process.env[key];

	if (defaultValue === undefined && !value) {
		throw new Error(`Required env variable ${key} not set`)
	}

	return value || defaultValue;
}

export type RepoAndImageEntry = { repoName: string; imageTag: string; };
export type TagKey = 'ageDeleted' |
	'whitelistedTagKept' |
	'whitelistedDigestKept' |
	'usedKept' |
	'whitelistedRepoKept' |
	'oneOfFewKept' |
	'ageKept';
export type ConfigEntry = {
	name: string
	tags: { [tag in TagKey]?: string[] }
}
const randomWords = require('random-words');

const getRandomRepoName = () => `faulty-crane-integration-testing-${readEnv('GITHUB_RUN_ID')}-${randomWords({ exactly: 1, join: '' })}`

export const config: ConfigEntry[] = [
	{
		name: getRandomRepoName(),
		tags: {
			ageDeleted: ['will-be-deleted-due-to-age'],
			whitelistedTagKept: ['will-be-kept-due-to-whitelisted-tag'],
			whitelistedDigestKept: ['will-be-kept-due-to-whitelisted-digest'],
			usedKept: ['will-be-kept-due-to-usage1', 'will-be-kept-due-to-usage2']
		}
	},
	{
		name: getRandomRepoName(),
		tags: {
			whitelistedRepoKept: ['will-be-kept-due-to-whitelisted-repo1', 'will-be-kept-due-to-whitelisted-repo2']
		}
	},
	{
		name: getRandomRepoName(),
		tags: {
			oneOfFewKept: ['will-be-kept-due-to-being-one-of-few-in-repo']
		}
	},
	{
		name: getRandomRepoName(),
		tags: {
			usedKept: ['will-be-kept-due-to-usage1', 'will-be-kept-due-to-usage2']
		}
	},
	{
		name: getRandomRepoName(),
		tags: {
			usedKept: ['will-be-kept-due-to-usage1', 'will-be-kept-due-to-usage2']
		}
	},
	{
		name: getRandomRepoName(),
		tags: {
			usedKept: ['will-be-kept-due-to-usage1', 'will-be-kept-due-to-usage2', 'will-be-kept-due-to-usage3'],
			ageKept: ['will-be-kept-due-to-age']
		}
	},
]

export const abc = 'abddd'

export const CONTAINER_REGISTRY_URL = readEnv('CONTAINER_REGISTRY_URL', 'docker.io')
export const CONTAINER_REGISTRY_REPO_PREFIX = {
	'docker.io': `/${readEnv('CONTAINER_REGISTRY_USERNAME')}`
}[CONTAINER_REGISTRY_URL]