SHELL:=/bin/bash -O extglob
APP = stock-manager
VERSION = 0.2.0

MAIN_DIR = cmd/main
MAIN_FILE = ${MAIN_DIR}/c-main.go

REL_DIR = release

EXE = e-${APP}-${VERSION} 
EXE_PATH = ${REL_DIR}/${EXE}
WIN_EXE_PATH = $(addsuffix .exe, ${EXE_PATH})
ZIP = $(addsuffix .zip,${APP}-${VERSION}-windows-386 )
ZIP_PATH = ${REL_DIR}/${ZIP}
TRANS = *.all.yaml
TRANS_PATH = ${MAIN_DIR}/c-int/${TRANS}
TRANS_DIR = ${REL_DIR}/c-int

rm-app : ${ZIP_PATH}
	cd ${REL_DIR} && \
	shopt extglob && \
	rm -r !(*.zip) && \
	shopt -u extglob

${REL_DIR} :
	mkdir ${REL_DIR}

${TRANS_DIR} : ${REL_DIR}
	mkdir ${TRANS_DIR}

${EXE_PATH} : ${REL_DIR}
	go build -o ${EXE_PATH} ${MAIN_FILE}

${WIN_EXE_PATH} : ${REL_DIR}
	GOOS=windows GOARCH=386	go build -o ${WIN_EXE_PATH} ${MAIN_FILE}

run-app : ${EXE_PATH} ${REL_DIR} ${TRANS_DIR}
	cp cmd/main/2016-01-29-n1-sortie.csv ${REL_DIR} && \
	cp cmd/main/2016-01-28-n2-entree.csv ${REL_DIR} && \
	cp cmd/main/2016-01-28-n1-entree.csv ${REL_DIR} && \
	cp ${TRANS_PATH} release/c-int && \
	cd ${REL_DIR} && \
	./${EXE} && \
	rm ${EXE}
 
${ZIP_PATH} : ${REL_DIR} ${EXE_PATH} ${WIN_EXE_PATH} run-app
	shopt extglob && \
	cd ${REL_DIR} && \
	zip -r ${ZIP} !(*.zip) && \
	shopt -u extglob
