## TODO : support windows build on drone.io

SHELL:=/bin/bash -O extglob
APP = stock-manager
VERSION = 0.2.0

MAIN_DIR = cmd/main
MAIN_FILE = ${MAIN_DIR}/main.go

REL_DIR = release

EXE = ${APP}-${VERSION} 
EXE_PATH = ${REL_DIR}/${EXE}
ZIP = $(addsuffix .zip,${EXE})
ZIP_PATH = ${REL_DIR}/${ZIP}
TRANS = *.all.yaml
TRANS_PATH = ${MAIN_DIR}/${TRANS}

rm-app : ${ZIP_PATH}
	cd ${REL_DIR} && \
	shopt extglob && \
	rm -r !(*.zip) && \
	shopt -u extglob

${REL_DIR} :
	mkdir ${REL_DIR}

${EXE_PATH} : ${REL_DIR}
	go build -o ${EXE_PATH} ${MAIN_FILE}

run-app : ${EXE_PATH}
	cp ${TRANS_PATH} ${REL_DIR} && \
	cd ${REL_DIR} && \
	./${EXE}
 
${ZIP_PATH} : ${REL_DIR} ${EXE_PATH} run-app
	shopt extglob && \
	cd ${REL_DIR} && \
	zip -r ${ZIP} !(*.zip) && \
	shopt -u extglob

#.PHONY : clean
#clean:
#	rm ${}
