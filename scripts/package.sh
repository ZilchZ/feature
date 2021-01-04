#!/bin/bash

CUR_PATH=`pwd`
echo "current path: CUR_PATH"

COMMIT_ID=`git rev-parse --short HEAD`
echo "commit id:$COMMIT_ID"

POSTFIX=''

SERVER=127.0.0.1
USER=root
CODE_PATH=/opt
WORK_PATH=/opt
LOG_PATH=/var/log
BUILD_ARCH=x86_64

function package() {
    target=$1
    echo "target:$target"
    if [[ "$target" == "" ]]; then
        exit 1
    fi

    env=$2
    app=${target}${POSTFIX}
    vpath=${app}-${COMMIT_ID}
    vfpath=${CUR_PATH}/objs/${vpath}
    tomlDest=${vfpath}/conf/conf.toml
    rm -rf ${vfpath}
    mkdir -p ${vfpath} ${vfpath}/conf
    cp -av ${CUR_PATH}/${target} ${vfpath}/
    if [[ "$env" == "prd" ]]; then
        cp -av ${CUR_PATH}/conf/${target}.toml ${tomlDest}
    else
        cp -av ${CUR_PATH}/conf/conf.toml ${vfpath}/conf
    fi

    cp -av ${CUR_PATH}/conf/${target}.logrotate ${vfpath}/conf/
    cp -av ${CUR_PATH}/conf/${target}.service ${vfpath}/conf/
    chmod 644 ${vfpath}/conf/*

    cd ${CUR_PATH}/objs
    tar czvf ${vpath}.tar.gz ${vpath}
    rm -rf ${vpath}
    echo "package $target to end"
    echo "ls -lh $CUR_PATH/objs"

    printPktCmd
}

function printPktCmd() {
    echo "package commands:"
    echo "
        cd $CODE_PATH
        mv ~/$vpath.tar.gz .
        tar xzvf $vpath.tar.gz
        rm -rf $vpath.tar.gz
        ls -lh $WORK_PATH/$app
        ln -sfnv $CODE_PATH/$vpath $WORK_PATH/$app
        cd $WORK_PATH/$app
        cp -avf ./conf/$target.logrotate /etc/logrotate.d/$target
        cp -avf ./conf/$target.service /usr/lib/systemd/system/

        systemctl restart $target
    "
}

function init() {
    ssh -t ${USER}@${SERVER} \
    "
    mkdir -p $CODE_PAHT $WORK_PATH $LOG_PAHT/$TARGET
    "
}

function release() {
    target=$1
        echo "target:$target"
        if [[ "$target" == "" ]]; then
            exit 1
        fi

        app = ${target}${POSTFIX}
        vpath=${app}-${COMMIT_ID}
        vfpath=${CUR_PATH}/objs/${vpath}

        init

        scp ${vfpath}.tar.gz ${USER}@${SERVER}:${CODE_PATH}
        ssh -t ${USER}@${SERVER} \
        "
        cd $CODE_PATH
        tar xzvf $vpath.tar.gz
        ls -lh $WORK_PATH/$app
        ln -sfnv $CODE_PATH/$vpath $WORK_PATH/$app
        cd $WORK_PATH/$app
        cp -avf ./conf/$target.logrotate /etc/logrotate.d/$target
        cp -avf ./conf/$target.service /usr/lib/systemd/system/

        systemctl restart $target
        systemctl status $target
        "
}

function rpmbuild() {
    cd ${CUR_PATH}/objs
    app=$1
    release=$2
    target=$3
    if [[ "target" == "" ]]; then
        target=x86_64
    fi
    buildroot=${app}-${release}.${target}
    mkdir -p ${buildroot}/opt/${app} ${buildroot}/opt/${app}/conf ${buildroot}/usr/lib/systemd/system ${buildroot}/etc/logrotate.d
    rm -rf ${app}-${COMMIT_ID}
    tar xzvf ${app}-${COMMIT_ID}.tar.gz
    cp -av ${app}-${COMMIT_ID}/${app} ${buildroot}/opt/${app}/
    cp -av ${app}-${COMMIT_ID}/conf/conf.toml ${buildroot}/opt/${app}/conf

    cp -av ${app}-${COMMIT_ID}/conf/${app}.service ${buildroot}/usr/lib/systemd/system/
    cp -av ${app}-${COMMIT_ID}/conf/${app}.logrotate ${buildroot}/etc/logrotate.d/${app}

    cp -av ../conf/${app}.spec .

    echo "rpmbuild -vv --buildroot `pwd`/$buildroot --target $target --bb objs/$app.spec"

}

case "$1" in
    package)
        package $2 $3;;
    release)
        release $2;;
    rpmbuild)
        rpmbuild $2 $3 $4;;
    *)
        echo "usage: $0 package|release|rpmbuild";;

esac