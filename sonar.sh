#!/bin/bash

sonar-scanner -X \
    -Dsonar.host.url=http://quality1.ecd.axway.int \
    -Dsonar.language=go \
    -Dsonar.projectName=WebMethods_Agenta \
    -Dsonar.projectVersion=1.0 \
    -Dsonar.projectKey=WebMethodsAgents \
    -Dsonar.sourceEncoding=UTF-8 \
    -Dsonar.projectBaseDir=${WORKSPACE} \
    -Dsonar.sources=. \
    -Dsonar.tests=. \
    -Dsonar.exclusions=**/testdata/**,**/*.json,**/definitions.go,**/errors.go \
    -Dsonar.test.inclusions=**/*test*.go \
    -Dsonar.go.tests.reportPaths=goreport.json \
    -Dsonar.go.coverage.reportPaths=gocoverage.out \
    -Dsonar.issuesReport.console.enable=true \
    -Dsonar.report.export.path=sonar-report.json

exit 0
