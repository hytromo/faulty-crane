import { buildNPushImages } from './docker';
import { RepoAndImageEntry, sleep, TagKey, config, CONTAINER_REGISTRY_URL } from './local_config';

const main = async () => {
	console.log(`Generating images for registry ${CONTAINER_REGISTRY_URL}`)

	// only one image is to be kept due to age, let's upload it later so it's fresh and won't be deleted
	const olderImagesBuildAndPush: RepoAndImageEntry[] = []
	const newerImagesBuildAndPush: RepoAndImageEntry[] = []
	for (const repo of config) {
		for (const _tagTypeKey in repo.tags) {
			const tagTypeKey = _tagTypeKey as TagKey
			const allTags = repo.tags[tagTypeKey];
			if (tagTypeKey !== 'ageKept') {
				if (allTags) {
					for (const imageTag of allTags) {
						olderImagesBuildAndPush.push({
							repoName: repo.name,
							imageTag,
						})
					}
				}
			} else {
				if (allTags) {
					for (const imageTag of allTags) {
						newerImagesBuildAndPush.push({
							repoName: repo.name,
							imageTag,
						})
					}
				}
			}
		}
	}

	buildNPushImages(olderImagesBuildAndPush)

	await sleep(60 * 2 * 1000)

	buildNPushImages(newerImagesBuildAndPush)
}

main()

