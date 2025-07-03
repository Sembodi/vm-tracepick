help:
	@echo 'Build targets:'
	@echo 'prep'
	@echo 'maintest'
	@echo	'artifactgenerator (step 1): generate artifacts from VM'
	@echo	'merge (step 2): merge docker repos'
	@echo	'builddocker (step 3): build docker image from artifacts'
	@echo 'rundocker (step 4): deploy containers from created docker images'
	@echo 'removedocker: stop and delete containers'
	@echo 'export: export saved metrics for data processing'
	@echo 'csv: make CSV from exported yaml files'
	@echo ''
	@echo 'clean: wipe make binaries bin/'
	@echo ''
	@echo '(Visualize GoProject:)'
	@echo 'visualize'

prep:
	mkdir bin

maintest:
	go build -o bin/maintest cmd/maintest/maintest.go

artifactgenerator:
	go build -o bin/artifactgenerator cmd/artifactgenerator/artifactgenerator.go

merge:
	go build -o bin/mergedockerfiles cmd/docker/mergedockerfiles.go

builddocker:
	go build -o bin/builddockerimage cmd/docker/builddockerimage.go

rundocker:
	go build -o bin/rundocker cmd/docker/rundocker.go

removedocker:
	go build -o bin/removedocker cmd/docker/removedocker.go

export:
	go build -o bin/exportmetrics cmd/metricexport/metricexport.go

csv:
	go build -o bin/csv cmd/metricexport/makecsv.go

clean:
	rm -r bin

visualize:
	go build -o bin/visualizeGoProject visualizeGoProject/visualizeGoProject.go
