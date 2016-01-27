BINARY = stock-manager
VERSION = 0.2.0

MAIN_DIR = cmd/main
OUT_DIR = release
TRANS_FILE = *.all.yaml

FILENAME = ${BINARY}-${VERSION}
RELEASE_DIR = ${OUT_DIR}/${VERSION}

MKDIR_P = mkdir -p
CD = cd
CP = cp

all:
	${MKDIR_P} ${OUT_DIR}
	${MKDIR_P} ${RELEASE_DIR}
	${CP} ${MAIN_DIR}/${TRANS_FILE} ${RELEASE_DIR}
	go build -o ${RELEASE_DIR}/${FILENAME} ${MAIN_DIR}/main.go
	${CD} ${RELEASE_DIR} && \
	./${FILENAME}
	zip -j ${OUT_DIR}/${FILENAME}.zip ${RELEASE_DIR}/*
