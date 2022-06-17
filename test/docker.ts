import Promise from 'bluebird'
import fs from 'fs';
import os from 'os';
import path from 'path';
import util from 'util';
import { CONTAINER_REGISTRY_REPO_PREFIX, CONTAINER_REGISTRY_URL, RepoAndImageEntry } from "./local_config";
const exec = util.promisify(require('child_process').exec);
import Dockerode from 'dockerode';
const replaceInFile = require('replace-in-file');
const docker = new Dockerode();

const getImageName = (repoName: string, imageTag: string) => `${CONTAINER_REGISTRY_URL}${CONTAINER_REGISTRY_REPO_PREFIX}/${repoName}:${imageTag}`

export const buildNPushImages = async (entries: RepoAndImageEntry[]) => Promise.map(entries, async ({ repoName, imageTag }) => {
	// copy and template the dockerfile to produce a unique docker image
	const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'faulty-crane-integration-testing'));
	const sourceDockerfileTemplate = path.join(__dirname, 'Dockerfile.tpl')
	const destDockerfile = path.join(tmpDir, 'Dockerfile')
	fs.copyFileSync(sourceDockerfileTemplate, destDockerfile)
	await replaceInFile({
		files: [destDockerfile],
		from: /RANDOM_REPLACE/g,
		to: `${new Date().toISOString()}-${Math.random()}`
	})
	const imageName = getImageName(repoName, imageTag)

	console.log('building', imageName)

	const stream = await docker.buildImage({
		context: tmpDir,
		src: ['Dockerfile']
	}, { t: imageName });

	await new Promise((resolve, reject) => {
		docker.modem.followProgress(stream, (err: any, res: any) => err ? reject(err) : resolve(res));
	});

	// dockerode makes it harder than it should to just push an image, so let's just hack it
	await exec(`docker push ${imageName}`);

	return imageName
}, { concurrency: 5 })
