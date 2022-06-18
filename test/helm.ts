import { exec, getImageName } from "./docker"
import path from 'path';
import fs from 'fs';
import yaml from 'js-yaml'
import { TagKey, config } from "./local_config"

const getUsedImages = () => {
	const usedImages: string[] = []
	for (const repo of config) {
		for (const _tagTypeKey in repo.tags) {
			const tagTypeKey = _tagTypeKey as TagKey
			const allTags = repo.tags[tagTypeKey]
			if (tagTypeKey === 'usedKept' && allTags) {
				usedImages.push(...allTags.map(t => getImageName(repo.name, t)))
			}
		}
	}

	if (usedImages.length != 9) {
		throw new Error('Exactly 9 produced images need to be used inside the k8s cluster');
	}

	return usedImages;
}

/**
 * 
 * @param iteration first iteration = first set of images, 2nd iteration = second set of images
 */
export const generateHelmValuesFile = (iteration = 1) => {
	const usedImages = getUsedImages()
	const helmChartValues = yaml.dump({
		"namespaces": [
			"team1",
			"team2"
		],
		"pods": [
			{
				"namespace": "team1",
				"name": "pod1",
				"image": usedImages[0]
			}
		],
		"deployments": [
			{
				"namespace": "team1",
				"name": "dep1",
				"image": usedImages[2]
			},
			{
				"namespace": "team2",
				"name": "dep2",
				"image": usedImages[3]
			}
		],
		"cronjobs": [
			{
				"namespace": "team1",
				"name": "cron1",
				"image": usedImages[4]
			}
		],
		"statefulsets": [
			{
				"namespace": "team2",
				"name": "ss1",
				"image": usedImages[5]
			}
		],
		"replicasets": [
			{
				"namespace": "team1",
				"name": "rs1",
				"image": usedImages[6]
			}
		]
	})

	fs.writeFileSync(path.join(__dirname, 'helm', 'tester', 'values.yaml'), helmChartValues)
}

export const installHelmChart = async () => {
	await exec('helm upgrade --wait --install tester ./helm/tester')
}