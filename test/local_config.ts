import { exec, getImageName } from "./docker";
import fs from 'fs'
import path from "path";
const randomWords = require('random-words');

export const sleep = (ms: number) => new Promise(res => setTimeout(res, ms));
export const readEnv = (key: string, defaultValue: string | undefined = undefined): string => {
	const value = process.env[key];

	if (defaultValue === undefined && !value) {
		console.error(new Error(`Required env variable ${key} not set`))
		process.exit(1)
	}

	return value || defaultValue || '';
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

const getRandomRepoName = () => `faulty-crane-integration-testing-${readEnv('GITHUB_RUN_ID')}-${randomWords({ exactly: 1, join: '' })}`

export const config: ConfigEntry[] = process.env.FCTESTING_STATIC_CONFIG ? JSON.parse(process.env.FCTESTING_STATIC_CONFIG) : [
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

export const CONTAINER_REGISTRY_URL = readEnv('CONTAINER_REGISTRY_URL', 'docker.io')
export const CONTAINER_REGISTRY_REPO_PREFIX = {
	'docker.io': `/${readEnv('CONTAINER_REGISTRY_USERNAME')}`
}[CONTAINER_REGISTRY_URL]

export const createFaultyCraneConfiguration = async () => {
	process.env.FAULTY_CRANE_CONTAINER_REGISTRY_USERNAME = readEnv('CONTAINER_REGISTRY_USERNAME')
	process.env.FAULTY_CRANE_CONTAINER_REGISTRY_PASSWORD = readEnv('CONTAINER_REGISTRY_PASSWORD')

	const whitelistedRepos: string[] = []
	const whitelistedDigests: string[] = []
	for (const repo of config) {
		if ((Object.keys(repo.tags) as TagKey[]).includes('whitelistedRepoKept')) {
			whitelistedRepos.push(repo.name)
		}

		if ('whitelistedDigestKept' in repo.tags) {
			for (const image of repo.tags.whitelistedDigestKept?.map(t => getImageName(repo.name, t)) || []) {
				whitelistedDigests.push(JSON.parse((await exec(`docker image inspect ${image}`)).stdout)[0].RepoDigests[0].split('@')[1])
			}
		}
	}


	const faultyCraneConfig = {
		"Dockerhub": {
			"Namespace": readEnv('CONTAINER_REGISTRY_USERNAME')
		},
		"Keep": {
			"YoungerThan": "2m",
			"AtLeast": 1,
			"UsedIn": {
				"KubernetesClusters": [
					{
						"Context": "minikube",
						"Namespace": "",
						"RunningInside": false
					}
				]
			},
			"Image": {
				"Tags": [
					"will-be-kept-due-to-whitelisted-tag"
				],
				"Digests": whitelistedDigests,
				"Repositories": whitelistedRepos
			}
		}
	}

	fs.writeFileSync(path.join(__dirname, 'faulty-crane.json'), JSON.stringify(faultyCraneConfig, null, 2))

	return faultyCraneConfig
}