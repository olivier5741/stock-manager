## TODO : support windows build on drone.io

SHELL:=/bin/bash -O extglob
APP = stock-manager
VERSION = 0.2.0

MAIN_DIR = cmd/main
MAIN_FILE = ${MAIN_DIR}/c-main.go

REL_DIR = release

EXE = ${APP}-${VERSION} 
EXE_PATH = ${REL_DIR}/${EXE}
ZIP = $(addsuffix .zip,${EXE})
ZIP_PATH = ${REL_DIR}/${ZIP}
TRANS = *.all.yaml
TRANS_PATH = ${MAIN_DIR}/c-int/${TRANS}

rm-app : ${ZIP_PATH}
	cd ${REL_DIR} && \
	shopt extglob && \
	rm -r !(*.zip) && \
	shopt -u extglob

${REL_DIR} :
	mkdir ${REL_DIR} && \
	mkdir ${REL_DIR}/c-int

${EXE_PATH} : ${REL_DIR}
	go build -o ${EXE_PATH} ${MAIN_FILE}

run-app : ${EXE_PATH}
	cp cmd/main/2016-01-29-n°1-sortie.csv ${REL_DIR} && \
	cp cmd/main/2016-01-28-n°2-entrée.csv ${REL_DIR} && \
	cp cmd/main/2016-01-28-n°1-entrée.csv ${REL_DIR} && \
	cp ${TRANS_PATH} ${REL_DIR}/c-int && \
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
