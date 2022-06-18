import { buildAndPushConfigImages, exec, getImageName } from './docker';
import { generateHelmValuesFile, installHelmChart } from './helm';
import { config, CONTAINER_REGISTRY_URL, createFaultyCraneConfiguration, sleep, TagKey } from './local_config';

const main = async () => {
	// generate the random image names and populate them into the config
	generateHelmValuesFile()
	// build those same images and push them to the registry
	await buildAndPushConfigImages()
	// now that the images exist we can build the faulty crane configuration (we need the images built to extract their digest for the configuration, for example)
	await createFaultyCraneConfiguration()
	// install the helm chart referencing those images
	await installHelmChart()

	const { stdout } = await exec('kubectl get pods -A')
	console.log(stdout)
}

main()

