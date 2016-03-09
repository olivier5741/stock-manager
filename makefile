SHELL:=/bin/bash -O extglob
APP = woutstock
VERSION = 0.3.0

MAIN_DIR = cmd/main
MAIN_FILE = ${MAIN_DIR}/main.go

REL_DIR = release

RICKY_DIR = test-ricky
MICKY_DIR = test-micky
DASHBOARD_DATA_DIR = données_dashboard

EXE = ${APP}-${VERSION} 
EXE_PATH = ${REL_DIR}/${EXE}
WIN_EXE_PATH = $(addsuffix .exe, ${EXE_PATH})
ZIP = $(addsuffix .zip,${APP}-${VERSION}-windows-386 )
ZIP_PATH = ${REL_DIR}/${ZIP}

rm-app : ${ZIP_PATH}
	cd ${REL_DIR} && \
	shopt extglob && \
	rm -r !(*.zip) && \
	shopt -u extglob

${RICKY_DIR} :
	mkdir -p ${REL_DIR}/${RICKY_DIR}/${DASHBOARD_DATA_DIR} 

${MICKY_DIR} :
	mkdir -p ${REL_DIR}/${MICKY_DIR}/${DASHBOARD_DATA_DIR} 

${EXE_PATH} : ${RICKY_DIR} ${MICKY_DIR}
	go build -o ${EXE_PATH} ${MAIN_FILE}

${WIN_EXE_PATH} : ${RICKY_DIR} ${MICKY_DIR}
	GOOS=windows GOARCH=386	go build -o ${WIN_EXE_PATH} ${MAIN_FILE}

run-app : ${EXE_PATH} ${RICKY_DIR} ${MICKY_DIR}
	cp cmd/main/test-ricky/2016-02-17-n1-produit_màj.csv ${REL_DIR}/${RICKY_DIR} && \
	cp cmd/main/test-ricky/2016-02-18-n1-entree.csv ${REL_DIR}/${RICKY_DIR} && \
	cp cmd/main/test-ricky/2016-02-18-n2-sortie.csv ${REL_DIR}/${RICKY_DIR} && \
	cp cmd/main/test-ricky/2016-02-18-n3-inventaire.csv ${REL_DIR}/${RICKY_DIR} && \
	cd ${REL_DIR} && \
	./${EXE} && \
	rm ${EXE}
 
${ZIP_PATH} : ${RICKY_DIR} ${MICKY_DIR} ${EXE_PATH} ${WIN_EXE_PATH} run-app
	shopt extglob && \
	cd ${REL_DIR} && \
	zip -r ${ZIP} !(*.zip) && \
	shopt -u extglob
