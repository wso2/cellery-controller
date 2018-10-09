package org.wso2.vick.sts.core;

public class VickStsException extends Exception {

    VickStsException(String msg, Throwable e) {

        super(msg, e);
    }

    public VickStsException(Throwable e) {

        super(e);
    }
}
